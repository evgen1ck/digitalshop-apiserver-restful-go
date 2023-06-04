package handlers_v1

import (
	"github.com/go-chi/chi/v5"
	"strconv"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/freekassa"
	"test-server-go/internal/models"
	"test-server-go/internal/storage"
	"time"
)

const (
	timeout            = 50 * time.Second
	rateLimitRequests  = 1000
	rateLimitInterval  = 1 * time.Minute
	requestMaxSize     = 4 * 1024 * 1024 // 4 MB
	uriMaxLength       = 1024            // 1024 runes
	serviceUnavailable = false

	//csrfTokenLength    = 32
	//csrfHeaderName     = "X-CSRF-Token"
	//csrfCookieDuration = 30 * time.Minute
)

type Resolver struct {
	App *models.Application
}

func (rs *Resolver) SetupRouterApiVer1(pathPrefix string) {
	r := chi.NewRouter()

	var corsAllowedOrigins []string
	if rs.App.Config.App.Debug {
		corsAllowedOrigins = append(corsAllowedOrigins, "http://localhost:"+strconv.Itoa(rs.App.Config.App.Port))
	} else {
		corsAllowedOrigins = append(corsAllowedOrigins, rs.App.Config.App.Service.Url.Client)
	}

	// CORS settings
	r.Use(api_v1.CorsMiddleware(corsAllowedOrigins))

	// Error handling
	r.Use(api_v1.ServiceUnavailableMiddleware(serviceUnavailable))          // Error 503 - ServiceName Unavailable
	r.Use(api_v1.RateLimitMiddleware(rateLimitRequests, rateLimitInterval)) // Error 429 - Too Many Requests
	r.Use(api_v1.UriLengthMiddleware(uriMaxLength))                         // Error 414 - URI Too Long
	r.Use(api_v1.RequestSizeMiddleware(requestMaxSize))                     // Error 413 - Payload Too Large
	r.Use(api_v1.UnprocessableEntityMiddleware)                             // Error 422 - Unprocessable Entity
	r.Use(api_v1.MethodNotAllowedMiddleware)                                // Error 405 - Method Not Allowed
	r.Use(api_v1.GatewayTimeoutMiddleware(timeout))                         // Error 504 - Gateway Timeout
	r.NotFound(api_v1.NotFoundMiddleware())                                 // Error 404 - Not Found
	//r.Use(api_v1.CsrfMiddleware(rs.Client, csrfTokenLength, csrfHeaderName, csrfCookieDuration)) // Error 403 - Forbidden

	rs.registerRoutes(r)

	rs.App.Router.Mount(pathPrefix, r)
}

func (rs *Resolver) registerRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", rs.AuthSignup)
		r.Post("/signup-with-token", rs.AuthSignupWithToken)
		r.Post("/login", rs.AuthLogin)
		r.Post("/alogin", rs.AuthAlogin)
		//r.Post("/login-with-token", rs.AuthLoginWithToken)
		//r.Post("/recover-password", rs.AuthRecoverPassword)
		//r.Post("/recover-password-with-token", rs.AuthRecoverPasswordWithToken)
	})
	r.Route("/product", func(r chi.Router) {
		r.Get("/mainpage", rs.ProductsDataForMainpage)
	})
	r.Route("/user", func(r chi.Router) {
		r.Use(api_v1.JwtAuthMiddleware(rs.App.Postgres, rs.App.Redis, rs.App.Logger, rs.App.Config.App.Jwt, storage.AccountRoleUser))
		r.Get("/order", rs.UserProfileOrders)
		r.Post("/payment", rs.UserNewPayment)
		r.Route("/profile", func(r chi.Router) {
			r.Patch("/", rs.UserProfileUpdate)
			r.Delete("/", rs.UserProfileDelete)
		})
		r.Post("/logout", rs.AuthLogout)
		//r.Post("/dump", rs.UserProfileDump)
	})
	r.Route("/admin", func(r chi.Router) {
		r.Use(api_v1.JwtAuthMiddleware(rs.App.Postgres, rs.App.Redis, rs.App.Logger, rs.App.Config.App.Jwt, storage.AccountRoleAdmin))
		r.Route("/product", func(r chi.Router) {
			r.Get("/", rs.AdminGetProducts)
			r.Delete("/", rs.AdminDeleteProduct)
		})
		r.Route("/service", func(r chi.Router) {
			r.Get("/", rs.AdminGetServices)
			r.Post("/", rs.AdminAddService)
			r.Post("/svg", rs.AdminAddService)
			r.Delete("/", rs.AdminDeleteService)
		})
		r.Route("/state", func(r chi.Router) {
			r.Get("/", rs.AdminGetStates)
		})
		r.Route("/item", func(r chi.Router) {
			r.Get("/", rs.AdminGetItems)
		})
		r.Route("/type", func(r chi.Router) {
			r.Get("/", rs.AdminGetTypes)
			r.Post("/", rs.AdminAddType)
			r.Delete("/", rs.AdminDeleteType)
		})
		r.Route("/subtype", func(r chi.Router) {
			r.Get("/", rs.AdminGetSubtypes)
			r.Post("/", rs.AdminAddSubtype)
			r.Delete("/", rs.AdminDeleteSubtype)
		})
		r.Route("/variant", func(r chi.Router) {
			r.Get("/", rs.AdminGetVariants)
			r.Post("/", rs.AdminCreateVariant)
			r.Patch("/", rs.AdminUpdateVariant)
			r.Delete("/", rs.AdminDeleteVariant)
			r.Route("/upload", func(r chi.Router) {
				r.Get("/", rs.AdminGetVariantUploads)
				r.Post("/", rs.AdminUploadVariant)
				r.Delete("/", rs.AdminDeleteVariantUpload)
			})
		})
		r.Route("/database", func(r chi.Router) {
			r.Route("/postgres", func(r chi.Router) {
				r.Get("/info", rs.ServerDatabasesPostgresInfo)
				r.Post("/backup", rs.ServerDatabasesPostgresBackup)
			})
		})
		r.Post("/logout", rs.AuthLogout)
	})
	r.Route("/resources", func(r chi.Router) {
		r.Get("/product_image/{id}", rs.ResourcesGetProductImage)
		r.Get("/svg/{id}", rs.ResourcesGetSvgFile)
	})
	r.Route("/freekassa", func(r chi.Router) {
		r.Use(api_v1.FreekassaIpWhitelistMiddleware(freekassa.AllowedFreekassaIPs, rs.App.Config.App.Service.Url.Client+"/finish"))
		r.Get("/notification", rs.FreekassaNotification)
	})
}
