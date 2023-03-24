package collection

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/sirupsen/logrus"
)

// populatePreMarketDailyStats - corrects the daily high, low, and volume for the current day if in pre-market
func (c *Collection) populatePreMarketDailyStats() {
	calendarDayUpdate := c.calendarDayLive.Get()

	if calendarDayUpdate.Condition != calendars.PreMarket {
		return
	}

	start := calendarDayUpdate.PreviousMarketDate.AsTime()

	multiBars, err := c.mdClient.GetMultiBars(
		c.symbols, md.GetBarsRequest{
			TimeFrame:  md.OneDay,
			Adjustment: "split",
			Start:      start,
			End:        start,
			TotalLimit: 1,
			Feed:       md.SIP,
		},
	)
	if err != nil {
		logrus.Panicf("failed to get daily bars: %v", err)
	}

	c.collection.IterCb(
		func(symbol string, snapshot *snapshots.Snapshot) {
			if intradayBars, ok := multiBars[symbol]; ok {
				var high float64
				var low float64
				var volume float64

				for _, bar := range intradayBars {
					if bar.High > high {
						high = bar.High
					}

					if bar.Low < low {
						low = bar.Low
					}

					volume += float64(bar.Volume)
				}

				snapshot.DailyHigh = high
				snapshot.DailyLow = low
				snapshot.DailyVolume = volume
			}
		},
	)
}
