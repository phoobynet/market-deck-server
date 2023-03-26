package assets

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/sirupsen/logrus"
	"time"
)

// Populate is a function that populates the database with assets.
func Populate() {
	r := GetRepository()

	lastUpdated := r.LastUpdated()

	if lastUpdated == nil || time.Since(*lastUpdated) > ttl {
		r.DeleteAll()

		alpacaAssets, err := r.alpacaClient.GetAssets(
			alpaca.GetAssetsRequest{
				Status:     string(alpaca.AssetActive),
				AssetClass: string(alpaca.USEquity),
			},
		)

		if err != nil {
			logrus.Panic(err)
		}

		var assets []Asset

		for _, alpacaAsset := range alpacaAssets {
			asset := FromAlpacaAsset(alpacaAsset)
			assets = append(assets, asset)
		}

		r.BulkInsert(assets)
	}
}
