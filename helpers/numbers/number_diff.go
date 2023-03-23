package numbers

import (
	"math"
)

type NumberDiffResult struct {
	OriginalValue  float64 `json:"originalValue"`
	NewValue       float64 `json:"newValue"`
	Change         float64 `json:"change"`
	AbsoluteChange float64 `json:"absoluteChange"`
	Multiplier     int8    `json:"multiplier"`
	Sign           string  `json:"sign"`
	ChangePercent  float64 `json:"changePercent"`
}

func NumberDiff(originalValue, newValue float64) NumberDiffResult {
	change := newValue - originalValue
	var multiplier int8 = 0
	sign := ""

	if change > 0 {
		multiplier = 1
		sign = "+"
	} else if change < 0 {
		multiplier = -1
		sign = "-"
	}

	changePercent := float64(0)

	if change != 0 && originalValue != 0 {
		changePercent = change / originalValue
	}

	return NumberDiffResult{
		OriginalValue:  originalValue,
		NewValue:       newValue,
		Change:         change,
		AbsoluteChange: math.Abs(change),
		Multiplier:     multiplier,
		ChangePercent:  changePercent,
		Sign:           sign,
	}
}
