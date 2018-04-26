package main

import (
	"golang.org/x/crypto/chacha20poly1305"
)

// Encrypt uses authenticated encryption with ChaCha20Poly1305
// to encrypt a message
func Encrypt(key, nonce, associated, message []byte) (ciphertext []byte, err error) {

	chacha, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}

	ciphertext = chacha.Seal(nil, nonce, message, associated)
	return

}
