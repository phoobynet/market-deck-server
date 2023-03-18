package assets

type Asset struct {
	Symbol   string `gorm:"primaryKey" json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Status   string `json:"status"`
	Class    string `json:"class"`
	Query    string `json:"-"`
}
