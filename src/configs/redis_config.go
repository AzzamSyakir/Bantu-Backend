package configs

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Redis *RedisConnection
}

type RedisConnection struct {
	Connection *redis.Client
}

func NewRedisConfig(env *EnvConfig) *RedisConfig {
	redisConfig := &RedisConfig{
		Redis: NewRedisConnection(env),
	}
	return redisConfig
}

func NewRedisConnection(env *EnvConfig) *RedisConnection {
	options := &redis.Options{
		Addr:     env.Redis.Addr,
		Password: env.Redis.Password,
		DB:       env.Redis.DB,
	}
	client := redis.NewClient(options)
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return nil
	}
	connection := &RedisConnection{
		Connection: client,
	}
	return connection
}
