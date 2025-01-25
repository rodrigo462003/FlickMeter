package hashing

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func HashPassword(p []byte) ([]byte, error) {
	const (
		time    = 1
		memory  = 64 * 1024
		threads = 4
		keyLen  = 32
	)

	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey(p, salt, time, memory, threads, keyLen)

	return hash, nil
}

func generateSalt() ([]byte, error) {
	const saltSize = 16

	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate salt: %v", err)
	}

	return salt, nil
}
