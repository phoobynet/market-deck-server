package numbers

import (
	"fmt"
	"math"
)

type NumberDiffResult struct {
	OriginalValue  float64
	NewValue       float64
	Change         float64
	AbsoluteChange float64
	Sign           int8
	ChangePercent  float64
}

func (n *NumberDiffResult) String() string {
	return fmt.Sprintf(
		"Original Value: %f, New Value: %f, Change: %f, AbsoluteChange: %f, Sign: %d, ChangePercent: %f",
		n.OriginalValue,
		n.NewValue,
		n.Change,
		n.AbsoluteChange,
		n.Sign,
		n.ChangePercent,
	)
}

func NumberDiff(originalValue, newValue float64) NumberDiffResult {
	change := newValue - originalValue
	var sign int8 = 0

	if change > 0 {
		sign = 1
	} else if change < 0 {
		sign = -1
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
		Sign:           sign,
		ChangePercent:  changePercent,
	}
}
