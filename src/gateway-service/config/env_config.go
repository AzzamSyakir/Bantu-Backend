package config

import (
	"os"
)

type AppEnv struct {
	Host            string
	UserHost        string
	ProductHost     string
	OrderHost       string
	GatewayHttpPort string
	OrderPort       string
	UserPort        string
	ProductPort     string
	GatewayGrpcPort string
}

type PostgresEnv struct {
	Host     string
	Port     string
	Gateway  string
	Password string
	Database string
}

type EnvConfig struct {
	App       *AppEnv
	GatewayDB *PostgresEnv
}

func NewEnvConfig() *EnvConfig {
	envConfig := &EnvConfig{
		App: &AppEnv{
			Host:            os.Getenv("GATEWAY_HOST"),
			UserHost:        os.Getenv("USER_HOST"),
			ProductHost:     os.Getenv("TRANSACTION_HOST"),
			OrderHost:       os.Getenv("REVIEW_HOST"),
			GatewayHttpPort: os.Getenv("GATEWAY_PORT"),
			GatewayGrpcPort: os.Getenv("GATEWAY_GRPC_PORT"),
			OrderPort:       os.Getenv("REVIEW_PORT"),
			UserPort:        os.Getenv("USER_PORT"),
			ProductPort:     os.Getenv("TRANSACTION_PORT"),
		},
		GatewayDB: &PostgresEnv{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_GATEWAY_PORT"),
			Gateway:  os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: "auth_db",
		},
	}
	return envConfig
}
