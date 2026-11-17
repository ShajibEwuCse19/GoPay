package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/utils"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			slog.Info("request",
				"request_id", utils.GetRequestID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start).String(),
			)
		},
	)
}
