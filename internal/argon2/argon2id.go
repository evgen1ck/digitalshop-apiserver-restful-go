package argon2

import (
	"crypto/rand"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
)

// Define parameters for Argon2id
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-argon2-13#section-7.4
const (
	saltLength  = 16
	keyLength   = 32
	iterations  = 3
	memory      = 64 * 1024
	parallelism = 2
)

// generateSalt generates a salt for password hashing
func generateSalt(length int) string {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		logrus.Fatal("salt generation error", err)
	}

	return string(salt)
}

// HashPassword hashes the input password using the Argon2id algorithm.
func HashPassword(password string, salt string) (string, string) {
	if salt == "" {
		salt = generateSalt(saltLength)
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey([]byte(password), []byte(salt), iterations, memory, parallelism, keyLength)

	return string(hash), salt
}
