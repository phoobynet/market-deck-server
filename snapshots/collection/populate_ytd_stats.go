package collection

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/sirupsen/logrus"
)

func (c *Collection) populateYTDStats() {
	startOfDay := carbon.NewCarbon().SetTimezone(carbon.NewYork).StartOfDay()

	ytdBarsMonthly, err := c.mdClient.GetMultiBars(
		c.symbols, md.GetBarsRequest{
			TimeFrame:  md.OneMonth,
			Adjustment: md.Split,
			Start:      startOfDay.SubYear().ToStdTime(),
			End:        startOfDay.ToStdTime(),
			Feed:       md.SIP,
		},
	)

	if err != nil {
		logrus.Panicf("failed to get last month of bars: %v", err)
	}

	c.collection.IterCb(
		func(symbol string, snapshot *snapshots.Snapshot) {
			if barsForSymbol, ok := ytdBarsMonthly[symbol]; ok {
				snapshot.YtdBars = bars.FromMarketDataBars(symbol, barsForSymbol)
				snapshot.YtdChange = numbers.NumberDiff(snapshot.YtdBars[0].Close, snapshot.Price)
			}
		},
	)
}
