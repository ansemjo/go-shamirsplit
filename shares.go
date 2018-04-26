package main

import (
	"github.com/dsprenkels/sss-go"
)

// CreateKeyShares splits 32 bytes into keyshares using dsprenkels/sss
func CreateKeyShares(key []byte, threshold, total int) (shares [][]byte, err error) {
	shares, err = sss.CreateKeyshares(key, total, threshold)
	return
}
