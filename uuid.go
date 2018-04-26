package main

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/blake2b"
)

func RandomUUID() uuid.UUID {

	return uuid.New()

}

func HashToUUID(bytes []byte) uuid.UUID {

	hash := blake2b.Sum256(bytes)
	return uuid.Must(uuid.FromBytes(hash[:16]))

}
