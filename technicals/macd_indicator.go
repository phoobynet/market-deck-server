package technicals

import (
	"github.com/markcheno/go-talib"
	. "github.com/phoobynet/market-deck-server/bars"
	"github.com/samber/lo"
)

type MACDIndicator struct {
}

type MACDIndicatorInput struct {
	Bars         []Bar
	FastPeriod   int
	SlowPeriod   int
	SignalPeriod int
}

type MACDIndicatorOutput struct {
	MACD      []float64
	Signal    []float64
	Histogram []float64
}

func (m *MACDIndicator) Calculate(input *MACDIndicatorInput) MACDIndicatorOutput {
	closingPrices := lo.Map[Bar, float64](
		input.Bars, func(bar Bar, i int) float64 {
			return bar.Close
		},
	)

	macd, sign, histogram := talib.Macd(closingPrices, input.FastPeriod, input.SlowPeriod, input.SignalPeriod)

	return MACDIndicatorOutput{
		MACD:      macd,
		Signal:    sign,
		Histogram: histogram,
	}
}
