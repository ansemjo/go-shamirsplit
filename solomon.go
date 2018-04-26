package main

import (
	"github.com/klauspost/reedsolomon"
)

// RSEncode splits and encodes data into $total shards, of which $threshold are needed for reconstruction
func RSEncode(data []byte, threshold, total int) (rscoded [][]byte, err error) {

	rs, err := reedsolomon.New(threshold, total-threshold)
	if err != nil {
		return
	}

	rscoded, err = rs.Split(data)
	if err != nil {
		return
	}

	err = rs.Encode(rscoded)
	return

}
