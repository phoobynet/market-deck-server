package technicals

import (
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

type TechnicalIndicator interface {
	// Calculate takes the latest time series and returns the latest value of the indicator.
	Calculate(timeSeries *techan.TimeSeries) big.Decimal

	String() string
}
