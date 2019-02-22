package cryptography

import (
	"github.com/dsprenkels/sss-go"
)

// see https://dsprenkels.com/sss-34c3.html for description and implementation details

// ShamirSplit splits 32 bytes into keyshares using Shamir secret sharing.
// When later recombining the key, at least $threshold keyshares are required
// to construct the correct key.
func ShamirSplit(key []byte, threshold, shares int) (keyshares [][]byte, err error) {
	return sss.CreateKeyshares(key, shares, threshold)
}

// ShamirCombine reconstructs a 32 byte key from previously split keyshares.
// There is a minimum threshold required to recover the correct secret, however
// no error is emitted if you do not supply enough shares. You will simply
// receive a corrupted / wrong key.
func ShamirCombine(keyshares [][]byte) (key []byte, err error) {
	return sss.CombineKeyshares(keyshares)
}
