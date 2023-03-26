package concept

import (
	"encoding/json"
	"fmt"
	"github.com/phoobynet/market-deck-server/sec"
	"github.com/phoobynet/market-deck-server/sec/http"
	"github.com/phoobynet/market-deck-server/sec/tickers"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	sharesOutstandingTTL = 24 * 7 * time.Hour
)

func SharesOutstanding(ticker string) *sec.CompanyFacts {
	t := tickers.GetRepository().GetByTicker(ticker)

	url := fmt.Sprintf(
		"%sCIK%s/dei/EntityCommonStockSharesOutstanding.json",
		conceptBaseURL,
		t.FullCIK(),
	)

	cache, err := http.GetWithCache(url, sharesOutstandingTTL)

	if err != nil {
		logrus.Errorf("Error getting shares outstanding for %s: %v", ticker, err)
		return nil
	}

	var companyFacts sec.CompanyFacts

	if err = json.Unmarshal(cache, &companyFacts); err != nil {
		logrus.Errorf("Error unmarshalling shares outstanding for %s: %v", ticker, err)
	}

	return &companyFacts
}
