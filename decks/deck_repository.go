package decks

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

type DeckRepository struct {
	db *gorm.DB
}

func NewDeckRepository(db *gorm.DB) *DeckRepository {
	d := &DeckRepository{
		db: db,
	}

	if d.Count() == 0 {
		_, err := d.Create("default", []string{})

		if err != nil {
			logrus.Fatalf("error creating default deck: %v", err)
		}
	}

	return d
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

func (d *DeckRepository) FindByName(name string) (*Deck, error) {
	var deck Deck
	err := d.db.Where("name = ?", name).First(&deck).Error
	return &deck, err
}

func (d *DeckRepository) UpdateByName(name string, symbols []string) (*Deck, error) {
	var deck Deck

	d.db.Where("name = ?", name).First(&deck)

	deck.Symbols = strings.Join(symbols, ",")

	err := d.db.Save(&deck).Error

	return &deck, err
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

func (d *DeckRepository) Count() int {
	var count int64
	d.db.Model(&Deck{}).Count(&count)
	return int(count)
}
