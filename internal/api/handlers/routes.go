package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"test-server-go/internal/api"
	"test-server-go/internal/models"
	"time"
)

const (
	compressLevel      = 5
	rateLimitRequests  = 15
	rateLimitInterval  = 10 * time.Second
	requestMaxSize     = 12 * 1024 * 1024 // 12 MB
	maxHeaderSize      = 1 * 1024 * 1024  // 1 MB
	uriMaxLength       = 1024
	serviceUnavailable = false
)

var (
	allowedMethods      = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	allowedContentTypes = []string{"application/json"}
)

type RouteHandler struct {
	App *models.Application
}

func (rh *RouteHandler) SetupRouter() http.Handler {
	r := chi.NewRouter()

	// CORS settings
	r.Use(api.CorsMiddleware())

	// Default settings
	r.Use(middleware.Recoverer)                                                  // Prevents server from crashing
	r.Use(middleware.StripSlashes)                                               // Optimizes paths
	r.Use(middleware.Logger)                                                     // Logging
	r.Use(api.ContentTypeCompressMiddleware(compressLevel, allowedContentTypes)) // Supports compression

	// Error handling
	r.Use(api.ServiceUnavailableMiddleware(serviceUnavailable))          // Error 503 - Service Unavailable
	r.Use(api.RateLimitMiddleware(rateLimitRequests, rateLimitInterval)) // Error 429 - Too Many Requests
	r.Use(api.UriLengthMiddleware(uriMaxLength))                         // Error 414 - URI Too Long
	r.Use(api.RequestSizeMiddleware(requestMaxSize))                     // Error 413 - Payload Too Large
	r.Use(api.UnprocessableEntityMiddleware)                             // Error 422 - Unprocessable Entity
	r.Use(api.NotAcceptableMiddleware(allowedContentTypes))              // Error 406 - Not Acceptable
	r.Use(api.UnsupportedMediaTypeMiddleware(allowedContentTypes))       // Error 415 - Unsupported Media Type
	r.Use(api.NotImplementedMiddleware(allowedMethods))                  // Error 501 - Not implemented
	r.Use(api.MethodNotAllowedMiddleware)                                // Error 405 - Method Not Allowed
	r.NotFound(api.NotFoundMiddleware())                                 // Error 404 - Not Found

	rh.RegisterRoutes(r)

	maxHeaderBytesMW := api.MaxHeaderBytesMiddleware(maxHeaderSize)
	handlerWithMiddleware := maxHeaderBytesMW(r)

	return handlerWithMiddleware
}

func (rh *RouteHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/signup", rh.AuthSignup)
		r.Post("/signupWithToken", rh.SignupWithToken)
	})
}
