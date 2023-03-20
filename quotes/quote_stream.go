package quotes

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Stream struct {
	mu            sync.RWMutex
	sc            *stream.StocksClient
	symbols       []string
	quotes        cmap.ConcurrentMap[string, Quote]
	streamChan    chan stream.Quote
	quoteChan     chan Quote
	unpublished   bool
	publishTicker *time.Ticker
}

func NewQuoteStream(ctx context.Context, sc *stream.StocksClient, publishChan chan<- map[string]Quote) *Stream {
	s := &Stream{
		sc:            sc,
		symbols:       make([]string, 0),
		streamChan:    make(chan stream.Quote, 1_000),
		quoteChan:     make(chan Quote, 1_000),
		unpublished:   true,
		publishTicker: time.NewTicker(1 * time.Second),
		quotes:        cmap.New[Quote](),
	}

	go func(ctx context.Context, s *Stream, publishChan chan<- map[string]Quote) {
		for {
			select {
			case b := <-s.streamChan:
				s.quoteChan <- FromStreamQuote(b)
			case quote := <-s.quoteChan:
				s.mu.Lock()
				s.unpublished = true
				s.quotes.Set(quote.Symbol, quote)
				s.mu.Unlock()
			case <-s.publishTicker.C:
				s.mu.Lock()
				if s.unpublished {
					publishChan <- s.quotes.Items()
				}

				s.unpublished = false
				s.mu.Unlock()
			case <-ctx.Done():
				s.publishTicker.Stop()
				err := s.sc.UnsubscribeFromQuotes(s.quotes.Keys()...)

				if err != nil {
					logrus.Errorf("failed to unsubscribe from quotes: %s", err)
				}
			}
		}
	}(ctx, s, publishChan)

	return s
}

func (s *Stream) Update(symbols []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	removedSymbols, addedSymbols := lo.Difference(s.symbols, symbols)

	if len(removedSymbols) > 0 {
		err := s.sc.UnsubscribeFromQuotes(removedSymbols...)

		if err != nil {
			logrus.Errorf("failed to unsubscribe from bars: %s", err)
		}

		for _, symbol := range removedSymbols {
			s.quotes.Remove(symbol)
		}
	}

	if len(addedSymbols) > 0 {
		err := s.sc.SubscribeToQuotes(
			func(q stream.Quote) {
				s.streamChan <- q
			}, addedSymbols...,
		)

		if err != nil {
			logrus.Errorf("failed to subscribe to quotes: %s", err)
		}
	}

	s.unpublished = true
}
