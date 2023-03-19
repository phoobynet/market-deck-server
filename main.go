package main

import (
	"context"
	"embed"
	"flag"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/database"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/messages"
	"github.com/phoobynet/market-deck-server/server"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

var (
	//go:embed dist
	dist       embed.FS
	quitChan   = make(chan os.Signal, 1)
	messageBus = make(chan messages.Message, 100_00)
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

	stocksClient := stream.NewStocksClient(md.SIP)

	err = stocksClient.Connect(context.TODO())

	if err != nil {
		logrus.Fatalf("error connecting to stocks client: %v", err)
	}

	server.InitSSE()

	mdClient := md.NewClient(md.ClientOpts{})
	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{})

	assetRepository := assets.NewAssetRepository(database.GetDB(), alpacaClient)
	calendarDayRepository := calendars.NewCalendarDayRepository(database.GetDB(), alpacaClient)
	deckRepository := decks.NewDeckRepository(database.GetDB())
	snapshotRepository := snapshots.NewSnapshotRepository(mdClient)
	barRepository := bars.NewBarRepository(mdClient)

	ctx, cancel := context.WithCancel(context.Background())

	realTimeSymbols := snapshots.NewSnapshotStream(
		ctx,
		stocksClient,
		snapshotRepository,
		deckRepository,
		barRepository,
		messageBus,
	)

	calendars.NewCalendarDayLive(ctx, alpacaClient, calendarDayRepository, messageBus)

	webServer := server.NewServer(config, dist, realTimeSymbols, deckRepository, assetRepository)

	go func() {
		for {
			select {
			case message := <-messageBus:
				server.Publish(message)
			case <-quitChan:
				logrus.Info("Shutting down...")
				cancel()
				os.Exit(0)
			}
		}
	}()

	logrus.Infof("Listening on %d...", config.ServerPort)
	webServer.Listen()
}
