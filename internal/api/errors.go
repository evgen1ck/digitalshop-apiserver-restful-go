package api

import (
	"encoding/json"
	"net/http"
)

func respondWithTooManyRequests(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	response := struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		ErrorCode  string `json:"error_code"`
		ResponseID string `json:"response_id"`
	}{
		StatusCode: http.StatusTooManyRequests,
		Message:    "Too many requests",
		ErrorCode:  "too_many_requests",
		ResponseID: "15095f25-aac3-4d60-a788-96cb5136f186",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithResourceNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response := struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		ErrorCode  string `json:"error_code"`
		ResponseID string `json:"response_id"`
	}{
		StatusCode: http.StatusNotFound,
		Message:    "Resource not found",
		ErrorCode:  "not_found",
		ResponseID: "15095f25-aac3-4d60-a788-96cb5136f186",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		ErrorCode  string `json:"error_code"`
		ResponseID string `json:"response_id"`
	}{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		ErrorCode:  "internal_server_error",
		ResponseID: "15095f25-aac3-4d60-a788-96cb5136f186",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func respondWithBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	response := struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		ErrorCode  string `json:"error_code"`
		ResponseID string `json:"response_id"`
	}{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		ErrorCode:  "bad_request",
		ResponseID: "15095f25-aac3-4d60-a788-96cb5136f186",
	}
	_ = json.NewEncoder(w).Encode(response)
	return
}
