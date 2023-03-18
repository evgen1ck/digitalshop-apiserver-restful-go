package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type errorResponse struct {
	StatusCode  int    `json:"status_code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// Red responds

func RespondWithTooManyRequests(w http.ResponseWriter, requests int, interval time.Duration) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	secondsStr := fmt.Sprintf("%.0f", interval.Seconds())
	response := errorResponse{
		StatusCode:  http.StatusTooManyRequests,
		Message:     "Too many requests",
		Description: "You have exceeded the allowed number of requests. A maximum of " + strconv.Itoa(requests) + " requests per " + secondsStr + " seconds can be sent",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response := errorResponse{
		StatusCode:  http.StatusNotFound,
		Message:     "Not found",
		Description: "The requested resource could not be found",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := errorResponse{
		StatusCode:  http.StatusInternalServerError,
		Message:     "Internal server error",
		Description: "An internal server error has occurred. Please try again, and if it doesn't help, contact technical support",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithBadRequest(w http.ResponseWriter, description string) {
	w.Header().Set("Content-Type", "application/json")
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

func RespondWithPayloadTooLarge(w http.ResponseWriter, maxSize int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestEntityTooLarge)
	response := errorResponse{
		StatusCode:  http.StatusRequestEntityTooLarge,
		Message:     "Payload too large",
		Description: "The request payload is too large. A maximum of " + strconv.FormatInt(maxSize/1024/1024, 10) + " megabytes can be sent",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithURITooLong(w http.ResponseWriter, maxLen int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestURITooLong)
	response := errorResponse{
		StatusCode:  http.StatusRequestURITooLong,
		Message:     "URI too long",
		Description: "The requested URI is too long. A maximum of " + strconv.Itoa(maxLen) + " bytes can be sent",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithServiceUnavailable(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	response := errorResponse{
		StatusCode:  http.StatusServiceUnavailable,
		Message:     "Service unavailable",
		Description: "The service is currently unavailable. Please try again later",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithUnsupportedMediaType(w http.ResponseWriter, allowedContentTypes []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnsupportedMediaType)
	response := errorResponse{
		StatusCode:  http.StatusUnsupportedMediaType,
		Message:     "Unsupported media type",
		Description: "The request contains an unsupported media type. Please use " + strings.Join(allowedContentTypes, ", ") + " media type(s)",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithNotAcceptable(w http.ResponseWriter, allowedContentTypes []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotAcceptable)
	response := errorResponse{
		StatusCode:  http.StatusNotAcceptable,
		Message:     "Not acceptable",
		Description: "The provided 'Accept' header does not support the allowed content type. Please use " + strings.Join(allowedContentTypes, ", ") + " content type(s)",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithNotImplemented(w http.ResponseWriter, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	response := errorResponse{
		StatusCode:  http.StatusNotImplemented,
		Message:     "Method not implemented",
		Description: "The requested method (" + strings.ToUpper(method) + ") is not implemented by the server",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithMethodNotAllowed(w http.ResponseWriter, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	response := errorResponse{
		StatusCode:  http.StatusMethodNotAllowed,
		Message:     "Method not allowed",
		Description: "The requested method (" + strings.ToUpper(method) + ") is not allowed for the specified resource",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithUnprocessableEntity(w http.ResponseWriter, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	response := errorResponse{
		StatusCode:  http.StatusUnprocessableEntity,
		Message:     "Unprocessable entity",
		Description: description,
	}
	_ = json.NewEncoder(w).Encode(response)
}

// Green responds

func RespondWithCreated(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}
