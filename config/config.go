package config

import "time"

// Config holds all typed configuration sections for the application.
type Config struct {
	Version      string
	ServiceName  string
	LogLevel     string
	HTTP         HTTPConfig
	DB           DBConfig
	JWT          JWTConfig
	Redis        RedisConfig
	Timeouts     TimeoutConfig

	// Backward-compatibility aliases
	HttpPort     int
	JWTSecretKey string
}

// HTTPConfig holds HTTP server network and timeout settings.
type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DBConfig holds PostgreSQL connection settings and pool limits.
type DBConfig struct {
	Host          string
	Port          int
	Name          string
	User          string
	Password      string
	EnableSSlMode bool
	MaxConns      int
	MinConns      int
}

// JWTConfig holds authentication token secret and expiration durations.
type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// RedisConfig holds cache and distributed lock connection settings.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// TimeoutConfig holds application-level operation budgets.
type TimeoutConfig struct {
	CreateCustomer time.Duration // total time budget of one register request
	Shutdown       time.Duration // max wait for in-flight requests on shutdown
}
