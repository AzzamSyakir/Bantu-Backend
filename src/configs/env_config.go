package configs

import (
	"fmt"
	"os"
	"strconv"
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
type RedisEnv struct {
	Addr     string
	Password string
	DB       int
}

type EnvConfig struct {
	App             *AppEnv
	Db              *PostgresEnv
	RabbitMq        *RabbitMqEnv
	Redis           *RedisEnv
	SecretKey       string
	XenditSecretKey string
}

func NewEnvConfig() *EnvConfig {
	queueNamesEnv := os.Getenv("RABBITMQ_QUEUE_NAMES")
	queues := strings.Split(queueNamesEnv, ",")
	cleanedQueues := make([]string, len(queues))
	addressRedis := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	for i, q := range queues {
		cleanedQueues[i] = strings.TrimSpace(q)
	}
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
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
		Redis: &RedisEnv{
			Addr:     addressRedis,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		SecretKey:       os.Getenv("SECRET_KEY"),
		XenditSecretKey: os.Getenv("XENDIT_KEY"),
	}
	return envConfig
}
