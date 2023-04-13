package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// TokenExpirationTime specifies the token expiration time.
const TokenExpirationTime = time.Hour * 24 * 21

// JwtData represents the custom JWT claims, which includes the account UUID and standard claims.
type JwtData struct {
	AccountUuid string `json:"account_uuid"`
	jwt.RegisteredClaims
}

// GenerateJwt generates a JWT token for a given account UUID and secret key.
// It sets the token to expire in 21 days and includes the issued at time.
func GenerateJwt(accountUuid string, secret string) (string, error) {
	claims := JwtData{
		AccountUuid: accountUuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	return token.SignedString([]byte(secret))
}

// ParseJwtToken parses a JWT token string and returns the custom claims or an error.
// It verifies the token signature using the secret key and checks for token expiration.
func ParseJwtToken(tokenString string, secret string) (*JwtData, error) {
	if tokenString == "" {
		return nil, errors.New("missing token")
	}

	claims := &JwtData{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS384.Alg() {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
