package tickers

import "time"

// Populate - downloads and stores company ticker information from the size in the local database
// If the data is older than 7 days, it will be refreshed.
func Populate() {
	r := GetRepository()

	lastUpdated := r.LastUpdated()

	if lastUpdated == nil || time.Since(*lastUpdated) > ttl {
		r.DeleteAll()

		t := get()

		r.BulkInsert(t)
	}
}
