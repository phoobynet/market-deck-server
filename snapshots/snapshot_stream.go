package snapshots

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/server"
	"github.com/phoobynet/market-deck-server/trades"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Stream struct {
	deckRepository      *decks.DeckRepository
	barStream           *bars.Stream
	quoteStream         *quotes.Stream
	tradeStream         *trades.Stream
	barChan             chan map[string]bars.Bar
	quoteChan           chan map[string]quotes.Quote
	tradeChan           chan map[string]trades.Trade
	snapshotsRepository *Repository
	mu                  sync.RWMutex
	snapshots           map[string]*Snapshot
	publishTicker       *time.Ticker
	publishInterval     time.Duration
}

func NewSnapshotStream(
	sc *stream.StocksClient,
	snapshotsRepository *Repository,
	deckRepository *decks.DeckRepository,
	messageBus chan<- server.Message,
) *Stream {

	l := &Stream{
		deckRepository:      deckRepository,
		snapshotsRepository: snapshotsRepository,
		snapshots:           make(map[string]*Snapshot),
		publishInterval:     1 * time.Second,
		publishTicker:       time.NewTicker(1 * time.Second),
		barChan:             make(chan map[string]bars.Bar, 1_000),
		quoteChan:           make(chan map[string]quotes.Quote, 1_000),
		tradeChan:           make(chan map[string]trades.Trade, 1_000),
	}

	go func(l *Stream) {
		for {
			select {
			case latestBar := <-l.barChan:
				for symbol, bar := range latestBar {
					l.snapshots[symbol].LatestBar = bar
				}

			case latestQuote := <-l.quoteChan:
				for symbol, quote := range latestQuote {
					l.snapshots[symbol].LatestQuote = quote
				}

			case lastTrade := <-l.tradeChan:
				for symbol, trade := range lastTrade {
					l.snapshots[symbol].LatestTrade = trade
				}
			case <-l.publishTicker.C:
				l.mu.RLock()
				messageBus <- server.Message{
					Event: server.Snapshots,
					Data:  l.snapshots,
				}
				l.mu.RUnlock()
			}
		}
	}(l)

	return l
}

func (l *Stream) UpdateSymbols(symbols []string) {
	l.publishTicker.Stop()
	defer l.publishTicker.Reset(l.publishInterval)

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err := l.deckRepository.UpdateByName("default", symbols)

	if err != nil {
		logrus.Fatalf("failed to update symbols: %v", err)
	}

	oldSymbols := lo.Keys(l.snapshots)

	removedSymbols, addedSymbols := lo.Difference(oldSymbols, symbols)

	if len(removedSymbols) > 0 {
		for _, symbol := range removedSymbols {
			delete(l.snapshots, symbol)
		}
	}

	if len(addedSymbols) > 0 {
		snapshots, err := l.snapshotsRepository.GetMulti(symbols)

		if err != nil {
			logrus.Errorf("failed to get snapshots: %v\n", err)
		}

		for symbol, snapshot := range snapshots {
			l.snapshots[symbol] = &snapshot
		}
	}

	l.barStream.Update(symbols)
	l.quoteStream.Update(symbols)
	l.tradeStream.Update(symbols)
}
