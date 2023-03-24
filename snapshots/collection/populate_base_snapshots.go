package collection

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/sirupsen/logrus"
)

func (c *Collection) populateBaseSnapshots(symbols []string) {
	if len(symbols) == 0 {
		return
	}

	multiSnapshots, err := c.mdClient.GetSnapshots(symbols, md.GetSnapshotRequest{})

	if err != nil {
		logrus.Panicf("failed to load base snapshots: %v", err)
	}

	calendarDayLive := c.calendarDayLive.Get()

	var latestTrades map[string]md.Trade

	if calendarDayLive.Condition == calendars.PreMarket {
		latestTrades, err = c.mdClient.GetLatestTrades(
			symbols, md.GetLatestTradeRequest{
				Feed: md.SIP,
			},
		)

		if err != nil {
			logrus.Panicf("failed to load latest trades: %v", err)
		}
	}

	for symbol, snapshot := range multiSnapshots {
		asset := c.assetRepo.Get(symbol)

		prevDailyBar := snapshot.PrevDailyBar

		if calendarDayLive.PreviousMarketDate.Date == snapshot.DailyBar.Timestamp.Format("2006-01-02") {
			prevDailyBar = snapshot.DailyBar
		}

		latestPrice := snapshot.LatestTrade.Price

		if calendarDayLive.Condition == calendars.PreMarket {
			if latestTrade, ok := latestTrades[symbol]; ok {
				latestPrice = latestTrade.Price
			}
		}

		c.collection.Set(
			symbol, &snapshots.Snapshot{
				Symbol:        symbol,
				Name:          asset.Name,
				Exchange:      asset.Exchange,
				Class:         asset.Class,
				Price:         snapshot.LatestTrade.Price,
				PrevClose:     prevDailyBar.Close,
				PrevCloseDate: prevDailyBar.Timestamp.Format("2006-01-02"),
				DailyHigh:     snapshot.DailyBar.High,
				DailyLow:      snapshot.DailyBar.Low,
				DailyVolume:   float64(snapshot.DailyBar.Volume),
				Change:        numbers.NumberDiff(prevDailyBar.Close, latestPrice),
			},
		)
	}
}
