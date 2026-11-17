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

// SendErrorWithReqID writes an error with request ID: {"error": "message", "request_id": "..."}.
func SendErrorWithReqID(w http.ResponseWriter, status int, message string, reqID string) {
	resp := map[string]string{"error": message}
	if reqID != "" {
		resp["request_id"] = reqID
	}
	SendResponse(w, status, resp)
}

// WriteError centralizes all HTTP, JSON, Domain, and Context errors across the application.
// Handlers should pass any error returned by decoding or services directly to this function.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	reqID := GetRequestID(r.Context())

	var (
		tooLarge      *http.MaxBytesError
		syntaxErr     *json.SyntaxError
		unmarshalErr  *json.UnmarshalTypeError
		validationErr *domain.ValidationError
	)

	switch {
	// Request Body Size Limit Exceeded
	case errors.As(err, &tooLarge):
		SendErrorWithReqID(w, http.StatusRequestEntityTooLarge, "Request body is too large", reqID)

	// Invalid JSON, Unknown Fields, or Empty Body
	case errors.As(err, &syntaxErr),
		errors.As(err, &unmarshalErr),
		errors.Is(err, io.EOF),
		strings.Contains(err.Error(), "json:"):
		SendErrorWithReqID(w, http.StatusBadRequest, "Invalid request body. Give me valid json", reqID)

	// Domain Validation Error (e.g. invalid phone, short PIN)
	case errors.As(err, &validationErr):
		SendErrorWithReqID(w, http.StatusBadRequest, validationErr.Error(), reqID)

	// Resource Conflict (Duplicate phone number or email)
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrCustomerAlreadyExists):
		SendErrorWithReqID(w, http.StatusConflict, "Phone number or email already registered", reqID)

	// Request Timeout
	case errors.Is(err, context.DeadlineExceeded):
		slog.Warn("request timed out", "request_id", reqID, "path", r.URL.Path)
		SendErrorWithReqID(w, http.StatusGatewayTimeout, "Request timed out. Please try again", reqID)

	// Client Disconnected
	case errors.Is(err, context.Canceled):
		slog.Info("client cancelled request", "request_id", reqID, "path", r.URL.Path)

	// Unhandled Internal Server Errors
	default:
		slog.Error("unhandled internal error", "request_id", reqID, "path", r.URL.Path, "error", err)
		SendErrorWithReqID(w, http.StatusInternalServerError, "An unexpected error occurred", reqID)
	}
}
