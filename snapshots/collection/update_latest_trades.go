package collection

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/phoobynet/market-deck-server/trades"
	"math"
)

func (c *Collection) UpdateLatestTrades(latestTrades *cmap.ConcurrentMap[string, *trades.Trade]) {
	c.collection.IterCb(
		func(symbol string, snapshot *snapshots.Snapshot) {
			if latestTrade, ok := latestTrades.Get(symbol); ok {
				snapshot.Price = latestTrade.Price
				snapshot.DailyHigh = math.Max(snapshot.DailyHigh, latestTrade.Price)
				snapshot.DailyLow = math.Min(snapshot.DailyLow, latestTrade.Price)
				snapshot.Change = numbers.NumberDiff(snapshot.PrevClose, snapshot.Price)
				snapshot.YtdChange = numbers.NumberDiff(snapshot.YtdBars[0].Close, latestTrade.Price)
			}
		},
	)
}
