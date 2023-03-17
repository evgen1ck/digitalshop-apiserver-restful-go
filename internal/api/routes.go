package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"test-server-go/internal/models"
	"time"
)

type RouteHandler struct {
	App *models.Application
}

const (
	compressLevel     = 5
	rateLimitRequests = 10
	rateLimitInterval = 10 * time.Second
	requestMaxSize    = 12 * 1024 * 1024 // 12 MB
	maxHeaderSize     = 1 * 1024 * 1024  // 1 MB
	uriMaxLength      = 1024
	usingUtf8         = true
)

var (
	allowedMethods      = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	allowedContentTypes = []string{"application/json"}
)

func (rh *RouteHandler) SetupRouter() http.Handler {
	r := chi.NewRouter()

	// CORS settings
	r.Use(corsMiddleware())

	// Default settings
	r.Use(middleware.Recoverer)                                   // Prevents server from crashing
	r.Use(middleware.StripSlashes)                                // Optimizes paths
	r.Use(middleware.Logger)                                      // Logging
	r.Use(middleware.Compress(compressLevel, "application/json")) // Supports compression

	// Error handling
	r.Use(serviceUnavailableMiddleware(false))                       // Error 503 - Service Unavailable
	r.Use(uriLengthMiddleware(uriMaxLength))                         // Error 414 - URI Too Long
	r.Use(requestSizeMiddleware(requestMaxSize))                     // Error 413 - Payload Too Large
	r.Use(rateLimitMiddleware(rateLimitRequests, rateLimitInterval)) // Error 429 - Too Many Requests
	r.Use(unsupportedMediaTypeMiddleware(allowedContentTypes))       // Error 415 - Unsupported Media Type
	r.Use(notImplementedMiddleware(allowedMethods))                  // Error 501 - Not implemented
	r.Use(notAcceptableMiddleware(allowedContentTypes, usingUtf8))   // Error 406 - Not acceptable
	r.Use(methodNotAllowedMiddleware)                                // Error 405 - Method Not Allowed
	r.NotFound(notFoundMiddleware())                                 // Error 404 - Not Found

	rh.RegisterRoutes(r)

	maxHeaderBytesMW := maxHeaderBytesMiddleware(maxHeaderSize)
	handlerWithMiddleware := maxHeaderBytesMW(r)

	return handlerWithMiddleware
}

func (rh *RouteHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/signup", rh.AuthSignup)
		r.Post("/signupWithToken", rh.SignupWithToken)
	})
}
