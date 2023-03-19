package calendars

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/market-deck-server/helpers/date"
)

type CalendarDay struct {
	Date            string `gorm:"primaryKey" json:"date"`
	PreMarketOpen   int64  `json:"preMarketOpen"`
	Open            int64  `json:"open"`
	Close           int64  `json:"close"`
	PostMarketClose int64  `json:"postMarketClose"`
}

func NewMarketCalendarDay(calendarDay alpaca.CalendarDay) CalendarDay {
	preMarketOpen := carbon.Parse(fmt.Sprintf("%s %s:00", calendarDay.Date, "04:00"), date.MarketTimeZone)
	openingTime := carbon.Parse(fmt.Sprintf("%s %s:00", calendarDay.Date, calendarDay.Open), date.MarketTimeZone)
	closingTime := carbon.Parse(fmt.Sprintf("%s %s:00", calendarDay.Date, calendarDay.Close), date.MarketTimeZone)
	postMarketClose := carbon.Parse(fmt.Sprintf("%s %s:00", calendarDay.Date, "20:00"), date.MarketTimeZone)

	return CalendarDay{
		Date:            calendarDay.Date,
		PreMarketOpen:   preMarketOpen.ToStdTime().UnixMicro(),
		Open:            openingTime.ToStdTime().UnixMicro(),
		Close:           closingTime.ToStdTime().UnixMicro(),
		PostMarketClose: postMarketClose.ToStdTime().UnixMicro(),
	}
}
