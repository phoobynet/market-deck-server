package technicals

import (
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	"time"
)

type TechnicalBar struct {
	bars.Bar
	Technicals map[string]big.Decimal `json:"technicals"`
}

func NewTechnicalBar(bar bars.Bar, technicals map[string]big.Decimal) *TechnicalBar {
	return &TechnicalBar{bar, technicals}
}

type TechnicalBars struct {
	bars       []*TechnicalBar
	timeSeries *techan.TimeSeries
	indicators []TechnicalIndicator
}

func (t *TechnicalBars) Push(bar bars.Bar) []*TechnicalBar {
	period := techan.NewTimePeriod(time.UnixMilli(bar.Timestamp), time.Minute)
	candle := techan.NewCandle(period)
	candle.OpenPrice = big.NewDecimal(bar.Open)
	candle.ClosePrice = big.NewDecimal(bar.Close)
	candle.MaxPrice = big.NewDecimal(bar.High)
	candle.MinPrice = big.NewDecimal(bar.Low)
	candle.Volume = big.NewDecimal(bar.Volume)

	t.timeSeries.AddCandle(candle)

	technicals := make(map[string]big.Decimal)

	for _, indicator := range t.indicators {
		technicals[indicator.String()] = indicator.Calculate(t.timeSeries)
	}

	t.bars = append(t.bars, NewTechnicalBar(bar, technicals))

	return t.bars
}

// https://www.investopedia.com/top-7-technical-analysis-tools-4773275

func NewTechnicalBars(indicators ...TechnicalIndicator) *TechnicalBars {
	return &TechnicalBars{
		bars:       make([]*TechnicalBar, 0),
		timeSeries: techan.NewTimeSeries(),
		indicators: indicators,
	}
}
