package config

import (
	"log/slog"
	"os"
)

type Config struct {
	Port            string
	ShutdownTimeout string
	DBDsn           string
}

func Load() *Config {
	slog.Debug("Loading environment variables")
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		ShutdownTimeout: getEnv("SHUTDOWN_TIMEOUT", "10"),

		DBDsn: getEnv("DB_DSN", "postgres://postgres:Test1234@localhost:5432/usermanagementdb?sslmode=disable"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
