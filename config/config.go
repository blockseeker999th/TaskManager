package config

import "os"

type PostgreSQLConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
	JWTSecret string
}

var Envs = initConfig()

func initConfig() PostgreSQLConfig {
	return PostgreSQLConfig{
		Host:      getEnv("PG_HOST", "localhost"),
		Port:      getEnv("PG_PORT", "5432"),
		User:      getEnv("PG_USER", "aleks"),
		Password:  getEnv("PG_PASSWORD", "tyu899uyt"),
		DBName:    getEnv("PG_DBNAME", "taskmanager"),
		JWTSecret: getEnv("JWT_SECRET", "rand"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
