package repository

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-deck-server/clients"
	"sync"
)

var snapshotRepositoryOnce sync.Once

var snapshotRepository *Repository

type Repository struct {
	mdClient *md.Client
}

func Get() *Repository {
	snapshotRepositoryOnce.Do(
		func() {
			snapshotRepository = &Repository{
				mdClient: clients.GetMarketDataClient(),
			}
		},
	)

	return snapshotRepository
}
