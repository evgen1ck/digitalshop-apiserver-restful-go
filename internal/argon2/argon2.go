package argon2

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

func generateSalt(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreatePassword(password string) ([]byte, string) {

	// Define parameters for Argon2i
	var salt, _ = generateSalt(16)
	const memory = 64 * 1024
	const threads = 4
	const keyLen = 32

	// Hash the password using Argon2i
	hash := argon2.IDKey([]byte(password), []byte(salt), 1, memory, threads, keyLen)
	return hash, salt
}
