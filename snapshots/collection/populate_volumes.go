package collection

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/sirupsen/logrus"
)

// populateVolumes - populates the volumes for each snapshot for any trading days in the last month
func (c *Collection) populateVolumes() {
	startOfDay := carbon.NewCarbon().SetTimezone(carbon.NewYork).StartOfDay()

	historicBars, err := c.mdClient.GetMultiBars(
		c.collection.Keys(), md.GetBarsRequest{
			TimeFrame:  md.OneDay,
			Start:      startOfDay.SubDays(35).ToStdTime(),
			End:        startOfDay.ToStdTime(),
			Adjustment: md.Split,
		},
	)

	if err != nil {
		logrus.Panicf("failed to get last month of bars: %v", err)
	}

	volumes := make([]snapshots.SnapshotVolume, 0)

	c.collection.IterCb(
		func(symbol string, snapshot *snapshots.Snapshot) {
			if barsForSymbol, ok := historicBars[symbol]; ok {
				for _, bar := range barsForSymbol {
					volumes = append(
						volumes, snapshots.SnapshotVolume{
							Date:   bar.Timestamp.Format("2006-01-02"),
							Volume: float64(bar.Volume),
						},
					)
				}

				snapshot.Volumes = volumes
			}
		},
	)
}
