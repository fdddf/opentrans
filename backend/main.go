package main

import (
	"context"

	"go.uber.org/fx"

	"github.com/fdddf/xcstrings-translator/internal/auth"
	"github.com/fdddf/xcstrings-translator/internal/config"
	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/server"
	"github.com/fdddf/xcstrings-translator/internal/services"
	"github.com/fdddf/xcstrings-translator/internal/translator"
)

func main() {
	app := fx.New(
		// Configuration module
		config.Module,

		// Database module
		database.Module,

		// Auth module
		auth.Module,

		// Services module
		services.Module,

		// Translator module
		translator.TranslatorModule,

		// Server module
		server.ServerModule,

		// Application startup
		fx.Invoke(InitializeAuth),
	)

	app.Run()
}

// InitializeAuthParams holds the dependencies for InitializeAuth
type InitializeAuthParams struct {
	fx.In

	Config    *config.FXConfig
	DB        *database.Database
	Lifecycle fx.Lifecycle
}

// InitializeAuth initializes the auth package with JWT secret and sets up lifecycle
func InitializeAuth(p InitializeAuthParams) {
	auth.SetJWTSecret([]byte(p.Config.JWT.Secret))

	// Add lifecycle hook to properly initialize auth if needed
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Any auth initialization needed
			return nil
		},
	})
}
