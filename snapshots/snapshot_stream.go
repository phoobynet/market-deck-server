package snapshots

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	cmap "github.com/orcaman/concurrent-map/v2"
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
	refreshTicker       *time.Ticker
	barRepository       *bars.Repository

	mu sync.RWMutex

	// live and frequently updated data which is used to calculate snapshots
	latestBars     cmap.ConcurrentMap[string, bars.Bar]
	latestTrades   cmap.ConcurrentMap[string, trades.Trade]
	latestQuotes   cmap.ConcurrentMap[string, quotes.Quote]
	dailyBars      cmap.ConcurrentMap[string, bars.Bar]
	prevDailyBars  cmap.ConcurrentMap[string, bars.Bar]
	previousCloses cmap.ConcurrentMap[string, float64]
	bars           cmap.ConcurrentMap[string, []bars.Bar]

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
		bars:                cmap.New[[]bars.Bar](),
		barRepository:       barRepository,
		symbols:             make([]string, 0),
		latestBars:          cmap.New[bars.Bar](),
		latestTrades:        cmap.New[trades.Trade](),
		latestQuotes:        cmap.New[quotes.Quote](),
		dailyBars:           cmap.New[bars.Bar](),
		prevDailyBars:       cmap.New[bars.Bar](),
		previousCloses:      cmap.New[float64](),
		refreshTicker:       time.NewTicker(1 * time.Second),
	}

	s.barStream = bars.NewBarStream(ctx, sc, s.barChan)
	s.tradeStream = trades.NewTradeStream(ctx, sc, s.tradeChan)
	s.quoteStream = quotes.NewQuoteStream(ctx, sc, s.quoteChan)

	go func(s *Stream) {
		for {
			select {
			case barsMap := <-s.barChan:
				s.latestBars.MSet(barsMap)

				for _, symbol := range s.symbols {
					existingBars, _ := s.bars.Get(symbol)
					latestBar, _ := s.latestBars.Get(symbol)
					s.bars.Set(symbol, append(existingBars, latestBar))
				}
			case quotesMap := <-s.quoteChan:
				s.latestQuotes.MSet(quotesMap)
			case tradesMap := <-s.tradeChan:
				s.latestTrades.MSet(tradesMap)
			case <-s.publishTicker.C:
				s.mu.RLock()
				data := make(map[string]Snapshot)

				for _, symbol := range s.symbols {
					latestTrade, _ := s.latestTrades.Get(symbol)
					latestBar, _ := s.latestBars.Get(symbol)
					latestQuote, _ := s.latestQuotes.Get(symbol)
					dailyBar, _ := s.dailyBars.Get(symbol)
					prevDailyBar, _ := s.prevDailyBars.Get(symbol)
					previousClose, _ := s.previousCloses.Get(symbol)

					snapshot := Snapshot{
						LatestBar:        latestBar,
						LatestQuote:      latestQuote,
						LatestTrade:      latestTrade,
						DailyBar:         dailyBar,
						PreviousDailyBar: prevDailyBar,
						PreviousClose:    previousClose,
					}

					diff := numbers.NumberDiff(snapshot.PreviousClose, snapshot.LatestTrade.Price)

					snapshot.Change = diff.Change
					snapshot.ChangePercent = diff.ChangePercent
					snapshot.ChangeAbs = diff.AbsoluteChange
					snapshot.ChangeSign = diff.Sign

					data[symbol] = snapshot
				}

				s.mu.RUnlock()

				messageBus <- messages.Message{
					Event: messages.Snapshots,
					Data:  data,
				}
			case <-s.refreshTicker.C:
				sec := time.Now().Second()

				if sec == 0 {
					logrus.Info("refreshing snapshots")
					s.mu.Lock()
					s.refreshSnapshot(s.symbols)
					s.mu.Unlock()
				}
			case <-ctx.Done():
				s.publishTicker.Stop()
				s.refreshTicker.Stop()
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
			s.latestTrades.Remove(symbol)
			s.latestBars.Remove(symbol)
			s.latestQuotes.Remove(symbol)
			s.latestBars.Remove(symbol)
			s.prevDailyBars.Remove(symbol)
			s.previousCloses.Remove(symbol)
			s.dailyBars.Remove(symbol)
			s.bars.Remove(symbol)
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
		s.dailyBars.Set(symbol, snapshots[symbol].DailyBar)
		s.prevDailyBars.Set(symbol, snapshots[symbol].PreviousDailyBar)
		s.previousCloses.Set(symbol, snapshots[symbol].PreviousClose)
	}
}

func (s *Stream) fillIntradayBars(symbols []string) {
	intradayMulti, err := s.barRepository.GetIntradayMulti(symbols)

	if err != nil {
		logrus.Errorf("failed to get intraday bars: %v", err)
	}

	for symbol, intraday := range intradayMulti {
		s.bars.Set(symbol, intraday)
	}
}

func (s *Stream) GetIntradayBars() map[string][]bars.Bar {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.bars.Items()
}
