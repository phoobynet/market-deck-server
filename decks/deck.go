package decks

import (
	"gorm.io/gorm"
)

type Deck struct {
	gorm.Model
	Name    string
	Symbols string
}
