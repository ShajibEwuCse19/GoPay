package health

import (
	"context"
	"net/http"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/utils"
)

const defaultPingTimeout = 2 * time.Second

type Handler struct {
	db Pinger
}

func NewHandler(db Pinger) *Handler {
	return &Handler{db: db}
}

// Live handles GET /health/live - always returns 200 OK if the HTTP process is alive.
func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, http.StatusOK, map[string]string{
		"status": "UP",
	})
}

// Ready handles GET /health/ready - returns 200 if dependencies respond, 503 if any check fails.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultPingTimeout)
	defer cancel()

	if h.db != nil {
		if err := h.db.PingContext(ctx); err != nil {
			utils.SendResponse(w, http.StatusServiceUnavailable, map[string]any{
				"status": "DOWN",
				"error":  "database unavailable",
			})
			return
		}
	}

	utils.SendResponse(w, http.StatusOK, map[string]any{
		"status": "UP",
		"checks": map[string]string{
			"database": "UP",
		},
	})
}
