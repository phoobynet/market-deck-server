package snapshots

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/golang-module/carbon/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/messages"
	"github.com/phoobynet/market-deck-server/trades"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"math"
	"sync"
	"time"
)

type SnapshotStream struct {
	mu                     sync.RWMutex
	deckRepo               *decks.DeckRepository
	snapshotsRepo          *Repository
	publishDuration        time.Duration
	publishTicker          *time.Ticker
	snapshotsLite          cmap.ConcurrentMap[string, Snapshot]
	barStream              *bars.Stream
	barChan                chan map[string]bars.Bar
	intradayBars           cmap.ConcurrentMap[string, []bars.Bar]
	tradeStream            *trades.Stream
	tradeChan              chan map[string]trades.Trade
	latestTrades           cmap.ConcurrentMap[string, trades.Trade]
	assetRepo              *assets.Repository
	symbols                []string
	calendarDayLive        *calendars.CalendarDayLive
	refreshSnapshotsTicker *time.Ticker
	ytdBars                []bars.Bar
	barRepo                *bars.Repository
}

func NewSnapshotLiteStream(
	ctx context.Context,
	calendarDayLive *calendars.CalendarDayLive,
	messageBus chan<- messages.Message,
) *SnapshotStream {
	s := &SnapshotStream{
		deckRepo:               decks.GetRepository(),
		snapshotsRepo:          GetRepository(),
		publishDuration:        500 * time.Millisecond,
		barChan:                make(chan map[string]bars.Bar, 1_000),
		tradeChan:              make(chan map[string]trades.Trade, 1_000),
		snapshotsLite:          cmap.New[Snapshot](),
		latestTrades:           cmap.New[trades.Trade](),
		assetRepo:              assets.GetRepository(),
		calendarDayLive:        calendarDayLive,
		refreshSnapshotsTicker: time.NewTicker(1 * time.Second),
		ytdBars:                make([]bars.Bar, 0),
		barRepo:                bars.GetRepository(),
		intradayBars:           cmap.New[[]bars.Bar](),
	}

	s.barStream = bars.NewBarStream(ctx, s.barChan)
	s.tradeStream = trades.NewTradeStream(ctx, s.tradeChan)
	s.publishTicker = time.NewTicker(s.publishDuration)

	go func() {
		for {
			select {
			case tradesMap := <-s.tradeChan:
				s.latestTrades.MSet(tradesMap)
			case barsMap := <-s.barChan:
				for symbol, latestBar := range barsMap {
					if snapshotLite, found := s.snapshotsLite.Get(symbol); found {
						newSnapshotLite := Snapshot{
							Symbol:        snapshotLite.Symbol,
							Name:          snapshotLite.Name,
							Exchange:      snapshotLite.Exchange,
							Price:         snapshotLite.Price,
							PrevClose:     snapshotLite.PrevClose,
							PrevCloseDate: snapshotLite.PrevCloseDate,
							DailyHigh:     math.Max(snapshotLite.DailyHigh, latestBar.High),
							DailyLow:      math.Min(snapshotLite.DailyLow, latestBar.Low),
							DailyVolume:   snapshotLite.DailyVolume + latestBar.Volume,
							Change:        snapshotLite.Change,
							Class:         snapshotLite.Class,
							Volumes:       snapshotLite.Volumes,
							MonthlyBars:   snapshotLite.MonthlyBars,
							YtdBars:       snapshotLite.YtdBars,
							YtdChange:     snapshotLite.YtdChange,
						}

						s.snapshotsLite.Set(symbol, newSnapshotLite)
					}
				}
			case <-s.refreshSnapshotsTicker.C:
				if time.Now().Second() == 2 {
					s.refreshSnapshots()
				}
			case <-ctx.Done():
				s.publishTicker.Stop()
			case <-s.publishTicker.C:
				s.mu.Lock()

				for _, symbol := range s.symbols {
					if latestTrade, ok := s.latestTrades.Get(symbol); ok {
						if snapshotLite, found := s.snapshotsLite.Get(symbol); found {
							newSnapshotLite := Snapshot{
								Symbol:        snapshotLite.Symbol,
								Name:          snapshotLite.Name,
								Exchange:      snapshotLite.Exchange,
								Price:         latestTrade.Price,
								PrevClose:     snapshotLite.PrevClose,
								PrevCloseDate: snapshotLite.PrevCloseDate,
								DailyHigh:     math.Max(snapshotLite.DailyHigh, latestTrade.Price),
								DailyLow:      math.Min(snapshotLite.DailyLow, latestTrade.Price),
								DailyVolume:   snapshotLite.DailyVolume,
								Change:        numbers.NumberDiff(snapshotLite.PrevClose, latestTrade.Price),
								Class:         snapshotLite.Class,
								Volumes:       snapshotLite.Volumes,
								MonthlyBars:   snapshotLite.MonthlyBars,
								YtdBars:       snapshotLite.YtdBars,
								YtdChange: numbers.NumberDiff(
									snapshotLite.YtdBars[0].Close,
									latestTrade.Price,
								),
							}

							s.snapshotsLite.Set(symbol, newSnapshotLite)
						}
					}
				}

				messageBus <- messages.Message{
					Event: messages.SnapshotsLite,
					Data:  s.snapshotsLite.Items(),
				}
				s.mu.Unlock()
			}
		}
	}()

	return s
}

