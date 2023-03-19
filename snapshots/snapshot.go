package snapshots

import (
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
)

type Snapshot struct {
	LatestBar   bars.Bar     `json:"lb"`
	LatestTrade trades.Trade `json:"lt"`
	LatestQuote quotes.Quote `json:"lq"`
}
