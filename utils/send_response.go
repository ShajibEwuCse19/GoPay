package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

func SendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var validationErr *domain.ValidationError

	switch {
	case errors.As(err, &validationErr):
		SendError(w, http.StatusBadRequest, validationErr.Error())

	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrCustomerAlreadyExists):
		SendError(w, http.StatusConflict, "Phone number or email already registered")

	case errors.Is(err, context.DeadlineExceeded):
		slog.Warn("request timed out", "path", r.URL.Path)
		SendError(w, http.StatusGatewayTimeout, "Request timed out. Please try again")

	case errors.Is(err, context.Canceled):
		// Client disconnected, do not write a response body
		slog.Info("client cancelled request", "path", r.URL.Path)

	default:
		slog.Error("unhandled internal error", "path", r.URL.Path, "error", err)
		SendError(w, http.StatusInternalServerError, "An unexpected error occurred")
	}
}

func SendError(w http.ResponseWriter, status int, message string) {
	SendResponse(w, status, map[string]string{"error": message})
}
