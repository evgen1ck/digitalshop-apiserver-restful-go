package argon2

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
	"log"
)

// generateSalt generates a salt for password hashing
func generateSalt(length int) string {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal(err)
	}

	return string(salt)
}

// HashPassword hashes the input password using the Argon2id algorithm.
func HashPassword(password string, salt string) (string, string) {
	// Define parameters for Argon2id
	const saltLength = 16
	const keyLength = 32
	const iterations = 1
	const memory = 64 * 1024
	const parallelism = 2

	if salt == "" {
		salt = generateSalt(saltLength)
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey([]byte(password), []byte(salt), iterations, memory, parallelism, keyLength)

	return string(hash), salt
}
