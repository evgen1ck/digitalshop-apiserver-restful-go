package argon2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	base64Salt := "cOCWOnHGAaBu3p2SoEGreA"

	firstPassword := "password"
	firstBase64PasswordHash, _ := HashPassword(firstPassword, base64Salt)

	secondPassword := "password"
	secondBase64PasswordHash, _ := HashPassword(secondPassword, base64Salt)

	assert.Equalf(t, firstBase64PasswordHash, secondBase64PasswordHash, "result should be equal")
}

func TestHashPassword2(t *testing.T) {
	base64Salt := "cOCWOnHGAaBu3p2SoEGreA"

	firstPassword := "password"
	firstBase64PasswordHash, _ := HashPassword(firstPassword, base64Salt)

	secondPassword := "passw0rd"
	secondBase64PasswordHash, _ := HashPassword(secondPassword, base64Salt)

	assert.NotEqualf(t, firstBase64PasswordHash, secondBase64PasswordHash, "result should not be equal")
}

func TestHashPassword3(t *testing.T) {
	firstPassword := "password"
	firstBase64PasswordHash, _ := HashPassword(firstPassword, "")

	secondPassword := "password"
	secondBase64PasswordHash, _ := HashPassword(secondPassword, "")

	assert.NotEqualf(t, firstBase64PasswordHash, secondBase64PasswordHash, "result should not be equal")
}

func TestCompareHashPasswords(t *testing.T) {
	password := "password123"
	base64PasswordHash := "dQ6IluEKycDJTJ/4q5MddItNX78Tpvrj84Dex6Kgu18"
	base64Salt := "4F6M7rm8AdF8puDrwoGsPg"

	assert.Truef(t, CompareHashPasswords(password, base64PasswordHash, base64Salt), "result should be true")
}

func TestCompareHashPasswords2(t *testing.T) {
	password := "password1234"
	base64PasswordHash := "dQ6IlugycDJTJ/4q5MddItNX78Tpvrj84Dex6Kgu18"
	base64Salt := "4F6M7rm8gdF8purwoGsPg"

	assert.Falsef(t, CompareHashPasswords(password, base64PasswordHash, base64Salt), "result should be error")
}
