package controllers

import (
	"github.com/fdddf/xcstrings-translator/internal/context"
	"github.com/fdddf/xcstrings-translator/internal/services"
	"github.com/gofiber/fiber/v2"
)

// AppController handles app management requests
type AppController struct{}

// NewAppController creates a new AppController
func NewAppController() *AppController {
	return &AppController{}
}

// CreateAppRequest represents the create app request
type CreateAppRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	BundleID         string `json:"bundleId"`
	AppleID          string `json:"appleId"`
	PrimaryLocale    string `json:"primaryLocale"`
	ShortDescription string `json:"shortDescription"`
	LongDescription  string `json:"longDescription"`
	Keywords         string `json:"keywords"`
	SupportURL       string `json:"supportUrl"`
	MarketingURL     string `json:"marketingUrl"`
	PrivacyURL       string `json:"privacyUrl"`
	Version          string `json:"version"`
	AppCategory      string `json:"appCategory"`
}

// UpdateAppRequest represents the update app request
type UpdateAppRequest struct {
	Name             *string `json:"name"`
	Description      *string `json:"description"`
	BundleID         *string `json:"bundleId"`
	AppleID          *string `json:"appleId"`
	PrimaryLocale    *string `json:"primaryLocale"`
	ShortDescription *string `json:"shortDescription"`
	LongDescription  *string `json:"longDescription"`
	Keywords         *string `json:"keywords"`
	SupportURL       *string `json:"supportUrl"`
	MarketingURL     *string `json:"marketingUrl"`
	PrivacyURL       *string `json:"privacyUrl"`
	Version          *string `json:"version"`
	AppCategory      *string `json:"appCategory"`
	IsReadyForReview *bool   `json:"isReadyForReview"`
}

// AddUserToAppRequest represents the add user to app request
type AddUserToAppRequest struct {
	UserID uint   `json:"userId"`
	Role   string `json:"role"`
}

// UpdateUserAppRoleRequest represents the update user app role request
type UpdateUserAppRoleRequest struct {
	Role string `json:"role"`
}

// GetApps retrieves all apps for the authenticated user
func (ctrl *AppController) GetApps(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	apps, err := services.GetAppsByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"apps":    apps,
	})
}

// CreateApp creates a new app
func (ctrl *AppController) CreateApp(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req CreateAppRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	app, err := services.CreateApp(userID, req.Name, req.Description, req.BundleID, req.AppleID, req.PrimaryLocale)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"app":     app,
	})
}

// GetApp retrieves a single app by ID
func (ctrl *AppController) GetApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"app":     app,
	})
}

// UpdateApp updates an existing app
func (ctrl *AppController) UpdateApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req UpdateAppRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["Name"] = *req.Name
	}
	if req.Description != nil {
		updates["Description"] = *req.Description
	}
	if req.BundleID != nil {
		updates["BundleID"] = *req.BundleID
	}
	if req.AppleID != nil {
		updates["AppleID"] = *req.AppleID
	}
	if req.PrimaryLocale != nil {
		updates["PrimaryLocale"] = *req.PrimaryLocale
	}
	if req.ShortDescription != nil {
		updates["ShortDescription"] = *req.ShortDescription
	}
	if req.LongDescription != nil {
		updates["LongDescription"] = *req.LongDescription
	}
	if req.Keywords != nil {
		updates["Keywords"] = *req.Keywords
	}
	if req.SupportURL != nil {
		updates["SupportURL"] = *req.SupportURL
	}
	if req.MarketingURL != nil {
		updates["MarketingURL"] = *req.MarketingURL
	}
	if req.PrivacyURL != nil {
		updates["PrivacyURL"] = *req.PrivacyURL
	}
	if req.Version != nil {
		updates["Version"] = *req.Version
	}
	if req.AppCategory != nil {
		updates["AppCategory"] = *req.AppCategory
	}
	if req.IsReadyForReview != nil {
		updates["IsReadyForReview"] = *req.IsReadyForReview
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No fields to update")
	}

	err = services.UpdateApp(uint(appID), updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App updated successfully",
	})
}

// DeleteApp deletes an app
func (ctrl *AppController) DeleteApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	err = services.DeleteApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App deleted successfully",
	})
}

// GetAppUsers retrieves users for an app
func (ctrl *AppController) GetAppUsers(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	users, err := services.GetUsersForApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

// AddUserToApp adds a user to an app
func (ctrl *AppController) AddUserToApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner of the app
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can add users")
	}

	var req AddUserToAppRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate role
	validRoles := map[string]bool{
		"owner":  true,
		"admin":  true,
		"editor": true,
		"viewer": true,
	}
	if !validRoles[req.Role] {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role. Valid roles: owner, admin, editor, viewer")
	}

	err = services.AddUserToApp(uint(appID), req.UserID, req.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User added to app successfully",
	})
}

// UpdateUserAppRole updates a user's role in an app
func (ctrl *AppController) UpdateUserAppRole(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userIDParam, err := c.ParamsInt("userId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner of the app
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can update user roles")
	}

	var req UpdateUserAppRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate role
	validRoles := map[string]bool{
		"owner":  true,
		"admin":  true,
		"editor": true,
		"viewer": true,
	}
	if !validRoles[req.Role] {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role. Valid roles: owner, admin, editor, viewer")
	}

	err = services.UpdateUserAppRole(uint(appID), uint(userIDParam), req.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User role updated successfully",
	})
}

// RemoveUserFromApp removes a user from an app
func (ctrl *AppController) RemoveUserFromApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userIDParam, err := c.ParamsInt("userId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner of the app
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can remove users")
	}

	// Prevent removing the owner from their own app
	if uint(userIDParam) == app.UserID {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot remove the owner from the app")
	}

	err = services.RemoveUserFromApp(uint(appID), uint(userIDParam))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User removed from app successfully",
	})
}

// GetAppStats retrieves statistics for an app
func (ctrl *AppController) GetAppStats(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	stats, err := services.GetAppStats(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"stats":   stats,
	})
}

// BulkUpdateAppLocalizations bulk updates app localizations
func (ctrl *AppController) BulkUpdateAppLocalizations(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req struct {
		Updates []map[string]interface{} `json:"updates"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if len(req.Updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "updates are required")
	}

	err = services.BulkUpdateAppLocalizations(uint(appID), req.Updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App localizations updated successfully",
	})
}

// SyncAppToApple syncs an app to Apple Connect
func (ctrl *AppController) SyncAppToApple(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req struct {
		ConfigID    uint   `json:"configId"`
		LanguageCode string `json:"languageCode"` // Optional: if provided, only sync this language
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.ConfigID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "configId is required")
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

	err = services.SyncAppToAppleConnect(uint(appID),
		config.ConfigData["issuerID"].(string),
		config.ConfigData["keyID"].(string),
		config.ConfigData["privateKey"].(string),
		req.LanguageCode,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App synced to Apple Connect successfully",
	})
}