package calendars

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/database"
	"github.com/phoobynet/market-deck-server/decks"
	"testing"
)

func TestNewCalendarDayRepository(t *testing.T) {
	database.Connect()
	database.Migrate(&assets.Asset{})
	database.Migrate(&CalendarDay{})
	database.Migrate(&decks.Deck{})

	db := database.GetDB()
	client := alpaca.NewClient(alpaca.ClientOpts{})

	repo := NewCalendarDayRepository(db, client)

	daysAgo := []int{
		365,
		180,
		90,
		30,
		14,
		7,
		4,
	}

	actual := repo.PickByIntervals(
		daysAgo,
	)

	if len(actual) != len(daysAgo) {
		t.Errorf("expected %d calendar day, got %d", len(daysAgo), len(actual))
	}

	for _, cd := range actual {
		t.Logf("\n%s", cd.String())
	}
}
