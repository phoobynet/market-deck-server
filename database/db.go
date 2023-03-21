package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func checkErr(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func Connect() {
	// e.g. a path option
	_db, err := gorm.Open(
		sqlite.Open("market-deck.db"), &gorm.Config{
			CreateBatchSize: 1_000,
		},
	)

	checkErr(err)

	db = _db
}

func Migrate(model interface{}) {
	checkErr(db.AutoMigrate(model))
}

func GetDB() *gorm.DB {
	return db
}
