package controllers

import (
	"fmt"

	appcontext "github.com/fdddf/opentrans/internal/context"
	"github.com/fdddf/opentrans/internal/auth"
	"github.com/fdddf/opentrans/internal/database"
	"github.com/gofiber/fiber/v2"
)

// UserController handles user-related requests
type UserController struct{}

// NewUserController creates a new UserController
func NewUserController() *UserController {
	return &UserController{}
}

// UpdateUserRequest represents the update user request
type UpdateUserRequest struct {
	Email *string `json:"email"`
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// GetUser retrieves a user by ID
func (ctrl *UserController) GetUser(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDParam := c.Params("id")
		var userID uint

		_, err := fmt.Sscanf(userIDParam, "%d", &userID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
		}

		authInstance := auth.NewAuth(db)
		user, err := authInstance.GetUserByID(userID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"user":    user,
		})
	}
}

// UpdateUser updates the authenticated user's profile
func (ctrl *UserController) UpdateUser(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := appcontext.GetUserIDFromContext(c)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
		}

		var req UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		authInstance := auth.NewAuth(db)
		user, err := authInstance.GetUserByID(userID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		// Update email if provided
		if req.Email != nil {
			if !auth.ValidateEmail(*req.Email) {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
			}

			// Check if email is already taken by another user
			var existingUser database.User
			result := db.Where("email = ? AND id != ?", *req.Email, userID).First(&existingUser)
			if result.Error == nil {
				return fiber.NewError(fiber.StatusConflict, "Email already in use")
			}

			user.Email = *req.Email
		}

		result := db.Save(user)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
		}

		// Clear password for security
		user.Password = ""

		return c.JSON(fiber.Map{
			"success": true,
			"user":    user,
		})
	}
}

// UpdateCurrentUser updates the authenticated user's own profile (for use in Profile page)
func (ctrl *UserController) UpdateCurrentUser(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := appcontext.GetUserIDFromContext(c)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
		}

		var req UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		authInstance := auth.NewAuth(db)
		user, err := authInstance.GetUserByID(userID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		// Update email if provided
		if req.Email != nil {
			if !auth.ValidateEmail(*req.Email) {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
			}

			// Check if email is already taken by another user
			var existingUser database.User
			result := db.Where("email = ? AND id != ?", *req.Email, userID).First(&existingUser)
			if result.Error == nil {
				return fiber.NewError(fiber.StatusConflict, "Email already in use")
			}

			user.Email = *req.Email
		}

		result := db.Save(user)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
		}

		// Clear password for security
		user.Password = ""

		return c.JSON(fiber.Map{
			"success": true,
			"user":    user,
		})
	}
}

// ChangePassword changes the authenticated user's password
func (ctrl *UserController) ChangePassword(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := appcontext.GetUserIDFromContext(c)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
		}

		var req ChangePasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate new password
		if len(req.NewPassword) < 6 {
			return fiber.NewError(fiber.StatusBadRequest, "New password must be at least 6 characters")
		}

		// Get user
		var user database.User
		result := db.First(&user, userID)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		// Verify current password
		if err := auth.CheckPassword(req.CurrentPassword, user.Password); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Current password is incorrect")
		}

		// Hash new password
		hashedPassword, err := auth.HashPassword(req.NewPassword)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
		}

		// Update password
		user.Password = hashedPassword
		result = db.Save(&user)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update password")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Password changed successfully",
		})
	}
}