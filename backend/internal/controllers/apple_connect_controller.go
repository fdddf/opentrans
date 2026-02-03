package controllers

import (
	"fmt"

	"github.com/fdddf/opentrans/internal/context"
	"github.com/fdddf/opentrans/internal/services"
	"github.com/gofiber/fiber/v2"
)

// AppleConnectController handles Apple Connect requests
type AppleConnectController struct{}

// NewAppleConnectController creates a new AppleConnectController
func NewAppleConnectController() *AppleConnectController {
	return &AppleConnectController{}
}

// CreateAppleConnectConfigRequest represents the create Apple Connect config request
type CreateAppleConnectConfigRequest struct {
	IssuerID   string `json:"issuerId"`
	KeyID      string `json:"keyId"`
	PrivateKey string `json:"privateKey"`
	IsDefault  bool   `json:"isDefault"`
}

// UpdateAppleConnectConfigRequest represents the update Apple Connect config request
type UpdateAppleConnectConfigRequest struct {
	IssuerID   *string `json:"issuerId"`
	KeyID      *string `json:"keyId"`
	PrivateKey *string `json:"privateKey"`
	IsDefault  *bool   `json:"isDefault"`
}

// SyncAppsRequest represents the sync apps request
type SyncAppsRequest struct {
	ConfigID uint `json:"configId"`
}

// SyncAppLocalizationsRequest represents the sync app localizations request
type SyncAppLocalizationsRequest struct {
	ConfigID      uint     `json:"configId"`
	LanguageCodes []string `json:"languageCodes"`
	Direction     string   `json:"direction"`
	Strategy      string   `json:"strategy"`
}

// TestCredentialsRequest represents the test credentials request
type TestCredentialsRequest struct {
	IssuerID   string `json:"issuerId"`
	KeyID      string `json:"keyId"`
	PrivateKey string `json:"privateKey"`
}

// GetAppleConnectConfigs retrieves all Apple Connect configs for the authenticated user
func (ctrl *AppleConnectController) GetAppleConnectConfigs(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	configs, err := services.GetProviderConfigsByUserAndType(userID, "appleconnect")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Transform configs to a more intuitive format
	result := make([]fiber.Map, 0, len(configs))
	for _, config := range configs {
		result = append(result, fiber.Map{
			"id":          config.ID,
			"userId":      config.UserID,
			"providerType": config.ProviderType,
			"issuerId":    config.ConfigData["issuerID"],
			"keyId":       config.ConfigData["keyID"],
			"isDefault":   config.IsDefault,
			"createdAt":   config.CreatedAt,
			"updatedAt":   config.UpdatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// CreateAppleConnectConfig creates a new Apple Connect config
func (ctrl *AppleConnectController) CreateAppleConnectConfig(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req CreateAppleConnectConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.IssuerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "issuerId is required")
	}
	if req.KeyID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "keyId is required")
	}
	if req.PrivateKey == "" {
		return fiber.NewError(fiber.StatusBadRequest, "privateKey is required")
	}

	configData := map[string]interface{}{
		"issuerID":   req.IssuerID,
		"keyID":      req.KeyID,
		"privateKey": req.PrivateKey,
	}

	config, err := services.CreateProviderConfig(userID, "appleconnect", configData, req.IsDefault)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":        config.ID,
			"issuerId":  req.IssuerID,
			"keyId":     req.KeyID,
			"isDefault": config.IsDefault,
			"createdAt": config.CreatedAt,
			"updatedAt": config.UpdatedAt,
		},
	})
}

// GetAppleConnectConfig retrieves a single Apple Connect config by ID
func (ctrl *AppleConnectController) GetAppleConnectConfig(c *fiber.Ctx) error {
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
		return fiber.NewError(fiber.StatusNotFound, "Configuration not found")
	}

	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}

	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Not an Apple Connect configuration")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":        config.ID,
			"issuerId":  config.ConfigData["issuerID"],
			"keyId":     config.ConfigData["keyID"],
			"isDefault": config.IsDefault,
			"createdAt": config.CreatedAt,
			"updatedAt": config.UpdatedAt,
		},
	})
}

// UpdateAppleConnectConfig updates an existing Apple Connect config
func (ctrl *AppleConnectController) UpdateAppleConnectConfig(c *fiber.Ctx) error {
	configID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid config ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req UpdateAppleConnectConfigRequest
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

	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Not an Apple Connect configuration")
	}

	updates := make(map[string]interface{})
	newConfigData := make(map[string]interface{})

	// Copy existing config data
	for k, v := range config.ConfigData {
		newConfigData[k] = v
	}

	// Update fields if provided
	if req.IssuerID != nil {
		newConfigData["issuerID"] = *req.IssuerID
	}
	if req.KeyID != nil {
		newConfigData["keyID"] = *req.KeyID
	}
	if req.PrivateKey != nil {
		newConfigData["privateKey"] = *req.PrivateKey
	}

	if len(newConfigData) > 0 {
		updates["ConfigData"] = newConfigData
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
		"data": fiber.Map{
			"id":        updatedConfig.ID,
			"issuerId":  updatedConfig.ConfigData["issuerID"],
			"keyId":     updatedConfig.ConfigData["keyID"],
			"isDefault": updatedConfig.IsDefault,
			"createdAt": updatedConfig.CreatedAt,
			"updatedAt": updatedConfig.UpdatedAt,
		},
	})
}

