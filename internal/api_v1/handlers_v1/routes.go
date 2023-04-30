package handlers_v1

import (
	"github.com/go-chi/chi/v5"
	"strconv"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/models"
	"test-server-go/internal/storage"
	"time"
)

const (
	timeout            = 5 * time.Second
	rateLimitRequests  = 100
	rateLimitInterval  = 1 * time.Minute
	requestMaxSize     = 4 * 1024 * 1024 // 4 MB
	uriMaxLength       = 1024            // 1024 runes
	serviceUnavailable = false

	//csrfTokenLength    = 32
	//csrfHeaderName     = "X-CSRF-Token"
	//csrfCookieDuration = 30 * time.Minute
)

var (
	allowedContentTypes = []string{"application/json", "text/plain"}
)

type Resolver struct {
	App *models.Application
}

func (rs *Resolver) SetupRouterApiVer1(pathPrefix string) {
	r := chi.NewRouter()

	// CORS settings
	r.Use(api_v1.CorsMiddleware([]string{
		rs.App.Config.App.Service.Url.Api,
		rs.App.Config.App.Service.Url.App,
		"http://localhost:" + strconv.Itoa(rs.App.Config.App.Port)}))

	// Error handling
	r.Use(api_v1.ServiceUnavailableMiddleware(serviceUnavailable))          // Error 503 - Service Unavailable
	r.Use(api_v1.RateLimitMiddleware(rateLimitRequests, rateLimitInterval)) // Error 429 - Too Many Requests
	r.Use(api_v1.UriLengthMiddleware(uriMaxLength))                         // Error 414 - URI Too Long
	r.Use(api_v1.RequestSizeMiddleware(requestMaxSize))                     // Error 413 - Payload Too Large
	r.Use(api_v1.UnprocessableEntityMiddleware)                             // Error 422 - Unprocessable Entity
	r.Use(api_v1.UnsupportedMediaTypeMiddleware(allowedContentTypes))       // Error 415 - Unsupported Media Type
	r.Use(api_v1.MethodNotAllowedMiddleware)                                // Error 405 - Method Not Allowed
	r.Use(api_v1.GatewayTimeoutMiddleware(timeout))                         // Error 504 - Gateway Timeout
	r.NotFound(api_v1.NotFoundMiddleware())                                 // Error 404 - Not Found
	//r.Use(api_v1.CsrfMiddleware(rs.App, csrfTokenLength, csrfHeaderName, csrfCookieDuration)) // Error 403 - Forbidden

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
	})
	r.Route("/products", func(r chi.Router) {
		r.Get("/mainpage", rs.ProductsDataForMainpage)
		r.Get("/", rs.ProductsData)
	})
	r.Route("/user", func(r chi.Router) {
		r.Use(api_v1.JwtAuthMiddleware(rs.App.Postgres, rs.App.Redis, rs.App.Logger, rs.App.Config.App.Jwt, storage.AccountRoleUser))
		r.Route("/profile", func(r chi.Router) {
			r.Patch("/", rs.UserProfileUpdate)
			r.Delete("/", rs.UserProfileDelete)
			r.Get("/orders", rs.UserProfileOrders)
			r.Post("/dump", rs.UserProfileDump)
			r.Route("/image", func(r chi.Router) {
				r.Get("/", rs.UserProfileOrders)
				r.Post("/", rs.UserProfileOrders)
			})
		})
		r.Post("/logout", rs.AuthLogout)
	})
	r.Route("/admin", func(r chi.Router) {
		r.Use(api_v1.JwtAuthMiddleware(rs.App.Postgres, rs.App.Redis, rs.App.Logger, rs.App.Config.App.Jwt, storage.AccountRoleAdmin))
		r.Route("/products", func(r chi.Router) {
			r.Get("/", rs.AdminGetProducts)
			r.Post("/", rs.AdminCreateProduct)
			r.Patch("/{id}", rs.AdminProductsUpdate)
			r.Delete("/{id}", rs.AdminProductsDelete)
		})
		r.Route("/profile", func(r chi.Router) {
			// routes
		})
	})
	r.Route("/resources", func(r chi.Router) {
		r.Get("/product_image/{id}", rs.ResourcesGetProductImage)
		r.Get("/svg_file/{id}", rs.ResourcesGetSvgFile)
	})
}
