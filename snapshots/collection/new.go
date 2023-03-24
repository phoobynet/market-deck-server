package collection

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/calendars"
	"github.com/phoobynet/market-deck-server/clients"
	"github.com/phoobynet/market-deck-server/snapshots"
)

type collectionMap = cmap.ConcurrentMap[string, *snapshots.Snapshot]

type Collection struct {
	collection      collectionMap
	mdClient        *md.Client
	assetRepo       *assets.Repository
	calendarDayLive *calendars.CalendarDayLive
	symbols         []string
}

func New(calendarDayLive *calendars.CalendarDayLive) *Collection {
	s := &Collection{
		mdClient:        clients.GetMarketDataClient(),
		assetRepo:       assets.GetRepository(),
		calendarDayLive: calendarDayLive,
	}

	return s
}
