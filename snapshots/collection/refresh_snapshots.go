package collection

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/sirupsen/logrus"
	"math"
)

// Refresh gets the latest snapshots from Alpaca and applies them to the collection
func (c *Collection) Refresh() {

	calendarDayLive := c.calendarDayLive.Get()

	if calendarDayLive.Condition == calendars.ClosedToday {
		return
	}

	multiSnapshots, err := c.mdClient.GetSnapshots(c.symbols, md.GetSnapshotRequest{})

	if err != nil {
		logrus.Panicf("failed to load base snapshots: %v", err)
	}

	c.collection.IterCb(
		func(symbol string, snapshot *snapshots.Snapshot) {
			if mdSnapshot, ok := multiSnapshots[symbol]; ok {

				if calendarDayLive.PreviousMarketDate.Date == mdSnapshot.PrevDailyBar.Timestamp.Format("2006-01-02") {
					snapshot.PrevCloseDate = calendarDayLive.PreviousMarketDate.Date
					snapshot.PrevClose = mdSnapshot.PrevDailyBar.Close
				} else {
					snapshot.PrevCloseDate = calendarDayLive.CurrentMarketDate.Date
					snapshot.PrevClose = mdSnapshot.DailyBar.Close
				}

				snapshot.DailyHigh = math.Max(mdSnapshot.LatestTrade.Price, snapshot.Price)
				snapshot.DailyLow = math.Min(mdSnapshot.LatestTrade.Price, snapshot.Price)
				snapshot.DailyVolume = float64(mdSnapshot.DailyBar.Volume)
			}
		},
	)
}
