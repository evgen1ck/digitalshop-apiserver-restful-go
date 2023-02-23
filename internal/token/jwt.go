package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims represent JWT claims
type Claims struct {
	AccountUuid string `json:"account_uuid"`
	jwt.StandardClaims
}

// GenerateToken creates a JWT token with provided claims and secret key
func GenerateToken(claims Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken parses a JWT token with provided secret key and returns claims
func ParseToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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

// SetClaims sets claims for JWT token
func SetClaims(uuid string, issuer string) Claims {
	return Claims{
		AccountUuid: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 21).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}
}
