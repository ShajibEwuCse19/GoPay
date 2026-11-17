package middleware

import (
	"net/http"
)

func Preflight(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handlePreFlightRequest(w, r)
			next.ServeHTTP(w, r)
		},
	)
}

func handlePreFlightRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
}
