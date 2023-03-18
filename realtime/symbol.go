package realtime

import (
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
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
