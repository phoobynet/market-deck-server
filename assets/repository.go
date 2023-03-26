package assets

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/phoobynet/market-deck-server/clients"
	"github.com/phoobynet/market-deck-server/database"
	"github.com/schollz/closestmatch"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

var assetRepositoryOnce sync.Once

var assetRepository *Repository

type Repository struct {
	alpacaClient *alpaca.Client
	populated    bool
	db           *gorm.DB
	search       []string
	cm           *closestmatch.ClosestMatch
}

func GetRepository() *Repository {
	assetRepositoryOnce.Do(
		func() {
			assetRepository = &Repository{
				alpacaClient: clients.GetAlpacaClient(),
				populated:    false,
				db:           database.Get(),
			}
		},
	)

	return assetRepository
}

// GetBySymbol returns an asset by its symbol/ticker
func (r *Repository) GetBySymbol(symbol string) *Asset {
	var asset Asset

	if err := r.db.Model(&Asset{}).Where("symbol = ?", symbol).Preload("SECTicker").First(&asset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		logrus.Panicf("Error getting asset: %v", err)
	}

	return &asset
}

func (r *Repository) GetMulti(symbols []string) map[string]Asset {
	var assets map[string]Asset

	if err := r.db.Where("symbol IN ?", symbols).Preload("SECTicker").Find(&assets).Error; err != nil {
		logrus.Panicf("Error getting assets: %v", err)
	}

	return assets
}

func (r *Repository) GetAll() []Asset {
	var assets []Asset

	if err := r.db.Model(&Asset{}).Preload("SECTicker").Find(&assets).Error; err != nil {
		logrus.Panicf("Error getting assets: %v", err)
	}

	return assets
}

func (r *Repository) GetByClass(assetClass alpaca.AssetClass) []Asset {
	var assets []Asset

	if err := r.db.Model(&Asset{}).Find(&assets, "class = ?", assetClass).Error; err != nil {
		logrus.Panicf("Error getting assets: %v", err)
	}

	return assets
}

func (r *Repository) LastUpdated() *time.Time {
	var queryResult Asset

	if err := r.db.Model(&Asset{}).First(&queryResult).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
	}

	return &queryResult.UpdatedAt
}

func (r *Repository) DeleteAll() {
	//goland:noinspection SqlResolve,SqlWithoutWhere
	if err := r.db.Exec("DELETE FROM assets").Error; err != nil {
		logrus.Panicf("Error deleting assets: %v", err)
	}
}

func (r *Repository) BulkInsert(assets []Asset) {
	if err := r.db.Model(&Asset{}).CreateInBatches(assets, 1_000).Error; err != nil {
		logrus.Panicf("Error inserting assets: %s", err)
	}
}
