package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// TokenExpirationTime specifies the token expiration time.
const TokenExpirationTime = time.Hour * 24 * 21

// JwtClaims represents the custom JWT claims, which includes the account UUID and standard claims.
type JwtClaims struct {
	AccountUuid string `json:"account_uuid"`
	jwt.StandardClaims
}

// GenerateJwt generates a JWT token for a given account UUID and secret key.
// It sets the token to expire in 21 days and includes the issued at time.
func GenerateJwt(accountUuid string, secret string) (string, error) {
	claims := JwtClaims{
		AccountUuid: accountUuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpirationTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	return token.SignedString([]byte(secret))
}

// ParseJwtToken parses a JWT token string and returns the custom claims or an error.
// It verifies the token signature using the secret key and checks for token expiration.
func ParseJwtToken(tokenString string, secret string) (*JwtClaims, error) {
	claims := &JwtClaims{}
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

// JwtAuthenticate authenticates a request with a JWT token and extracts the account UUID.
// If the token is missing or invalid, it returns an HTTP error.
func JwtAuthenticate(token string, secret string) (*JwtClaims, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}

	claims, err := ParseJwtToken(token, secret)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}