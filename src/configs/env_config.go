package configs

import (
	"os"
)

type AppEnv struct {
	AppHost string
	AppPort string
}

type PostgresEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
type RabbitMqEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	Queue    string
}

type EnvConfig struct {
	App *AppEnv
	Db  *PostgresEnv
}

func NewEnvConfig() *EnvConfig {
	envConfig := &EnvConfig{
		App: &AppEnv{
			AppHost: os.Getenv("GATEWAY_APP_HOST"),
			AppPort: os.Getenv("GATEWAY_APP_PORT"),
		},
		Db: &PostgresEnv{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},
	}
	return envConfig
}
