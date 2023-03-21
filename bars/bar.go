package bars

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type Bar struct {
	Symbol     string  `json:"S"`
	Open       float64 `json:"o"`
	High       float64 `json:"h"`
	Low        float64 `json:"l"`
	Close      float64 `json:"c"`
	Volume     float64 `json:"v"`
	TradeCount uint64  `json:"n"`
	Timestamp  int64   `json:"t"`
}

// Date taken from the Timestamp field (ms) and formatted as YYYY-MM-DD
func (b *Bar) Date() string {
	return time.UnixMilli(b.Timestamp).Format("2006-01-02")
}

func (b *Bar) String() string {
	j := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := j.Marshal(b)

	if err != nil {
		return fmt.Sprintf("Error marshalling bar: %s", err)
	}

	return fmt.Sprintf("%s", data)
}
