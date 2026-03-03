package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fdddf/opentrans/internal/dao/query"
	"github.com/fdddf/opentrans/internal/database"
	"gorm.io/gorm"
)

// AppLocalizationService handles app localization operations
type AppLocalizationService struct {
	DB    *database.Database
	Query *query.Query
}

// CreateAppLocalization creates a new localization for an app
func (s *AppLocalizationService) CreateAppLocalization(appID uint, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, description, keywords, whatsNew, promotionalText, whatToTest string) (*database.AppLocalization, error) {
	if !IsSupportedLanguage(languageCode) {
		return nil, fmt.Errorf("unsupported language: %s", languageCode)
	}

	// Validate localization data
	validationResult := ValidateAppLocalization(languageCode, name, subtitle, "", keywords, privacyURL, marketingURL, supportURL, promotionalText, whatsNew)
	if !validationResult.IsValid {
		var errorMsgs []string
		for _, err := range validationResult.Errors {
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", err.Field, err.Message))
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(errorMsgs, "; "))
	}

	// Check if localization already exists for this app and language
	existingLocalization, err := s.Query.AppLocalization.Where(
		s.Query.AppLocalization.AppID.Eq(appID),
		s.Query.AppLocalization.LanguageCode.Eq(languageCode),
	).First()
	if err == nil && existingLocalization != nil {
		return nil, errors.New("localization already exists for this language and app")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing localization: %v", err)
	}

	now := time.Now()
	localization := &database.AppLocalization{
		AppID:           appID,
		LanguageCode:    languageCode,
		Name:            name,
		Subtitle:        subtitle,
		PrivacyURL:      privacyURL,
		MarketingURL:    marketingURL,
		SupportURL:      supportURL,
		Description:     description,
		Keywords:        keywords,
		WhatsNew:        whatsNew,
		PromotionalText: promotionalText,
		WhatToTest:      whatToTest,
		Locale:          languageCode, // Default the locale to the language code
		SyncedAt:        &now,
		Source:          "local",
		SyncStatus:      "pending",
	}

	if err := s.Query.AppLocalization.Create(localization); err != nil {
		return nil, fmt.Errorf("failed to create app localization: %v", err)
	}

	return localization, nil
}

// GetAppLocalization retrieves a specific localization for an app
func (s *AppLocalizationService) GetAppLocalization(appID uint, languageCode string) (*database.AppLocalization, error) {
	localization, err := s.Query.AppLocalization.Where(
		s.Query.AppLocalization.AppID.Eq(appID),
		s.Query.AppLocalization.LanguageCode.Eq(languageCode),
	).Preload(s.Query.AppLocalization.App).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("app localization not found")
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	return localization, nil
}

// GetAppLocalizations retrieves all localizations for an app
func (s *AppLocalizationService) GetAppLocalizations(appID uint) ([]database.AppLocalization, error) {
	localizations, err := s.Query.AppLocalization.Where(
		s.Query.AppLocalization.AppID.Eq(appID),
	).Order(s.Query.AppLocalization.LanguageCode.Asc()).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve app localizations: %v", err)
	}

	// Convert slice of pointers to slice of values
	result := make([]database.AppLocalization, len(localizations))
	for i, loc := range localizations {
		result[i] = *loc
	}

	return result, nil
}

// GetAppLocalizationByLanguageCode retrieves a localization by language code only
func (s *AppLocalizationService) GetAppLocalizationByLanguageCode(languageCode string) ([]database.AppLocalization, error) {
	localizations, err := s.Query.AppLocalization.Where(
		s.Query.AppLocalization.LanguageCode.Eq(languageCode),
	).Preload(s.Query.AppLocalization.App).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve app localizations by language code: %v", err)
	}

	// Convert slice of pointers to slice of values
	result := make([]database.AppLocalization, len(localizations))
	for i, loc := range localizations {
		result[i] = *loc
	}

	return result, nil
}

