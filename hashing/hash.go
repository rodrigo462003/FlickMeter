package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base32"
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
	salt, hash := hashSplit[4], hashSplit[5]

	newHash := argon2.IDKey([]byte(password), []byte(salt), time, memory, threads, keyLen)
	encodedHash := base32.StdEncoding.EncodeToString(newHash)

	return subtle.ConstantTimeCompare([]byte(encodedHash), []byte(hash)) == 1
}

func HashPassword(p []byte) string {
	salt := rand.Text()
	hash := argon2.IDKey(p, []byte(salt), time, memory, threads, keyLen)

	encodedHash := base32.StdEncoding.EncodeToString(hash)
	return fmt.Sprintf(
		"argon2id$%d,%d,%d,%d,%s,%s",
		argon2.Version, memory, time, threads,
		salt, encodedHash,
	)
}
