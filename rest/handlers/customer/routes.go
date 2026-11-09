package customer

import (
	"net/http"

	middleware "github.com/ShajibEwuCse19/gopay-core/rest/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	mux.Handle("POST /api/v1/auth/customers/register",
		manager.AddMiddleware(
			http.HandlerFunc(h.CreateCustomers),
			middleware.Logger, // need to remove
		),
	)
}
