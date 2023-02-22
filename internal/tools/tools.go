package tools

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateConfirmCode() (string, error) {
	// Generate 3 random bytes
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert bytes to a number
	n := new(big.Int).SetBytes(b)

	// Limit the number to 6 digits
	n.Mod(n, big.NewInt(1000000))

	// Again limit the number to 6 digits
	return fmt.Sprintf("%06d", n), nil
}
