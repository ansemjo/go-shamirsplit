package main

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/ed25519"
)

// keyFromBytes generates an Ed25519 private key from 32 bytes of entropy (a key)
func keyFromBytes(key []byte) (pubkey ed25519.PublicKey, seckey ed25519.PrivateKey, err error) {

	if len(key) != 32 {
		err = ErrKeySize
		return
	}

	buf := bytes.NewBuffer(key)
	pubkey, seckey, err = ed25519.GenerateKey(buf)
	return

}

// Sign signs a message with a Ed25519 private key generated from 32 bytes of entropy
func Sign(key, message []byte) (signature []byte, pubkey ed25519.PublicKey, err error) {

	pubkey, seckey, err := keyFromBytes(key)
	if err != nil {
		return
	}

	signature = ed25519.Sign(seckey, message)
	return

}

// Verify simply verifies a previously signed message
var Verify = ed25519.Verify

// ErrKeySize indicates that the byte array used as entropy source is not 32 bytes long
var ErrKeySize = errors.New("key is not 32 bytes long")
