package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	base64Salt := "cOCWOnHGAaBu3p2SoEGreA"

	firstPassword := "password"
	firstBase64PasswordHash, _, err := HashPassword(firstPassword, base64Salt)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	secondPassword := "password"
	secondBase64PasswordHash, _, err := HashPassword(secondPassword, base64Salt)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	assert.Equalf(t, firstBase64PasswordHash, secondBase64PasswordHash, "result should be equal")
}

func TestCompareHashPasswords(t *testing.T) {
	password := "password123"
	base64PasswordHash := "dQ6IluEKycDJTJ/4q5MddItNX78Tpvrj84Dex6Kgu18"
	base64Salt := "4F6M7rm8AdF8puDrwoGsPg"

	result, err := CompareHashPasswords(password, base64PasswordHash, base64Salt)
	if err != nil {
		t.Errorf("CompareHashPasswords failed: %v", err)
	}

	assert.Truef(t, result, "result should be true")
}
