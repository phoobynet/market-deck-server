package tickers

import (
	"github.com/phoobynet/market-deck-server/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Repository struct {
	db *gorm.DB
}

var tickerRepositoryOnce sync.Once
var tickerRepository *Repository

// GetRepository returns a pointer to the ticker repository
// The repository is created on the first call to this function
func GetRepository() *Repository {
	tickerRepositoryOnce.Do(
		func() {
			tickerRepository = &Repository{
				db: database.Get(),
			}
		},
	)

	return tickerRepository
}

// GetByTicker returns a ticker by its ticker symbol
func (r *Repository) GetByTicker(ticker string) *Ticker {
	var queryResult Ticker

	if err := r.db.Model(&Ticker{}).Where("ticker = ?", ticker).First(&queryResult).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
	}

	return &queryResult
}

// GetAsMap returns a map of tickers
func (r *Repository) GetAsMap() map[string]Ticker {
	var queryResult []Ticker

	if err := r.db.Model(&Ticker{}).Find(&queryResult).Error; err != nil {
		logrus.Panicf("Error getting tickers: %v", err)
	}

	result := make(map[string]Ticker, len(queryResult))

	for _, ticker := range queryResult {
		result[ticker.Ticker] = ticker
	}

	return result
}

// BulkInsert inserts a slice of tickers into the database
func (r *Repository) BulkInsert(tickers []Ticker) {
	if err := r.db.Model(&Ticker{}).CreateInBatches(tickers, 1_000).Error; err != nil {
		logrus.Panicf("Error inserting tickers: %v", err)
	}
}

// LastUpdated returns the last updated time of the ticker table
// a nil pointer is returned if the table is empty
func (r *Repository) LastUpdated() *time.Time {
	var queryResult Ticker

	if err := r.db.Model(&Ticker{}).First(&queryResult).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
	}

	return &queryResult.UpdatedAt
}

func (r *Repository) DeleteAll() {
	//goland:noinspection SqlResolve,SqlWithoutWhere
	if err := r.db.Exec("DELETE FROM tickers").Error; err != nil {
		logrus.Panicf("Error deleting tickers: %v", err)
	}
}
