//go:build !windows || !gui || !cgo
// +build !windows !gui !cgo

package main

import (
	"os"

	"github.com/fdddf/xcstrings-translator/cmd"
	"github.com/fdddf/xcstrings-translator/internal/auth"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// If .env file doesn't exist, that's OK - we'll use defaults or environment variables
	}

	// Set a default JWT secret if not provided
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "default-secret-key-for-development")
	}

	// Initialize auth package with the JWT secret
	auth.JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	cmd.Execute()
}
