package main

import (
	"encoding/base64"
	"encoding/pem"
	"strconv"

	"golang.org/x/crypto/ed25519"
)

type Shard struct {
	index     int
	nonce     []byte
	keyshare  []byte
	pubkey    ed25519.PublicKey
	signature []byte
	data      []byte
}

func (s *Shard) toBlock() (block *pem.Block) {

	base64e := base64.StdEncoding.EncodeToString

	headers := make(map[string]string)
	headers["Index"] = strconv.Itoa(s.index)
	headers["Nonce"] = base64e(s.nonce)
	headers["Keyshare"] = base64e(s.keyshare)
	headers["Public Key"] = base64e(s.pubkey)
	headers["Signature"] = base64e(s.signature)

	return &pem.Block{Type: pemtype, Headers: headers, Bytes: s.data}

}

func EncodePEM(shard Shard) (p []byte) {

	block := shard.toBlock()
	return pem.EncodeToMemory(block)

}

const pemtype = "SHAMIR SHARED MESSAGE"
