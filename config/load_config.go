package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var configuration *Config

// Load loads the local .env file (if present) and parses environment variables.
// It reports ALL missing required variables at once instead of failing on the first one.
func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg, err := loadFromLookup(os.Getenv) // now, lookup = os.Getenv
	if err != nil {
		return nil, err
	}
	configuration = cfg
	return cfg, nil
}

// LoadFromLookup loads configuration using a custom lookup function.
// Useful for isolated unit testing without depending on os.Getenv or local .env files.
func LoadFromLookup(lookup func(string) string) (*Config, error) {
	return loadFromLookup(lookup)
}

// loadFromLookup parses configuration from any lookup provider (e.g. os.Getenv or a test map).
func loadFromLookup(lookup func(string) string) (*Config, error) {
	var missing []string
	var invalid []string

	//1. validation annonyous function to check require string, requireInt, getEnv, intEnv, durationEnv
	require := func(key string) string {
		val := strings.TrimSpace(lookup(key))
		if val == "" {
			missing = append(missing, key)
			return ""
		}
		return val
	}

	requireInt := func(key string) int {
		valStr := require(key)
		if valStr == "" {
			return 0
		}
		val, err := strconv.Atoi(valStr)
		if err != nil {
			invalid = append(invalid, fmt.Sprintf("%s=%q (must be an integer)", key, valStr))
			return 0
		}
		return val
	}

	getEnv := func(key, def string) string {
		val := strings.TrimSpace(lookup(key))
		if val == "" {
			return def
		}
		return val
	}

	intEnv := func(key string, def int) int {
		raw := strings.TrimSpace(lookup(key))
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

	durationEnv := func(key string, def time.Duration) time.Duration {
		raw := strings.TrimSpace(lookup(key))
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

	// 2. Service metadata
	version := getEnv("VERSION", "1.0.0")
	serviceName := getEnv("SERVICE_NAME", "FINTECH")
	logLevel := getEnv("LOG_LEVEL", "info")

	// 3. Required variables from environment
	httpPort := requireInt("HTTP_PORT")
	jwtSecretKey := require("JWT_SECRET_KEY")
	dbHost := require("DB_HOST")
	dbPort := requireInt("DB_PORT")
	dbName := require("DB_NAME")
	dbUser := require("DB_USER")
	dbPassword := require("DB_PASSWORD")

	enableSllMode, err := strconv.ParseBool(getEnv("DB_ENABLE_SSL_MODE", "false"))
	if err != nil {
		invalid = append(invalid, fmt.Sprintf("DB_ENABLE_SSL_MODE=%q (must be true or false)", lookup("DB_ENABLE_SSL_MODE")))
	}

	// If any required variables are missing or invalid, report all of them at once!
	var errMsgs []string
	if len(missing) > 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("missing required environment variables: %s", strings.Join(missing, ", ")))
	}
	if len(invalid) > 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("invalid environment variables: %s", strings.Join(invalid, ", ")))
	}
	if len(errMsgs) > 0 {
		return nil, errors.New(strings.Join(errMsgs, "; "))
	}

	// 4. Set the global Config instance if there is missing or invalid environment variables
	cfg := &Config{
		Version:     version,
		ServiceName: serviceName,
		LogLevel:    logLevel,

		HTTP: HTTPConfig{
			Port:         httpPort,
			ReadTimeout:  durationEnv("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: durationEnv("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  durationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},

		DB: DBConfig{
			Host:           dbHost,
			Port:           dbPort,
			Name:           dbName,
			User:           dbUser,
			Password:       dbPassword,
			EnableSSlMode:  enableSllMode,
			MaxConns:       intEnv("DB_MAX_CONNS", 25),
			MinConns:       intEnv("DB_MIN_CONNS", 5),
			MigrationsPath: getEnv("DB_MIGRATIONS_PATH", "./migrations"),
		},

		JWT: JWTConfig{
			SecretKey:       jwtSecretKey,
			AccessTokenTTL:  durationEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: durationEnv("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     intEnv("REDIS_PORT", 6379),
			Password: lookup("REDIS_PASSWORD"),
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

	return cfg, nil
}

// GetConfig returns the singleton configuration instance.
// If not loaded, it attempts to load via Load() and logs fatal on error.
func GetConfig() *Config {
	if configuration == nil {
		cfg, err := Load()
		if err != nil {
			log.Fatalf("failed to load configuration: %v", err)
		}
		configuration = cfg
	}
	return configuration
}
