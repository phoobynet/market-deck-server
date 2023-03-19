package trades

import md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"

type Repository struct {
	mdClient *md.Client
}

func NewTradeRepository(mdClient *md.Client) *Repository {
	return &Repository{mdClient}
}

func (r *Repository) GetLatestMulti(symbols []string) (map[string]Trade, error) {
	rawTrades, err := r.mdClient.GetLatestTrades(
		symbols, md.GetLatestTradeRequest{
			Feed: md.SIP,
		},
	)

	if err != nil {
		return nil, err
	}

	trades := make(map[string]Trade)

	for symbol, rawTrade := range rawTrades {
		trades[symbol] = FromMarketDataTrade(symbol, rawTrade)
	}

	return trades, nil
}
