package api_v1

import (
	"context"
	"errors"
	"net/http"
	"test-server-go/internal/auth"
)

const (
	AuthenticatedContextKey = "jwt_data"
)

func ContextSetAuthenticated(r *http.Request, data *auth.JwtClaims) error {
	context.WithValue(r.Context(), AuthenticatedContextKey, data)
	return nil
}

func ContextGetAuthenticated(r *http.Request) (*auth.JwtClaims, error) {
	value, ok := r.Context().Value(AuthenticatedContextKey).(*auth.JwtClaims)
	if !ok {
		return nil, errors.New("user context key not found")
	}
	if value == nil {
		return nil, errors.New("user context key is null")
	}

	return value, nil
}
