package clients

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"sync"
)

var alpacaClientOnce sync.Once

var alpacaClient *alpaca.Client

func GetAlpacaClient() *alpaca.Client {
	alpacaClientOnce.Do(
		func() {
			alpacaClient = alpaca.NewClient(alpaca.ClientOpts{})
		},
	)

	return alpacaClient
}
