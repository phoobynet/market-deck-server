package realtime

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Symbols struct {
	mu                    sync.RWMutex
	symbolMap             map[string]*Symbol
	alpacaClient          *alpaca.Client
	marketDataClient      *marketdata.Client
	stocksClient          *stream.StocksClient
	publishTicker         *time.Ticker
	publishInterval       time.Duration
	publishChan           chan<- map[string]*Symbol
	tradeChan             chan stream.Trade
	quoteChan             chan stream.Quote
	barChan               chan stream.Bar
	calendarDayRepository *calendars.CalendarDayRepository
	assetRepository       *assets.AssetRepository
}

func NewLiveSymbols(
	realtimeSymbolsChan chan<- map[string]*Symbol,
	alpacaClient *alpaca.Client,
	marketDataClient *marketdata.Client,
	stocksClient *stream.StocksClient,
	publishInterval time.Duration,
	calendarDayRepository *calendars.CalendarDayRepository,
	assetRepository *assets.AssetRepository,
) *Symbols {
	l := &Symbols{
		symbolMap:             make(map[string]*Symbol),
		alpacaClient:          alpacaClient,
		marketDataClient:      marketDataClient,
		stocksClient:          stocksClient,
		publishTicker:         time.NewTicker(publishInterval),
		publishInterval:       publishInterval,
		publishChan:           realtimeSymbolsChan,
		tradeChan:             make(chan stream.Trade, 100_000),
		quoteChan:             make(chan stream.Quote, 100_000),
		barChan:               make(chan stream.Bar, 100_000),
		calendarDayRepository: calendarDayRepository,
		assetRepository:       assetRepository,
	}

	go func() {
		for range l.publishTicker.C {
			realtimeSymbolsChan <- l.symbolMap
		}
	}()

	return l
}

func (l *Symbols) updateTrade(streamTrade stream.Trade) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if liveSymbol, ok := l.symbolMap[streamTrade.Symbol]; ok {
		liveSymbol.Trade = trades.FromStreamTrade(streamTrade)
	}
}

func (l *Symbols) updateQuote(streamQuote stream.Quote) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if liveSymbol, ok := l.symbolMap[streamQuote.Symbol]; ok {
		liveSymbol.Quote = quotes.FromStreamQuote(streamQuote)
	}
}

func (l *Symbols) updateBar(streamBar stream.Bar) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if liveSymbol, ok := l.symbolMap[streamBar.Symbol]; ok {
		liveSymbol.Bar = bars.FromStreamBar(streamBar)
		liveSymbol.IntradayBars = append(liveSymbol.IntradayBars, liveSymbol.Bar)
	}
}

func (l *Symbols) UpdateSymbols(symbols []string) {
	logrus.Infof("updating symbols: %v...", symbols)
	l.publishTicker.Stop()
	defer l.publishTicker.Reset(l.publishInterval)

	oldSymbols := lo.Keys(l.symbolMap)

	removedSymbols, addedSymbols := lo.Difference(oldSymbols, symbols)

	if len(removedSymbols) > 0 {
		_ = l.stocksClient.UnsubscribeFromTrades(removedSymbols...)
		_ = l.stocksClient.UnsubscribeFromQuotes(removedSymbols...)
		_ = l.stocksClient.UnsubscribeFromBars(removedSymbols...)

		for _, symbol := range removedSymbols {
			delete(l.symbolMap, symbol)
		}
	}

	if len(addedSymbols) > 0 {
		checkErr := func(err error) {
			if err != nil {
				logrus.Error(err)
			}
		}

		logrus.Infof("Pre-loading data for %v...", addedSymbols)
		latestBars := l.getLatestBars(addedSymbols)
		latestQuotes := l.getLatestQuotes(addedSymbols)
		latestTrades := l.getLatestTrades(addedSymbols)
		previousDailyBars := l.getPreviousDailyBars(addedSymbols)

		logrus.Infof("Pre-loading data for %v...DONE", addedSymbols)
		//intradayBars := l.getIntradayBars(addedSymbols)
		//ytdDailyBars := l.getYtdDailyBars(addedSymbols)

		logrus.Infof("latest bars: %v", len(latestBars))

		for _, symbol := range addedSymbols {
			l.symbolMap[symbol] = &Symbol{
				Asset:        l.assetRepository.Get(symbol),
				Bar:          latestBars[symbol],
				Trade:        latestTrades[symbol],
				Quote:        latestQuotes[symbol],
				PrevDailyBar: previousDailyBars[symbol],
				//IntradayBars: intradayBars[symbol],
				//YtdDailyBars: ytdDailyBars[symbol],
			}
		}

		logrus.Infof("Subscribing to %v...", addedSymbols)
		checkErr(l.stocksClient.SubscribeToTrades(l.updateTrade, addedSymbols...))
		checkErr(l.stocksClient.SubscribeToQuotes(l.updateQuote, addedSymbols...))
		checkErr(l.stocksClient.SubscribeToBars(l.updateBar, addedSymbols...))
		logrus.Infof("Subscribing to %v...DONE", addedSymbols)
	}
}

