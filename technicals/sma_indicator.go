package technicals

import (
	"fmt"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

type SMAIndicator struct {
	window int
}

func NewSMAIndicator(window int) *SMAIndicator {
	return &SMAIndicator{window: window}
}

func (s *SMAIndicator) String() string {
	return fmt.Sprintf("SMA->{\"window\":%d}", s.window)
}

func (s *SMAIndicator) Calculate(timeSeries *techan.TimeSeries) big.Decimal {
	closePrices := techan.NewClosePriceIndicator(timeSeries)

	return techan.NewSimpleMovingAverage(closePrices, s.window).Calculate(timeSeries.LastIndex())
}
