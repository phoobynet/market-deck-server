package quotes

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

func FromStreamQuote(q stream.Quote) Quote {
	return Quote{
		Symbol:      q.Symbol,
		AskPrice:    q.AskPrice,
		AskSize:     float64(q.AskSize),
		AskExchange: q.AskExchange,
		BidPrice:    q.BidPrice,
		BidSize:     float64(q.BidSize),
		BidExchange: q.BidExchange,
		Conditions:  q.Conditions,
		Tape:        q.Tape,
		Timestamp:   q.Timestamp.UnixMilli(),
	}
}

func FromMarketDataQuote(symbol string, q marketdata.Quote) Quote {
	return Quote{
		Symbol:      symbol,
		AskPrice:    q.AskPrice,
		AskSize:     float64(q.AskSize),
		AskExchange: q.AskExchange,
		BidPrice:    q.BidPrice,
		BidSize:     float64(q.BidSize),
		BidExchange: q.BidExchange,
		Conditions:  q.Conditions,
		Tape:        q.Tape,
		Timestamp:   q.Timestamp.UnixMilli(),
	}
}
