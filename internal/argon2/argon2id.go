package argon2

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"test-server-go/internal/logger"
)

// Define parameters for Argon2id
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-argon2-13#section-7.4
const (
	argonVersion = 0x13
	saltLength   = 16
	keyLength    = 32
	iterations   = 3
	memory       = 64 * 1024
	parallelism  = 2
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
	if argon2.Version != argonVersion {
		logger.NewError("error in hashes password", errors.New("version argon is not 19"))
	}

	if salt == "" {
		salt = generateSalt(saltLength, logger)
	}

	// Hash the password using Argon2id
	passwordHash := argon2.IDKey([]byte(password), []byte(salt), iterations, memory, parallelism, keyLength)

	// Base64 encode the salt and hashed password
	base64PasswordHash := base64.RawStdEncoding.EncodeToString(passwordHash)
	base64Salt := base64.RawStdEncoding.EncodeToString([]byte(salt))

	return base64PasswordHash, base64Salt
}

func CompareHashPasswords(password string, base64PasswordHash string, base64Salt string, logger *logger.Logger) bool {
	if argon2.Version != argonVersion {
		logger.NewError("error in hashes password", errors.New("version argon is not 19"))
	}

	localBase64Salt, err := base64.RawStdEncoding.DecodeString(base64Salt)
	if err != nil {
		logger.NewError("error in decode base64", err)
	}

	localPasswordHash := argon2.IDKey([]byte(password), localBase64Salt, iterations, memory, parallelism, keyLength)

	localBase64PasswordHash := base64.RawStdEncoding.EncodeToString(localPasswordHash)
	if localBase64PasswordHash == base64PasswordHash {
		return true
	}
	return false
}
