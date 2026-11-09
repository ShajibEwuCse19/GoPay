package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
	middleware "github.com/ShajibEwuCse19/gopay-core/rest/middlewares"
)

type Server struct {
	customerHandler *customer.Handler
	config          *config.Config
}

func NewServer(
	customerHandler *customer.Handler,
	config *config.Config,
) *Server {
	return &Server{
		customerHandler: customerHandler,
		config:          config,
	}
}

func (server *Server) StartServer() {
	manager := middleware.NewManager()
	manager.AppendGlobalMiddlewares(
		middleware.Cors,
		middleware.Preflight,
		middleware.Logger,
	)

	mux := http.NewServeMux()
	wrapperMux := manager.WrapMuxWithGlobalMiddleware(mux)
	server.customerHandler.RegisterRoutes(mux, manager)

	err := http.ListenAndServe(":"+strconv.Itoa(server.config.HttpPort), wrapperMux)
	if err != nil {
		fmt.Println("Error starting the server", err)
	}
}
