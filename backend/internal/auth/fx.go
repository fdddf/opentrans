package auth

import (
	"go.uber.org/fx"

	"github.com/fdddf/opentrans/internal/database"
)

// Module is the FX module for auth
var Module = fx.Module("auth",
	fx.Provide(NewAuthFromFX),
)

// AuthParams holds the dependencies for Auth
type AuthParams struct {
	fx.In

	DB *database.Database
}

// NewAuthFromFX creates a new Auth instance with dependency injection
func NewAuthFromFX(p AuthParams) *Auth {
	auth := &Auth{
		DB: p.DB,
	}
	// Set the global auth instance for backward compatibility
	SetAuthInstance(auth)
	return auth
}
