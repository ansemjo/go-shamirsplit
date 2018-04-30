package cryptography

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ed25519"
)

// keyFromBytes generates Ed25519 keys from 32 bytes of entropy (a key)
func keysFromBytes(key []byte) (pubkey ed25519.PublicKey, seckey ed25519.PrivateKey, err error) {

	if len(key) != ed25519PrivateSeedSize {
		err = ErrKeySize
		return
	}

	return ed25519.GenerateKey(bytes.NewBuffer(key))

}

// EdSign signs a message with a Ed25519 private key generated from 32 bytes of entropy.
func EdSign(key, message []byte) (signature []byte, pubkey ed25519.PublicKey, err error) {

	pubkey, seckey, err := keysFromBytes(key)
	if err != nil {
		return
	}

	signature = ed25519.Sign(seckey, message)
	return

}

// EdVerify verifies a previously signed message with an Ed25519 signature.
func EdVerify(pubkey, message, signature []byte) (ok bool) {

	p := ed25519.PublicKey(pubkey)
	ok = ed25519.Verify(p, message, signature)
	return

}

// the amount of entropy needed to generate keys, i.e. the length of needed byte array
const ed25519PrivateSeedSize = 32

// ErrKeySize indicates that the byte array used as entropy source is not 32 bytes long
var ErrKeySize = fmt.Errorf("key is not %d bytes long", ed25519PrivateSeedSize)
