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
	compressLevel      = 5
	rateLimitRequests  = 10
	rateLimitInterval  = 10 * time.Second
	requestMaxSize     = 4 * 1024 * 1024 // 4MB
	maxHeaderSize      = 1024            // 1MB
	uriMaxLength       = 1024
	supportedMediaType = "application/json"
)

var allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

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
	r.Use(unsupportedMediaTypeMiddleware(supportedMediaType))        // Error 415 - Unsupported Media Type
	r.Use(notImplementedMiddleware(allowedMethods))                  // Error 501 - (Method) Not implemented
	r.Use(methodNotAllowedHandler)                                   // Error 405 - Method Not Allowed
	r.NotFound(notFoundMiddleware())                                 // Error 404 - Not Found

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/signup", rh.AuthSignup)
		r.Post("/signupWithToken", rh.SignupWithToken)
	})

	//r.Get("/api/v1/getAllAlbums", handlers.getAllAlbums)

	maxHeaderBytesMW := maxHeaderBytesMiddleware(maxHeaderSize)
	handlerWithMiddleware := maxHeaderBytesMW(r)

	return handlerWithMiddleware
}