func (s *SnapshotStream) UpdateSymbols(symbols []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.publishTicker.Stop()
	defer s.publishTicker.Reset(s.publishDuration)

	s.updateDeck(symbols)

	removedSymbols, addedSymbols := lo.Difference(s.symbols, symbols)

	s.removeSymbols(removedSymbols)
	s.addSymbols(addedSymbols)
	s.tradeStream.Update(s.symbols)
}

func (s *SnapshotStream) removeSymbols(removedSymbols []string) {
	if len(removedSymbols) == 0 {
		return
	}

	s.symbols = lo.Filter(
		s.symbols, func(symbol string, _ int) bool {
			return !lo.Contains(removedSymbols, symbol)
		},
	)

	for _, symbol := range removedSymbols {
		s.snapshotsLite.Remove(symbol)
	}
}

func (s *SnapshotStream) addSymbols(addedSymbols []string) {
	if len(addedSymbols) == 0 {
		return
	}

	s.symbols = append(s.symbols, addedSymbols...)

	multiSnapshots, err := s.snapshotsRepo.GetMulti(addedSymbols)

	if err != nil {
		logrus.Panicf("failed to get snapshots: %v", err)
	}

	lastMonthOfBars, err := s.barRepo.GetHistoricMulti(
		addedSymbols,
		marketdata.OneDay,
		carbon.Now(carbon.NewYork).StartOfDay().SubMonths(1).ToStdTime(),
		carbon.Now(carbon.NewYork).StartOfDay().ToStdTime(),
	)

	if err != nil {
		logrus.Panicf("failed to get last month of bars: %v", err)
	}

	yearToDateBars, err := s.barRepo.GetHistoricMulti(
		addedSymbols,
		marketdata.OneMonth,
		carbon.Now(carbon.NewYork).StartOfDay().SubMonths(12).ToStdTime(),
		carbon.Now(carbon.NewYork).StartOfDay().ToStdTime(),
	)

	if err != nil {
		logrus.Panicf("failed to get year to date bars: %v", err)
	}

	logrus.Infof("last month of bars: %v", yearToDateBars)

	calendarDayUpdate := s.calendarDayLive.Get()

	intradayBarsBySymbol := cmap.New[[]bars.Bar]()

	// daily high, low and volume must be acquired from the minute bars
	if calendarDayUpdate.Condition == calendars.PreMarket {
		intradayBars, err := s.barRepo.GetIntradayMulti(addedSymbols)

		if err != nil {
			logrus.Panicf("failed to get intraday bars: %v", err)
		}

		intradayBarsBySymbol.MSet(intradayBars)
	}

	var asset *assets.Asset

	for symbol, snapshot := range multiSnapshots {
		asset = s.assetRepo.Get(symbol)

		if asset == nil {
			logrus.Panicf("failed to get asset: %v", err)
		}

		volumes := make([]SnapshotVolume, 0)

		if barsForSymbol, ok := lastMonthOfBars[symbol]; ok {
			for _, bar := range barsForSymbol {
				volumes = append(
					volumes, SnapshotVolume{
						Date:   bar.Date(),
						Volume: bar.Volume,
					},
				)
			}
		}

		change := numbers.NumberDiff(snapshot.ActualPreviousDailyBar.Close, snapshot.LatestTrade.Price)

		var dailyHigh float64
		var dailyLow float64
		var dailyVolume float64

		if calendarDayUpdate.Condition == calendars.PreMarket {
			if intradayBars, ok := intradayBarsBySymbol.Get(symbol); ok {
				highs := make([]float64, 0)
				lows := make([]float64, 0)
				volume := 0.0

				for _, bar := range intradayBars {
					highs = append(highs, bar.High)
					lows = append(lows, bar.Low)
					volume += bar.Volume
				}

				dailyHigh = lo.Max(highs)
				dailyLow = lo.Min(lows)
				dailyVolume = volume
			}
		} else {
			dailyHigh = snapshot.DailyBar.High
			dailyLow = snapshot.DailyBar.Low
			dailyVolume = snapshot.DailyBar.Volume
		}

		s.snapshotsLite.Set(
			symbol, Snapshot{
				Symbol:        symbol,
				Name:          asset.Name,
				Exchange:      asset.Exchange,
				Price:         snapshot.LatestTrade.Price,
				PrevClose:     snapshot.ActualPreviousDailyBar.Close,
				PrevCloseDate: snapshot.ActualPreviousDailyBar.Date(),
				DailyHigh:     dailyHigh,
				DailyLow:      dailyLow,
				DailyVolume:   dailyVolume,
				Change:        change,
				Class:         asset.Class,
				Volumes:       volumes,
				MonthlyBars:   yearToDateBars[symbol],
				YtdBars:       yearToDateBars[symbol],
				YtdChange:     numbers.NumberDiff(yearToDateBars[symbol][0].Close, snapshot.LatestTrade.Price),
			},
		)
	}
}

