package cryptography

// I know this is not strictly "cryptography" but it fits
// best in this package.

import (
	"bytes"

	"github.com/alexbakker/pkcs7"
	"github.com/klauspost/reedsolomon"
)

// TODO: what if shares == threshold? encoder complains about 0 parity shards

// ReedSolomonEncode splits and encodes data into $total shards,
// out of which $threshold are needed for reconstruction.
func ReedSolomonEncode(data []byte, threshold, total int) (split [][]byte, err error) {

	// reversibly pad to a multiple of threshold,
	// i.e. the number of data blocks
	padded, err := pkcs7.Pad(data, threshold)
	if err != nil {
		return
	}

	// instantiate encoder
	// data shares = threshold; parity shares = total - threshold
	rs, err := reedsolomon.New(threshold, total-threshold)
	if err != nil {
		return
	}

	// split data into equal-sized blocks
	split, err = rs.Split(padded)
	if err != nil {
		return
	}

	// calculate parity blocks
	err = rs.Encode(split)
	return

}

// ReedSolomonReconstruct does
func ReedSolomonReconstruct(split [][]byte, threshold, total int) (data []byte, err error) {

	// instantiate decoder
	// data shares = threshold; parity shares = total - threshold
	rs, err := reedsolomon.New(threshold, total-threshold)
	if err != nil {
		return
	}

	// reconstruct data blocks
	err = rs.ReconstructData(split)
	if err != nil {
		return
	}

	// get maximum block length
	max := 0
	for _, block := range split {
		l := len(block)
		if l > max {
			max = l
		}
	}

	// join data blocks. we know that the size must be exactly that
	// of all data blocks combined, since we padded it earlier
	buf := bytes.NewBuffer(make([]byte, 0, max*threshold))
	err = rs.Join(buf, split, max*threshold)
	if err != nil {
		return
	}

	// remove padding
	data, err = pkcs7.Unpad(buf.Bytes(), threshold)
	return

}
