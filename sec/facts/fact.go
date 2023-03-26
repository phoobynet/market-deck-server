package facts

import "gorm.io/gorm"

type Fact struct {
	gorm.Model
	Ticker          string      `json:"ticker"`
	CIK             int         `json:"cik"`
	Root            string      `json:"root"`     // Root - e.g. dei, us-gaap
	Concept         string      `json:"concept"`  // e.g. EntityCommonStockSharesOutstanding
	Label           string      `json:"label"`    // e.g. Label provided for the concept
	UnitType        string      `json:"unitType"` // e.g. Shares, USD - describes what the value (val) is
	EndDate         string      `json:"end"`
	Value           interface{} `json:"val"`
	AccessionNumber string      `json:"accn"`
	FinancialYear   int         `json:"fy"`
	FinancialPeriod string      `json:"fp"`
	Form            string      `json:"form"`
	Filed           string      `json:"filed"`
	Frame           string      `json:"frame"`
}
