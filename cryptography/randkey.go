package cryptography

import "crypto/rand"

// RandomKey reads randomness to generate a 32 byte key and
// panics if something went wrong.
func RandomKey() (key []byte) {

	k := make([]byte, 32)
	_, err := rand.Read(k)

	if err != nil {
		panic("could not read randomness")
	}

	return k

}
