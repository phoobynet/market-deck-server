package technicals

import (
	"encoding/json"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	bars2 "github.com/phoobynet/market-deck-server/bars"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"testing"
)

var testBars []bars2.Bar

func getBars() []bars2.Bar {
	if len(testBars) == 0 {
		file, err := os.Open("technical_bars_test.json")

		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		if err != nil {
			logrus.Error(err)
		}

		data, err := io.ReadAll(file)

		if err != nil {
			logrus.Error(err)
		}

		var bars []marketdata.Bar

		err = json.Unmarshal(data, &bars)

		if err != nil {
			logrus.Error(err)
		}

		for _, bar := range bars {
			testBars = append(testBars, bars2.FromMarketDataBar("AAPL", bar))
		}
	}

	return testBars
}

func TestTechnicalBars_SMA(t *testing.T) {
	smaIndicator := NewSMAIndicator(5)

	technicalBars := NewTechnicalBars(smaIndicator)

	for _, bar := range getBars() {
		technicalBars.Push(bar)
	}

	var result []*TechnicalBar

	for _, bar := range result {
		for indicatorKey, technical := range bar.Technicals {
			t.Logf("%s: %v\n", indicatorKey, technical)
		}
	}
}

func BenchmarkNewTechnicalBarsSMA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smaIndicator := NewSMAIndicator(5)

		technicalBars := NewTechnicalBars(smaIndicator)

		for _, bar := range getBars() {
			technicalBars.Push(bar)
		}
	}
}
