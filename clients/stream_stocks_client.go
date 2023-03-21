package clients

import (
	"context"
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/sirupsen/logrus"
	"sync"
)

var streamClientOnce sync.Once

var stocksClient *stream.StocksClient

func GetStocksClient() *stream.StocksClient {
	streamClientOnce.Do(
		func() {
			stocksClient = stream.NewStocksClient(md.SIP)

			err := stocksClient.Connect(context.TODO())

			if err != nil {
				logrus.Panicf("error connecting to stocks client: %v", err)
			}
		},
	)

	return stocksClient
}
