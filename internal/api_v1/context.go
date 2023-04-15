package api_v1

import (
	"context"
	"errors"
	"net/http"
	"test-server-go/internal/auth"
)

const (
	AuthenticatedJwtTokenContextKey = "jwt_token"
	AuthenticatedJwtDataContextKey  = "jwt_data"
)

func ContextSetAuthenticated(r *http.Request, token string, data *auth.JwtData) error {
	context.WithValue(r.Context(), AuthenticatedJwtTokenContextKey, token)
	context.WithValue(r.Context(), AuthenticatedJwtDataContextKey, data)
	return nil
}

func ContextGetAuthenticated(r *http.Request) (string, *auth.JwtData, error) {
	token, ok := r.Context().Value(AuthenticatedJwtTokenContextKey).(string)
	if !ok {
		return "", nil, errors.New("user context jwt token key not found")
	}
	if token == "" {
		return "", nil, errors.New("user context jwt token key is null")
	}

	data, ok := r.Context().Value(AuthenticatedJwtDataContextKey).(*auth.JwtData)
	if !ok {
		return "", nil, errors.New("user context jwt data key not found")
	}
	if data == nil {
		return "", nil, errors.New("user context jwt data key is null")
	}

	return token, data, nil
}
