package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// Module is the FX module for config
var Module = fx.Module("config",
	fx.Provide(NewConfig),
)

// NewConfig creates a new configuration instance
func NewConfig() (*FXConfig, error) {
	// Try multiple paths for .env file
	paths := []string{"../.env", "../../.env", ".env"}
	loaded := false

	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			fmt.Printf("Loaded .env file from: %s\n", path)
			loaded = true
			break
		}
	}

	if !loaded {
		// If .env file doesn't exist, that's OK - we'll use defaults or environment variables
		fmt.Println("No .env file found, using environment variables")
	}

	cfg := &FXConfig{
		Database: GetDatabaseConfig(),
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "3000"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret-key-for-development"),
		},
		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromAddress:  getEnv("EMAIL_FROM", "noreply@example.com"),
			BaseURL:      getEnv("BASE_URL", "http://localhost:3000"),
		},
	}

	return cfg, nil
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// StripeConfig holds Stripe configuration
type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromAddress  string
	BaseURL      string
}

// FXConfig is the main configuration struct for FX
type FXConfig struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Stripe   StripeConfig
	Email    EmailConfig
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
