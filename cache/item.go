package cache

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	URL    string `gorm:"primaryKey"`
	Data   string `gorm:"type:text"`
	TTLSec int64
}
