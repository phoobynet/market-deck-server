package trades

type Trade struct {
	Symbol    string  `json:"S"`
	Price     float64 `json:"p"`
	Size      float64 `json:"s"`
	Timestamp int64   `json:"t"`
}
