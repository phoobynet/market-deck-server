package sec

type Fact struct {
	EndDate         string  `json:"end"`
	Value           float64 `json:"val"`
	AccessionNumber string  `json:"accn"`
	FinancialYear   int     `json:"fy"`
	FinancialPeriod string  `json:"fp"`
	Form            string  `json:"form"`
	Filed           string  `json:"filed"`
	Frame           string  `json:"frame"`
}

type CompanyFacts struct {
	Ticker      string            `json:"ticker"`
	CIK         int               `json:"cik"`
	Taxonomy    string            `json:"taxonomy"`
	Tag         string            `json:"tag"`
	Label       string            `json:"label"`
	Description string            `json:"description"`
	EntityName  string            `json:"entityName"`
	Units       map[string][]Fact `json:"units"`
}
