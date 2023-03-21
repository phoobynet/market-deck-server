package snapshots

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
)

type Snapshot struct {
	LatestBar   bars.Bar     `json:"lb"`
	LatestTrade trades.Trade `json:"lt"`
	LatestQuote quotes.Quote `json:"lq"`
	// PreviousDailyBar is the previous day's bar, if the market is open.
	ActualPreviousDailyBar bars.Bar                  `json:"apdb"`
	PreviousDailyBar       bars.Bar                  `json:"pdb"`
	DailyBar               bars.Bar                  `json:"db"`
	PreviousClose          float64                   `json:"pc"`
	Changes                map[string]SnapshotChange `json:"changes"`
}

func (s *Snapshot) String() string {
	j := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := j.Marshal(s)

	if err != nil {
		return ""
	}

	return string(data)
}

type SnapshotChange struct {
	Since         int64   `json:"since"`
	Label         string  `json:"label"`
	Change        float64 `json:"c"`
	ChangePercent float64 `json:"cp"`
	ChangeSign    int8    `json:"cs"`
	ChangeAbs     float64 `json:"ca"`
}

func (s *SnapshotChange) String() string {
	j := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := j.Marshal(s)

	if err != nil {
		return ""
	}

	return string(data)
}
