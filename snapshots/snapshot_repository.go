package snapshots

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/golang-module/carbon/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/clients"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
	"sync"
)

var snapshotRepositoryOnce sync.Once

var snapshotRepository *Repository

type Repository struct {
	mdClient *md.Client
}

func GetRepository() *Repository {
	snapshotRepositoryOnce.Do(
		func() {
			snapshotRepository = &Repository{
				mdClient: clients.GetMarketDataClient(),
			}
		},
	)

	return snapshotRepository
}

func (r *Repository) GetMulti(symbols []string) (map[string]Snapshot, error) {
	mdSnapshots, err := r.mdClient.GetSnapshots(symbols, md.GetSnapshotRequest{})

	if err != nil {
		return nil, err
	}

	result := make(map[string]Snapshot)

	now := carbon.Now(carbon.NewYork).Format("Y-m-d")

	for symbol, mdSnapshot := range mdSnapshots {
		dailyBar := bars.FromMarketDataBar(symbol, *mdSnapshot.DailyBar)
		previousDailyBar := bars.FromMarketDataBar(symbol, *mdSnapshot.PrevDailyBar)

		actualPreviousDailyBar := previousDailyBar

		if dailyBar.Date() < now {
			actualPreviousDailyBar = dailyBar
		}

		previousClose := actualPreviousDailyBar.Close

		s := Snapshot{
			LatestBar:              bars.FromMarketDataBar(symbol, *mdSnapshot.MinuteBar),
			LatestQuote:            quotes.FromMarketDataQuote(symbol, *mdSnapshot.LatestQuote),
			LatestTrade:            trades.FromMarketDataTrade(symbol, *mdSnapshot.LatestTrade),
			DailyBar:               dailyBar,
			PreviousDailyBar:       previousDailyBar,
			ActualPreviousDailyBar: actualPreviousDailyBar,
			PreviousClose:          previousClose,
		}

		diff := numbers.NumberDiff(s.PreviousClose, s.LatestTrade.Price)
		changes := cmap.New[SnapshotChange]()

		changes.Set(
			"Since Previous", SnapshotChange{
				Change:        diff.Change,
				ChangePercent: diff.ChangePercent,
				ChangeSign:    diff.Multiplier,
				ChangeAbs:     diff.AbsoluteChange,
			},
		)

		s.Changes = changes.Items()

		result[symbol] = s
	}

	return result, nil
}
