package api_v1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"io"
	"net/http"
	"strconv"
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

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)

		requestDuration.Observe(elapsed.Seconds())
		requestsProcessed.Inc()
	})
}

func NotFoundMiddleware() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RedRespond(w,
			http.StatusNotFound,
			"Not found",
			"The requested resource could not be found")
	}
}

func UriLengthMiddleware(maxLength int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.RequestURI) > maxLength {
				RedRespond(w,
					http.StatusRequestURITooLong,
					"Request URI too long",
					"The requested URI is too long. A maximum of "+strconv.Itoa(maxLength)+" bytes can be sent")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequestSizeMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Calculate total header size
			totalHeaderSize := int64(0)
			for key, values := range r.Header {
				totalHeaderSize += int64(len(key) + 2) // Adding the size of the key, colon and space
				for _, value := range values {
					totalHeaderSize += int64(len(value) + 2) // Adding value size and newline characters (CRLF)
				}
			}
			// Check if the sum of header and body size exceeds maxSize
			contentLength := r.ContentLength
			if contentLength < 0 {
				contentLength = 0
			}
			totalSize := totalHeaderSize + contentLength
			if totalSize > maxSize {
				RedRespond(w,
					http.StatusRequestEntityTooLarge,
					"Request entity too large",
					"The request payload is too large. A maximum of "+strconv.FormatInt(maxSize/1024/1024, 10)+" megabytes can be sent. You have sent "+strconv.FormatInt(totalSize, 10)+" bytes")
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
			RedRespond(w,
				http.StatusTooManyRequests,
				"Too many requests",
				"You have exceeded the allowed number of requests. A maximum of "+strconv.Itoa(requests)+" requests per "+fmt.Sprintf("%.0f", interval.Seconds())+" seconds can be sent")
			return
		}),
	)
}

func ServiceUnavailableMiddleware(enable bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enable {
				RedRespond(w,
					http.StatusServiceUnavailable,
					"Service unavailable",
					"The service is currently unavailable. Please try again later")
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
				RedRespond(w,
					http.StatusUnsupportedMediaType,
					"Unsupported media type",
					"The request contains an unsupported media type. Please use one of "+strings.Join(allowedContentTypes, ", ")+" allowed media types")
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
				RedRespond(w,
					http.StatusNotImplemented,
					"Not implemented",
					"The requested method ("+strings.ToUpper(method)+") is not implemented by the server")
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
			RedRespond(w,
				http.StatusMethodNotAllowed,
				"Method not allowed",
				"The requested method ("+strings.ToUpper(r.Method)+") is not allowed for the specified resource")
			return
		}
	})
}

//func NotAcceptableMiddleware(allowedContentTypes []string) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			acceptHeader := r.Header.Get("Accept")
//			if acceptHeader != "" {
//				acceptsAllowedContentType := false
//				for _, allowedContentType := range allowedContentTypes {
//					if strings.Contains(acceptHeader, allowedContentType) {
//						acceptsAllowedContentType = true
//						break
//					}
//				}
//				if !acceptsAllowedContentType {
//					RedRespond(w,
//						http.StatusNotAcceptable,
//						"Not acceptable",
//						"The provided 'Accept' header does not support the allowed content type. Please use one of "+strings.Join(allowedContentTypes, ", ")+" allowed content types")
//					return
//				}
//			}
//			next.ServeHTTP(w, r)
//		})
//	}
//}

func UnprocessableEntityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, values := range r.Header {
			if !utf8.ValidString(key) {
				RedRespond(w,
					http.StatusUnprocessableEntity,
					"Unprocessable entity",
					"Header keys must be valid UTF-8. Please use UTF-8 encoding")
				return
			}
			for _, value := range values {
				if !utf8.ValidString(value) {
					RedRespond(w,
						http.StatusUnprocessableEntity,
						"Unprocessable entity",
						"Header values must be valid UTF-8. Please use UTF-8 encoding")
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
				RedRespond(w,
					http.StatusUnprocessableEntity,
					"Unprocessable entity",
					"Request body must be valid UTF-8. Please use UTF-8 encoding")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func GatewayTimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					RedRespond(w,
						http.StatusGatewayTimeout,
						"Gateway timeout",
						"The server did not receive a response within "+strconv.FormatFloat(timeout.Seconds(), 'f', 0, 64)+" seconds. Please try again later")
				}
			case <-done:
			}
		})
	}
}

func HttpVersionCheckMiddleware(supportedHttpVersions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, version := range supportedHttpVersions {
				if r.Proto == version {
					next.ServeHTTP(w, r)
					return
				}
			}
			RedRespond(w,
				http.StatusHTTPVersionNotSupported,
				"HTTP version not supported",
				"HTTP version not supported. Please use one of "+strings.Join(supportedHttpVersions, ", ")+" supported http versions")
		})
	}
}
