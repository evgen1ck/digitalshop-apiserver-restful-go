package auth

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strconv"
	"time"
)

const (
	// SecretLength is the length of the secret key in bytes
	SecretLength = 10
	// Base32Length is the length of the secret key in base32 encoding
	Base32Length = 16
)

// User represents a user in the system
type User struct {
	ID       int
	Username string
	Secret   string
	Enabled  bool
}

// GenerateSecret generates a new secret for 2FA
func GenerateSecret() (string, error) {
	secret := make([]byte, SecretLength)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("could not generate secret: %w", err)
	}

	encoded := base32.StdEncoding.EncodeToString(secret)
	return encoded[:Base32Length], nil
}

// ValidateCode checks if the provided code is valid for a given user
func (u *User) ValidateCode(code string) bool {
	hash := sha1.New()
	counter := uint64(time.Now().Unix() / 30)

	for i := 0; i < 8; i++ {
		b := make([]byte, 8)
		for j := 7; j >= 0; j-- {
			b[j] = byte(counter & 0xff)
			counter >>= 8
		}

		hash.Write(b)
		hash.Write([]byte(u.Secret))
		sum := hash.Sum(nil)

		offset := sum[len(sum)-1] & 0xf
		truncated := uint32(sum[offset]&0x7f)<<24 |
			uint32(sum[offset+1]&0xff)<<16 |
			uint32(sum[offset+2]&0xff)<<8 |
			uint32(sum[offset+3]&0xff)

		if strconv.Itoa(int(truncated%1000000)) == code {
			return true
		}
	}

	return false
}
