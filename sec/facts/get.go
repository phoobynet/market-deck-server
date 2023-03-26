package facts

import (
	"fmt"
	sechttp "github.com/phoobynet/market-deck-server/sec/http"
	"github.com/phoobynet/market-deck-server/sec/tickers"
	"github.com/sirupsen/logrus"
)

func Get(query FactQuery) []Fact {
	tickersRepo := tickers.GetRepository()

	secTicker := tickersRepo.GetByTicker(query.Ticker)

	if secTicker == nil {
		logrus.Panicf("No ticker found for %s", query.Ticker)
	}

	r := GetRepository()

	if r.HasTicker(query.Ticker, ttl) {
		return r.Find(query)
	}

	url := fmt.Sprintf("https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json", secTicker.FullCIK())

	logrus.Infof("Getting facts for %s from %s\n", query.Ticker, url)

	data, err := sechttp.GetWithCache(url, ttl)

	if err != nil {
		logrus.Panicf("Error getting facts for %s: %v", query.Ticker, err)
	}

	facts := parseFacts(secTicker, data)

	r.BulkInsert(facts)

	return Get(query)
}
