package snapshots

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
)

type Repository struct {
	mdClient        *md.Client
	assetRepository *assets.AssetRepository
}

func NewSnapshotRepository(mdClient *md.Client, assetRepository *assets.AssetRepository) *Repository {
	return &Repository{
		mdClient:        mdClient,
		assetRepository: assetRepository,
	}
}

func (r *Repository) GetMulti(symbols []string) (map[string]Snapshot, error) {
	mdSnapshots, err := r.mdClient.GetSnapshots(symbols, md.GetSnapshotRequest{})

	if err != nil {
		return nil, err
	}

	symbolsAssets := r.assetRepository.GetMulti(symbols)

	result := make(map[string]Snapshot)

	for symbol, mdSnapshot := range mdSnapshots {
		result[symbol] = Snapshot{
			Asset:       symbolsAssets[symbol],
			LatestBar:   bars.FromMarketDataBar(symbol, *mdSnapshot.MinuteBar),
			LatestQuote: quotes.FromMarketDataQuote(symbol, *mdSnapshot.LatestQuote),
			LatestTrade: trades.FromMarketDataTrade(symbol, *mdSnapshot.LatestTrade),
		}
	}

	return result, nil
}