// DeleteAppleConnectConfig deletes an Apple Connect config
func (ctrl *AppleConnectController) DeleteAppleConnectConfig(c *fiber.Ctx) error {
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

	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Not an Apple Connect configuration")
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

// TestAppleConnectConnection tests an Apple Connect connection
func (ctrl *AppleConnectController) TestAppleConnectConnection(c *fiber.Ctx) error {
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

	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Not an Apple Connect configuration")
	}

	// Test by attempting to sync apps (this will validate the credentials)
	_, err = services.SyncAppsFromAppleConnect(uint(userID), config.ConfigData["issuerID"].(string), config.ConfigData["keyID"].(string), config.ConfigData["privateKey"].(string))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Connection successful",
	})
}

// TestAppleConnectCredentials tests Apple Connect credentials
func (ctrl *AppleConnectController) TestAppleConnectCredentials(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req TestCredentialsRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.IssuerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "issuerId is required")
	}
	if req.KeyID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "keyId is required")
	}
	if req.PrivateKey == "" {
		return fiber.NewError(fiber.StatusBadRequest, "privateKey is required")
	}

	// Test by attempting to sync apps (this will validate the credentials)
	_, err := services.SyncAppsFromAppleConnect(uint(userID), req.IssuerID, req.KeyID, req.PrivateKey)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Connection successful",
	})
}

// SyncAppleApps syncs apps from Apple Connect
func (ctrl *AppleConnectController) SyncAppleApps(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req SyncAppsRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if req.ConfigID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "configId is required")
	}

	// Get the user's Apple Connect configuration
	config, err := services.GetProviderConfig(req.ConfigID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Provider configuration not found")
	}
	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}
	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Provider configuration is not Apple Connect")
	}

	_, ok = config.ConfigData["issuerID"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "issuerID is missing from configuration")
	}
	_, ok = config.ConfigData["keyID"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "keyID is missing from configuration")
	}
	_, ok = config.ConfigData["privateKey"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "privateKey is missing from configuration")
	}

	queueService := services.GetQueueService()
	if queueService.AppService == nil || queueService.AppLocalizationService == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Apple Connect service unavailable")
	}
	apps, err := services.NewAppleConnectService(services.AppleConnectServiceDeps{
		AppService:             queueService.AppService,
		AppLocalizationService: queueService.AppLocalizationService,
	}).SyncApps(userID,
		config.ConfigData["issuerID"].(string),
		config.ConfigData["keyID"].(string),
		"",
		config.ConfigData["privateKey"].(string),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, syncedApp := range apps {
		_, _ = services.BindAppProviderConfig(syncedApp.ID, config.ID, "appleconnect", false)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Synced %d apps from Apple Connect", len(apps)),
		"count":   len(apps),
	})
}

// SyncAppleAppLocalizations syncs localizations from Apple Connect
func (ctrl *AppleConnectController) SyncAppleAppLocalizations(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("appId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req SyncAppLocalizationsRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if req.ConfigID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "configId is required")
	}
	// Set defaults
	if req.Direction == "" {
		req.Direction = "pull"
	}
	if req.Strategy == "" {
		req.Strategy = "apple_first"
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	config, err := services.GetProviderConfig(req.ConfigID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Provider configuration not found")
	}
	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}
	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Provider configuration is not Apple Connect")
	}

	_, err = services.BindAppProviderConfig(uint(appID), config.ID, "appleconnect", true)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	binding, err := services.GetAppProviderConfig(uint(appID), config.ID)
	if err != nil || binding.ProviderConfigID != config.ID {
		return fiber.NewError(fiber.StatusForbidden, "Apple Connect configuration not bound to this app")
	}

	queueService := services.GetQueueService()
	if queueService.AppService == nil || queueService.AppLocalizationService == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Apple Connect service unavailable")
	}

	// Use the new sync method with direction and strategy
	syncService := services.NewAppleConnectService(services.AppleConnectServiceDeps{
		AppService:             queueService.AppService,
		AppLocalizationService: queueService.AppLocalizationService,
	})

	localizations, conflicts, err := syncService.SyncWithDirectionAndStrategy(
		uint(appID),
		services.SyncDirection(req.Direction),
		services.ConflictResolutionStrategy(req.Strategy),
		config.ConfigData["issuerID"].(string),
		config.ConfigData["keyID"].(string),
		"",
		config.ConfigData["privateKey"].(string),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   fmt.Sprintf("Synced %d localizations from Apple Connect", len(localizations)),
		"count":     len(localizations),
		"conflicts": conflicts,
	})
}