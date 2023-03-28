package tickers

import (
	secclient "github.com/phoobynet/market-deck-server/sec/sechttp"
	"github.com/sirupsen/logrus"
	"time"
)

// companyTickersUrl - the URL to the SEC's list of company tickers
const companyTickersUrl = "https://www.sec.gov/files/company_tickers_exchange.json"

// tickers - in memory cache of the list of tickers
var tickers []Ticker

// ttl - time to live for the cached version of raw tickers response
var ttl = time.Hour * 24 * 7

// get - obtains the list of tickers from the SEC, or returns the cached version if the TTL has not expired
func get() []Ticker {
	if tickers != nil {
		return tickers
	}

	data, err := secclient.GetWithCache(companyTickersUrl, ttl)

	if err != nil {
		logrus.Fatalf("Error getting company tickers: %v", err)
	}

	r, err := unmarshallResponse(data)

	if err != nil {
		logrus.Panicf("Error unmarshalling response: %v", err)
	}

	tickers = responseToTickers(r)

	return tickers
}
