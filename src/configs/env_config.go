package configs

import (
	"os"
	"strings"
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
	Queues   []string
}

type EnvConfig struct {
	App      *AppEnv
	Db       *PostgresEnv
	RabbitMq *RabbitMqEnv
}

func NewEnvConfig() *EnvConfig {
	queueNamesEnv := os.Getenv("RABBITMQ_QUEUE_NAMES")
	queues := strings.Split(queueNamesEnv, ",")
	cleanedQueues := make([]string, len(queues))
	for i, q := range queues {
		cleanedQueues[i] = strings.TrimSpace(q)
	}
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
			Queues:   cleanedQueues,
		},
	}
	return envConfig
}
