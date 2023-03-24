package collection

import "github.com/phoobynet/market-deck-server/snapshots"

func (c *Collection) Items() map[string]*snapshots.Snapshot {
	return c.collection.Items()
}
