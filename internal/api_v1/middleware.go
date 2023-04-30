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
	"test-server-go/internal/auth"
	"test-server-go/internal/logger"
	"test-server-go/internal/storage"
	tl "test-server-go/internal/tools"
	"time"
	"unicode/utf8"
)

func CorsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler
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

func ServiceUnavailableMiddleware(unavailable bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if unavailable {
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
			if !tl.ContainsStringInSlice(contentType, allowedContentTypes) {
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

func MethodNotAllowedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		if ww.Status() == http.StatusMethodNotAllowed {
			RedRespond(w,
				http.StatusMethodNotAllowed,
				"Method not allowed",
				"The requested method ("+r.Method+") is not allowed for the specified resource")
			return
		}
	})
}

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

//func CsrfMiddleware(app *models.Application, csrfTokenLength int, csrfName string, csrfCookieDuration time.Duration) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			safeMethods := []string{"GET", "HEAD", "OPTIONS", "TRACE"}
//
//			if !tl.StringInSlice(r.Method, safeMethods) {
//				// Get CSRF token from request header
//				requestToken := r.Header.Get(csrfName)
//				if requestToken == "" {
//					RedRespond(w, http.StatusForbidden, "Forbidden", "CSRF token not found in header")
//					return
//				}
//
//				// Get CSRF token from request cookie
//				cookie, err := r.Cookie(csrfName)
//				if err != nil {
//					RedRespond(w, http.StatusForbidden, "Forbidden", "CSRF token not found in cookie")
//					return
//				}
//
//				// Check if the CSRF tokens match and is valid in the server-side store
//				existsCsrfToken, err := storage.CheckCsrfTokenExists(r.Context(), app.Postgres, requestToken)
//				if err != nil {
//					RespondWithInternalServerError(w)
//				}
//				if requestToken != cookie.Value || !existsCsrfToken {
//					RedRespond(w, http.StatusForbidden, "Forbidden", "Invalid CSRF token")
//					return
//				}
//
//				// Remove the used token from the server-side store
//				err = storage.DeleteCsrfToken(r.Context(), app.Postgres, requestToken)
//				if err != nil {
//					RespondWithInternalServerError(w)
//				}
//			}
//
//			// Generate a new CSRF token
//			token, err := tl.GenerateRandomString(csrfTokenLength)
//			if err != nil {
//				RespondWithInternalServerError(w)
//				return
//			}
//
//			// Save the generated token in the server-side store
//			err = storage.CreateCsrfToken(r.Context(), app.Postgres, token)
//			if err != nil {
//				RespondWithInternalServerError(w)
//			}
//
//			// Create a CSRF cookie with the new token
//			csrfCookie := &http.Cookie{
//				Name:     csrfName,
//				Value:    token,
//				Path:     "/",
//				Expires:  time.Now().Add(csrfCookieDuration),
//				HttpOnly: true,
//				Secure:   true, // Send the cookie only over HTTPS
//			}
//
//			// Set the CSRF cookie
//			http.SetCookie(w, csrfCookie)
//
//			// Set the CSRF token header for the response
//			w.Header().Set(csrfName, token)
//
//			// Call the next handler in the chain
//			next.ServeHTTP(w, r)
//		})
//	}
//}

func JwtAuthMiddleware(pdb *storage.Postgres, rdb *storage.Redis, logger *logger.Logger, secret, role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				RedRespond(w, http.StatusUnauthorized, "Unauthorized", "No authorization header provided")
				return
			}
			if !strings.HasPrefix(authHeader, "Bearer ") {
				RedRespond(w, http.StatusUnauthorized, "Unauthorized", "No bearer in header provided")
				return
			}

			// Get token
			tokenString := strings.Replace(strings.TrimSpace(authHeader), "Bearer ", "", 1)
			if tokenString == "" {
				RedRespond(w, http.StatusUnauthorized, "Unauthorized", "No token header provided")
				return
			}

			// Parse jwt token
			jwtData, err := auth.ParseJwtToken(tokenString, secret)
			if err != nil {
				RedRespond(w, http.StatusUnauthorized, "Unauthorized", "Invalid token")
				return
			}

			// Check token expiration
			if jwtData.ExpiresAt.Time.Before(time.Now()) {
				RedRespond(w, http.StatusForbidden, "Unauthorized", "Token expired")
				return
			}

			// Check if token is in blacklist
			tokenInStopList, err := storage.CheckBlockedTokenExists(r.Context(), rdb, tokenString)
			if err != nil {
				RespondWithInternalServerError(w)
				logger.NewWarn("Error in founding account in the list", err)
				return
			}
			if tokenInStopList {
				RedRespond(w, http.StatusForbidden, "Forbidden", "This token in stop-list")
				return
			}

			// Get account state and check on exists
			state, err := storage.GetStateAccount(r.Context(), pdb, jwtData.AccountUuid, role)
			if state == "" {
				RedRespond(w, http.StatusUnauthorized, "Unauthorized", "The account was not found in the list of "+role+"s")
				return
			} else if err != nil {
				RespondWithInternalServerError(w)
				logger.NewWarn("Error in founding account in the list", err)
				return
			}

			// Check account on state (blocked, deleted...)
			switch state {
			case storage.AccountStateBlocked:
				RedRespond(w, http.StatusForbidden, "Forbidden", "This account has been blocked")
				return
			case storage.AccountStateDeleted:
				RedRespond(w, http.StatusForbidden, "Forbidden", "This account has been deleted")
				return
			}

			// Set context key
			r = ContextSetAuthenticated(r, tokenString, jwtData)
			if err != nil {
				RespondWithInternalServerError(w)
				logger.NewWarn("Error in setting auth context key", err)
				return
			}

			// Update last account activity
			storage.UpdateLastAccountActivity(r.Context(), pdb, jwtData.AccountUuid)

			next.ServeHTTP(w, r)
		})
	}
}
