package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	StatusCode  int    `json:"status_code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func respondWithTooManyRequests(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	response := ErrorResponse{
		StatusCode: http.StatusTooManyRequests,
		Message:    "Too many requests",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{
		StatusCode: http.StatusNotFound,
		Message:    "Not found",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	response := ErrorResponse{
		StatusCode:  http.StatusBadRequest,
		Message:     "Bad request",
		Description: message,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithPayloadTooLarge(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestEntityTooLarge)
	response := ErrorResponse{
		StatusCode: http.StatusRequestEntityTooLarge,
		Message:    "Payload too large",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithURITooLong(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestURITooLong)
	response := ErrorResponse{
		StatusCode: http.StatusRequestURITooLong,
		Message:    "URI too long",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithServiceUnavailable(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	response := ErrorResponse{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service unavailable",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithUnsupportedMediaType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnsupportedMediaType)
	response := ErrorResponse{
		StatusCode: http.StatusUnsupportedMediaType,
		Message:    "Unsupported media type",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithNotImplemented(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	response := ErrorResponse{
		StatusCode: http.StatusNotImplemented,
		Message:    "Method not implemented",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithMethodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	response := ErrorResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Message:    "Method not allowed",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithCreated(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(result)
}
