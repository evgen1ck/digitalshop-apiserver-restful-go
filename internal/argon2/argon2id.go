package argon2

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"log"
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
func generateSalt(length int) string {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatalf("error in salt generation %e", err)
	}

	return string(salt)
}

// HashPassword hashes the input password using the Argon2id algorithm.
func HashPassword(password string, salt string) (string, string) {
	if salt == "" {
		salt = generateSalt(saltLength)
	}

	// Hash the password using Argon2id
	passwordHash := argon2.IDKey([]byte(password), []byte(salt), iterations, memory, parallelism, keyLength)

	// Base64 encode the salt and hashed password
	base64PasswordHash := base64.RawStdEncoding.EncodeToString(passwordHash)
	base64Salt := base64.RawStdEncoding.EncodeToString([]byte(salt))

	return base64PasswordHash, base64Salt
}

func CompareHashPasswords(password string, base64PasswordHash string, base64Salt string) bool {
	localBase64Salt, err := base64.RawStdEncoding.DecodeString(base64Salt)
	if err != nil {
		log.Fatalf("error in decode base64 %e", err)
	}

	localPasswordHash := argon2.IDKey([]byte(password), localBase64Salt, iterations, memory, parallelism, keyLength)

	localBase64PasswordHash := base64.RawStdEncoding.EncodeToString(localPasswordHash)
	if localBase64PasswordHash == base64PasswordHash {
		return true
	}
	return false
}
