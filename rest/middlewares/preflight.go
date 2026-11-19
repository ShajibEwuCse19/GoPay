package middleware

import (
	"net/http"
)

// Preflight intercepts HTTP OPTIONS requests, responds with 200 OK, and short-circuits the middleware chain.
func Preflight(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}
