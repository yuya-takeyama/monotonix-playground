package config

import "os"

type Config struct {
	Port     string
	LogLevel string
	DBHost   string
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		DBHost:   getEnv("DB_HOST", "localhost"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}