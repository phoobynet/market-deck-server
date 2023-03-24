package bars

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"time"
)

type BarCollection struct {
	bars          []Bar
	closingPrices []float64
	timeframe     marketdata.TimeFrame
}

func NewBarCollection(bars []Bar, timeframe marketdata.TimeFrame) *BarCollection {
	b := &BarCollection{
		bars:          make([]Bar, 0),
		closingPrices: make([]float64, 0),
		timeframe:     timeframe,
	}

	for _, bar := range bars {
		b.Push(bar)
	}

	return b
}

func (b *BarCollection) Push(bar Bar) {
	b.bars = append(b.bars, bar)
	b.closingPrices = append(b.closingPrices, bar.Close)
}

func (b *BarCollection) ClosingPrices() []float64 {
	return b.closingPrices
}

func (b *BarCollection) Bars() []Bar {
	return b.bars
}

func (b *BarCollection) Timeframe() marketdata.TimeFrame {
	return b.timeframe
}

func (b *BarCollection) LastBar() *Bar {
	if len(b.bars) == 0 {
		return nil
	}

	lastBar := b.bars[len(b.bars)-1]

	return &lastBar
}

func (b *BarCollection) FirstBar() *Bar {
	if len(b.bars) == 0 {
		return nil
	}

	firstBar := b.bars[0]

	return &firstBar
}

func (b *BarCollection) DateRange() (*time.Time, *time.Time) {
	firstBar := b.FirstBar()

	if firstBar == nil {
		return nil, nil
	}

	lastBar := b.LastBar()

	if lastBar == nil {
		return nil, nil
	}

	start := time.UnixMilli(firstBar.Timestamp)
	end := time.UnixMilli(lastBar.Timestamp)

	return &start, &end
}
