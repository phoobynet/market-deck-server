package tickers

// responseToTickers - converts the response from the SEC API to a slice of Ticker structs
func responseToTickers(r *response) []Ticker {
	tickers := make([]Ticker, 0)

	var cikPosition int
	var namePosition int
	var tickerPosition int
	var exchangePosition int

	for i, field := range r.Fields {
		if field == "cik" {
			cikPosition = i
		} else if field == "name" {
			namePosition = i
		} else if field == "ticker" {
			tickerPosition = i
		} else if field == "exchange" {
			exchangePosition = i
		}
	}

	for _, data := range r.Data {
		ticker := Ticker{
			CIK:      int(data[cikPosition].(float64)),
			Ticker:   data[tickerPosition].(string),
			Name:     data[namePosition].(string),
			Exchange: data[exchangePosition].(string),
		}

		tickers = append(tickers, ticker)
	}

	return tickers
}
