package bars

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Stream struct {
	mu            sync.RWMutex
	sc            *stream.StocksClient
	symbols       []string
	bars          map[string]Bar
	streamChan    chan stream.Bar
	barChan       chan Bar
	unpublished   bool
	publishTicker *time.Ticker
}

func NewBarStream(ctx context.Context, sc *stream.StocksClient, publishChan chan<- map[string]Bar) *Stream {
	s := &Stream{
		sc:            sc,
		symbols:       make([]string, 0),
		streamChan:    make(chan stream.Bar, 1_000),
		barChan:       make(chan Bar, 1_000),
		unpublished:   true,
		publishTicker: time.NewTicker(5 * time.Second),
	}

	go func(ctx context.Context, s *Stream, publishChan chan<- map[string]Bar) {
		for {
			select {
			case b := <-s.streamChan:
				s.barChan <- FromStreamBar(b)
			case bar := <-s.barChan:
				s.mu.Lock()
				s.unpublished = true
				s.bars[bar.Symbol] = bar
				s.mu.Unlock()
			case <-s.publishTicker.C:
				s.mu.Lock()
				if s.unpublished {
					publishChan <- s.bars
				}

				s.unpublished = false
				s.mu.Unlock()
			case <-ctx.Done():
				s.publishTicker.Stop()
				err := s.sc.UnsubscribeFromBars(lo.Keys(s.bars)...)
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
			delete(s.bars, symbol)
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
