package argon2

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
	"test-server-go/internal/logger"
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
func generateSalt(length int, logger *logger.Logger) string {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		logger.NewError("error in salt generation", err)
	}

	return string(salt)
}

// HashPassword hashes the input password using the Argon2id algorithm.
func HashPassword(password string, salt string, logger *logger.Logger) (string, string) {
	if salt == "" {
		salt = generateSalt(saltLength, logger)
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey([]byte(password), []byte(salt), iterations, memory, parallelism, keyLength)

	return string(hash), salt
}
