package tickers

import (
	"encoding/json"
)

// response - the raw response from the SEC - the precursor to the Ticker struct
type response struct {
	Fields []string        `json:"fields"`
	Data   [][]interface{} `json:"data"`
}

func unmarshallResponse(data []byte) (*response, error) {
	var r response

	err := json.Unmarshal(data, &r)

	return &r, err
}
