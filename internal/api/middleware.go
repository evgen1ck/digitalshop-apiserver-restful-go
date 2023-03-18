package api

import (
	"bytes"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func CorsMiddleware() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler
}

func ContentTypeCompressMiddleware(level int, allowedContentTypes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, allowedContentType := range allowedContentTypes {
				if strings.Contains(r.Header.Get("Accept"), allowedContentType) {
					compressMiddleware := middleware.Compress(level, allowedContentType)
					compressMiddleware(next).ServeHTTP(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func NotFoundMiddleware() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RespondWithNotFound(w)
	}
}

func UriLengthMiddleware(maxLength int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.RequestURI) > maxLength {
				RespondWithURITooLong(w, maxLength)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequestSizeMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			err := r.ParseMultipartForm(maxSize)
			if err != nil && err != http.ErrNotMultipart {
				RespondWithPayloadTooLarge(w, maxSize)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimitMiddleware(requests int, interval time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(
		requests,
		interval,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			RespondWithTooManyRequests(w, requests, interval)
			return
		}),
	)
}

func ServiceUnavailableMiddleware(enable bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enable {
				RespondWithServiceUnavailable(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UnsupportedMediaTypeMiddleware(allowedContentTypes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			allowed := false
			if contentType != "" {
				for _, allowedContentType := range allowedContentTypes {
					if strings.HasPrefix(contentType, allowedContentType) {
						allowed = true
						break
					}
				}
			}
			if !allowed {
				RespondWithUnsupportedMediaType(w, allowedContentTypes)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func NotImplementedMiddleware(allowedMethods []string) func(http.Handler) http.Handler {
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
				RespondWithNotImplemented(w, method)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func MaxHeaderBytesMiddleware(maxHeaderSize int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			totalHeaderSize := 0
			for key, values := range r.Header {
				totalHeaderSize += len(key) + 2
				for _, value := range values {
					totalHeaderSize += len(value) + 2
				}
			}
			if totalHeaderSize > maxHeaderSize {
				RespondWithBadRequest(w, "Request headers too large")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func MethodNotAllowedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		if ww.Status() == http.StatusMethodNotAllowed {
			RespondWithMethodNotAllowed(w, r.Method)
			return
		}
	})
}

func NotAcceptableMiddleware(allowedContentTypes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptHeader := r.Header.Get("Accept")
			if acceptHeader != "" {
				acceptsAllowedContentType := false
				for _, allowedContentType := range allowedContentTypes {
					if strings.Contains(acceptHeader, allowedContentType) {
						acceptsAllowedContentType = true
						break
					}
				}
				if !acceptsAllowedContentType {
					RespondWithNotAcceptable(w, allowedContentTypes)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UnprocessableEntityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, values := range r.Header {
			if !utf8.ValidString(key) {
				RespondWithUnprocessableEntity(w, "Header keys must be valid UTF-8. Please use UTF-8 encoding")
				return
			}
			for _, value := range values {
				if !utf8.ValidString(value) {
					RespondWithUnprocessableEntity(w, "Header values must be valid UTF-8. Please use UTF-8 encoding")
					return
				}
			}
		}
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				RespondWithInternalServerError(w)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			if !utf8.Valid(bodyBytes) {
				RespondWithUnprocessableEntity(w, "Request body must be valid UTF-8. Please use UTF-8 encoding")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
