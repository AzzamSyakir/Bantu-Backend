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
	DBNeed   string
}
type RabbitMqEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	Queue    string
}

type EnvConfig struct {
	App       *AppEnv
	Db        *PostgresEnv
	RabbitMq  *RabbitMqEnv
	SecretKey string
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
			DBNeed:   os.Getenv("DB_NEED"),
		},
		RabbitMq: &RabbitMqEnv{
			Host:     os.Getenv("RABBITMQ_HOST"),
			Port:     os.Getenv("RABBITMQ_PORT"),
			User:     os.Getenv("RABBITMQ_USER"),
			Password: os.Getenv("RABBITMQ_PASSWORD"),
			Queue:    os.Getenv("RABBITMQ_QUEUE_NAMES"),
		},
		SecretKey: os.Getenv("SECRET_KEY"),
	}
	return envConfig
}
