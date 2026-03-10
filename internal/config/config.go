package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
}

func MustLoad() *Config {
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbUser := requireEnv("POSTGRES_USER")
	dbPassword := requireEnv("POSTGRES_PASSWORD")
	dbName := getEnv("POSTGRES_DB", "todo_db")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: databaseURL,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}
