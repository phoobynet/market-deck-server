package snapshots

import (
	"context"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/messages"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Stream struct {
	deckRepository        *decks.DeckRepository
	barStream             *bars.Stream
	quoteStream           *quotes.Stream
	tradeStream           *trades.Stream
	barChan               chan map[string]bars.Bar
	quoteChan             chan map[string]quotes.Quote
	tradeChan             chan map[string]trades.Trade
	snapshotsRepository   *Repository
	publishTicker         *time.Ticker
	publishInterval       time.Duration
	refreshTicker         *time.Ticker
	barRepository         *bars.Repository
	calendarDayRepository *calendars.Repository

	mu sync.RWMutex

	// live and frequently updated data which is used to calculate snapshots
	latestBars     cmap.ConcurrentMap[string, bars.Bar]
	latestTrades   cmap.ConcurrentMap[string, trades.Trade]
	latestQuotes   cmap.ConcurrentMap[string, quotes.Quote]
	dailyBars      cmap.ConcurrentMap[string, bars.Bar]
	prevDailyBars  cmap.ConcurrentMap[string, bars.Bar]
	previousCloses cmap.ConcurrentMap[string, float64]
	intradayBars   cmap.ConcurrentMap[string, []bars.Bar]
	ytdBars        cmap.ConcurrentMap[string, []bars.Bar]

	symbols []string

	ytdPickDaysAgo  []int
	calendarDayLive *calendars.CalendarDayLive
}

func NewSnapshotStream(
	ctx context.Context,
	calendarDayLive *calendars.CalendarDayLive,
	messageBus chan<- messages.Message,
) *Stream {
	s := &Stream{
		deckRepository:        decks.GetRepository(),
		snapshotsRepository:   GetRepository(),
		calendarDayRepository: calendars.GetRepository(),
		calendarDayLive:       calendarDayLive,
		barRepository:         bars.GetRepository(),
		publishInterval:       1 * time.Second,
		publishTicker:         time.NewTicker(1 * time.Second),
		barChan:               make(chan map[string]bars.Bar, 1_000),
		quoteChan:             make(chan map[string]quotes.Quote, 1_000),
		tradeChan:             make(chan map[string]trades.Trade, 1_000),
		intradayBars:          cmap.New[[]bars.Bar](),
		ytdBars:               cmap.New[[]bars.Bar](),
		symbols:               make([]string, 0),
		latestBars:            cmap.New[bars.Bar](),
		latestTrades:          cmap.New[trades.Trade](),
		latestQuotes:          cmap.New[quotes.Quote](),
		dailyBars:             cmap.New[bars.Bar](),
		prevDailyBars:         cmap.New[bars.Bar](),
		previousCloses:        cmap.New[float64](),
		refreshTicker:         time.NewTicker(1 * time.Second),
	}

	s.barStream = bars.NewBarStream(ctx, s.barChan)
	s.tradeStream = trades.NewTradeStream(ctx, s.tradeChan)
	s.quoteStream = quotes.NewQuoteStream(ctx, s.quoteChan)

	s.ytdPickDaysAgo = []int{
		365,
		180,
		90,
		30,
		14,
		7,
	}

	go func(s *Stream) {
		for {
			select {
			case barsMap := <-s.barChan:
				s.latestBars.MSet(barsMap)

				for _, symbol := range s.symbols {
					existingBars, _ := s.intradayBars.Get(symbol)
					latestBar, _ := s.latestBars.Get(symbol)
					s.intradayBars.Set(symbol, append(existingBars, latestBar))
				}
			case quotesMap := <-s.quoteChan:
				s.latestQuotes.MSet(quotesMap)
			case tradesMap := <-s.tradeChan:
				s.latestTrades.MSet(tradesMap)
			case <-s.publishTicker.C:
				s.mu.RLock()

				snapshots := cmap.New[Snapshot]()

				today := time.Now().Format("2006-01-02")

				for _, symbol := range s.symbols {
					latestTrade, _ := s.latestTrades.Get(symbol)
					latestBar, _ := s.latestBars.Get(symbol)
					latestQuote, _ := s.latestQuotes.Get(symbol)
					dailyBar, _ := s.dailyBars.Get(symbol)
					prevDailyBar, _ := s.prevDailyBars.Get(symbol)
					previousClose, _ := s.previousCloses.Get(symbol)
					intradayBars, _ := s.intradayBars.Get(symbol)

					var high float64
					low := 999_999_999.00

					for _, bar := range intradayBars {
						if bar.High > high {
							high = bar.High
						}

						if bar.Low < low && bar.Low > 0 {
							low = bar.Low
						}
					}

					snapshot := Snapshot{
						LatestBar:        latestBar,
						LatestQuote:      latestQuote,
						LatestTrade:      latestTrade,
						DailyBar:         dailyBar,
						PreviousDailyBar: prevDailyBar,
						PreviousClose:    previousClose,
						IntradayBars:     intradayBars,
						IntradayLow:      low,
						IntradayHigh:     high,
					}

					actualPreviousDailyBar := snapshot.PreviousDailyBar

					if snapshot.DailyBar.Date() < today {
						actualPreviousDailyBar = snapshot.DailyBar
					}

					snapshot.ActualPreviousDailyBar = actualPreviousDailyBar

					changes := cmap.New[SnapshotChange]()

					// previous close
					diff := numbers.NumberDiff(snapshot.PreviousClose, snapshot.LatestTrade.Price)

					changes.Set(
						snapshot.ActualPreviousDailyBar.Date(), SnapshotChange{
							Since:         snapshot.ActualPreviousDailyBar.Timestamp,
							Label:         "Prev Close",
							Change:        diff.Change,
							ChangePercent: diff.ChangePercent,
							ChangeAbs:     diff.AbsoluteChange,
							ChangeSign:    diff.Multiplier,
						},
					)

					snapshot.Changes = changes.Items()

					snapshots.Set(symbol, snapshot)
				}

				messageBus <- messages.Message{
					Event: messages.Snapshots,
					Data:  snapshots.Items(),
				}
				s.mu.RUnlock()
			case <-s.refreshTicker.C:
				sec := time.Now().Second()

				if sec == 0 {
					logrus.Info("refreshing snapshots")
					s.mu.Lock()
					s.refreshSnapshot(s.symbols, false)
					s.mu.Unlock()
				}
			case <-ctx.Done():
				s.publishTicker.Stop()
				s.refreshTicker.Stop()
			}
		}
	}(s)

	return s
}

