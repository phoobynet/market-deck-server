package facts

import (
	"github.com/phoobynet/market-deck-server/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

type Repository struct {
	db *gorm.DB
}

var onceRepository sync.Once

var repository *Repository

func GetRepository() *Repository {
	onceRepository.Do(
		func() {
			repository = &Repository{
				db: database.Get(),
			}
		},
	)

	return repository
}

func (r *Repository) GetByTicker(ticker string) []Fact {
	var factUnits []Fact

	if err := r.db.Model(&Fact{}).Where("ticker = ?", strings.ToUpper(ticker)).Find(&factUnits).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil
		}

		logrus.Panicf("Error getting fact units for ticker %s: %v", ticker, err)
	}

	return factUnits
}

func (r *Repository) HasTicker(ticker string, ttl time.Duration) bool {
	var factUnit Fact

	if err := r.db.Model(&Fact{}).First(&factUnit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		} else {
			logrus.Panicf("Error checking if ticker %s exists: %v", ticker, err)
		}
	}

	if time.Now().Sub(factUnit.UpdatedAt) > ttl {
		if err := r.db.Model(&Fact{}).Where(
			"ticker = ?",
			strings.ToUpper(ticker),
		).Delete(&Fact{}).Error; err != nil {
			logrus.Panicf("Error deleting ticker %s: %v", ticker, err)
		}

		return false
	}

	return true
}

func (r *Repository) BulkInsert(factUnits []Fact) {
	if err := r.db.Model(&Fact{}).CreateInBatches(&factUnits, 1_000).Error; err != nil {
		logrus.Panicf("Error inserting fact units: %v", err)
	}
}
