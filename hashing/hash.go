package hashing

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func MustHashPassword(p []byte) string {
	const (
		time    = 1
		memory  = 64 * 1024
		threads = 4
		keyLen  = 32
	)

	salt := mustGenerateSalt()
	hash := argon2.IDKey(p, salt, time, memory, threads, keyLen)
	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedHash := base64.StdEncoding.EncodeToString(hash)

	result := fmt.Sprintf("argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, threads, encodedSalt, encodedHash)

	return result
}

func mustGenerateSalt() []byte {
	const saltSize = 16

	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		panic(fmt.Sprintf("rand.Read failed to generate salt: %v", err))
	}

	return salt
}
