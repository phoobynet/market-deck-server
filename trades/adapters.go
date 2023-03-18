package trades

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

func FromStreamTrade(trade stream.Trade) Trade {
	return Trade{
		Symbol:    trade.Symbol,
		Price:     trade.Price,
		Size:      float64(trade.Size),
		Timestamp: trade.Timestamp.UnixMilli(),
	}
}

func FromMarketDataTrade(symbol string, trade marketdata.Trade) Trade {
	return Trade{
		Symbol:    symbol,
		Price:     trade.Price,
		Size:      float64(trade.Size),
		Timestamp: trade.Timestamp.UnixMilli(),
	}
}
