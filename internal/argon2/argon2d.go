package argon2

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
	"log"
)

// generateSalt generates a salt for password hashing
func generateSalt(length int) []byte {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal(err)
	}

	return salt
}

// HashPassword hashes the input password using the Argon2d algorithm.
func HashPassword(password string) (string, string) {
	// Define parameters for Argon2d
	const saltLength = 16
	const keyLength = 32
	const iterations = 3
	const memory = 64 * 1024
	const parallelism = 2

	salt := generateSalt(saltLength)

	// Hash the password using Argon2d
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	return string(hash), string(salt)
}
