package common

import (
	"encoding/json"
	"net/http"
	"time"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Timestamp time.Time    `json:"timestamp"`
	Status    int          `json:"status"`
	Error     string       `json:"error"`
	Message   string       `json:"message"`
	Path      string       `json:"path,omitempty"`
	Details   []FieldError `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, message string, r *http.Request) {
	WriteDetailedError(w, status, message, nil, r)
}

func WriteDetailedError(w http.ResponseWriter, status int, message string, details []FieldError, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Timestamp: time.Now(),
		Status:    status,
		Error:     http.StatusText(status),
		Message:   message,
		Path:      r.URL.Path,
		Details:   details,
	}

	json.NewEncoder(w).Encode(resp)
}
