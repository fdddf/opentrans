package controllers

import (
	"github.com/fdddf/opentrans/internal/context"
	"github.com/fdddf/opentrans/internal/database"
	"github.com/fdddf/opentrans/internal/services"
	"github.com/gofiber/fiber/v2"
)

// ProviderController handles provider configuration requests
type ProviderController struct{}

// NewProviderController creates a new ProviderController
func NewProviderController() *ProviderController {
	return &ProviderController{}
}

// CreateProviderConfigRequest represents the create provider config request
type CreateProviderConfigRequest struct {
	ProviderType string                 `json:"providerType"`
	ConfigData   map[string]interface{} `json:"configData"`
	IsDefault    bool                   `json:"isDefault"`
}

// UpdateProviderConfigRequest represents the update provider config request
type UpdateProviderConfigRequest struct {
	ProviderType *string                `json:"providerType"`
	ConfigData   map[string]interface{} `json:"configData"`
	IsDefault    *bool                  `json:"isDefault"`
}

// GetProviderConfigs retrieves all provider configs for the authenticated user
func (ctrl *ProviderController) GetProviderConfigs(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	configs, err := services.GetProviderConfigsByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Sanitize configs before returning
	sanitizedConfigs := make([]database.ProviderConfig, 0, len(configs))
	for _, config := range configs {
		sanitizedConfig := services.SanitizeProviderConfig(&config)
		if sanitizedConfig != nil {
			sanitizedConfigs = append(sanitizedConfigs, *sanitizedConfig)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"configs": sanitizedConfigs,
	})
}

// CreateProviderConfig creates a new provider config
func (ctrl *ProviderController) CreateProviderConfig(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req CreateProviderConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate the config data
	err := services.ValidateProviderConfigData(req.ProviderType, req.ConfigData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	config, err := services.CreateProviderConfig(userID, req.ProviderType, req.ConfigData, req.IsDefault)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"config":  services.SanitizeProviderConfig(config),
	})
}

// GetProviderConfig retrieves a single provider config by ID
func (ctrl *ProviderController) GetProviderConfig(c *fiber.Ctx) error {
	configID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid config ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	config, err := services.GetProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"config":  services.SanitizeProviderConfig(config),
	})
}

// UpdateProviderConfig updates an existing provider config
func (ctrl *ProviderController) UpdateProviderConfig(c *fiber.Ctx) error {
	configID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid config ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req UpdateProviderConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing config to check ownership
	config, err := services.GetProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Configuration not found")
	}

	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}

	updates := make(map[string]interface{})
	if req.ProviderType != nil {
		updates["ProviderType"] = *req.ProviderType
	}
	if req.ConfigData != nil {
		// Validate the config data using the provider type from request or existing config
		providerType := config.ProviderType
		if req.ProviderType != nil {
			providerType = *req.ProviderType
		}
		err := services.ValidateProviderConfigData(providerType, req.ConfigData)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		updates["ConfigData"] = req.ConfigData
	}
	if req.IsDefault != nil {
		updates["IsDefault"] = *req.IsDefault
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No fields to update")
	}

	err = services.UpdateProviderConfig(uint(configID), updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Get the updated config to return
	updatedConfig, err := services.GetProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"config":  services.SanitizeProviderConfig(updatedConfig),
	})
}

// DeleteProviderConfig deletes a provider config
func (ctrl *ProviderController) DeleteProviderConfig(c *fiber.Ctx) error {
	configID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid config ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get existing config to check ownership
	config, err := services.GetProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Configuration not found")
	}

	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}

	err = services.DeleteProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Configuration deleted successfully",
	})
}

// GetDefaultProviderConfig retrieves the default provider config for a given type
func (ctrl *ProviderController) GetDefaultProviderConfig(c *fiber.Ctx) error {
	providerType := c.Params("type")

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	config, err := services.GetDefaultProviderConfig(userID, providerType)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"config":  services.SanitizeProviderConfig(config),
	})
}