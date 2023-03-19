package bars

import "time"

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

func (b *Bar) Date() string {
	return time.UnixMicro(b.Timestamp).Format("2006-01-02")
}
