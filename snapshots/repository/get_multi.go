package repository

import (
	md "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

func (r *Repository) GetMulti(symbols []string) (map[string]*md.Snapshot, error) {
	return r.mdClient.GetSnapshots(symbols, md.GetSnapshotRequest{})
}
