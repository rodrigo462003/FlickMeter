package hashing

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func HashPassword(p []byte) (string, error) {
	const (
		time    = 1
		memory  = 64 * 1024
		threads = 4
		keyLen  = 32
	)

	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(p, salt, time, memory, threads, keyLen)
	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedHash := base64.StdEncoding.EncodeToString(hash)

	result := fmt.Sprintf("argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, threads, encodedSalt, encodedHash)

	return result, nil
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
