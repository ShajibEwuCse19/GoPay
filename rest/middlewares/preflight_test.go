package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPreflight_StopsOptionsRequest(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusTeapot)
	})

	handler := Preflight(nextHandler)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/customers", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if nextCalled {
		t.Fatalf("expected next handler NOT to be called on OPTIONS, but it was invoked")
	}
}

func TestPreflight_PassesNonOptionsRequest(t *testing.T) {
	testMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range testMethods {
		t.Run(method, func(t *testing.T) {
			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := Preflight(nextHandler)

			req := httptest.NewRequest(method, "/api/v1/customers", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if !nextCalled {
				t.Fatalf("expected next handler to be called for method %s, but it was not", method)
			}
			if rec.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
			}
		})
	}
}
