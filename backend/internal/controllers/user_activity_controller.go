package controllers

import (
	"github.com/fdddf/opentrans/internal/context"
	"github.com/fdddf/opentrans/internal/database"
	"github.com/gofiber/fiber/v2"
)

// UserActivityController handles user activity logging and retrieval
type UserActivityController struct{}

// NewUserActivityController creates a new user activity controller
func NewUserActivityController() *UserActivityController {
	return &UserActivityController{}
}

// GetActivities returns the recent activities for the authenticated user
func (c *UserActivityController) GetActivities(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID, ok := context.GetUserIDFromContext(ctx)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
		}

		var activities []database.UserActivity
		result := db.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&activities)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}

		return ctx.JSON(fiber.Map{
			"success":    true,
			"activities": activities,
		})
	}
}

// GetAdminActivities returns activities for all users (admin only)
func (c *UserActivityController) GetAdminActivities(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var activities []database.UserActivity
		result := db.Order("created_at DESC").Limit(100).Find(&activities)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}

		return ctx.JSON(fiber.Map{
			"success":    true,
			"activities": activities,
		})
	}
}