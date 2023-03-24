package stream

import (
	"context"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/messages"
	"github.com/phoobynet/market-deck-server/snapshots/collection"
	"github.com/phoobynet/market-deck-server/trades"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type SnapshotStream struct {
	mu                     sync.RWMutex
	snapshotsCollection    *collection.Collection
	deckRepo               *decks.DeckRepository
	publishDuration        time.Duration
	publishTicker          *time.Ticker
	tradeStream            *trades.Stream
	tradeChan              chan map[string]*trades.Trade
	tradesMap              cmap.ConcurrentMap[string, *trades.Trade]
	refreshSnapshotsTicker *time.Ticker
}

func New(
	ctx context.Context,
	calendarDayLive *calendars.CalendarDayLive,
	messageBus chan<- messages.Message,
) *SnapshotStream {
	tradesMap := cmap.New[*trades.Trade]()
	s := &SnapshotStream{
		snapshotsCollection:    collection.New(calendarDayLive),
		deckRepo:               decks.GetRepository(),
		publishDuration:        1 * time.Second,
		tradeChan:              make(chan map[string]*trades.Trade, 1_000),
		refreshSnapshotsTicker: time.NewTicker(time.Second),
		tradesMap:              tradesMap,
	}

	s.tradeStream = trades.NewTradeStream(ctx, s.tradeChan)
	s.publishTicker = time.NewTicker(s.publishDuration)

	go func() {
		for {
			select {
			case tradesMap := <-s.tradeChan:
				s.tradesMap.MSet(tradesMap)
			case <-s.refreshSnapshotsTicker.C:
				if time.Now().Second() == 2 {
					s.snapshotsCollection.Refresh()
				}
			case <-ctx.Done():
				s.publishTicker.Stop()
			case <-s.publishTicker.C:
				s.mu.Lock()
				s.snapshotsCollection.UpdateLatestTrades(s.tradesMap)
				messageBus <- messages.Message{
					Event: messages.Snapshots,
					Data:  s.snapshotsCollection.Items(),
				}
				s.mu.Unlock()
			}
		}
	}()

	s.loadDeck()

	return s
}

func (s *SnapshotStream) loadDeck() {
	s.mu.Lock()
	defer s.mu.Unlock()

	deck, err := s.deckRepo.FindByName("default")

	if err != nil {
		logrus.Panicf("failed to load deck: %v", err)
	}

	if deck.Symbols == "" {
		logrus.Warnf("no symbols found in deck")
		return
	}

	symbols := strings.Split(deck.Symbols, ",")

	s.tradeStream.UpdateSymbols(symbols)
	s.snapshotsCollection.UpdateSymbols(symbols)
}

func (s *SnapshotStream) UpdateSymbols(symbols []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.publishTicker.Stop()
	defer s.publishTicker.Reset(s.publishDuration)

	s.updateDeck(symbols)

	s.tradeStream.UpdateSymbols(symbols)
	s.snapshotsCollection.UpdateSymbols(symbols)
}

func (s *SnapshotStream) updateDeck(symbols []string) {
	if len(symbols) == 0 {
		_, err := s.deckRepo.ClearByName("default")

		if err != nil {
			logrus.Panicf("failed to clear symbols: %v", err)
		}
	} else {

		_, err := s.deckRepo.UpdateByName("default", symbols)

		if err != nil {
			logrus.Panicf("failed to update symbols: %v", err)
		}
	}
}
