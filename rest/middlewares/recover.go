package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recover catches a panic in a handler, logs it, and answers 500 instead of dropping the connection.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					slog.Error("panic recovered",
						"panic", rec,
						"method", r.Method,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal server error"}`))
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}
