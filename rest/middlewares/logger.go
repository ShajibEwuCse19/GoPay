package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			end := time.Now()
			timeDiff := end.Sub(start)
			fmt.Printf("Logger: %s %s %s\n", r.Method, r.URL.Path, timeDiff.String())
		},
	)
}