// refreshSnapshots updates the snapshots lite from source
// should be run every minute on or just after the minute
func (s *SnapshotStream) refreshSnapshots() {
	// TODO: Duplicate code - refactor
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.symbols) == 0 {
		return
	}

	multiSnapshots, err := s.snapshotsRepo.GetMulti(s.symbols)

	if err != nil {
		logrus.Panicf("failed to get snapshots: %v", err)
	}

	calendarDayUpdate := s.calendarDayLive.Get()

	intradayBarsBySymbol := cmap.New[[]bars.Bar]()

	// daily high, low and volume must be acquired from the minute bars
	if calendarDayUpdate.Condition == calendars.PreMarket {
		intradayBars, err := s.barRepo.GetIntradayMulti(s.symbols)

		if err != nil {
			logrus.Panicf("failed to get intraday bars: %v", err)
		}

		intradayBarsBySymbol.MSet(intradayBars)
	}

	for symbol, snapshot := range multiSnapshots {
		if snapshotLite, found := s.snapshotsLite.Get(symbol); found {

			var dailyHigh float64
			var dailyLow float64
			var dailyVolume float64

			if calendarDayUpdate.Condition == calendars.PreMarket {
				if intradayBars, ok := intradayBarsBySymbol.Get(symbol); ok {
					highs := make([]float64, 0)
					lows := make([]float64, 0)
					volume := 0.0

					for _, bar := range intradayBars {
						highs = append(highs, bar.High)
						lows = append(lows, bar.Low)
						volume += bar.Volume
					}

					dailyHigh = lo.Max(highs)
					dailyLow = lo.Min(lows)
					dailyVolume = volume
				}
			} else {
				dailyHigh = snapshot.DailyBar.High
				dailyLow = snapshot.DailyBar.Low
				dailyVolume = snapshot.DailyBar.Volume
			}
			change := numbers.NumberDiff(snapshot.ActualPreviousDailyBar.Close, snapshot.LatestTrade.Price)

			newSnapshotLite := Snapshot{
				Symbol:        snapshotLite.Symbol,
				Name:          snapshotLite.Name,
				Exchange:      snapshotLite.Exchange,
				Price:         snapshotLite.Price,
				PrevClose:     snapshot.ActualPreviousDailyBar.Close,
				PrevCloseDate: snapshot.ActualPreviousDailyBar.Date(),
				DailyHigh:     dailyHigh,
				DailyLow:      dailyLow,
				DailyVolume:   dailyVolume,
				Change:        change,
				Class:         snapshotLite.Class,
				YtdBars:       snapshotLite.YtdBars,
				Volumes:       snapshotLite.Volumes,
			}

			s.snapshotsLite.Set(symbol, newSnapshotLite)
		}
	}
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
