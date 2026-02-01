package server

import (
	"context"
	"fmt"

	"github.com/fdddf/xcstrings-translator/internal/config"
	"github.com/fdddf/xcstrings-translator/internal/controllers"
	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// ServerModule is the FX module for server
var ServerModule = fx.Module("server",
	fx.Provide(NewFiberApp),
	fx.Invoke(StartServer),
)

// ServerParams holds the dependencies for the server
type ServerParams struct {
	fx.In

	Config    *config.FXConfig
	DB        *database.Database
	Lifecycle fx.Lifecycle
}

// ServerStateWithDB extends ServerState to include database dependency
type ServerStateWithDB struct {
	*controllers.ServerState
	DB *database.Database
}

// NewServerStateWithDB creates a new server state with database
func NewServerStateWithDB(db *database.Database) *ServerStateWithDB {
	return &ServerStateWithDB{
		ServerState: &controllers.ServerState{},
		DB:          db,
	}
}

// NewFiberApp creates a new Fiber application
func NewFiberApp(p ServerParams) (*fiber.App, error) {
	// Create the Fiber app using the NewApp function with database
	app, err := NewApp(p.DB)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// StartServer starts the Fiber server
func StartServer(lc fx.Lifecycle, app *fiber.App, cfg *config.FXConfig) {
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Printf("Starting server on %s\n", addr)
			go func() {
				if err := app.Listen(addr); err != nil {
					fmt.Printf("Server error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down server")
			return app.Shutdown()
		},
	})
}
