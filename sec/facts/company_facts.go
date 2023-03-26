package facts

type FactUnit struct {
	EndDate         string  `json:"end"`
	Value           float64 `json:"val"`
	AccessionNumber string  `json:"accn"`
	FinancialYear   int     `json:"fy"`
	FinancialPeriod string  `json:"fp"`
	Form            string  `json:"form"`
	Filed           string  `json:"filed"`
	Frame           string  `json:"frame"`
}

type Fact struct {
	Label       string              `json:"label"`
	Description string              `json:"description"`
	Units       map[string]FactUnit `json:"units"`
}

type CompanyFacts struct {
	CIK        int               `json:"cik"`
	EntityName string            `json:"entityName"`
	Facts      map[string][]Fact `json:"facts"`
}
