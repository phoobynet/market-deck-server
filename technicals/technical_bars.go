package technicals

import (
	"fmt"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	"time"
)

type TechnicalBars struct {
}

func NewTechnicalBars(bars []bars.Bar) *TechnicalBars {
	series := techan.NewTimeSeries()

	for _, bar := range bars {
		period := techan.NewTimePeriod(time.UnixMicro(bar.Timestamp), time.Minute)
		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewDecimal(bar.Open)
		candle.ClosePrice = big.NewDecimal(bar.Close)
		candle.MaxPrice = big.NewDecimal(bar.High)
		candle.MinPrice = big.NewDecimal(bar.Low)
		candle.Volume = big.NewDecimal(bar.Volume)

		series.AddCandle(candle)
	}

	// https://www.investopedia.com/top-7-technical-analysis-tools-4773275
	closePrices := techan.NewClosePriceIndicator(series)
	ema := techan.NewEMAIndicator(closePrices, 10)
	rsi := techan.NewRelativeStrengthIndicator(closePrices, 10)
	macd := techan.NewMACDIndicator(closePrices, 9, 26)
	slowStoch := techan.NewSlowStochasticIndicator(closePrices, 14)
	fastStoch := techan.NewFastStochasticIndicator(series, 14)

	fmt.Println(ema.Calculate(0))
	fmt.Println(macd.Calculate(0))
	fmt.Println(rsi.Calculate(0))
	fmt.Println(slowStoch.Calculate(0))
	fmt.Println(fastStoch.Calculate(0))

	return &TechnicalBars{}
}
