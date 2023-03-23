package snapshots

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/golang-module/carbon/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/clients"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"time"
)

type Snapshots struct {
	snapshots       cmap.ConcurrentMap[string, *Snapshot]
	mdClient        *md.Client
	assetRepo       *assets.Repository
	barRepo         *bars.Repository
	calendarDayLive *calendars.CalendarDayLive
	symbols         []string
}

func NewSnapshots(calendarDayLive *calendars.CalendarDayLive) *Snapshots {
	s := &Snapshots{
		mdClient:        clients.GetMarketDataClient(),
		assetRepo:       assets.GetRepository(),
		barRepo:         bars.GetRepository(),
		calendarDayLive: calendarDayLive,
	}

	return s
}

func (s *Snapshots) Update(symbols []string) {
	removedSymbols, addedSymbols := lo.Difference(s.symbols, symbols)

	if len(removedSymbols) > 0 {
		for _, symbol := range removedSymbols {
			s.snapshots.Remove(symbol)
		}

		s.symbols = lo.Filter(
			s.symbols, func(symbol string, _ int) bool {
				return !lo.Contains(removedSymbols, symbol)
			},
		)
	}

	if len(addedSymbols) > 0 {
		s.symbols = append(s.symbols, addedSymbols...)
		s.populateBaseSnapshots(addedSymbols) // l
		s.populateVolumes()
		s.populateDailyStats()
	}
}

func (s *Snapshots) populateBaseSnapshots(symbols []string) {
	snapshots, err := s.mdClient.GetSnapshots(symbols, md.GetSnapshotRequest{})

	if err != nil {
		logrus.Panicf("failed to load base snapshots: %v", err)
	}

	for symbol, snapshot := range snapshots {
		asset := s.assetRepo.Get(symbol)

		s.snapshots.Set(
			symbol, &Snapshot{
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

// populateVolumes - populates the volumes for each snapshot for any trading days in the last month
func (s *Snapshots) populateVolumes() {
	startOfDay := carbon.NewCarbon().SetTimezone(carbon.NewYork).StartOfDay()

	// TODO: Just use mdClient directly
	historicBars, err := s.barRepo.GetHistoricMulti(
		s.snapshots.Keys(),
		md.OneDay,
		startOfDay.SubDays(35).ToStdTime(),
		startOfDay.ToStdTime(),
	)

	if err != nil {
		logrus.Panicf("failed to get last month of bars: %v", err)
	}

	volumes := make([]SnapshotVolume, 0)

	s.snapshots.IterCb(
		func(symbol string, snapshot *Snapshot) {
			if barsForSymbol, ok := historicBars[symbol]; ok {
				for _, bar := range barsForSymbol {
					volumes = append(
						volumes, SnapshotVolume{
							Date:   bar.Date(),
							Volume: bar.Volume,
						},
					)
				}

				snapshot.Volumes = volumes
			}
		},
	)
}

// populateDailyStats - corrects the daily high, low, and volume for the current day if in pre-market
func (s *Snapshots) populateDailyStats() {
	calendarDayUpdate := s.calendarDayLive.Get()

	if calendarDayUpdate.Condition != calendars.PreMarket {
		return
	}

	start, _ := time.Parse("2006-01-02", calendarDayUpdate.PreviousMarketDate.Date)

	multiBars, err := s.mdClient.GetMultiBars(
		s.symbols, md.GetBarsRequest{
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

	var dailyHigh float64
	var dailyLow float64
	var dailyVolume float64

	s.snapshots.IterCb(
		func(symbol string, snapshot *Snapshot) {
			if intradayBars, ok := multiBars[symbol]; ok {
				highs := make([]float64, 0)
				lows := make([]float64, 0)
				volume := 0.0

				for _, bar := range intradayBars {
					highs = append(highs, bar.High)
					lows = append(lows, bar.Low)
					volume += float64(bar.Volume)
				}

				dailyHigh = lo.Max(highs)
				dailyLow = lo.Min(lows)
				dailyVolume = volume

				snapshot.DailyHigh = dailyHigh
				snapshot.DailyLow = dailyLow
				snapshot.DailyVolume = dailyVolume
			}
		},
	)
}
