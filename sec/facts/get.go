package facts

import (
	"fmt"
	sechttp "github.com/phoobynet/market-deck-server/sec/http"
	"github.com/phoobynet/market-deck-server/sec/tickers"
	"github.com/sirupsen/logrus"
)

func Get(ticker string) []Fact {
	tickersRepo := tickers.GetRepository()

	secTicker := tickersRepo.GetByTicker(ticker)

	if secTicker == nil {
		logrus.Panicf("No ticker found for %s", ticker)
	}

	r := GetRepository()

	if r.HasTicker(ticker, ttl) {
		return r.GetByTicker(ticker)
	}

	url := fmt.Sprintf("https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json", secTicker.FullCIK())

	data, err := sechttp.GetWithCache(url, ttl)

	if err != nil {
		logrus.Panicf("Error getting facts for %s: %v", ticker, err)
	}

	facts := parseFacts(data)

	r.BulkInsert(facts)

	return facts
}
