package calendars

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type CalendarDayRepository struct {
	alpacaClient *alpaca.Client
	populated    bool
	db           *gorm.DB
}

func NewCalendarDayRepository(db *gorm.DB, alpacaClient *alpaca.Client) *CalendarDayRepository {
	return &CalendarDayRepository{
		alpacaClient: alpacaClient,
		populated:    false,
		db:           db,
	}
}

func (r *CalendarDayRepository) populate() {
	if r.populated {
		return
	}

	if r.Count() > 0 {
		r.populated = true
		return
	}

	calendarDays, err := r.alpacaClient.GetCalendar(
		alpaca.GetCalendarRequest{
			Start: time.Now().AddDate(-1, 0, 0),
			End:   time.Now().AddDate(1, 0, 0),
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
	r.populated = true
}

func (r *CalendarDayRepository) Get(date string) CalendarDay {
	var marketDate CalendarDay

	r.db.Where("date = ?", date).First(&marketDate)

	return marketDate
}

func (r *CalendarDayRepository) GetPrevious() *CalendarDay {
	r.populate()

	var marketDate CalendarDay

	r.db.Where("date < ?", time.Now().Format("2006-01-02")).
		Order("date desc").
		First(&marketDate)

	return &marketDate
}

func (r *CalendarDayRepository) GetNext(date string) *CalendarDay {
	r.populate()

	var marketDate CalendarDay

	r.db.Where("date > ?", date).Order("date asc").First(&marketDate)

	return &marketDate
}

func (r *CalendarDayRepository) GetToday() *CalendarDay {
	var marketDate CalendarDay

	r.db.Where("date = ?", time.Now().Format("2006-01-02")).First(&marketDate)

	return &marketDate
}

func (r *CalendarDayRepository) GetBetween(start, end time.Time) []CalendarDay {
	r.populate()

	var marketDate []CalendarDay

	r.db.Where(
		"date >= ? AND date <= ?",
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	).Order("date asc").Find(&marketDate)

	return marketDate
}

func (r *CalendarDayRepository) Count() int64 {
	var count int64

	r.db.Model(&CalendarDay{}).Count(&count)

	return count
}
