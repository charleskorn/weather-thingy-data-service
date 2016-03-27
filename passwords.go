package main

import (
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/unicode/norm"
	"io"
)

const hashIterations = 100000
const saltBytes int = 32

func computePasswordHash(password string, salt []byte, iterations int) []byte {
	passwordBytes := norm.NFC.Bytes([]byte(password))

	return pbkdf2.Key(passwordBytes, salt, iterations, sha256.Size, sha256.New)
}

func generateHashingSalt() ([]byte, error) {
	salt := make([]byte, saltBytes)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	return salt, nil
}
