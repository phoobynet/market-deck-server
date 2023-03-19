package assets

type Asset struct {
	Symbol   string `gorm:"primaryKey" json:"S"`
	Name     string `json:"n"`
	Exchange string `json:"x"`
	Status   string `json:"-"`
	Class    string `json:"-"`
	Query    string `json:"-"`
}
