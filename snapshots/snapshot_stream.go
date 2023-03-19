package snapshots

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/go-co-op/gocron"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/messages"
	"github.com/phoobynet/market-deck-server/quotes"
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
	snapshotScheduler   *gocron.Scheduler
}

func NewSnapshotStream(
	ctx context.Context,
	sc *stream.StocksClient,
	snapshotsRepository *Repository,
	deckRepository *decks.DeckRepository,
	messageBus chan<- messages.Message,
) *Stream {
	s := &Stream{
		deckRepository:      deckRepository,
		snapshotsRepository: snapshotsRepository,
		snapshots:           make(map[string]*Snapshot),
		publishInterval:     1 * time.Second,
		publishTicker:       time.NewTicker(1 * time.Second),
		barChan:             make(chan map[string]bars.Bar, 1_000),
		quoteChan:           make(chan map[string]quotes.Quote, 1_000),
		tradeChan:           make(chan map[string]trades.Trade, 1_000),
	}

	s.barStream = bars.NewBarStream(ctx, sc, s.barChan)
	s.tradeStream = trades.NewTradeStream(ctx, sc, s.tradeChan)
	s.quoteStream = quotes.NewQuoteStream(ctx, sc, s.quoteChan)
	s.snapshotScheduler = gocron.NewScheduler(time.UTC)

	snapshotRefreshStartAt := carbon.
		Now().
		StartOfMinute().
		AddMinutes(1).
		ToStdTime()

	_, err := s.snapshotScheduler.
		Every(1).
		Minute().
		StartAt(snapshotRefreshStartAt).
		Do(s.refreshSnapshot)

	if err != nil {
		logrus.Errorf("failed to create snapshot job: %v", err)
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
				messageBus <- messages.Message{
					Event: messages.Snapshots,
					Data:  l.snapshots,
				}
				l.mu.RUnlock()
			case <-ctx.Done():
				l.snapshotScheduler.Stop()
				l.publishTicker.Stop()
			}
		}
	}(s)

	return s
}

func (s *Stream) UpdateSymbols(symbols []string) {
	s.publishTicker.Stop()
	defer s.publishTicker.Reset(s.publishInterval)

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.deckRepository.UpdateByName("default", symbols)

	if err != nil {
		logrus.Fatalf("failed to update symbols: %v", err)
	}

	oldSymbols := lo.Keys(s.snapshots)

	removedSymbols, addedSymbols := lo.Difference(oldSymbols, symbols)

	if len(removedSymbols) > 0 {
		for _, symbol := range removedSymbols {
			delete(s.snapshots, symbol)
		}
	}

	if len(addedSymbols) > 0 {
		snapshots, err := s.snapshotsRepository.GetMulti(symbols)

		if err != nil {
			logrus.Errorf("failed to get snapshots: %v\n", err)
		}

		for symbol, snapshot := range snapshots {
			s.snapshots[symbol] = &snapshot
		}
	}

	s.barStream.Update(symbols)
	s.quoteStream.Update(symbols)
	s.tradeStream.Update(symbols)
}

func (s *Stream) refreshSnapshot() {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshots, err := s.snapshotsRepository.GetMulti(lo.Keys(s.snapshots))

	if err != nil {
		logrus.Errorf("failed to get snapshots: %v", err)
	}

	for symbol, snapshot := range snapshots {
		s.snapshots[symbol].DailyBar = snapshot.DailyBar
		s.snapshots[symbol].PreviousDailyBar = snapshot.PreviousDailyBar
	}
}
