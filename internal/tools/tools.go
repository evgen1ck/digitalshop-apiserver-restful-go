package tools

import (
	"crypto/rand"
	"math/big"
)

func GenerateConfirmCode() (int64, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return 0, err
	}
	return num.Int64() + 100000, nil
}
