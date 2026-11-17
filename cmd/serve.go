package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/infra/db"
	"github.com/ShajibEwuCse19/gopay-core/infra/logger"
	"github.com/ShajibEwuCse19/gopay-core/repo"
	"github.com/ShajibEwuCse19/gopay-core/rest"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
	services "github.com/ShajibEwuCse19/gopay-core/services/customer"
)

func Serve() {
	config, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error:\n%v\n", err)
		os.Exit(1)
	}

	logger.Init(config.ServiceName, config.LogLevel)
	slog.Info(
		"Starting server...",
		slog.String("version", config.Version),
		slog.Int("port", config.HTTP.Port),
	)

	dbCon, err := db.NewConnection(config)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer dbCon.Close()

	customerRepo := repo.NewCustomerRepo(dbCon)
	serviceCustomer := services.NewService(customerRepo)
	customerHandler := customer.NewHandler(config, serviceCustomer)
	server := rest.NewServer(customerHandler, config)
	server.StartServer()
}
