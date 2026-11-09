package db

import (
	"fmt"
	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/jmoiron/sqlx"
)

func NewConnection(config *config.Config) (*sqlx.DB, error) {
	dbSource := getConnectionString(config)
	dbConnection, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return dbConnection, nil
}

func getConnectionString(dbConfig *config.Config) string {
	sslMode := "disable"

	if dbConfig.DB.EnableSSlMode {
		sslMode = "require"
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		dbConfig.DB.User,
		dbConfig.DB.Password,
		dbConfig.DB.Host,
		dbConfig.DB.Port,
		dbConfig.DB.Name,
		sslMode,
	)
}