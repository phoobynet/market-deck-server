package bars

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/market-deck-server/clients"
	"sync"
	"time"
)

type Repository struct {
	mdClient *md.Client
}

var barRepositoryOnce sync.Once

var barRepository *Repository

func GetRepository() *Repository {
	barRepositoryOnce.Do(
		func() {
			barRepository = &Repository{
				mdClient: clients.GetMarketDataClient(),
			}
		},
	)

	return barRepository
}

func (c *Repository) GetLatestMulti(symbols []string) (map[string]Bar, error) {
	rawBars, err := c.mdClient.GetLatestBars(
		symbols, md.GetLatestBarRequest{
			Feed: md.SIP,
		},
	)

	if err != nil {
		return nil, err
	}

	bars := make(map[string]Bar)

	for symbol, rawBar := range rawBars {
		bars[symbol] = FromMarketDataBar(symbol, rawBar)
	}

	return bars, nil
}

func (c *Repository) GetHistoricMulti(symbols []string, timeframe md.TimeFrame, start, end time.Time) (
	map[string][]Bar, error,
) {
	rawBars, err := c.mdClient.GetMultiBars(
		symbols, md.GetBarsRequest{
			TimeFrame:  timeframe,
			Adjustment: md.Split,
			Start:      start,
			End:        end,
			Feed:       md.SIP,
		},
	)

	if err != nil {
		return nil, err
	}

	bars := make(map[string][]Bar)

	for symbol, rawBars := range rawBars {
		for _, rawBar := range rawBars {
			bars[symbol] = append(bars[symbol], FromMarketDataBar(symbol, rawBar))
		}
	}

	return bars, nil
}

func (c *Repository) GetIntradayMulti(symbols []string) (map[string][]Bar, error) {
	start := now().StartOfDay().ToStdTime()
	end := now().EndOfDay().ToStdTime()

	return c.GetHistoricMulti(symbols, md.OneMin, start, end)
}

func (c *Repository) GetYtdDailyMulti(symbols []string) (map[string][]Bar, error) {
	start := now().SubYears(1).StartOfDay().ToStdTime()
	end := now().EndOfDay().ToStdTime()

	return c.GetHistoricMulti(symbols, md.OneDay, start, end)
}

func now() carbon.Carbon {
	return carbon.Now(carbon.NewYork)
}
