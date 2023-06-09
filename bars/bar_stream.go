package bars

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/clients"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Stream struct {
	mu            sync.RWMutex
	sc            *stream.StocksClient
	symbols       []string
	bars          cmap.ConcurrentMap[string, Bar]
	streamChan    chan stream.Bar
	barChan       chan Bar
	unpublished   bool
	publishTicker *time.Ticker
}

func NewBarStream(ctx context.Context, publishChan chan<- map[string]Bar) *Stream {
	s := &Stream{
		sc:            clients.GetStocksClient(),
		symbols:       make([]string, 0),
		streamChan:    make(chan stream.Bar, 1_000),
		barChan:       make(chan Bar, 1_000),
		unpublished:   true,
		publishTicker: time.NewTicker(5 * time.Second),
		bars:          cmap.New[Bar](),
	}

	go func(ctx context.Context, s *Stream, publishChan chan<- map[string]Bar) {
		for {
			select {
			case b := <-s.streamChan:
				s.barChan <- FromStreamBar(b)
			case bar := <-s.barChan:
				s.mu.Lock()
				s.unpublished = true
				s.bars.Set(bar.Symbol, bar)
				s.mu.Unlock()
			case <-s.publishTicker.C:
				s.mu.Lock()
				if s.unpublished {
					publishChan <- s.bars.Items()
				}

				s.unpublished = false
				s.mu.Unlock()
			case <-ctx.Done():
				s.publishTicker.Stop()
				err := s.sc.UnsubscribeFromBars(s.bars.Keys()...)
				if err != nil {
					logrus.Errorf("failed to unsubscribe from bars: %s", err)
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
		err := s.sc.UnsubscribeFromBars(removedSymbols...)

		if err != nil {
			logrus.Errorf("failed to unsubscribe from bars: %s", err)
		}

		for _, symbol := range removedSymbols {
			s.bars.Remove(symbol)
		}
	}

	if len(addedSymbols) > 0 {
		err := s.sc.SubscribeToBars(
			func(b stream.Bar) {
				s.streamChan <- b
			}, addedSymbols...,
		)

		if err != nil {
			logrus.Errorf("failed to subscribe to bars: %s", err)
		}
	}

	s.unpublished = true
}
