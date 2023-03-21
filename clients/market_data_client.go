package clients

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"sync"
)

var mdClientOnce sync.Once

var mdClient *md.Client

func GetMarketDataClient() *md.Client {
	mdClientOnce.Do(
		func() {
			mdClient = md.NewClient(
				md.ClientOpts{},
			)
		},
	)

	return mdClient
}
