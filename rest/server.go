package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
	middleware "github.com/ShajibEwuCse19/gopay-core/rest/middlewares"
)

const shutdownTimeout = 10 * time.Second

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
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(server.config.HttpPort),
		Handler: server.buildHandler(),
	}

	// ctx will get cancelled signal on Ctrl+C (SIGINT) or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go server.listen(srv, stop)

	<-ctx.Done() // wait for context cancelation(signal)
	shutdown(srv)
}

// buildHandler registers routes and wraps the mux with global middlewares.
func (server *Server) buildHandler() http.Handler {
	manager := middleware.NewManager()
	manager.AppendGlobalMiddlewares(
		middleware.Cors,
		middleware.Preflight,
		middleware.Logger,
	)

	mux := http.NewServeMux()
	wrapperMux := manager.WrapMuxWithGlobalMiddleware(mux)
	server.customerHandler.RegisterRoutes(mux, manager)

	return wrapperMux
}

// listen blocks serving requests; on an unexpected error it calls stop to trigger shutdown.
func (server *Server) listen(srv *http.Server, stop context.CancelFunc) {
	fmt.Println("Server listening on", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Error starting the server", err)
		stop()
	}
}

// shutdown stops accepting new requests, waits for in-flight ones, then releases the port.
func shutdown(srv *http.Server) {
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Graceful shutdown failed:", err)
		return
	}
	fmt.Println("Server stopped")
}
