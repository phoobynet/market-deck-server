package calendars

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/clients"
	"github.com/phoobynet/market-deck-server/database"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

var calendarRepositoryOnce sync.Once

var calendarRepository *Repository

type Repository struct {
	alpacaClient    *alpaca.Client
	populated       bool
	calendarDays    []CalendarDay
	calendarDaysMap cmap.ConcurrentMap[string, CalendarDay]
	calendarDates   []string
	db              *gorm.DB
}

func GetRepository() *Repository {
	calendarRepositoryOnce.Do(
		func() {
			calendarRepository = &Repository{
				alpacaClient:    clients.GetAlpacaClient(),
				populated:       false,
				db:              database.Get(),
				calendarDays:    make([]CalendarDay, 0),
				calendarDaysMap: cmap.New[CalendarDay](),
				calendarDates:   make([]string, 0),
			}
		},
	)

	return calendarRepository
}

// populate - ensure the calendar days are loaded into memory, if not, load them from the database
// if the table is empty, then load them from the alpaca api, then load them into memory
func (r *Repository) populate() {
	if r.populated {
		return
	}

	if r.Count() > 0 && !r.populated {
		var calendarDays []CalendarDay
		r.db.Model(&CalendarDay{}).Find(&calendarDays) // load into memory

		for _, calendarDay := range calendarDays {
			r.calendarDaysMap.Set(calendarDay.Date, calendarDay)
			r.calendarDays = append(r.calendarDays, calendarDay)
			r.calendarDates = append(r.calendarDates, calendarDay.Date)
		}

		r.populated = true
		return
	}

	calendarDays, err := r.alpacaClient.GetCalendar(
		alpaca.GetCalendarRequest{
			Start: time.Now().AddDate(-1, -1, 0),
			End:   time.Now().AddDate(1, 1, 0),
		},
	)

	if err != nil {
		logrus.Panic(err)
	}

	var marketDates []CalendarDay

	for _, calendarDay := range calendarDays {
		marketDate := NewMarketCalendarDay(calendarDay)
		marketDates = append(marketDates, marketDate)
	}

	r.db.Create(&marketDates)

	r.populate()
}

func (r *Repository) Get(date string) CalendarDay {
	var marketDate CalendarDay

	r.db.Where("date = ?", date).First(&marketDate)

	return marketDate
}

func (r *Repository) GetPrevious() *CalendarDay {
	r.populate()

	var marketDate CalendarDay

	r.db.Where("date < ?", time.Now().Format("2006-01-02")).
		Order("date desc").
		First(&marketDate)

	return &marketDate
}

func (r *Repository) GetNext(date string) *CalendarDay {
	r.populate()

	var marketDate CalendarDay

	r.db.Where("date > ?", date).Order("date asc").First(&marketDate)

	return &marketDate
}

func (r *Repository) GetToday() *CalendarDay {
	var marketDate CalendarDay

	r.db.Where("date = ?", time.Now().Format("2006-01-02")).First(&marketDate)

	return &marketDate
}

func (r *Repository) GetBetween(start, end time.Time) []CalendarDay {
	r.populate()

	var marketDate []CalendarDay

	r.db.Where(
		"date >= ? AND date <= ?",
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	).Order("date asc").Find(&marketDate)

	return marketDate
}

func (r *Repository) Count() int64 {
	var count int64

	r.db.Model(&CalendarDay{}).Count(&count)

	return count
}

// PickByIntervals - pick calendar days that are on or after a date relative to now
func (r *Repository) PickByIntervals(daysAgo []int) []CalendarDay {
	r.populate()

	pickedCalendarDays := make([]CalendarDay, 0)

	now := time.Now()

	pickDates := lo.Map[int, string](
		daysAgo, func(days int, _ int) string {
			return now.AddDate(0, 0, -days).Format("2006-01-02")
		},
	)

	for _, pickDate := range pickDates {
		calendarDay, ok := lo.Find[CalendarDay](
			r.calendarDays, func(calendarDay CalendarDay) bool {
				return calendarDay.Date >= pickDate
			},
		)

		if ok {
			pickedCalendarDays = append(pickedCalendarDays, calendarDay)
		} else {
			logrus.Panicf("could not find calendar day for date %s", pickDate)
		}
	}

	return pickedCalendarDays
}
