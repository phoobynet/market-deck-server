package snapshots

import (
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
)

type Snapshot struct {
	LatestBar        bars.Bar                  `json:"lb"`
	LatestTrade      trades.Trade              `json:"lt"`
	LatestQuote      quotes.Quote              `json:"lq"`
	PreviousDailyBar bars.Bar                  `json:"pdb"`
	DailyBar         bars.Bar                  `json:"db"`
	PreviousClose    float64                   `json:"pc"`
	Changes          map[string]SnapshotChange `json:"changes"`
}

type SnapshotChange struct {
	Since         int64   `json:"since"`
	Change        float64 `json:"c"`
	ChangePercent float64 `json:"cp"`
	ChangeSign    int8    `json:"cs"`
	ChangeAbs     float64 `json:"ca"`
}
