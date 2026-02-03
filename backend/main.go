package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"

	"github.com/fdddf/opentrans/internal/auth"
	"github.com/fdddf/opentrans/internal/config"
	"github.com/fdddf/opentrans/internal/database"
	"github.com/fdddf/opentrans/internal/server"
	"github.com/fdddf/opentrans/internal/services"
	"github.com/fdddf/opentrans/internal/translator"
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

	Config      *config.FXConfig
	DB          *database.Database
	Lifecycle   fx.Lifecycle
	QueueService *services.QueueService
}

// InitializeAuth initializes the auth package with JWT secret and sets up lifecycle
func InitializeAuth(p InitializeAuthParams) {
	auth.SetJWTSecret([]byte(p.Config.JWT.Secret))

	// Add lifecycle hook to properly initialize auth if needed
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Initialize admin user if not exists
			if p.Config.Admin.CreateIfNotExists {
				if err := EnsureAdminUser(p.DB, p.Config.Admin); err != nil {
					fmt.Printf("Warning: failed to initialize admin user: %v\n", err)
				}
			}

			// Start queue processing in background
			if p.QueueService != nil {
				go p.QueueService.PollForNextJob(2 * time.Second)
				fmt.Println("Translation queue processor started")
			}

			return nil
		},
	})
}

// EnsureAdminUser ensures an admin user exists, creating one if necessary
func EnsureAdminUser(db *database.Database, adminConfig config.AdminConfig) error {
	// Check if any admin user already exists
	var existingAdmin database.User
	result := db.Where("role = ?", "admin").First(&existingAdmin)

	if result.Error == nil {
		// Admin already exists, no need to create
		fmt.Printf("Admin user already exists: %s\n", existingAdmin.Username)
		return nil
	}

	// No admin exists, create the default admin
	hashedPassword, err := auth.HashPassword(adminConfig.Password)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %v", err)
	}

	adminUser := &database.User{
		Username:       adminConfig.Username,
		Email:          adminConfig.Email,
		Password:       hashedPassword,
		IsActive:       true,
		IsActivated:    true,
		Role:           "admin",
		SubscriptionType: "premium",
		MaxApps:        100,
		MaxTranslations: 100000,
		CurrentUsage:   0,
		CurrentAppCount: 0,
	}

	if err := db.Create(adminUser).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	fmt.Printf("Created admin user: %s (%s)\n", adminUser.Username, adminUser.Email)
	fmt.Printf("Admin password: %s\n", adminConfig.Password)
	fmt.Println("Please change the admin password after first login!")

	return nil
}
