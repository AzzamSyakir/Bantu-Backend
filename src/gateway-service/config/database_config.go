package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	GatewayDB *PostgresDatabase
}

type PostgresDatabase struct {
	Connection *sql.DB
}

func NewGatewayDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		GatewayDB: NewGatewayDB(envConfig),
	}
	return databaseConfig
}

func NewGatewayDB(envConfig *EnvConfig) *PostgresDatabase {
	var url string
	if envConfig.GatewayDB.Password == "" {
		url = fmt.Sprintf(
			"postgresql://%s@%s:%s/%s",
			envConfig.GatewayDB.Gateway,
			envConfig.GatewayDB.Host,
			envConfig.GatewayDB.Port,
			envConfig.GatewayDB.Database,
		)
	} else {
		url = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			envConfig.GatewayDB.Gateway,
			envConfig.GatewayDB.Password,
			envConfig.GatewayDB.Host,
			envConfig.GatewayDB.Port,
			envConfig.GatewayDB.Database,
		)
	}

	connection, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	connection.SetConnMaxLifetime(300 * time.Second)
	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(10)
	userDB := &PostgresDatabase{
		Connection: connection,
	}
	return userDB
}
