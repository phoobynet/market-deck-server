package assets

import (
	"github.com/phoobynet/market-deck-server/sec/tickers"
	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model
	Symbol    string         `gorm:"primaryKey" json:"S"`
	Name      string         `json:"n"`
	Exchange  string         `json:"x"`
	SECTicker tickers.Ticker `json:"secTicker" gorm:"foreignKey:Symbol;references:Ticker"`
	Status    string         `json:"-"`
	Class     string         `json:"-"`
	Query     string         `json:"-"`
}
