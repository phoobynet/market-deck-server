package decks

import (
	"gorm.io/gorm"
	"strings"
)

type DeckRepository struct {
	db *gorm.DB
}

func NewDeckRepository(db *gorm.DB) *DeckRepository {
	return &DeckRepository{
		db: db,
	}
}

func (d *DeckRepository) Create(name string, symbols []string) (*Deck, error) {
	var deck Deck

	deck.Name = name
	deck.Symbols = strings.Join(symbols, ",")

	err := d.db.Create(&deck).Error

	return &deck, err
}

func (d *DeckRepository) FindAll() ([]Deck, error) {
	var decks []Deck
	err := d.db.Find(&decks).Error
	return decks, err
}

func (d *DeckRepository) Update(id uint, name string, symbols []string) (*Deck, error) {
	var deck Deck

	d.db.First(&deck, id)

	deck.Name = name
	deck.Symbols = strings.Join(symbols, ",")

	err := d.db.Save(&deck).Error

	return &deck, err
}

func (d *DeckRepository) Delete(id uint) error {
	return d.db.
		Delete(&Deck{}, id).
		Error
}