// UpdateAppLocalization updates an existing localization
func (s *AppLocalizationService) UpdateAppLocalization(appID uint, languageCode string, updates map[string]interface{}) error {
	if _, ok := updates["SyncStatus"]; !ok {
		updates["SyncStatus"] = "pending"
	}
	
	// Log updates for debugging
	if desc, ok := updates["description"]; ok && desc != nil {
		descStr := desc.(string)
		newlineCount := strings.Count(descStr, "\n")
		fmt.Printf("[UpdateAppLocalization] Updating %s: description %d chars, %d newlines\n", languageCode, len(descStr), newlineCount)
	}
	
	result, err := s.Query.AppLocalization.Where(
		s.Query.AppLocalization.AppID.Eq(appID),
		s.Query.AppLocalization.LanguageCode.Eq(languageCode),
	).Updates(updates)
	if err != nil {
		return fmt.Errorf("failed to update app localization: %v", err)
	}

	if result.RowsAffected == 0 {
		return errors.New("app localization not found")
	}

	return nil
}

// UpdateAppLocalizationWithValidation updates an existing localization with validation
func (s *AppLocalizationService) UpdateAppLocalizationWithValidation(appID uint, languageCode string, updates map[string]interface{}) error {
	// Extract fields for validation (using snake_case keys to match updates map)
	name := getStringFromMap(updates, "name")
	subtitle := getStringFromMap(updates, "subtitle")
	shortDescription := getStringFromMap(updates, "short_description")
	keywords := getStringFromMap(updates, "keywords")
	privacyURL := getStringFromMap(updates, "privacy_url")
	marketingURL := getStringFromMap(updates, "marketing_url")
	supportURL := getStringFromMap(updates, "support_url")
	promotionalText := getStringFromMap(updates, "promotional_text")
	whatsNew := getStringFromMap(updates, "whats_new")

	// Validate localization data
	validationResult := ValidateAppLocalization(languageCode, name, subtitle, shortDescription, keywords, privacyURL, marketingURL, supportURL, promotionalText, whatsNew)
	if !validationResult.IsValid {
		var errorMsgs []string
		for _, err := range validationResult.Errors {
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", err.Field, err.Message))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errorMsgs, "; "))
	}

	return s.UpdateAppLocalization(appID, languageCode, updates)
}

// getStringFromMap safely extracts a string from a map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// DeleteAppLocalization soft deletes an app localization
func (s *AppLocalizationService) DeleteAppLocalization(appID uint, languageCode string) error {
	result, err := s.Query.AppLocalization.Where(
		s.Query.AppLocalization.AppID.Eq(appID),
		s.Query.AppLocalization.LanguageCode.Eq(languageCode),
	).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete app localization: %v", err)
	}

	if result.RowsAffected == 0 {
		return errors.New("app localization not found")
	}

	return nil
}

// BulkCreateAppLocalizations creates multiple localizations for an app
func (s *AppLocalizationService) BulkCreateAppLocalizations(appID uint, localizations []database.AppLocalization) error {
	if len(localizations) == 0 {
		return nil
	}

	for i := range localizations {
		localizations[i].AppID = appID
	}

	// Use a transaction to ensure all localizations are created together
	return s.Query.Transaction(func(tx *query.Query) error {
		for _, localization := range localizations {
			// Check if localization already exists for this app and language
			existing, err := tx.AppLocalization.Where(
				tx.AppLocalization.AppID.Eq(localization.AppID),
				tx.AppLocalization.LanguageCode.Eq(localization.LanguageCode),
			).First()

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new localization if it doesn't exist
				if err := tx.AppLocalization.Create(&localization); err != nil {
					return fmt.Errorf("failed to create localization: %v", err)
				}
			} else if err != nil {
				return fmt.Errorf("database error: %v", err)
			} else {
				// Update existing localization
				_, err = tx.AppLocalization.Where(
					tx.AppLocalization.ID.Eq(existing.ID),
				).Updates(map[string]interface{}{
					"Name":             localization.Name,
					"Subtitle":         localization.Subtitle,
					"PrivacyURL":       localization.PrivacyURL,
					"MarketingURL":     localization.MarketingURL,
					"SupportURL":       localization.SupportURL,
					"Description":      localization.Description,
					"Keywords":         localization.Keywords,
					"WhatsNew":         localization.WhatsNew,
					"PromotionalText":  localization.PromotionalText,
				})
				if err != nil {
					return fmt.Errorf("failed to update localization: %v", err)
				}
			}
		}
		return nil
	})
}

