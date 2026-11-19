package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors_AllowedOrigin(t *testing.T) {
	allowed := []string{"http://localhost:3000", "https://wallet.gopay.com"}
	corsMiddleware := Cors(allowed)

	dummyHandler := func(called *bool) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*called = true
			w.WriteHeader(http.StatusOK)
		})
	}

	t.Run("sets allow origin and credentials when origin is in whitelist", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rec := httptest.NewRecorder()

		corsMiddleware(dummyHandler(&called)).ServeHTTP(rec, req)

		if !called {
			t.Fatal("expected downstream handler to be called")
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:3000")
		}
		if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
			t.Errorf("Access-Control-Allow-Credentials = %q, want %q", got, "true")
		}
		if got := rec.Header().Get("Vary"); got != "Origin" {
			t.Errorf("Vary = %q, want %q", got, "Origin")
		}
		if got := rec.Header().Get("Access-Control-Expose-Headers"); got != "X-Request-ID" {
			t.Errorf("Access-Control-Expose-Headers = %q, want %q", got, "X-Request-ID")
		}
	})

	t.Run("does not set allow origin when origin is not in whitelist", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "https://evil-hacker.com")
		rec := httptest.NewRecorder()

		corsMiddleware(dummyHandler(&called)).ServeHTTP(rec, req)

		if !called {
			t.Fatal("expected downstream handler to be called")
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("expected empty Access-Control-Allow-Origin, got %q", got)
		}
		if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "" {
			t.Errorf("expected empty Access-Control-Allow-Credentials, got %q", got)
		}
	})

	t.Run("handles request without Origin header", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		corsMiddleware(dummyHandler(&called)).ServeHTTP(rec, req)

		if !called {
			t.Fatal("expected downstream handler to be called")
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("expected empty Access-Control-Allow-Origin, got %q", got)
		}
	})

	t.Run("matches wildcard in allowed origins", func(t *testing.T) {
		wildcardCors := Cors([]string{"*"})
		called := false
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "https://random-site.com")
		rec := httptest.NewRecorder()

		wildcardCors(dummyHandler(&called)).ServeHTTP(rec, req)

		if !called {
			t.Fatal("expected downstream handler to be called")
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://random-site.com" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "https://random-site.com")
		}
	})
}
