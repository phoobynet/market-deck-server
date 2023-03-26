package http

import "github.com/phoobynet/market-deck-server/cache"

var cacheRepository *cache.Repository

func init() {
	cacheRepository = cache.GetRepository()
}
