package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/infra/db"
	"github.com/ShajibEwuCse19/gopay-core/repo"
	"github.com/ShajibEwuCse19/gopay-core/rest"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
	services "github.com/ShajibEwuCse19/gopay-core/services/customer"
	"github.com/ShajibEwuCse19/gopay-core/utils"
)

func Serve() {
	config := config.GetConfig()
	fmt.Println("Version:", config.Version)
	fmt.Println("Service Name:", config.ServiceName)
	fmt.Println("Http Port:", config.HttpPort)
	utils.FreePort(":" + strconv.Itoa(config.HttpPort))

	dbCon, err := db.NewConnection(config)
	if err != nil {
		fmt.Println("db connection err:", err)
		os.Exit(1)
	}
	err = db.MigrateDB(dbCon, "./migrations") 
	if err != nil {
		fmt.Println("db migration err:", err)
		os.Exit(1)
	}

	customerRepo := repo.NewCustomerRepo(dbCon)
	serviceCustomer := services.NewService(customerRepo)
	customerHandler := customer.NewHandler(serviceCustomer)
	server := rest.NewServer(customerHandler, config)
	server.StartServer()
}
