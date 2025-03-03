package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	time    = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
)

func PasswordsMatch(password, dbHash string) bool {
	hashSplit := strings.Split(dbHash, ",")

	encodedSalt, encodedHash := hashSplit[4], hashSplit[5]
	salt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		panic(err)
	}

	hash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		panic(err)
	}
	newHash := argon2.IDKey([]byte(password), []byte(salt), time, memory, threads, keyLen)

	return subtle.ConstantTimeCompare(newHash, hash) == 1
}

func HashPassword(p []byte) string {
	salt := generateSalt()
	hash := argon2.IDKey(p, salt, time, memory, threads, keyLen)

	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedHash := base64.StdEncoding.EncodeToString(hash)
	return fmt.Sprintf(
		"argon2id$%d,%d,%d,%d,%s,%s",
		argon2.Version, memory, time, threads,
		encodedSalt, encodedHash,
	)
}

func generateSalt() []byte {
	const saltSize = 16

	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}

	return salt
}
