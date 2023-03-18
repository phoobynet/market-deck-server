package realtime

import (
	"market-deck/assets"
	"market-deck/bars"
	"market-deck/quotes"
	"market-deck/trades"
)

type Symbol struct {
	Asset        assets.Asset `json:"asset"`
	Bar          bars.Bar     `json:"bar"`
	Trade        trades.Trade `json:"trade"`
	Quote        quotes.Quote `json:"quote"`
	PrevDailyBar bars.Bar     `json:"prevDailyBar"`
	IntradayBars []bars.Bar   `json:"intradayBars"`
	YtdDailyBars []bars.Bar   `json:"ytdDailyBars"`
}
