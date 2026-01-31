package controllers

import (
	"github.com/fdddf/xcstrings-translator/internal/context"
	"github.com/fdddf/xcstrings-translator/internal/services"
	"github.com/gofiber/fiber/v2"
)

// AppLocalizationController handles app localization requests
type AppLocalizationController struct{}

// NewAppLocalizationController creates a new AppLocalizationController
func NewAppLocalizationController() *AppLocalizationController {
	return &AppLocalizationController{}
}

// CreateAppLocalizationRequest represents the create app localization request
type CreateAppLocalizationRequest struct {
	LanguageCode        string `json:"languageCode"`
	Name                string `json:"name"`
	Subtitle            string `json:"subtitle"`
	PrivacyURL          string `json:"privacyUrl"`
	MarketingURL        string `json:"marketingUrl"`
	SupportURL          string `json:"supportUrl"`
	DownloadDescription string `json:"downloadDescription"`
	ShortDescription    string `json:"shortDescription"`
	LongDescription     string `json:"longDescription"`
	Keywords            string `json:"keywords"`
	ReleaseNotes        string `json:"releaseNotes"`
}

// UpdateAppLocalizationRequest represents the update app localization request
type UpdateAppLocalizationRequest struct {
	Name                string `json:"name"`
	Subtitle            string `json:"subtitle"`
	PrivacyURL          string `json:"privacyUrl"`
	MarketingURL        string `json:"marketingUrl"`
	SupportURL          string `json:"supportUrl"`
	DownloadDescription string `json:"downloadDescription"`
	ShortDescription    string `json:"shortDescription"`
	LongDescription     string `json:"longDescription"`
	Keywords            string `json:"keywords"`
	ReleaseNotes        string `json:"releaseNotes"`
}

// AddAppLanguageRequest represents the add app language request
type AddAppLanguageRequest struct {
	Language string `json:"language"`
}

// GetAppLocalizations retrieves all localizations for an app
func (ctrl *AppLocalizationController) GetAppLocalizations(c *fiber.Ctx) error {
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

	localizations, err := services.GetAppLocalizations(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"localizations": localizations,
	})
}

// CreateAppLocalization creates a new app localization
func (ctrl *AppLocalizationController) CreateAppLocalization(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	role, _ := context.GetUserRoleFromContext(c)
	if role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Only admin can create localizations")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req CreateAppLocalizationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := services.ValidateLanguages([]string{req.LanguageCode}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	localization, err := services.CreateAppLocalization(
		uint(appID),
		req.LanguageCode,
		req.Name,
		req.Subtitle,
		req.PrivacyURL,
		req.MarketingURL,
		req.SupportURL,
		req.DownloadDescription,
		req.ShortDescription,
		req.LongDescription,
		req.Keywords,
		req.ReleaseNotes,
		"", // PromotionalText
		"", // WhatToTest
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	_ = services.UpdateAppLocalization(uint(appID), req.LanguageCode, map[string]interface{}{
		"SyncStatus": "pending",
		"Source":     "local",
	})

	return c.JSON(fiber.Map{
		"success":      true,
		"localization": localization,
	})
}

// GetAppLocalization retrieves a single app localization
func (ctrl *AppLocalizationController) GetAppLocalization(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	languageCode := c.Params("language")

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

	localization, err := services.GetAppLocalization(uint(appID), languageCode)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"localization": localization,
	})
}

// UpdateAppLocalization updates an existing app localization
func (ctrl *AppLocalizationController) UpdateAppLocalization(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	languageCode := c.Params("language")

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	role, _ := context.GetUserRoleFromContext(c)
	if role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Only admin can update localizations")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req UpdateAppLocalizationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	updates := map[string]interface{}{
		"Name":                req.Name,
		"Subtitle":            req.Subtitle,
		"PrivacyURL":          req.PrivacyURL,
		"MarketingURL":        req.MarketingURL,
		"SupportURL":          req.SupportURL,
		"DownloadDescription": req.DownloadDescription,
		"ShortDescription":    req.ShortDescription,
		"LongDescription":     req.LongDescription,
		"Keywords":            req.Keywords,
		"ReleaseNotes":        req.ReleaseNotes,
		"SyncStatus":          "pending",
		"Source":              "local",
	}

	err = services.UpdateAppLocalization(uint(appID), languageCode, updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App localization updated successfully",
	})
}

// DeleteAppLocalization deletes an app localization
func (ctrl *AppLocalizationController) DeleteAppLocalization(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	languageCode := c.Params("language")

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	role, _ := context.GetUserRoleFromContext(c)
	if role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Only admin can delete localizations")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	err = services.DeleteAppLocalization(uint(appID), languageCode)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App localization deleted successfully",
	})
}

// GetAppLanguages retrieves languages for an app
func (ctrl *AppLocalizationController) GetAppLanguages(c *fiber.Ctx) error {
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

	languages, err := services.GetAppSupportedLanguages(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
}

// AddAppLanguage adds a language to an app
func (ctrl *AppLocalizationController) AddAppLanguage(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Only owner may manage languages
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}
	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can manage languages")
	}

	var req AddAppLanguageRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := services.ValidateLanguages([]string{req.Language}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = services.GetOrCreateAppLocalization(
		uint(appID),
		req.Language,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	languages, err := services.GetAppSupportedLanguages(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
}

// RemoveAppLanguage removes a language from an app
func (ctrl *AppLocalizationController) RemoveAppLanguage(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	language := c.Params("language")

	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Only owner may manage languages
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}
	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can manage languages")
	}

	if err := services.ValidateLanguages([]string{language}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := services.DeleteAppLocalization(uint(appID), language); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	languages, err := services.GetAppSupportedLanguages(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
}