// GetOrCreateAppLocalization gets an existing localization or creates a new one
func (s *AppLocalizationService) GetOrCreateAppLocalization(appID uint, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, description, keywords, whatsNew string) (*database.AppLocalization, error) {
	// First try to get existing localization
	localization, err := s.GetAppLocalization(appID, languageCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// If there was an error other than "not found", return it
		return nil, err
	}

	if err == nil {
		// Localization exists, return it
		return localization, nil
	}

	// Localization doesn't exist, create it with empty values for new fields
	return s.CreateAppLocalization(appID, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, description, keywords, whatsNew, "", "")
}

// GetAppSupportedLanguages retrieves all languages supported by an app
func (s *AppLocalizationService) GetAppSupportedLanguages(appID uint) ([]string, error) {
	var localizationCodes []string
	err := s.Query.AppLocalization.Where(
		s.Query.AppLocalization.AppID.Eq(appID),
	).Pluck(s.Query.AppLocalization.LanguageCode, &localizationCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve supported languages: %v", err)
	}

	// Also add the app's primary locale if not already included
	// We'll need to pass the AppService here, but for now using global function
	app, err := GetApp(appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %v", err)
	}

	if app.PrimaryLocale != "" {
		found := false
		for _, code := range localizationCodes {
			if code == app.PrimaryLocale {
				found = true
				break
			}
		}
		if !found {
			localizationCodes = append(localizationCodes, app.PrimaryLocale)
		}
	}

	return localizationCodes, nil
}

// Global functions for backward compatibility
var appLocalizationServiceInstance *AppLocalizationService

func SetAppLocalizationService(db *database.Database) {
	appLocalizationServiceInstance = &AppLocalizationService{
		DB:    db,
		Query: query.Use(db.DB),
	}
}

func CreateAppLocalization(appID uint, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, description, keywords, whatsNew, promotionalText, whatToTest string) (*database.AppLocalization, error) {
	return appLocalizationServiceInstance.CreateAppLocalization(appID, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, description, keywords, whatsNew, promotionalText, whatToTest)
}

func GetAppLocalization(appID uint, languageCode string) (*database.AppLocalization, error) {
	return appLocalizationServiceInstance.GetAppLocalization(appID, languageCode)
}

func GetAppLocalizations(appID uint) ([]database.AppLocalization, error) {
	return appLocalizationServiceInstance.GetAppLocalizations(appID)
}

func GetAppLocalizationByLanguageCode(languageCode string) ([]database.AppLocalization, error) {
	return appLocalizationServiceInstance.GetAppLocalizationByLanguageCode(languageCode)
}

func UpdateAppLocalization(appID uint, languageCode string, updates map[string]interface{}) error {
	return appLocalizationServiceInstance.UpdateAppLocalization(appID, languageCode, updates)
}

func DeleteAppLocalization(appID uint, languageCode string) error {
	return appLocalizationServiceInstance.DeleteAppLocalization(appID, languageCode)
}

func BulkCreateAppLocalizations(appID uint, localizations []database.AppLocalization) error {
	return appLocalizationServiceInstance.BulkCreateAppLocalizations(appID, localizations)
}

func GetOrCreateAppLocalization(appID uint, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, description, keywords, whatsNew string) (*database.AppLocalization, error) {
	return appLocalizationServiceInstance.GetOrCreateAppLocalization(appID, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, description, keywords, whatsNew)
}

func GetAppSupportedLanguages(appID uint) ([]string, error) {
	return appLocalizationServiceInstance.GetAppSupportedLanguages(appID)
}