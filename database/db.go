package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once

var db *gorm.DB

func GetDB() *gorm.DB {
	// e.g. a path option
	once.Do(
		func() {
			_db, err := gorm.Open(
				sqlite.Open("market-deck.db"), &gorm.Config{
					CreateBatchSize: 1_000,
				},
			)

			if err != nil {
				logrus.Panicf("Failed to get DB: %v\n", err)
			}

			db = _db
		},
	)

	return db
}
