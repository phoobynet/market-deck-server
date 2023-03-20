package snapshots

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/golang-module/carbon/v2"
	"github.com/phoobynet/market-deck-server/bars"
	"github.com/phoobynet/market-deck-server/helpers/date"
	"github.com/phoobynet/market-deck-server/helpers/numbers"
	"github.com/phoobynet/market-deck-server/quotes"
	"github.com/phoobynet/market-deck-server/trades"
)

type Repository struct {
	mdClient *md.Client
}

func NewSnapshotRepository(mdClient *md.Client) *Repository {
	return &Repository{
		mdClient: mdClient,
	}
}

func (r *Repository) GetMulti(symbols []string) (map[string]Snapshot, error) {
	mdSnapshots, err := r.mdClient.GetSnapshots(symbols, md.GetSnapshotRequest{})

	if err != nil {
		return nil, err
	}

	result := make(map[string]Snapshot)

	now := carbon.Now(date.MarketTimeZone).Format("Y-m-d")

	for symbol, mdSnapshot := range mdSnapshots {
		dailyBar := bars.FromMarketDataBar(symbol, *mdSnapshot.DailyBar)
		previousDailyBar := bars.FromMarketDataBar(symbol, *mdSnapshot.PrevDailyBar)

		previousClose := previousDailyBar.Close

		if dailyBar.Date() < now {
			previousClose = dailyBar.Close
		}

		s := Snapshot{
			LatestBar:        bars.FromMarketDataBar(symbol, *mdSnapshot.MinuteBar),
			LatestQuote:      quotes.FromMarketDataQuote(symbol, *mdSnapshot.LatestQuote),
			LatestTrade:      trades.FromMarketDataTrade(symbol, *mdSnapshot.LatestTrade),
			DailyBar:         dailyBar,
			PreviousDailyBar: previousDailyBar,
			PreviousClose:    previousClose,
		}

		diff := numbers.NumberDiff(s.PreviousClose, s.LatestTrade.Price)

		s.Change = diff.Change
		s.ChangePercent = diff.ChangePercent
		s.ChangeSign = diff.Sign
		s.ChangeAbs = diff.AbsoluteChange

		result[symbol] = s
	}

	return result, nil
}
