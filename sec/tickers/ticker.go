package tickers

import (
	"fmt"
	"gorm.io/gorm"
)

// Ticker - SEC company ticker source directly from the SEC
type Ticker struct {
	gorm.Model `json:"-"`
	CIK        int    `json:"cik"`
	Ticker     string `json:"ticker" gorm:"primaryKey"`
	Name       string `json:"title"`
	Exchange   string `json:"exchange"`
}

// FullCIK returns the CIK with a fixed length of 10 digits (leading zeroes).
func (c *Ticker) FullCIK() string {
	return fmt.Sprintf("%010d", c.CIK)
}
