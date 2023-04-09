package api_v1

import (
	"context"
	"errors"
	"net/http"
	"test-server-go/internal/auth"
)

const (
	AuthenticatedUserContextKey = "jwt_data"
)

func ContextSetAuthenticatedUser(r *http.Request, data *auth.JwtClaims) error {
	if _, ok := r.Context().Value(AuthenticatedUserContextKey).(*auth.JwtClaims); ok {
		return errors.New("authenticated user context key already exists")
	}

	context.WithValue(r.Context(), AuthenticatedUserContextKey, data)

	return nil
}

func ContextGetAuthenticatedUser(r *http.Request) (*auth.JwtClaims, error) {
	value, ok := r.Context().Value(AuthenticatedUserContextKey).(*auth.JwtClaims)
	if !ok {
		return nil, errors.New("user context key not found")
	}
	if value == nil {
		return nil, errors.New("user context key is null")
	}

	return value, nil
}