func (s *Stream) UpdateSymbols(symbols []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.publishTicker.Stop()
	defer s.publishTicker.Reset(s.publishInterval)

	if len(symbols) == 0 {
		_, err := s.deckRepository.ClearByName("default")

		if err != nil {
			logrus.Panicf("failed to clear symbols: %v", err)
		}
	} else {

		_, err := s.deckRepository.UpdateByName("default", symbols)

		if err != nil {
			logrus.Panicf("failed to update symbols: %v", err)
		}
	}

	removedSymbols, addedSymbols := lo.Difference(s.symbols, symbols)

	if len(removedSymbols) > 0 {
		s.symbols = lo.Filter(
			s.symbols, func(symbol string, _ int) bool {
				return !lo.Contains(removedSymbols, symbol)
			},
		)

		for _, symbol := range removedSymbols {
			s.latestTrades.Remove(symbol)
			s.latestBars.Remove(symbol)
			s.latestQuotes.Remove(symbol)
			s.latestBars.Remove(symbol)
			s.prevDailyBars.Remove(symbol)
			s.previousCloses.Remove(symbol)
			s.dailyBars.Remove(symbol)
			s.intradayBars.Remove(symbol)
		}
	}

	if len(addedSymbols) > 0 {
		s.symbols = append(s.symbols, addedSymbols...)

		s.refreshSnapshot(addedSymbols, true)

		go s.fillIntradayBars(addedSymbols)
		go s.fillYtdBarsAtIntervals(addedSymbols)
	}

	s.barStream.Update(symbols)
	s.quoteStream.Update(symbols)
	s.tradeStream.Update(symbols)
}

// refreshSnapshot refreshes the daily bar, previous daily bar, and previous close.
// when initializeLatestValues is true, it will also initialize the latest bar, quote, and trade.
func (s *Stream) refreshSnapshot(symbols []string, initializeLatestValues bool) {
	if len(symbols) == 0 {
		return
	}

	snapshots, err := s.snapshotsRepository.GetMulti(symbols)

	if err != nil {
		logrus.Errorf("failed to get snapshots: %v", err)
	}

	for _, symbol := range symbols {
		if initializeLatestValues {
			s.latestBars.Set(symbol, snapshots[symbol].LatestBar)
			s.latestQuotes.Set(symbol, snapshots[symbol].LatestQuote)
			s.latestTrades.Set(symbol, snapshots[symbol].LatestTrade)
		}
		s.dailyBars.Set(symbol, snapshots[symbol].DailyBar)
		s.prevDailyBars.Set(symbol, snapshots[symbol].PreviousDailyBar)
		s.previousCloses.Set(symbol, snapshots[symbol].PreviousClose)
	}
}

func (s *Stream) fillIntradayBars(symbols []string) {
	if len(symbols) == 0 {
		return
	}

	intradayMulti, err := s.barRepository.GetIntradayMulti(symbols)

	if err != nil {
		logrus.Errorf("failed to get intraday bars: %v", err)
	}

	for symbol, intraday := range intradayMulti {
		s.intradayBars.Set(symbol, intraday)
	}
}

func (s *Stream) fillYtdBarsAtIntervals(symbols []string) {
	if len(symbols) == 0 {
		return
	}

	ytdMulti, err := s.barRepository.GetYtdDailyMulti(symbols)

	if err != nil {
		logrus.Errorf("failed to get YTD bars: %v", err)
	}

	pickCalendarDays := s.calendarDayRepository.PickByIntervals(s.ytdPickDaysAgo)

	// For YTD bars, select at intervals
	for symbol, ytdBars := range ytdMulti {
		pickedBars := make([]bars.Bar, 0)

		for _, pickCalendarDay := range pickCalendarDays {
			bar, ok := lo.Find(
				ytdBars, func(bar bars.Bar) bool {
					if bar == (bars.Bar{}) {
						logrus.Panicf("bar is empty for %s on %s", symbol, pickCalendarDay.Date)
					}

					return bar.Date() == pickCalendarDay.Date
				},
			)

			if ok {
				pickedBars = append(pickedBars, bar)
			} else {
				logrus.Errorf("no bar found for %s on %s", symbol, pickCalendarDay.Date)
			}
		}

		if len(pickedBars) > 0 {
			s.ytdBars.Set(symbol, pickedBars)
		}
	}
}

func (s *Stream) GetIntradayBars() map[string][]bars.Bar {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.intradayBars.Items()
}

func (s *Stream) GetYtdBars() map[string][]bars.Bar {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ytdBars.Items()
}
