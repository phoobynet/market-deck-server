package cache

import (
	"github.com/phoobynet/market-deck-server/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

var cacheRepositoryOnce sync.Once
var cacheRepository *Repository

type Repository struct {
	db *gorm.DB
}

func GetRepository() *Repository {
	cacheRepositoryOnce.Do(
		func() {
			cacheRepository = &Repository{
				db: database.Get(),
			}
		},
	)
	return cacheRepository
}

func (r *Repository) Get(url string) *Item {
	var item Item

	err := r.db.Model(&Item{}).First(&item, "url = ?", strings.ToLower(url)).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
	} else {
		threshold := item.UpdatedAt.Add(item.TTL)

		if time.Now().After(threshold) {
			if err := r.db.Model(&Item{}).Delete(&item).Error; err != nil {
				logrus.Errorf("Error deleting cache item: %s", err.Error())
			}
		}
	}

	return &item
}

func (r *Repository) Set(url string, data []byte, ttl time.Duration) {
	item := &Item{
		URL:  strings.ToLower(url),
		Data: data,
		TTL:  ttl,
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
