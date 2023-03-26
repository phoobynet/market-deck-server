package http

import (
	"github.com/sirupsen/logrus"
	"time"
)

// GetWithCache - Attempts to retrieve the data from the cache repository; otherwise, a direct call is made to the SEC.
// The SEC response will be stored with a TTL
func GetWithCache(url string, ttl time.Duration) ([]byte, error) {
	result := cacheRepository.Get(url)

	if result != nil {
		logrus.Printf("Cache hit for %s", url)
		return result.Data, nil
	}

	data, err := Get(url)

	if err != nil {
		return nil, err
	}

	cacheRepository.Set(url, data, ttl)

	return data, nil
}
