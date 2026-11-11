package config

import "time"

type Config struct {
	Version      string
	ServiceName  string
	HttpPort     int
	JWTSecretKey string
	DB           DBConfig
	Timeouts     TimeoutConfig
}

type DBConfig struct {
	Host          string
	Port          int
	Name          string
	User          string
	Password      string
	EnableSSlMode bool
}

type TimeoutConfig struct {
	CreateCustomer time.Duration // total time budget of one register request
	Shutdown       time.Duration // max wait for in-flight requests on shutdown
}
