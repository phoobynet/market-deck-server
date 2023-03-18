package main

import (
	"embed"
	"flag"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/database"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/events"
	"github.com/phoobynet/market-deck-server/realtime"
	"github.com/phoobynet/market-deck-server/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

var (
	//go:embed dist
	dist     embed.FS
	quitChan = make(chan os.Signal, 1)
)

func main() {
	signal.Notify(quitChan, os.Interrupt)
	logrus.SetFormatter(&logrus.TextFormatter{})

	var configPath string

	flag.StringVar(&configPath, "config", "config.toml", "Path to config file")

	config, err := server.LoadConfig(configPath)

	if err != nil {
		logrus.Fatalf("Error loading config: %s", err)
	}

	stocksClient := stream.NewStocksClient(marketdata.SIP)
	marketDataClient := marketdata.NewClient(marketdata.ClientOpts{})
	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{})

	assetRepository := assets.NewAssetRepository(database.GetDB(), alpacaClient)
	calendarDayRepository := calendars.NewCalendarDayRepository(database.GetDB(), alpacaClient)
	deckRepository := decks.NewDeckRepository(database.GetDB())

	realtimeSymbolsChan := make(chan map[string]*realtime.Symbol, 100)
	calendarDayUpdateChan := make(chan calendars.CalendarDayUpdate, 100)

	realTimeSymbols := realtime.NewLiveSymbols(
		realtimeSymbolsChan,
		alpacaClient,
		marketDataClient,
		stocksClient,
		1*time.Second,
		calendarDayRepository,
		assetRepository,
	)

	calendars.NewCalendarDayLive(calendarDayUpdateChan, alpacaClient, calendarDayRepository)

	go func() {
		for realtimeSymbols := range realtimeSymbolsChan {
			server.Publish(events.RealtimeSymbols, realtimeSymbols)
		}
	}()

	go func() {
		for calendarDayUpdate := range calendarDayUpdateChan {
			server.Publish(events.CalendarDayUpdate, calendarDayUpdate)
		}
	}()

	server.InitSSE()

	webServer := server.NewServer(config, dist, realTimeSymbols, deckRepository)

	webServer.Listen()
}
