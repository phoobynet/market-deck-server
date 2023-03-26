package facts

import "gorm.io/gorm"

type Fact struct {
	gorm.Model      `json:"-"`
	Ticker          string  `json:"ticker" gorm:"index"`
	CIK             int     `json:"cik" gorm:"index"`
	Root            string  `json:"root"`     // Root - e.g. dei, us-gaap
	Concept         string  `json:"concept"`  // e.g. EntityCommonStockSharesOutstanding
	UnitType        string  `json:"unitType"` // e.g. Shares, USD - describes what the value (val) is
	EndDate         string  `json:"end" gorm:"index"`
	Value           float64 `json:"val"`
	AccessionNumber string  `json:"accn"`
	FinancialYear   int     `json:"fy" gorm:"index"`
	FinancialPeriod string  `json:"fp" gorm:"index"`
	Form            string  `json:"form"`
	Filed           string  `json:"filed"`
	Frame           string  `json:"frame"`
}
