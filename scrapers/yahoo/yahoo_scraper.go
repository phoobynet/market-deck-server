package yahoo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/phoobynet/market-deck-server/cache"
	"github.com/phoobynet/market-deck-server/scrapers"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetTickerSummary(ticker string) (*scrapers.ScrapedSummary, error) {
	cacheRepository := cache.GetRepository()

	var body string

	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s", ticker)

	if item := cacheRepository.Get(url); item != nil {
		body = string(item.Data)
	} else {
		response, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(response.Body)

		rawBody, err := io.ReadAll(response.Body)

		if err != nil {
			return nil, err
		}

		cacheRepository.Set(url, rawBody, 24*time.Hour)

		body = string(rawBody)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))

	if err != nil {
		return nil, err
	}

	summary := &scrapers.ScrapedSummary{
		Ticker: ticker,
	}

	summary.Name = getName(ticker, doc)
	summary.Data = getSummaryData(doc)

	return summary, nil
}

func getName(ticker string, doc *goquery.Document) string {
	nameSelector := fmt.Sprintf("h1:contains('(%s)')", ticker)

	nameSelection := doc.Find(nameSelector)

	if nameSelection != nil {
		nameText := nameSelection.Text()

		index := strings.Index(nameText, fmt.Sprintf("(%s)", ticker))

		return strings.TrimSpace(nameText[:index])
	}

	return ""
}

func getSummaryData(doc *goquery.Document) map[string]string {
	data := map[string]string{}

	doc.Find("#quote-summary table").Each(
		func(i int, s *goquery.Selection) {
			s.Find("tr").Each(
				func(i int, s *goquery.Selection) {
					var label string
					var v string
					s.Find("td").Each(
						func(i int, s *goquery.Selection) {
							if i == 0 {
								label = s.Text()
							} else if i == 1 { //value
								v = s.Text()
							} else {
								label = ""
								v = ""
							}

							if label != "" || v != "" {
								data[label] = v
							}
						},
					)
				},
			)
		},
	)

	fmt.Printf("%v", data)

	return data
}
