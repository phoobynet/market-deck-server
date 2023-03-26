package assets

import "time"

// ttl - the lifetime of assets data before requiring a refresh
const ttl = 7 * 24 * time.Hour
