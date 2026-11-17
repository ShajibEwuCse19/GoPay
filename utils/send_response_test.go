package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

func TestWriteError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "request body too large",
			err:        &http.MaxBytesError{Limit: 1024},
			wantStatus: http.StatusRequestEntityTooLarge,
			wantBody:   `"error":"Request body is too large"`,
		},
		{
			name:       "invalid json syntax",
			err:        &json.SyntaxError{Offset: 5},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"Invalid request body. Give me valid json"`,
		},
		{
			name:       "unknown field error",
			err:        errors.New("json: unknown field \"unknown\""),
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"Invalid request body. Give me valid json"`,
		},
		{
			name:       "validation error",
			err:        &domain.ValidationError{Field: "email", Message: "is invalid"},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"email: is invalid"`,
		},
		{
			name:       "conflict error",
			err:        domain.ErrConflict,
			wantStatus: http.StatusConflict,
			wantBody:   `"error":"Phone number or email already registered"`,
		},
		{
			name:       "wrapped conflict error",
			err:        errors.Join(errors.New("db insert"), domain.ErrConflict),
			wantStatus: http.StatusConflict,
			wantBody:   `"error":"Phone number or email already registered"`,
		},
		{
			name:       "timeout error",
			err:        context.DeadlineExceeded,
			wantStatus: http.StatusGatewayTimeout,
			wantBody:   `"error":"Request timed out. Please try again"`,
		},
		{
			name:       "client canceled",
			err:        context.Canceled,
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
		{
			name:       "unknown internal error",
			err:        errors.New("unexpected database crash"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `"error":"An unexpected error occurred"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/test", nil)

			WriteError(rec, req, tt.err)

			if tt.err == context.Canceled {
				if rec.Body.Len() != 0 {
					t.Errorf("expected empty body for context.Canceled, got %q", rec.Body.String())
				}
				return
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body = %q, want substring %q", rec.Body.String(), tt.wantBody)
			}
		})
	}

	t.Run("includes request_id from context in error response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req = req.WithContext(WithRequestID(req.Context(), "req-uuid-test-999"))

		WriteError(rec, req, errors.New("boom"))

		if !strings.Contains(rec.Body.String(), `"request_id":"req-uuid-test-999"`) {
			t.Errorf("expected body to contain request_id, got %q", rec.Body.String())
		}
	})
}
