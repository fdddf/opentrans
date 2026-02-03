package context

import (
	"strings"

	"github.com/fdddf/opentrans/internal/auth"
	"github.com/fdddf/opentrans/internal/database"
	"github.com/fdddf/opentrans/internal/services"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware authenticates requests using JWT
func AuthMiddleware(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	// Check if it's a Bearer token
	tokenString := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	// Parse and validate the token
	claims, err := auth.ParseJWT(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Fetch user from database
	user, err := auth.GetUserByID(claims.UserID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Add user info to context
	c.Locals("user", user)
	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("role", user.Role)

	return c.Next()
}

// OptionalAuthMiddleware allows requests to proceed even if not authenticated
func OptionalAuthMiddleware(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// No token provided, continue without authentication
		c.Locals("user", nil)
		c.Locals("userID", uint(0))
		c.Locals("username", "")
		c.Locals("role", "")
		return c.Next()
	}

	// Check if it's a Bearer token
	tokenString := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// Invalid format, continue without authentication
		c.Locals("user", nil)
		c.Locals("userID", uint(0))
		c.Locals("username", "")
		c.Locals("role", "")
		return c.Next()
	}

	// Parse and validate the token
	claims, err := auth.ParseJWT(tokenString)
	if err != nil {
		// Invalid token, continue without authentication
		c.Locals("user", nil)
		c.Locals("userID", uint(0))
		c.Locals("username", "")
		c.Locals("role", "")
		return c.Next()
	}

	// Fetch user from database
	user, err := auth.GetUserByID(claims.UserID)
	if err != nil {
		// User not found, continue without authentication
		c.Locals("user", nil)
		c.Locals("userID", uint(0))
		c.Locals("username", "")
		c.Locals("role", "")
		return c.Next()
	}

	// Add user info to context
	c.Locals("user", user)
	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("role", user.Role)

	return c.Next()
}

// GetUserFromContext retrieves the authenticated user from the context
func GetUserFromContext(c *fiber.Ctx) (*database.User, bool) {
	user, ok := c.Locals("user").(*database.User)
	return user, ok
}

// GetUserIDFromContext retrieves the authenticated user ID from the context
func GetUserIDFromContext(c *fiber.Ctx) (uint, bool) {
	userID, ok := c.Locals("userID").(uint)
	return userID, ok
}

// GetUserRoleFromContext retrieves the authenticated user role from the context
func GetUserRoleFromContext(c *fiber.Ctx) (string, bool) {
	role, ok := c.Locals("role").(string)
	return role, ok
}

// AdminOnly middleware ensures only admin role can access the route
func AdminOnly(c *fiber.Ctx) error {
	role, ok := GetUserRoleFromContext(c)
	if !ok || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "admin role required",
		})
	}
	return c.Next()
}

// SubscriptionRequired ensures user has active subscription/quota
func SubscriptionRequired(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	if err := services.RequireActiveSubscriptionAndQuota(userID); err != nil {
		return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Next()
}