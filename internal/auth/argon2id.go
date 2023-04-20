package auth

import (
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"test-server-go/internal/tools"
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

// HashPassword hashes the input password using the Argon2id algorithm.
func HashPassword(password string, salt string) (string, string, error) {
	if salt == "" {
		var err error
		salt, err = tools.GenerateRandomString(saltLength)
		if err != nil {
			return "", "", err
		}
	}

	// Hash the password using Argon2id
	passwordHash := argon2.IDKey([]byte(password), []byte(salt), iterations, memory, parallelism, keyLength)

	// Base64 encode the salt and hashed password
	base64PasswordHash := base64.RawStdEncoding.EncodeToString(passwordHash)
	base64Salt := base64.RawStdEncoding.EncodeToString([]byte(salt))

	return base64PasswordHash, base64Salt, nil
}

// CompareHashPasswords compares the input password with the hashed password and salt.
func CompareHashPasswords(password string, base64PasswordHash string, base64Salt string) (bool, error) {
	newBase64Salt, err := base64.RawStdEncoding.DecodeString(base64Salt)
	if err != nil {
		return false, err
	}

	newPasswordHash := argon2.IDKey([]byte(password), newBase64Salt, iterations, memory, parallelism, keyLength)
	newBase64PasswordHash := base64.RawStdEncoding.EncodeToString(newPasswordHash)

	return newBase64PasswordHash == base64PasswordHash, nil
}
