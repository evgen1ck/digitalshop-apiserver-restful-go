package api

import (
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"net/http"
	"strings"
	"time"
)

func corsMiddleware() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler
}

func notFoundMiddleware() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithNotFound(w)
	}
}

func uriLengthMiddleware(maxLength int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.RequestURI) > maxLength {
				respondWithURITooLong(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func requestSizeMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			err := r.ParseMultipartForm(maxSize)
			if err != nil && err != http.ErrNotMultipart {
				respondWithPayloadTooLarge(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func rateLimitMiddleware(requests int, interval time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(
		requests,
		interval,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			respondWithTooManyRequests(w)
			return
		}),
	)
}

func serviceUnavailableMiddleware(enable bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !enable {
				respondWithServiceUnavailable(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func unsupportedMediaTypeMiddleware(allowedContentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			if contentType != "" && !strings.HasPrefix(contentType, allowedContentType) {
				respondWithUnsupportedMediaType(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func notImplementedMiddleware(allowedMethods []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			allowed := false
			for _, allowedMethod := range allowedMethods {
				if method == allowedMethod {
					allowed = true
					break
				}
			}
			if !allowed {
				respondWithNotImplemented(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func maxHeaderBytesMiddleware(maxHeaderSize int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			totalHeaderSize := 0
			for key, values := range r.Header {
				totalHeaderSize += len(key) + 2 // For colon and space after key
				for _, value := range values {
					totalHeaderSize += len(value) + 2 // For newlines
				}
			}
			if totalHeaderSize > maxHeaderSize {
				respondWithBadRequest(w, "Request headers too large")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
