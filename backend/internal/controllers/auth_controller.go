package controllers

import (
	"fmt"
	"os"

	"github.com/fdddf/opentrans/internal/auth"
	"github.com/fdddf/opentrans/internal/database"
	"github.com/fdddf/opentrans/internal/email"
	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication-related requests
type AuthController struct{}

// NewAuthController creates a new AuthController
func NewAuthController() *AuthController {
	return &AuthController{}
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register handles user registration
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Username, email, and password are required",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	user, err := auth.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	// Log registration activity
	ctrl.logActivity(c, user.ID, "user_registered", fmt.Sprintf("User %s registered", user.Username), "")

	// Send activation email
	go func() {
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:3000"
		}

		err := email.SendActivationEmail(
			user.Email,
			user.Username,
			user.ActivationCode,
			baseURL,
		)
		if err != nil {
			fmt.Printf("Failed to send activation email: %v\n", err)
		}
	}()

	response := struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		User    *database.User `json:"user"`
	}{
		Success: true,
		Message: "Registration successful. Please check your email for activation.",
		User:    user,
	}

	return c.JSON(response)
}

// Login handles user login
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	user, token, err := auth.LoginUser(req.Username, req.Password)
	if err != nil {
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: err.Error(),
		}
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	// Log login activity
	ctrl.logActivity(c, user.ID, "user_logged_in", fmt.Sprintf("User %s logged in", user.Username), "")

	response := struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		User    *database.User `json:"user"`
		Token   string         `json:"token"`
	}{
		Success: true,
		Message: "Login successful",
		User:    user,
		Token:   token,
	}

	return c.JSON(response)
}

// Logout handles user logout
func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "Logged out successfully",
	}

	return c.JSON(response)
}

// ActivateUser handles user account activation
func (ctrl *AuthController) ActivateUser(c *fiber.Ctx) error {
	activationCode := c.Params("code")

	err := auth.ActivateUser(activationCode)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account activated successfully. You can now log in.",
	})
}

// logActivity logs a user activity
func (ctrl *AuthController) logActivity(c *fiber.Ctx, userID uint, action, details, useragent string) {
	// For now, we'll just print the activity since we need to properly inject database
	fmt.Printf("Activity logged (no DB): userID=%d, action=%s, details=%s\n", userID, action, details)
}