package collection

import "github.com/samber/lo"

func (c *Collection) UpdateSymbols(symbols []string) {
	removedSymbols, addedSymbols := lo.Difference(c.symbols, symbols)

	if len(removedSymbols) > 0 {
		for _, symbol := range removedSymbols {
			c.collection.Remove(symbol)
		}

		c.symbols = lo.Filter(
			c.symbols, func(symbol string, _ int) bool {
				return !lo.Contains(removedSymbols, symbol)
			},
		)
	}

	if len(addedSymbols) > 0 {
		c.symbols = append(c.symbols, addedSymbols...)
		c.populateBaseSnapshots(addedSymbols) // l
		c.populateVolumes()
		c.populatePreMarketDailyStats()
	}
}
