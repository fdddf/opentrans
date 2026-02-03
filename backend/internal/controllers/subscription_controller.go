package controllers

import (
	"fmt"

	"github.com/fdddf/opentrans/internal/context"
	"github.com/fdddf/opentrans/internal/services"
	"github.com/gofiber/fiber/v2"
)

// SubscriptionController handles subscription management requests
type SubscriptionController struct{}

// NewSubscriptionController creates a new SubscriptionController
func NewSubscriptionController() *SubscriptionController {
	return &SubscriptionController{}
}

// GetUserSubscription retrieves the subscription info for the authenticated user
func (ctrl *SubscriptionController) GetUserSubscription(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	user, err := services.GetUserSubscriptionInfo(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

// SubscriptionWebhook handles subscription webhooks (e.g., from Stripe)
func (ctrl *SubscriptionController) SubscriptionWebhook(c *fiber.Ctx) error {
	// For now, just log the webhook payload - in a real application, you'd process specific Stripe events
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	fmt.Printf("Received subscription webhook: %+v\n", req)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Webhook received",
	})
}

// GetUsage retrieves usage statistics for the authenticated user
func (ctrl *SubscriptionController) GetUsage(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	overLimit, usage, limit, err := services.CheckUserUsage(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"overLimit":  overLimit,
		"usage":      usage,
		"limit":      limit,
		"percentage": float64(usage) / float64(limit) * 100,
	})
}