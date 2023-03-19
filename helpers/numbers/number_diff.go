package numbers

import "math"

type NumberDiffResult struct {
	Change         float64
	AbsoluteChange float64
	Sign           int8
	ChangePercent  float64
}

func NumberDiff(originalValue, newValue float64) NumberDiffResult {
	change := newValue - originalValue
	var sign int8 = 0

	if change > 0 {
		sign = 1
	} else if change < 0 {
		sign = -1
	}

	return NumberDiffResult{
		Change:         change,
		AbsoluteChange: math.Abs(change),
		Sign:           sign,
		ChangePercent:  change / originalValue,
	}
}
