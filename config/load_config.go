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
	// Load local .env file if available (ignore error in container/production environments)
	_ = godotenv.Load()

	version := getEnv("VERSION", "1.0.0")
	serviceName := getEnv("SERVICE_NAME", "FINTECH")

	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatal("HTTP_PORT is required and must be an integer")
		os.Exit(1)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is required")
		os.Exit(1)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("DB_PORT is required and must be an integer")
		os.Exit(1)
	}
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	enableSllMode, err := strconv.ParseBool(getEnv("DB_ENABLE_SSL_MODE", "false"))
	if err != nil {
		log.Fatal("DB_ENABLE_SSL_MODE must be a boolean (true/false)")
		os.Exit(1)
	}

	configuration = &Config{
		Version:     version,
		ServiceName: serviceName,

		HTTP: HTTPConfig{
			Port:         httpPort,
			ReadTimeout:  durationEnv("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: durationEnv("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  durationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},

		DB: DBConfig{
			Host:          dbHost,
			Port:          dbPort,
			Name:          dbName,
			User:          dbUser,
			Password:      dbPassword,
			EnableSSlMode: enableSllMode,
			MaxConns:      intEnv("DB_MAX_CONNS", 25),
			MinConns:      intEnv("DB_MIN_CONNS", 5),
		},

		JWT: JWTConfig{
			SecretKey:       jwtSecretKey,
			AccessTokenTTL:  durationEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: durationEnv("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     intEnv("REDIS_PORT", 6379),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       intEnv("REDIS_DB", 0),
		},

		Timeouts: TimeoutConfig{
			CreateCustomer: durationEnv("CREATE_CUSTOMER_TIMEOUT", 5*time.Second),
			Shutdown:       durationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		},

		// Backward-compatibility aliases
		HttpPort:     httpPort,
		JWTSecretKey: jwtSecretKey,
	}
}

// getEnv reads a string environment variable with a default fallback.
func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

// intEnv reads an optional integer environment variable; it falls back to def if unset or invalid.
func intEnv(key string, def int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return def
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("%s=%q is not a valid integer, using default %d", key, raw, def)
		return def
	}
	return val
}

// durationEnv reads an optional duration like "5s" or "500ms"; it falls back to def if unset or invalid.
func durationEnv(key string, def time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return def
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		log.Printf("%s=%q is not a valid duration, using default %s", key, raw, def)
		return def
	}
	return d
}

// GetConfig returns the singleton configuration instance.
func GetConfig() *Config {
	if configuration == nil {
		loadConfig()
	}
	return configuration
}
