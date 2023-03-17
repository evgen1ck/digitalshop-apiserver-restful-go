package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ErrorResponse struct {
	StatusCode  int    `json:"status_code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// Red responds
func respondWithTooManyRequests(w http.ResponseWriter, requests int, interval time.Duration) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	secondsStr := fmt.Sprintf("%.0f", interval.Seconds())
	response := ErrorResponse{
		StatusCode:  http.StatusTooManyRequests,
		Message:     "Too many requests",
		Description: "You have exceeded the allowed number of requests. A maximum of " + strconv.Itoa(requests) + " requests per " + secondsStr + " seconds can be sent",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{
		StatusCode:  http.StatusNotFound,
		Message:     "Not found",
		Description: "The requested resource could not be found",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := ErrorResponse{
		StatusCode:  http.StatusInternalServerError,
		Message:     "Internal server error",
		Description: "An internal server error has occurred",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithBadRequest(w http.ResponseWriter, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	response := ErrorResponse{
		StatusCode:  http.StatusBadRequest,
		Message:     "Bad request",
		Description: description,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithPayloadTooLarge(w http.ResponseWriter, maxSize int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestEntityTooLarge)
	response := ErrorResponse{
		StatusCode:  http.StatusRequestEntityTooLarge,
		Message:     "Payload too large",
		Description: "The request payload is too large. A maximum of " + strconv.FormatInt(maxSize/1024/1024, 10) + " megabytes can be sent",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithURITooLong(w http.ResponseWriter, maxLen int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestURITooLong)
	response := ErrorResponse{
		StatusCode:  http.StatusRequestURITooLong,
		Message:     "URI too long",
		Description: "The requested URI is too long. A maximum of " + strconv.Itoa(maxLen) + " bytes can be sent",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithServiceUnavailable(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	response := ErrorResponse{
		StatusCode:  http.StatusServiceUnavailable,
		Message:     "Service unavailable",
		Description: "The service is currently unavailable. Please try again later",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithUnsupportedMediaType(w http.ResponseWriter, allowedContentType string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnsupportedMediaType)
	response := ErrorResponse{
		StatusCode:  http.StatusUnsupportedMediaType,
		Message:     "Unsupported media type",
		Description: "The request contains an unsupported media type. Please use " + allowedContentType + " media type",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithNotImplemented(w http.ResponseWriter, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	response := ErrorResponse{
		StatusCode:  http.StatusNotImplemented,
		Message:     "Method not implemented",
		Description: "The requested method (" + strings.ToUpper(method) + ") is not implemented by the server",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithMethodNotAllowed(w http.ResponseWriter, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	response := ErrorResponse{
		StatusCode:  http.StatusMethodNotAllowed,
		Message:     "Method not allowed",
		Description: "The requested method (" + strings.ToUpper(method) + ") is not allowed for the specified resource",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithNotAcceptable(w http.ResponseWriter, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotAcceptable)
	response := ErrorResponse{
		StatusCode:  http.StatusNotAcceptable,
		Message:     "Not acceptable",
		Description: description,
	}
	_ = json.NewEncoder(w).Encode(response)
}

// Green responds
func respondWithCreated(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}
