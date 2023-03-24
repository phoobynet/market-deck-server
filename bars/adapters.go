package bars

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/samber/lo"
)

func FromStreamBar(bar stream.Bar) Bar {
	return Bar{
		Symbol:     bar.Symbol,
		Open:       bar.Open,
		High:       bar.High,
		Low:        bar.Low,
		Close:      bar.Close,
		Volume:     float64(bar.Volume),
		TradeCount: bar.TradeCount,
		Timestamp:  bar.Timestamp.UnixMilli(),
	}
}

func FromMarketDataBar(symbol string, bar marketdata.Bar) Bar {
	return Bar{
		Symbol:     symbol,
		Open:       bar.Open,
		High:       bar.High,
		Low:        bar.Low,
		Close:      bar.Close,
		Volume:     float64(bar.Volume),
		TradeCount: bar.TradeCount,
		Timestamp:  bar.Timestamp.UnixMilli(),
	}
}

func FromMarketDataBars(symbol string, bars []marketdata.Bar) []Bar {
	return lo.Map[marketdata.Bar, Bar](
		bars, func(bar marketdata.Bar, _ int) Bar {
			return FromMarketDataBar(symbol, bar)
		},
	)
}
