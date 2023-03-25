package cache

import (
	"github.com/phoobynet/market-deck-server/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

var once sync.Once

type Repository struct {
	db *gorm.DB
}

func GetRepository() *Repository {
	var r *Repository
	once.Do(
		func() {
			r = &Repository{
				db: database.GetDB(),
			}
		},
	)
	return r
}

func (r *Repository) Get(url string) *Item {
	var item Item

	err := r.db.Model(&Item{}).First(&item, "url = ?", strings.ToLower(url)).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
	} else {
		threshold := item.UpdatedAt.Add(time.Duration(item.TTLSec))

		if time.Now().After(threshold) {
			if err := r.db.Model(&Item{}).Delete(&item).Error; err != nil {
				logrus.Errorf("Error deleting cache item: %s", err.Error())
			}
		}
	}

	return &item
}

func (r *Repository) Set(url string, data string, ttlSec int64) {
	item := &Item{
		URL:    strings.ToLower(url),
		Data:   data,
		TTLSec: ttlSec,
	}

	r.db.Model(&Item{}).Create(item)
}

func (r *Repository) Delete(url string) {
	item := r.Get(url)

	if item != nil {
		if err := r.db.Model(&Item{}).Delete(item).Error; err != nil {
			logrus.Errorf("Error deleting cache item: %s", err.Error())
		}
	}
}
