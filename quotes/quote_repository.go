package quotes

import md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"

type Repository struct {
	mdClient *md.Client
}

func NewQuoteRepository(mdClient *md.Client) *Repository {
	return &Repository{mdClient}
}

func (r *Repository) GetLatestMulti(symbols []string) (map[string]Quote, error) {
	rawQuotes, err := r.mdClient.GetLatestQuotes(
		symbols, md.GetLatestQuoteRequest{
			Feed: md.SIP,
		},
	)

	if err != nil {
		return nil, err
	}

	quotes := make(map[string]Quote)

	for symbol, rawQuote := range rawQuotes {
		quotes[symbol] = FromMarketDataQuote(symbol, rawQuote)
	}

	return quotes, nil
}
