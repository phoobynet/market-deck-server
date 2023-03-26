package facts

type FactSummary struct {
	Year                         int     `json:"year"`
	Period                       string  `json:"period"`
	CommonStockSharesOutstanding float64 `json:"outstandingShares"` // concept: EntityCommonStockSharesOutstanding
	PublicFloat                  float64 `json:"publicFloat"`       // concept: EntityPublicFloat
	OperatingIncome              float64 `json:"operatingIncome"`   // concept: OperatingIncomeLoss
	NetIncome                    float64 `json:"netIncome"`         // concept: NetIncomeLoss
	GrossProfit                  float64 `json:"grossProfit"`       // concept: GrossProfit
	EPSBasic                     float64 `json:"epsBasic"`          // concept: EarningsPerShareBasic
	EPSDiluted                   float64 `json:"epsDiluted"`        // concept: EarningsPerShareDiluted
}
