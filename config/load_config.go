package config

import (
	"log"
	"os"
	"strconv"
	"time"
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
		Timeouts: TimeoutConfig{
			CreateCustomer: durationEnv("CREATE_CUSTOMER_TIMEOUT", 5*time.Second),
			Shutdown:       durationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		},
	}
}

// durationEnv reads an optional duration like "5s" or "500ms"; it falls back to def if unset or invalid.
func durationEnv(key string, def time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return def
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		log.Printf("%s=%q is not a valid duration, using %s", key, raw, def)
		return def
	}
	return d
}

func GetConfig() *Config {
	if configuration == nil {
		loadConfig()
	}
	return configuration
}
