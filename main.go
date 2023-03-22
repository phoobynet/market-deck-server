package main

import (
	"context"
	"embed"
	"flag"
	"github.com/phoobynet/market-deck-server/assets"
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

func fatal(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	signal.Notify(quitChan, os.Interrupt)
	logrus.SetFormatter(&logrus.TextFormatter{})

	config := loadConfig()

	migrateDatabase()

	ctx, cancel := context.WithCancel(context.Background())

	calendarDayLive := calendars.GetCalendarDayLive(ctx, messageBus)

	snapshotStream := snapshots.NewSnapshotStream(
		ctx,
		calendarDayLive,
		messageBus,
	)

	server.InitSSE()
	webServer := server.NewServer(config, dist, snapshotStream)

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

func migrateDatabase() {
	db := database.GetDB()
	fatal(db.AutoMigrate(&assets.Asset{}))
	fatal(db.AutoMigrate(&calendars.CalendarDay{}))
	fatal(db.AutoMigrate(&decks.Deck{}))
}

func loadConfig() *server.Config {
	var configPath string

	flag.StringVar(&configPath, "config", "config.toml", "Path to config file")

	config, err := server.LoadConfig(configPath)

	fatal(err)

	return config
}
