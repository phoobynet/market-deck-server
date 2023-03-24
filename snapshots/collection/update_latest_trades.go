package collection

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/phoobynet/market-deck-server/trades"
	"github.com/sirupsen/logrus"
	"math"
)

func (c *Collection) UpdateLatestTrades(latestTrades cmap.ConcurrentMap[string, *trades.Trade]) {
	c.collection.IterCb(
		func(symbol string, snapshot *snapshots.Snapshot) {
			if latestTrade, ok := latestTrades.Get(symbol); ok {

				dailyLow := snapshot.DailyLow

				if dailyLow == 0 || latestTrade.Price < dailyLow {
					dailyLow = latestTrade.Price
				}

				snapshot.Price = latestTrade.Price
				snapshot.DailyHigh = math.Max(snapshot.DailyHigh, latestTrade.Price)
				snapshot.DailyLow = dailyLow
				snapshot.Change = numbers.NumberDiff(snapshot.PrevClose, latestTrade.Price)

				if len(snapshot.YtdBars) > 0 {
					snapshot.YtdChange = numbers.NumberDiff(snapshot.YtdBars[0].Close, latestTrade.Price)
				} else {
					logrus.Warnf("No YTD bars for %v", symbol)
				}
			}
		},
	)
}
