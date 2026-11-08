package customer

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"POST /api/v1/auth/customers/register",
		h.CreateCustomers,
	)
}