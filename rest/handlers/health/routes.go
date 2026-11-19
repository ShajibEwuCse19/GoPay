package health

import "net/http"

// RegisterRoutes registers liveness and readiness probe routes.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health/live", h.Live)
	mux.HandleFunc("GET /health/ready", h.Ready)
}
