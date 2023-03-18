package assets

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"strings"
)

func FromAlpacaAsset(asset alpaca.Asset) Asset {
	return Asset{
		Symbol:   asset.Symbol,
		Name:     asset.Name,
		Exchange: asset.Exchange,
		Status:   string(asset.Status),
		Class:    string(asset.Class),
		Query:    strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s %s", asset.Symbol, asset.Name)), " ", ""),
	}
}
