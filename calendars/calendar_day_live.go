package calendars

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/golang-module/carbon/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/helpers/date"
	"github.com/phoobynet/market-deck-server/messages"
	"sync"
	"time"
)

type CalendarDayLive struct {
	mu                sync.RWMutex
	alpacaClient      *alpaca.Client
	repository        *Repository
	publishTicker     *time.Ticker
	calendarDays      []CalendarDay
	calendarDaysMap   cmap.ConcurrentMap[string, CalendarDay]
	calendarDayUpdate CalendarDayUpdate
	ctx               context.Context
	nyTimezone        carbon.Carbon
}

func NewCalendarDayLive(
	ctx context.Context,
	alpacaClient *alpaca.Client,
	calendarDayRepository *Repository,
	messageBus chan<- messages.Message,
) *CalendarDayLive {
	l := &CalendarDayLive{
		alpacaClient:    alpacaClient,
		publishTicker:   time.NewTicker(1 * time.Second),
		nyTimezone:      date.GetNewYorkZone(),
		calendarDays:    make([]CalendarDay, 0),
		calendarDaysMap: cmap.New[CalendarDay](),
		repository:      calendarDayRepository,
	}

	l.populateMarketDates()

	go func() {
		for {
			select {
			case <-l.publishTicker.C:
				l.update()
				messageBus <- messages.Message{
					Event: messages.CalendarDayUpdate,
					Data:  l.calendarDayUpdate,
				}
			case <-ctx.Done():
				l.publishTicker.Stop()
			}
		}
	}()

	return l
}

func (l *CalendarDayLive) populateMarketDates() {
	end := l.nyTimezone.Now().StartOfDay().AddDays(7)
	start := end.SubDays(14).ToStdTime()

	l.calendarDays = l.repository.GetBetween(start, end.ToStdTime())

	for _, marketDate := range l.calendarDays {
		l.calendarDaysMap.Set(marketDate.Date, marketDate)
	}
}

func (l *CalendarDayLive) update() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := date.GetNewYorkZone().Now()

	nowUtcMicro := now.ToStdTime().UnixMicro()

	dateKey := now.Format("Y-m-d")

	if marketDate, ok := l.calendarDaysMap.Get(dateKey); ok {
		if nowUtcMicro >= marketDate.PostMarketClose {
			l.calendarDayUpdate.Condition = ClosedForTheDay
		} else if nowUtcMicro >= marketDate.Close {
			l.calendarDayUpdate.Condition = PostMarket
		} else if nowUtcMicro >= marketDate.Open {
			l.calendarDayUpdate.Condition = Open
		} else if nowUtcMicro >= marketDate.PreMarketOpen {
			l.calendarDayUpdate.Condition = PreMarket
		} else {

			l.calendarDayUpdate.Condition = ClosedOpeningLater
		}
		l.calendarDayUpdate.CurrentMarketDate = marketDate
	} else {
		l.calendarDayUpdate.Condition = ClosedToday
	}

	var previousMarketDate CalendarDay

	for _, marketDate := range l.calendarDays {
		if marketDate.Date >= dateKey {
			break
		}
		previousMarketDate = marketDate
	}

	l.calendarDayUpdate.PreviousMarketDate = previousMarketDate

	var nextMarketDate CalendarDay

	for _, marketDate := range l.calendarDays {
		if marketDate.Date > dateKey {
			nextMarketDate = marketDate
			break
		}
	}

	l.calendarDayUpdate.NextMarketDate = nextMarketDate
	l.calendarDayUpdate.At = now.ToStdTime().UnixMilli()
}

func (l *CalendarDayLive) Get() CalendarDayUpdate {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.calendarDayUpdate
}
