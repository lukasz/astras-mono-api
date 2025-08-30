// Package database provides database configuration and utilities.
package database

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config represents database configuration settings
type Config struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Database     string        `json:"database"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	SSLMode      string        `json:"ssl_mode"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := &Config{
		Host:         getEnvOrDefault("DB_HOST", "localhost"),
		Port:         getEnvIntOrDefault("DB_PORT", 5432),
		Database:     getEnvOrDefault("DB_NAME", "astras"),
		Username:     getEnvOrDefault("DB_USER", "postgres"),
		Password:     getEnvOrDefault("DB_PASSWORD", ""),
		SSLMode:      getEnvOrDefault("DB_SSL_MODE", "disable"),
		MaxOpenConns: getEnvIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: getEnvIntOrDefault("DB_MAX_IDLE_CONNS", 5),
		MaxLifetime:  getEnvDurationOrDefault("DB_MAX_LIFETIME", 5*time.Minute),
	}

	return config
}

// DSN returns the PostgreSQL data source name
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

// DatabaseURL returns the database URL format (for migrations)
func (c *Config) DatabaseURL() string {
	if c.Password != "" {
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.Username, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
		)
	}
	return fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=%s",
		c.Username, c.Host, c.Port, c.Database, c.SSLMode,
	)
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault returns environment variable as int or default
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDurationOrDefault returns environment variable as duration or default
func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}