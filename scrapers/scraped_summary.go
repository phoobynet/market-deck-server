package scrapers

type ScrapedSummary struct {
	Ticker string            `json:"ticker"`
	Name   string            `json:"name"`
	Data   map[string]string `json:"data"`
}
