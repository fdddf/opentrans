package auth

import (
	"go.uber.org/fx"

	"github.com/fdddf/xcstrings-translator/internal/database"
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
	return &Auth{
		DB: p.DB,
	}
}
