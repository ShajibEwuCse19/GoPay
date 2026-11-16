package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

// SendResponse marshals data to JSON and sets the HTTP status code.
func SendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// SendError writes an error in one consistent JSON shape: {"error": "message"}.
func SendError(w http.ResponseWriter, status int, message string) {
	SendResponse(w, status, map[string]string{"error": message})
}

// WriteError centralizes all HTTP, JSON, Domain, and Context errors across the application.
// Handlers should pass any error returned by decoding or services directly to this function.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		tooLarge      *http.MaxBytesError
		syntaxErr     *json.SyntaxError
		unmarshalErr  *json.UnmarshalTypeError
		validationErr *domain.ValidationError
	)

	switch {
	// Request Body Size Limit Exceeded
	case errors.As(err, &tooLarge):
		SendError(w, http.StatusRequestEntityTooLarge, "Request body is too large")

	// Invalid JSON, Unknown Fields, or Empty Body
	case errors.As(err, &syntaxErr),
		errors.As(err, &unmarshalErr),
		errors.Is(err, io.EOF),
		strings.HasPrefix(err.Error(), "json:"):
		SendError(w, http.StatusBadRequest, "Invalid request body. Give me valid json")

	// Domain Validation Error (e.g. invalid phone, short PIN)
	case errors.As(err, &validationErr):
		SendError(w, http.StatusBadRequest, validationErr.Error())

	// Resource Conflict (Duplicate phone number or email)
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrCustomerAlreadyExists):
		SendError(w, http.StatusConflict, "Phone number or email already registered")

	// Request Timeout
	case errors.Is(err, context.DeadlineExceeded):
		slog.Warn("request timed out", "path", r.URL.Path)
		SendError(w, http.StatusGatewayTimeout, "Request timed out. Please try again")

	// Client Disconnected
	case errors.Is(err, context.Canceled):
		slog.Info("client cancelled request", "path", r.URL.Path)

	// Unhandled Internal Server Errors
	default:
		slog.Error("unhandled internal error", "path", r.URL.Path, "error", err)
		SendError(w, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
