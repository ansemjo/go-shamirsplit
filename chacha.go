package main

import (
	"crypto/rand"

	"golang.org/x/crypto/chacha20poly1305"
)

// Encrypt takes a 32 byte key, associated data and a message, generates a random
// nonce and encrypts using ChaCha20Poly1305
func Encrypt(key, ad, message []byte) (nonce, ciphertext []byte, err error) {

	chacha, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}

	nonce = make([]byte, chacha20poly1305.NonceSize)
	_, err = rand.Read(nonce)
	if err != nil {
		return
	}

	ciphertext = chacha.Seal(nil, nonce, message, ad)
	return

}
