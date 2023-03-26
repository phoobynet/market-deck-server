package facts

import (
	"github.com/phoobynet/market-deck-server/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

type FactQuery struct {
	Ticker          string
	FinancialYear   int
	FinancialPeriod string
	Form            string
	Concept         string
}

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

func (r *Repository) Find(query FactQuery) []Fact {
	var queryResults []Fact

	sql := "select * from facts where ticker = @Ticker"

	if query.FinancialYear != 0 {
		sql += " and financial_year = @FinancialYear"
	}

	if query.FinancialPeriod != "" {
		sql += " and financial_period = @FinancialPeriod"
	}

	if query.Form != "" {
		sql += " and form = @Form"
	}

	if query.Concept != "" {
		sql += " and concept = @Concept"
	}

	sql += " order by end_date desc"

	if err := r.db.Model(&Fact{}).Raw(sql, query).Find(&queryResults).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			logrus.Panicf("Error gett: %v", err)
		}
	}

	return queryResults
}

func (r *Repository) HasTicker(ticker string, ttl time.Duration) bool {
	var fact Fact

	if err := r.db.Model(&Fact{}).Where("ticker = ?", strings.ToUpper(ticker)).First(&fact).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		} else {
			logrus.Panicf("Error checking if ticker %s exists: %v", ticker, err)
		}
	}

	if time.Now().Sub(fact.UpdatedAt) > ttl {
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