func (l *Symbols) getLatestBars(symbols []string) map[string]bars.Bar {
	latestBars, err := l.marketDataClient.GetLatestBars(
		symbols, marketdata.GetLatestBarRequest{},
	)

	if err != nil {
		logrus.Fatalf("failed to get latest bars for %v:\n %v", symbols, err)
	}

	result := make(map[string]bars.Bar)

	for symbol, bar := range latestBars {
		result[symbol] = bars.FromMarketDataBar(symbol, bar)
	}

	return result
}

func (l *Symbols) getLatestTrades(symbols []string) map[string]trades.Trade {
	multiTrades, err := l.marketDataClient.GetLatestTrades(symbols, marketdata.GetLatestTradeRequest{})

	if err != nil {
		logrus.Fatalf("failed to get latest trades: %v: \n %v", symbols, err)
	}

	result := make(map[string]trades.Trade)

	for symbol, trade := range multiTrades {
		result[symbol] = trades.FromMarketDataTrade(symbol, trade)
	}

	return result
}

func (l *Symbols) getLatestQuotes(symbols []string) map[string]quotes.Quote {
	multiQuotes, err := l.marketDataClient.GetLatestQuotes(symbols, marketdata.GetLatestQuoteRequest{})

	if err != nil {
		logrus.Fatalf("failed to get quotes %v:\n %v", symbols, err)
	}

	result := make(map[string]quotes.Quote)

	for symbol, quote := range multiQuotes {
		result[symbol] = quotes.FromMarketDataQuote(symbol, quote)
	}

	return result
}

func (l *Symbols) getPreviousDailyBars(symbols []string) map[string]bars.Bar {
	previousMarketDate := l.calendarDayRepository.GetPrevious()

	result := make(map[string]bars.Bar)

	start := carbon.
		FromStdTime(time.UnixMicro(previousMarketDate.PreMarketOpen)).
		SetTimezone("America/New_York").
		StartOfDay().
		ToStdTime()

	multiBars, err := l.marketDataClient.GetMultiBars(
		symbols, marketdata.GetBarsRequest{
			TimeFrame: marketdata.TimeFrame{
				Unit: marketdata.Day,
				N:    1,
			},
			Start: start,
		},
	)

	if err != nil {
		logrus.Panic(err)
	}

	for symbol, symbolBars := range multiBars {
		for _, bar := range symbolBars {
			result[symbol] = bars.FromMarketDataBar(symbol, bar)
		}
	}

	return result
}

func (l *Symbols) getIntradayBars(symbols []string) map[string][]bars.Bar {
	marketDate := l.calendarDayRepository.GetToday()

	if marketDate == nil {
		marketDate = l.calendarDayRepository.GetPrevious()
	}

	start := carbon.
		FromStdTime(time.UnixMicro(marketDate.PreMarketOpen)).
		SetTimezone("America/New_York").
		ToStdTime()

	end := carbon.FromStdTime(time.UnixMicro(marketDate.PostMarketClose)).ToStdTime()

	return l.getMultiBars(
		symbols, marketdata.TimeFrame{
			Unit: marketdata.Min,
			N:    1,
		},
		start,
		end,
	)
}

func (l *Symbols) getYtdDailyBars(symbols []string) map[string][]bars.Bar {
	start := time.Now().AddDate(0, 0, -365)

	end := time.Now()

	return l.getMultiBars(
		symbols, marketdata.TimeFrame{
			Unit: marketdata.Min,
			N:    1,
		}, start, end,
	)
}

func (l *Symbols) getMultiBars(
	symbols []string,
	timeframe marketdata.TimeFrame,
	start, end time.Time,
) map[string][]bars.Bar {
	multiBars, err := l.marketDataClient.GetMultiBars(
		symbols, marketdata.GetBarsRequest{
			TimeFrame: timeframe,
			Start:     start,
			End:       end,
		},
	)

	if err != nil {
		logrus.Fatalf("failed to get muliti bars for %v\n %v", symbols, err)
	}

	result := make(map[string][]bars.Bar)

	for symbol, symbolBars := range multiBars {
		for _, bar := range symbolBars {
			result[symbol] = append(result[symbol], bars.FromMarketDataBar(symbol, bar))
		}
	}

	return result
}
