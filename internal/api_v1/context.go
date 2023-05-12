package api_v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"test-server-go/internal/auth"
)

const (
	AuthenticatedJwtTokenContextKey = "jwt_token"
	AuthenticatedJwtDataContextKey  = "jwt_data"
)

func ContextSetAuthenticated(r *http.Request, token string, data *auth.JwtData) *http.Request {
	ctx := context.WithValue(r.Context(), AuthenticatedJwtTokenContextKey, token)
	ctx = context.WithValue(ctx, AuthenticatedJwtDataContextKey, data)

	fmt.Println(data)
	return r.WithContext(ctx)
}

func ContextGetAuthenticated(r *http.Request) (string, *auth.JwtData, error) {
	tokenValue := r.Context().Value(AuthenticatedJwtTokenContextKey)
	if tokenValue == nil {
		return "", nil, errors.New("user context jwt token key not found")
	}
	token, ok := tokenValue.(string)
	if !ok {
		return "", nil, errors.New("user context jwt token key is not a string")
	}

	dataValue := r.Context().Value(AuthenticatedJwtDataContextKey)
	if dataValue == nil {
		return "", nil, errors.New("user context jwt data key not found")
	}
	data, ok := dataValue.(*auth.JwtData)
	if !ok {
		return "", nil, errors.New("user context jwt data key is not of expected type")
	}

	return token, data, nil
}
