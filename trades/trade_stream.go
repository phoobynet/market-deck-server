package trades

import (
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
	trades        map[string]Trade
	streamChan    chan stream.Trade
	tradeChan     chan Trade
	unpublished   bool
	publishTicker *time.Ticker
	publishChan   chan<- map[string]Trade
}

func NewTradeStream(sc *stream.StocksClient, publishChan chan<- map[string]Trade) *Stream {
	s := &Stream{
		sc:            sc,
		symbols:       make([]string, 0),
		streamChan:    make(chan stream.Trade, 1_000),
		tradeChan:     make(chan Trade, 1_000),
		unpublished:   true,
		publishTicker: time.NewTicker(1 * time.Second),
		publishChan:   publishChan,
	}

	go func(s *Stream) {
		for {
			select {
			case b := <-s.streamChan:
				s.tradeChan <- FromStreamTrade(b)
			case trade := <-s.tradeChan:
				s.mu.Lock()
				s.unpublished = true
				s.trades[trade.Symbol] = trade
				s.mu.Unlock()
			case <-s.publishTicker.C:
				s.mu.Lock()
				if s.unpublished {
					s.publishChan <- s.trades
				}

				s.unpublished = false
				s.mu.Unlock()
			}
		}
	}(s)

	return s
}

func (s *Stream) Update(symbols []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	removedSymbols, addedSymbols := lo.Difference(s.symbols, symbols)

	if len(removedSymbols) > 0 {
		err := s.sc.UnsubscribeFromTrades(removedSymbols...)

		if err != nil {
			logrus.Errorf("failed to unsubscribe from bars: %s", err)
		}

		for _, symbol := range removedSymbols {
			delete(s.trades, symbol)
		}
	}

	if len(addedSymbols) > 0 {
		err := s.sc.SubscribeToTrades(
			func(t stream.Trade) {
				s.streamChan <- t
			}, addedSymbols...,
		)

		if err != nil {
			logrus.Errorf("failed to subscribe to trades: %s", err)
		}
	}

	s.unpublished = true
}
