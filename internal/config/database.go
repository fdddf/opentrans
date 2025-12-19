package config

import (
	"os"
)

// DatabaseConfig holds database configuration settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

// GetDatabaseConfig returns the database configuration from environment variables or defaults
func GetDatabaseConfig() DatabaseConfig {
	config := DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "xcstrings"),
		Password: getEnvOrDefault("DB_PASSWORD", "xcstrings"),
		DBName:   getEnvOrDefault("DB_NAME", "xcstrings"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
	}

	// If DATABASE_URL is provided, use it instead of individual settings
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		config.URL = databaseURL
	} else {
		// Construct the connection string from individual settings
		config.URL = "host=" + config.Host + " user=" + config.User + " password=" + config.Password + " dbname=" + config.DBName + " port=" + config.Port + " sslmode=" + config.SSLMode
	}

	return config
}

// getEnvOrDefault returns the environment variable value or the default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
