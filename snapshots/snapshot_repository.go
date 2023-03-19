package snapshots

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
)

type Repository struct {
	mdClient *md.Client
}

func NewSnapshotRepository(mdClient *md.Client) *Repository {
	return &Repository{
		mdClient: mdClient,
	}
}

func (r *Repository) GetMulti(symbols []string) (map[string]Snapshot, error) {
	mdSnapshots, err := r.mdClient.GetSnapshots(symbols, md.GetSnapshotRequest{})

	if err != nil {
		return nil, err
	}

	result := make(map[string]Snapshot)

	for symbol, mdSnapshot := range mdSnapshots {
		result[symbol] = Snapshot{
			LatestBar:   bars.FromMarketDataBar(symbol, *mdSnapshot.MinuteBar),
			LatestQuote: quotes.FromMarketDataQuote(symbol, *mdSnapshot.LatestQuote),
			LatestTrade: trades.FromMarketDataTrade(symbol, *mdSnapshot.LatestTrade),
		}
	}

	return result, nil
}
