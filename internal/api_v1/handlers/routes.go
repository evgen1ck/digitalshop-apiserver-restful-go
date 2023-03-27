package handlers

import (
	"github.com/go-chi/chi/v5"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/models"
	"time"
)

const (
	timeout            = 5 * time.Second
	rateLimitRequests  = 80
	rateLimitInterval  = 1 * time.Minute
	requestMaxSize     = 12 * 1024 * 1024 // 12 MB
	uriMaxLength       = 1024             // 1024 runes
	serviceUnavailable = false
)

var (
	allowedMethods        = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	allowedContentTypes   = []string{"application/json", "text/plain"}
	supportedHttpVersions = []string{"HTTP/1.1", "HTTP/2.0"}
)

type Resolver struct {
	App *models.Application
}

func (rs *Resolver) SetupRouterApiVer1(pathPrefix string) {
	r := chi.NewRouter()

	// CORS settings
	r.Use(api_v1.CorsMiddleware())

	// Error handling
	r.Use(api_v1.ServiceUnavailableMiddleware(serviceUnavailable))          // Error 503 - Service Unavailable
	r.Use(api_v1.RateLimitMiddleware(rateLimitRequests, rateLimitInterval)) // Error 429 - Too Many Requests
	r.Use(api_v1.UriLengthMiddleware(uriMaxLength))                         // Error 414 - URI Too Long
	r.Use(api_v1.RequestSizeMiddleware(requestMaxSize))                     // Error 413 - Payload Too Large
	r.Use(api_v1.UnprocessableEntityMiddleware)                             // Error 422 - Unprocessable Entity
	r.Use(api_v1.UnsupportedMediaTypeMiddleware(allowedContentTypes))       // Error 415 - Unsupported Media Type
	r.Use(api_v1.NotImplementedMiddleware(allowedMethods))                  // Error 501 - Not implemented
	r.Use(api_v1.MethodNotAllowedMiddleware)                                // Error 405 - Method Not Allowed
	r.Use(api_v1.HttpVersionCheckMiddleware(supportedHttpVersions))         // Error 505 - HTTP Version Not Supported
	r.Use(api_v1.GatewayTimeoutMiddleware(timeout))                         // Error 504 - Gateway Timeout
	r.NotFound(api_v1.NotFoundMiddleware())                                 // Error 404 - Not Found

	rs.registerRoutes(r)

	rs.App.Router.Mount(pathPrefix, r)
}

func (rs *Resolver) registerRoutes(r chi.Router) {
	r.Route("/server", func(r chi.Router) {
		r.Route("/databases", func(r chi.Router) {
			r.Route("/postgres", func(r chi.Router) {
				r.Post("/info", rs.ServerDatabasesPostgresInfo)
				r.Post("/backup", rs.ServerDatabasesPostgresBackup)
			})
		})
	})
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", rs.AuthSignup)
		r.Post("/signup-with-token", rs.AuthSignupWithToken)
		r.Post("/login", rs.AuthLogin)
		r.Post("/login-with-token", rs.AuthLoginWithToken)
		r.Post("/recover-password", rs.AuthRecoverPassword)
		r.Post("/recover-password-with-token", rs.AuthRecoverPasswordWithToken)
		r.Post("/logout", rs.AuthLogout)
	})
	r.Route("/products", func(r chi.Router) {
		r.Get("/data", rs.ProductsData)
	})
	r.Route("/user", func(r chi.Router) {
		r.Use(api_v1.AuthUserMiddleware(rs.App.Config.App.Jwt))
		r.Route("/profile", func(r chi.Router) {
			r.Get("/data", rs.UserProfileData)
			r.Post("/dump", rs.UserProfileDump)
			r.Patch("/update", rs.UserProfileUpdate)
			r.Delete("/delete", rs.UserProfileDelete)
			r.Get("/orders", rs.UserProfileOrders)
		})
	})
	r.Route("/seller", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Get("/data", rs.SellerProductsData)
			r.Post("/add", rs.SellerProductsAdd)
			r.Patch("/update/{id}", rs.SellerProductsUpdate)
			r.Delete("/delete/{id}", rs.SellerProductsDelete)
		})
		r.Route("/profile", func(r chi.Router) {
			// routes
		})
	})
	r.Route("/admin", func(r chi.Router) {
		// routes
	})
}
