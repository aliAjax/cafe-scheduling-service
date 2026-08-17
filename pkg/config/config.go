package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	ServerPort    string
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() Config {
	return Config{
		DBHost:        getenv("DB_HOST", "127.0.0.1"),
		DBPort:        getenv("DB_PORT", "3306"),
		DBUser:        getenv("DB_USER", "cafe"),
		DBPassword:    getenv("DB_PASSWORD", "cafe_password"),
		DBName:        getenv("DB_NAME", "cafe_scheduling"),
		ServerPort:    getenv("SERVER_PORT", "8080"),
		JWTSecret:     getenv("JWT_SECRET", "change-me-in-production"),
		JWTExpiration: time.Duration(getenvInt("JWT_EXPIRES_HOURS", 24)) * time.Hour,
	}
}

func (c Config) DSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName +
		"?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true"
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
