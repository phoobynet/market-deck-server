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
	publishTicker       *time.Ticker
	publishInterval     time.Duration
	snapshotScheduler   *gocron.Scheduler
	barRepository       *bars.Repository

	mu sync.RWMutex

	// live and frequently updated data which is used to calculate snapshots
	latestBars     map[string]bars.Bar
	latestTrades   map[string]trades.Trade
	latestQuotes   map[string]quotes.Quote
	dailyBars      map[string]bars.Bar
	prevDailyBars  map[string]bars.Bar
	previousCloses map[string]float64
	bars           map[string][]bars.Bar

	symbols []string
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
		publishInterval:     1 * time.Second,
		publishTicker:       time.NewTicker(1 * time.Second),
		barChan:             make(chan map[string]bars.Bar, 1_000),
		quoteChan:           make(chan map[string]quotes.Quote, 1_000),
		tradeChan:           make(chan map[string]trades.Trade, 1_000),
		bars:                make(map[string][]bars.Bar, 0),
		barRepository:       barRepository,
		symbols:             make([]string, 0),
		latestBars:          make(map[string]bars.Bar, 0),
		latestTrades:        make(map[string]trades.Trade, 0),
		latestQuotes:        make(map[string]quotes.Quote, 0),
		dailyBars:           make(map[string]bars.Bar, 0),
		prevDailyBars:       make(map[string]bars.Bar, 0),
		previousCloses:      make(map[string]float64, 0),
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
		Do(
			func() {
				s.refreshSnapshot(s.symbols)
			},
		)

	if err != nil {
		logrus.Errorf("failed to create snapshot job: %v", err)
	}

	go func(s *Stream) {
		for {
			select {
			case barsMap := <-s.barChan:
				s.mu.Lock()

				s.latestBars = barsMap

				for _, symbol := range s.symbols {
					s.bars[symbol] = append(s.bars[symbol], barsMap[symbol])
				}

				s.mu.Unlock()
			case quotesMap := <-s.quoteChan:
				s.mu.Lock()
				s.latestQuotes = quotesMap
				s.mu.Unlock()
			case tradesMap := <-s.tradeChan:
				s.mu.Lock()
				s.latestTrades = tradesMap
				s.mu.Unlock()
			case <-s.publishTicker.C:
				s.mu.RLock()
				data := make(map[string]Snapshot)

				for _, symbol := range s.symbols {
					snapshot := Snapshot{
						LatestBar:        s.latestBars[symbol],
						LatestQuote:      s.latestQuotes[symbol],
						LatestTrade:      s.latestTrades[symbol],
						DailyBar:         s.dailyBars[symbol],
						PreviousDailyBar: s.prevDailyBars[symbol],
						PreviousClose:    s.previousCloses[symbol],
					}

					diff := numbers.NumberDiff(snapshot.LatestTrade.Price, snapshot.PreviousClose)

					snapshot.Change = diff.Change
					snapshot.ChangePercent = diff.ChangePercent
					snapshot.ChangeAbs = diff.AbsoluteChange
					snapshot.ChangeSign = diff.Sign

					data[symbol] = snapshot
				}

				messageBus <- messages.Message{
					Event: messages.Snapshots,
					Data:  data,
				}
				s.mu.RUnlock()
			case <-ctx.Done():
				s.snapshotScheduler.Stop()
				s.publishTicker.Stop()
			}
		}
	}(s)

	return s
}

func (s *Stream) UpdateSymbols(symbols []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.publishTicker.Stop()
	defer s.publishTicker.Reset(s.publishInterval)

	_, err := s.deckRepository.UpdateByName("default", symbols)

	if err != nil {
		logrus.Fatalf("failed to update symbols: %v", err)
	}

	removedSymbols, addedSymbols := lo.Difference(s.symbols, symbols)

	if len(removedSymbols) > 0 {
		s.symbols = lo.Filter(
			s.symbols, func(symbol string, _ int) bool {
				return !lo.Contains(removedSymbols, symbol)
			},
		)

		for _, symbol := range removedSymbols {
			delete(s.latestBars, symbol)
			delete(s.latestQuotes, symbol)
			delete(s.latestBars, symbol)
			delete(s.prevDailyBars, symbol)
			delete(s.previousCloses, symbol)
			delete(s.dailyBars, symbol)
			delete(s.bars, symbol)
		}
	}

	if len(addedSymbols) > 0 {
		s.symbols = append(s.symbols, addedSymbols...)

		s.refreshSnapshot(addedSymbols)

		// run in the background is fine
		go s.fillIntradayBars(addedSymbols)
	}

	s.barStream.Update(symbols)
	s.quoteStream.Update(symbols)
	s.tradeStream.Update(symbols)

}

func (s *Stream) refreshSnapshot(symbols []string) {
	snapshots, err := s.snapshotsRepository.GetMulti(symbols)

	if err != nil {
		logrus.Errorf("failed to get snapshots: %v", err)
	}

	for _, symbol := range symbols {
		s.dailyBars[symbol] = snapshots[symbol].DailyBar
		s.prevDailyBars[symbol] = snapshots[symbol].PreviousDailyBar
		s.previousCloses[symbol] = snapshots[symbol].PreviousClose
	}
}

func (s *Stream) fillIntradayBars(symbols []string) {
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
