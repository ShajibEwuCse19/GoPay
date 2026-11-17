package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/ShajibEwuCse19/gopay-core/utils"
)

// Recover catches a panic in a handler, logs it, and answers 500 instead of dropping the connection.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					reqID := utils.GetRequestID(r.Context())
					slog.Error("panic recovered",
						"request_id", reqID,
						"panic", rec,
						"method", r.Method,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					res := map[string]string{
						"error": "Internal server error",
					}
					if reqID != "" {
						res["request_id"] = reqID
					}
					_ = json.NewEncoder(w).Encode(res)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}
