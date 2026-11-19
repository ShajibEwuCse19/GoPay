package rest

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/health"
	middleware "github.com/ShajibEwuCse19/gopay-core/rest/middlewares"
)

type Server struct {
	customerHandler *customer.Handler
	healthHandler   *health.Handler
	config          *config.Config
}

func NewServer(
	customerHandler *customer.Handler,
	healthHandler *health.Handler,
	config *config.Config,
) *Server {
	return &Server{
		customerHandler: customerHandler,
		healthHandler:   healthHandler,
		config:          config,
	}
}

func (server *Server) StartServer() {
	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(server.config.HTTP.Port),
		Handler:           server.buildHandler(),
		ReadHeaderTimeout: server.config.HTTP.ReadHeaderTimeout,
		ReadTimeout:       server.config.HTTP.ReadTimeout,
		WriteTimeout:      server.config.HTTP.WriteTimeout,
		IdleTimeout:       server.config.HTTP.IdleTimeout,
	}

	// ctx will get cancelled signal on Ctrl+C (SIGINT) or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go server.listen(srv, stop)

	<-ctx.Done() // wait for context cancelation(signal)
	shutdown(srv, server.config.Timeouts.Shutdown)
}

// buildHandler registers routes and wraps the mux with global middlewares.
func (server *Server) buildHandler() http.Handler {
	manager := middleware.NewManager()
	manager.AppendGlobalMiddlewares(
		middleware.RequestID,
		middleware.Recover, // catches panics from every inner middleware and handler
		middleware.Cors(server.config.HTTP.AllowedOrigins),
		middleware.Preflight,
		middleware.Logger,
	)

	mux := http.NewServeMux()
	wrapperMux := manager.WrapMuxWithGlobalMiddleware(mux)
	server.customerHandler.RegisterRoutes(mux, manager)
	if server.healthHandler != nil {
		server.healthHandler.RegisterRoutes(mux)
	}

	return wrapperMux
}

// listen blocks serving requests; on an unexpected error it calls stop to trigger shutdown.
func (server *Server) listen(srv *http.Server, stop context.CancelFunc) {
	slog.Info("server listening", "addr", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("error starting the server", "error", err)
		stop()
	}
}

// shutdown stops accepting new requests, waits for in-flight ones, then releases the port.
func shutdown(srv *http.Server, timeout time.Duration) {
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		return
	}
	slog.Info("server stopped")
}
