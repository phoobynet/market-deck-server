package snapshots

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/go-co-op/gocron"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
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
	snapshots           map[string]Snapshot
	bars                map[string][]bars.Bar
	publishTicker       *time.Ticker
	publishInterval     time.Duration
	snapshotScheduler   *gocron.Scheduler
	barRepository       *bars.Repository
}

func NewSnapshotStream(
	ctx context.Context,
	sc *stream.StocksClient,
	snapshotsRepository *Repository,
	deckRepository *decks.DeckRepository,
	barRepository *bars.Repository,
	messageBus chan<- messages.Message,
) *Stream {
	s := &Stream{
		deckRepository:      deckRepository,
		snapshotsRepository: snapshotsRepository,
		snapshots:           make(map[string]Snapshot),
		publishInterval:     1 * time.Second,
		publishTicker:       time.NewTicker(1 * time.Second),
		barChan:             make(chan map[string]bars.Bar, 1_000),
		quoteChan:           make(chan map[string]quotes.Quote, 1_000),
		tradeChan:           make(chan map[string]trades.Trade, 1_000),
		bars:                make(map[string][]bars.Bar, 0),
		barRepository:       barRepository,
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
			case barsMap := <-l.barChan:
				l.mu.Lock()
				for symbol, bar := range barsMap {
					if snapshot, ok := l.snapshots[symbol]; ok {
						snapshot.LatestBar = bar
						l.snapshots[symbol] = snapshot
					}
					l.bars[symbol] = append(l.bars[symbol], bar)
				}
				l.mu.Unlock()
			case quotesMap := <-l.quoteChan:
				l.mu.Lock()
				for symbol, quote := range quotesMap {
					if snapshot, ok := l.snapshots[symbol]; ok {
						snapshot.LatestQuote = quote
						l.snapshots[symbol] = snapshot
					}
				}
				l.mu.Unlock()
			case tradesMap := <-l.tradeChan:
				l.mu.Lock()

				for symbol, latestTrade := range tradesMap {
					if snapshot, ok := l.snapshots[symbol]; ok {
						diff := numbers.NumberDiff(l.snapshots[symbol].PreviousClose, latestTrade.Price)

						snapshot.Change = diff.Change
						snapshot.ChangePercent = diff.ChangePercent
						snapshot.LatestTrade = latestTrade
						l.snapshots[symbol] = snapshot
					}

				}

				l.mu.Unlock()
			case <-l.publishTicker.C:
				l.mu.Lock()
				messageBus <- messages.Message{
					Event: messages.Snapshots,
					Data:  l.snapshots,
				}
				l.mu.Unlock()
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
			delete(s.bars, symbol)
		}
	}

	if len(addedSymbols) > 0 {
		snapshots, err := s.snapshotsRepository.GetMulti(symbols)

		if err != nil {
			logrus.Errorf("failed to get snapshots: %v\n", err)
		}

		for symbol, snapshot := range snapshots {
			s.snapshots[symbol] = snapshot
			logrus.Infof("snapshot: %v", s.snapshots)
		}

		go s.fillIntradayBars(symbols)
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
		if existingSnapshot, ok := s.snapshots[symbol]; ok {
			existingSnapshot.DailyBar = snapshot.DailyBar
			existingSnapshot.PreviousDailyBar = snapshot.PreviousDailyBar
			existingSnapshot.PreviousClose = snapshot.PreviousClose
			s.snapshots[symbol] = existingSnapshot
		}
	}
}

func (s *Stream) fillIntradayBars(symbols []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	intradayMulti, err := s.barRepository.GetIntradayMulti(symbols)

	if err != nil {
		logrus.Errorf("failed to get intraday bars: %v", err)
	}

	for symbol, intraday := range intradayMulti {
		s.bars[symbol] = intraday
	}
}

func (s *Stream) GetIntradayBars() map[string][]bars.Bar {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.bars
}
