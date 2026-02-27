package services

import (
	"encoding/json"

	"github.com/fdddf/opentrans/internal/dao/query"
	"github.com/fdddf/opentrans/internal/database"
	"github.com/gofiber/fiber/v2"
)

// ActivityService handles user activity logging
type ActivityService struct {
	db    *database.Database
	Query *query.Query
}

// NewActivityService creates a new activity service
func NewActivityService(db *database.Database, q *query.Query) *ActivityService {
	return &ActivityService{
		db:    db,
		Query: q,
	}
}

// LogActivity records a user activity
func (s *ActivityService) LogActivity(ctx *fiber.Ctx, userID uint, action string, details interface{}) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		detailsJSON = []byte("{}")
	}

	activity := &database.UserActivity{
		UserID:    userID,
		Action:    action,
		Details:   string(detailsJSON),
		IPAddress: ctx.IP(),
		UserAgent: ctx.Get("User-Agent"),
	}

	return s.Query.UserActivity.Create(activity)
}

// Helper functions for logging common activities

// LogUserLogin logs user login
func (s *ActivityService) LogUserLogin(ctx *fiber.Ctx, userID uint, username string) error {
	return s.LogActivity(ctx, userID, "user_logged_in", map[string]interface{}{
		"username": username,
	})
}

// LogUserLogout logs user logout
func (s *ActivityService) LogUserLogout(ctx *fiber.Ctx, userID uint, username string) error {
	return s.LogActivity(ctx, userID, "user_logged_out", map[string]interface{}{
		"username": username,
	})
}

// LogUserRegister logs user registration
func (s *ActivityService) LogUserRegister(ctx *fiber.Ctx, userID uint, username, email string) error {
	return s.LogActivity(ctx, userID, "user_registered", map[string]interface{}{
		"username": username,
		"email":    email,
	})
}

// LogAppCreated logs app creation
func (s *ActivityService) LogAppCreated(ctx *fiber.Ctx, userID uint, appID uint, name, bundleID string) error {
	return s.LogActivity(ctx, userID, "app_created", map[string]interface{}{
		"appId":    appID,
		"name":     name,
		"bundleId": bundleID,
	})
}

// LogAppUpdated logs app update
func (s *ActivityService) LogAppUpdated(ctx *fiber.Ctx, userID uint, appID uint, name string) error {
	return s.LogActivity(ctx, userID, "app_updated", map[string]interface{}{
		"appId": appID,
		"name":  name,
	})
}

// LogAppDeleted logs app deletion
func (s *ActivityService) LogAppDeleted(ctx *fiber.Ctx, userID uint, appID uint, name string) error {
	return s.LogActivity(ctx, userID, "app_deleted", map[string]interface{}{
		"appId": appID,
		"name":  name,
	})
}

// LogProjectCreated logs project creation
func (s *ActivityService) LogProjectCreated(ctx *fiber.Ctx, userID uint, projectID uint, name, fileName string) error {
	return s.LogActivity(ctx, userID, "project_created", map[string]interface{}{
		"projectId": projectID,
		"name":      name,
		"fileName":  fileName,
	})
}

// LogProjectUpdated logs project update
func (s *ActivityService) LogProjectUpdated(ctx *fiber.Ctx, userID uint, projectID uint, name string) error {
	return s.LogActivity(ctx, userID, "project_updated", map[string]interface{}{
		"projectId": projectID,
		"name":      name,
	})
}

// LogProjectDeleted logs project deletion
func (s *ActivityService) LogProjectDeleted(ctx *fiber.Ctx, userID uint, projectID uint, name string) error {
	return s.LogActivity(ctx, userID, "project_deleted", map[string]interface{}{
		"projectId": projectID,
		"name":      name,
	})
}

// LogTranslationStarted logs translation job start
func (s *ActivityService) LogTranslationStarted(ctx *fiber.Ctx, userID uint, projectID *uint, appID *uint, providerType string, targetLanguages []string) error {
	return s.LogActivity(ctx, userID, "translation_started", map[string]interface{}{
		"projectId":        projectID,
		"appId":            appID,
		"providerType":     providerType,
		"targetLanguages":  targetLanguages,
	})
}

// LogTranslationCompleted logs translation job completion
func (s *ActivityService) LogTranslationCompleted(ctx *fiber.Ctx, userID uint, projectID *uint, appID *uint, jobType string, total int) error {
	return s.LogActivity(ctx, userID, "translation_completed", map[string]interface{}{
		"projectId": projectID,
		"appId":     appID,
		"jobType":   jobType,
		"total":     total,
	})
}

// LogAppSyncToApple logs app sync to Apple
func (s *ActivityService) LogAppSyncToApple(ctx *fiber.Ctx, userID uint, appID uint, name string) error {
	return s.LogActivity(ctx, userID, "app_synced_to_apple", map[string]interface{}{
		"appId": appID,
		"name":  name,
	})
}

// LogAppleAppsSynced logs Apple apps sync
func (s *ActivityService) LogAppleAppsSynced(ctx *fiber.Ctx, userID uint, configID uint, count int) error {
	return s.LogActivity(ctx, userID, "apple_apps_synced", map[string]interface{}{
		"configId": configID,
		"count":    count,
	})
}

// LogAppleLocalizationsSynced logs Apple localizations sync
func (s *ActivityService) LogAppleLocalizationsSynced(ctx *fiber.Ctx, userID uint, appID uint, direction, strategy string, languageCodes []string) error {
	return s.LogActivity(ctx, userID, "apple_localizations_synced", map[string]interface{}{
		"appId":         appID,
		"direction":     direction,
		"strategy":      strategy,
		"languageCodes": languageCodes,
	})
}

// LogProviderConfigCreated logs provider config creation
func (s *ActivityService) LogProviderConfigCreated(ctx *fiber.Ctx, userID uint, configID uint, providerType string) error {
	return s.LogActivity(ctx, userID, "provider_config_created", map[string]interface{}{
		"configId":     configID,
		"providerType": providerType,
	})
}

// LogProviderConfigUpdated logs provider config update
func (s *ActivityService) LogProviderConfigUpdated(ctx *fiber.Ctx, userID uint, configID uint, providerType string) error {
	return s.LogActivity(ctx, userID, "provider_config_updated", map[string]interface{}{
		"configId":     configID,
		"providerType": providerType,
	})
}

// LogProviderConfigDeleted logs provider config deletion
func (s *ActivityService) LogProviderConfigDeleted(ctx *fiber.Ctx, userID uint, configID uint, providerType string) error {
	return s.LogActivity(ctx, userID, "provider_config_deleted", map[string]interface{}{
		"configId":     configID,
		"providerType": providerType,
	})
}

// LogSubscriptionChanged logs subscription changes
func (s *ActivityService) LogSubscriptionChanged(ctx *fiber.Ctx, userID uint, eventType, subscriptionType string) error {
	return s.LogActivity(ctx, userID, "subscription_changed", map[string]interface{}{
		"eventType":         eventType,
		"subscriptionType":  subscriptionType,
	})
}