package cryptography

import (
	"golang.org/x/crypto/chacha20poly1305"
)

// Encrypt uses authenticated encryption with ChaCha20Poly1305 to encrypt a message.
func Encrypt(key, nonce, associated, message []byte) (ciphertext []byte, err error) {

	chacha, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}

	ciphertext = chacha.Seal(nil, nonce, message, associated)
	return

}

// Decrypt authenticates and decrypts a message that was encrypted with ChaCha20Poly1305.
func Decrypt(key, nonce, associated, ciphertext []byte) (message []byte, err error) {

	chacha, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}

	message, err = chacha.Open(nil, nonce, ciphertext, associated)
	if err != nil {
		message = []byte{}
	}
	return

}
