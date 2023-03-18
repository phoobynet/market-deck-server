package quotes

type Quote struct {
	Symbol      string   `json:"S"`
	AskPrice    float64  `json:"ap"`
	AskSize     float64  `json:"as"`
	AskExchange string   `json:"ax"`
	BidPrice    float64  `json:"bp"`
	BidSize     float64  `json:"bs"`
	BidExchange string   `json:"bx"`
	Conditions  []string `json:"c"`
	Tape        string   `json:"z"`
	Timestamp   int64    `json:"t"`
}
