package config

type Config struct {
	Version      string
	ServiceName  string
	HttpPort     int
	JWTSecretKey string
	DB           DBConfig
}

type DBConfig struct {
	Host          string
	Port          int
	Name          string
	User          string
	Password      string
	EnableSSlMode bool
}