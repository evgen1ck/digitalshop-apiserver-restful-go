package api_v1

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	StatusCode  int    `json:"status_code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func setStandardHeadersForJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

// Red responds

func RedRespond(w http.ResponseWriter, statusCode int, message, description string) {
	setStandardHeadersForJson(w)
	w.WriteHeader(statusCode)
	response := errorResponse{
		StatusCode:  statusCode,
		Message:     message,
		Description: description,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithInternalServerError(w http.ResponseWriter) {
	setStandardHeadersForJson(w)
	w.WriteHeader(http.StatusInternalServerError)
	response := errorResponse{
		StatusCode:  http.StatusInternalServerError,
		Message:     "Internal server error",
		Description: "An internal server error has occurred. Please try again, and if it doesn't help, contact technical support",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithBadRequest(w http.ResponseWriter, description string) {
	setStandardHeadersForJson(w)
	w.WriteHeader(http.StatusBadRequest)
	if description == "" {
		description = "Invalid request payload. Please double-check the data you are sending, and if this doesn't help, contact technical support"
	}
	response := errorResponse{
		StatusCode:  http.StatusBadRequest,
		Message:     "Bad request",
		Description: description,
	}
	_ = json.NewEncoder(w).Encode(response)
}

// Green responds

func RespondWithCreated(w http.ResponseWriter, result interface{}) {
	setStandardHeadersForJson(w)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}
