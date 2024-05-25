package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

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
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env variables", err)
	}

	return PostgreSQLConfig{
		Host:      "db",
		Port:      getEnv("PG_PORT", "5432"),
		User:      getEnv("PG_USER", "default"),
		Password:  getEnv("PG_PASSWORD", "default"),
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
