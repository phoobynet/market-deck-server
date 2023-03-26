package cache

import (
	"gorm.io/gorm"
	"time"
)

type Item struct {
	gorm.Model
	URL  string `gorm:"primaryKey"`
	Data []byte
	TTL  time.Duration
}
