package config

import (
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

var configuration *Config

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	version := os.Getenv("VERSION")
	serviceName := os.Getenv("SERVICE_NAME")
	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatal("HTTP_PORT is not set")
		os.Exit(1)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
		os.Exit(1)
	}

	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("DB_PORT is not set")
		os.Exit(1)
	}
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	enableSllMode, err := strconv.ParseBool(os.Getenv("DB_ENABLE_SSL_MODE"))
	if err != nil {
		log.Fatal("DB_ENABLE_SSL_MODE is not set")
		os.Exit(1)
	}

	configuration = &Config{
		Version:      version,
		ServiceName:  serviceName,
		HttpPort:     httpPort,
		JWTSecretKey: jwtSecretKey,
		DB: DBConfig{
			Host:          host,
			Port:          port,
			Name:          name,
			User:          user,
			Password:      password,
			EnableSSlMode: enableSllMode,
		},
	}
}

func GetConfig() *Config {
	if configuration == nil {
		loadConfig()
	}
	return configuration
}