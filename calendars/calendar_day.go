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
		PreMarketOpen:   preMarketOpen.ToStdTime().UnixMicro(),
		Open:            openingTime.ToStdTime().UnixMicro(),
		Close:           closingTime.ToStdTime().UnixMicro(),
		PostMarketClose: postMarketClose.ToStdTime().UnixMicro(),
	}
}

func (c *CalendarDay) String() string {
	return fmt.Sprintf(
		"Date: %s\n PreMarketOpen: %s\n Open: %s\n Close: %s\n PostMarketClose: %s\n",
		c.Date,
		time.UnixMicro(c.PreMarketOpen).Format("15:04"),
		time.UnixMicro(c.Open).Format("15:04"),
		time.UnixMicro(c.Close).Format("15:04"),
		time.UnixMicro(c.PostMarketClose).Format("15:04"),
	)
}
