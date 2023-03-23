package snapshots

import (
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
)

type SnapshotLite struct {
	Class         string                   `json:"class"`
	Symbol        string                   `json:"symbol"`
	Name          string                   `json:"name"`
	Exchange      string                   `json:"exchange"`
	Price         float64                  `json:"price"`
	PrevClose     float64                  `json:"prevClose"`
	PrevCloseDate string                   `json:"prevCloseDate"`
	DailyHigh     float64                  `json:"dailyHigh"`
	DailyLow      float64                  `json:"dailyLow"`
	DailyVolume   float64                  `json:"dailyVolume"`
	Change        numbers.NumberDiffResult `json:"change"`
	Volumes       []SnapshotLiteVolume     `json:"volumes"`
	MonthlyBars   []bars.Bar               `json:"monthlyBars"`
	YtdChange     numbers.NumberDiffResult `json:"ytdChange"`
}

type SnapshotLiteVolume struct {
	Date   string  `json:"date"`
	Volume float64 `json:"vol"`
}
