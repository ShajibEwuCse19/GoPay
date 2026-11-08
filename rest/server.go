package rest

import (
	"net/http"

	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
)

type Server struct {
	customerHandler *customer.Handler
}

func NewServer(customerHandler *customer.Handler) *Server {
	return &Server{
		customerHandler: customerHandler,
	}
}

func (server *Server) StartServer() {
	mux := http.NewServeMux()

	server.customerHandler.RegisterRoutes(mux)
}