package calendars

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/golang-module/carbon/v2"
	"time"
)

type CalendarDay struct {
	Date            string `gorm:"primaryKey" json:"date"`
	PreMarketOpen   int64  `json:"preMarketOpen"`
	Open            int64  `json:"open"`
	Close           int64  `json:"close"`
	PostMarketClose int64  `json:"postMarketClose"`
}

func NewMarketCalendarDay(calendarDay alpaca.CalendarDay) CalendarDay {
	preMarketOpen := carbon.Parse(fmt.Sprintf("%s %s:00", calendarDay.Date, "04:00"), carbon.NewYork)
	openingTime := carbon.Parse(fmt.Sprintf("%s %s:00", calendarDay.Date, calendarDay.Open), carbon.NewYork)
	closingTime := carbon.Parse(fmt.Sprintf("%s %s:00", calendarDay.Date, calendarDay.Close), carbon.NewYork)
	postMarketClose := carbon.Parse(fmt.Sprintf("%s %s:00", calendarDay.Date, "20:00"), carbon.NewYork)

	return CalendarDay{
		Date:            calendarDay.Date,
		PreMarketOpen:   preMarketOpen.ToStdTime().UnixMilli(),
		Open:            openingTime.ToStdTime().UnixMilli(),
		Close:           closingTime.ToStdTime().UnixMilli(),
		PostMarketClose: postMarketClose.ToStdTime().UnixMilli(),
	}
}

func (c *CalendarDay) AsTime() time.Time {
	return carbon.Parse(c.Date, carbon.NewYork).ToStdTime()
}

func (c *CalendarDay) String() string {
	return fmt.Sprintf(
		"Date: %s\n PreMarketOpen: %s\n Open: %s\n Close: %s\n PostMarketClose: %s\n",
		c.Date,
		time.UnixMilli(c.PreMarketOpen).Format("15:04"),
		time.UnixMilli(c.Open).Format("15:04"),
		time.UnixMilli(c.Close).Format("15:04"),
		time.UnixMilli(c.PostMarketClose).Format("15:04"),
	)
}
