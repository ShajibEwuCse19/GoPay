package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShajibEwuCse19/gopay-core/utils"
)

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("generates new request ID when none provided", func(t *testing.T) {
		var capturedCtxID string

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtxID = utils.GetRequestID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		handler := RequestID(nextHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if capturedCtxID == "" {
			t.Fatal("expected request ID to be set in context, got empty")
		}

		headerID := rec.Header().Get(HeaderXRequestID)
		if headerID == "" {
			t.Fatal("expected X-Request-ID response header to be set, got empty")
		}

		if capturedCtxID != headerID {
			t.Errorf("context ID %q does not match header ID %q", capturedCtxID, headerID)
		}
	})

	t.Run("preserves incoming X-Request-ID header", func(t *testing.T) {
		incomingID := "client-trace-id-999"
		var capturedCtxID string

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtxID = utils.GetRequestID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		handler := RequestID(nextHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(HeaderXRequestID, incomingID)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if capturedCtxID != incomingID {
			t.Errorf("context ID = %q, want preserved %q", capturedCtxID, incomingID)
		}

		headerID := rec.Header().Get(HeaderXRequestID)
		if headerID != incomingID {
			t.Errorf("header ID = %q, want preserved %q", headerID, incomingID)
		}
	})
}
