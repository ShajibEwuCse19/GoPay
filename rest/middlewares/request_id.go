package middleware

import (
	"net/http"

	"github.com/ShajibEwuCse19/gopay-core/utils"
)

const HeaderXRequestID = "X-Request-ID"

// RequestID extracts incoming X-Request-ID or generates a new one,
// injects it into context, and sets the response header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(HeaderXRequestID)
		if reqID == "" {
			reqID = utils.GenerateRequestID()
		}

		ctx := utils.WithRequestID(r.Context(), reqID)
		w.Header().Set(HeaderXRequestID, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
