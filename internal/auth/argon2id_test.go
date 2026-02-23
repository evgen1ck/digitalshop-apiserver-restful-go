package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	password := "dawdawd$$@#wada$%DD33"
	base64PasswordHash := "aYbu27f3ubbO61MCg2Mu2tUswepdX/0HVZ3olapSiBc"
	base64Salt := "jRho+4SIsja+bEf72Z6AfQ"

	result, err := CompareHashPasswords(password, base64PasswordHash, base64Salt)
	if err != nil {
		t.Errorf("CompareHashPasswords failed: %v", err)
	}

	assert.Truef(t, result, "result should be true")
}
