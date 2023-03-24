package collection

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/sirupsen/logrus"
)

func (c *Collection) populateBaseSnapshots(symbols []string) {
	multiSnapshots, err := c.mdClient.GetSnapshots(symbols, md.GetSnapshotRequest{})

	if err != nil {
		logrus.Panicf("failed to load base snapshots: %v", err)
	}

	for symbol, snapshot := range multiSnapshots {
		asset := c.assetRepo.Get(symbol)

		c.collection.Set(
			symbol, &snapshots.Snapshot{
				Symbol:        symbol,
				Name:          asset.Name,
				Exchange:      asset.Exchange,
				Class:         asset.Class,
				Price:         snapshot.LatestTrade.Price,
				PrevClose:     snapshot.PrevDailyBar.Close,
				PrevCloseDate: snapshot.PrevDailyBar.Timestamp.Format("2006-01-02"),
				DailyHigh:     snapshot.DailyBar.High,
				DailyLow:      snapshot.DailyBar.Low,
				DailyVolume:   float64(snapshot.DailyBar.Volume),
				Change:        numbers.NumberDiff(snapshot.PrevDailyBar.Close, snapshot.LatestTrade.Price),
			},
		)
	}
}
