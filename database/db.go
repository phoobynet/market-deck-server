package database

import (
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/decks"
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

func init() {
	// TODO: Provide some sort of configuration option for determining the location of local data,
	// e.g. a path option
	_db, err := gorm.Open(
		sqlite.Open("market-deck.db"), &gorm.Config{
			CreateBatchSize: 1_000,
		},
	)

	checkErr(err)

	checkErr(_db.AutoMigrate(&assets.Asset{}))
	checkErr(_db.AutoMigrate(&calendars.CalendarDay{}))
	checkErr(_db.AutoMigrate(&decks.Deck{}))

	db = _db
}

func GetDB() *gorm.DB {
	return db
}
