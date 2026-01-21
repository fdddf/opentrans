package services

import (
	"errors"
	"fmt"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"gorm.io/gorm"
)

// AppLocalizationService handles app localization operations
type AppLocalizationService struct {
	DB *database.Database
}

// CreateAppLocalization creates a new localization for an app
func (s *AppLocalizationService) CreateAppLocalization(appID uint, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, downloadDescription, shortDescription, longDescription, keywords, releaseNotes string) (*database.AppLocalization, error) {
	if !IsSupportedLanguage(languageCode) {
		return nil, fmt.Errorf("unsupported language: %s", languageCode)
	}
	// Check if localization already exists for this app and language
	var existingLocalization database.AppLocalization
	result := s.DB.Where("app_id = ? AND language_code = ?", appID, languageCode).First(&existingLocalization)
	if result.Error == nil {
		return nil, errors.New("localization already exists for this language and app")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing localization: %v", result.Error)
	}

	localization := &database.AppLocalization{
		AppID:               appID,
		LanguageCode:        languageCode,
		Name:                name,
		Subtitle:            subtitle,
		PrivacyURL:          privacyURL,
		MarketingURL:        marketingURL,
		SupportURL:          supportURL,
		DownloadDescription: downloadDescription,
		ShortDescription:    shortDescription,
		LongDescription:     longDescription,
		Keywords:            keywords,
		ReleaseNotes:        releaseNotes,
	}

	result = s.DB.Create(localization)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create app localization: %v", result.Error)
	}

	return localization, nil
}

// GetAppLocalization retrieves a specific localization for an app
func (s *AppLocalizationService) GetAppLocalization(appID uint, languageCode string) (*database.AppLocalization, error) {
	var localization database.AppLocalization
	result := s.DB.Preload("App").Where("app_id = ? AND language_code = ?", appID, languageCode).First(&localization)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("app localization not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &localization, nil
}

// GetAppLocalizations retrieves all localizations for an app
func (s *AppLocalizationService) GetAppLocalizations(appID uint) ([]database.AppLocalization, error) {
	var localizations []database.AppLocalization
	result := s.DB.Where("app_id = ?", appID).Order("language_code ASC").Find(&localizations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve app localizations: %v", result.Error)
	}

	return localizations, nil
}

// GetAppLocalizationByLanguageCode retrieves a localization by language code only
func (s *AppLocalizationService) GetAppLocalizationByLanguageCode(languageCode string) ([]database.AppLocalization, error) {
	var localizations []database.AppLocalization
	result := s.DB.Preload("App").Where("language_code = ?", languageCode).Find(&localizations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve app localizations by language code: %v", result.Error)
	}

	return localizations, nil
}

// UpdateAppLocalization updates an existing localization
func (s *AppLocalizationService) UpdateAppLocalization(appID uint, languageCode string, updates map[string]interface{}) error {
	result := s.DB.Model(&database.AppLocalization{}).Where("app_id = ? AND language_code = ?", appID, languageCode).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update app localization: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("app localization not found")
	}

	return nil
}

// DeleteAppLocalization soft deletes an app localization
func (s *AppLocalizationService) DeleteAppLocalization(appID uint, languageCode string) error {
	result := s.DB.Where("app_id = ? AND language_code = ?", appID, languageCode).Delete(&database.AppLocalization{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete app localization: %v", result.Error)
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
	tx := s.DB.Begin()
	defer tx.Rollback()

	for _, localization := range localizations {
		// Check if localization already exists for this app and language
		var existing database.AppLocalization
		result := tx.Where("app_id = ? AND language_code = ?",
			localization.AppID, localization.LanguageCode).First(&existing)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new localization if it doesn't exist
			result = tx.Create(&localization)
			if result.Error != nil {
				return fmt.Errorf("failed to create localization: %v", result.Error)
			}
		} else if result.Error != nil {
			return fmt.Errorf("database error: %v", result.Error)
		} else {
			// Update existing localization
			result = tx.Model(&database.AppLocalization{}).
				Where("app_id = ? AND language_code = ?",
					localization.AppID, localization.LanguageCode).
				Updates(map[string]interface{}{
					"Name":                localization.Name,
					"Subtitle":            localization.Subtitle,
					"PrivacyURL":          localization.PrivacyURL,
					"MarketingURL":        localization.MarketingURL,
					"SupportURL":          localization.SupportURL,
					"DownloadDescription": localization.DownloadDescription,
					"ShortDescription":    localization.ShortDescription,
					"LongDescription":     localization.LongDescription,
					"Keywords":            localization.Keywords,
					"ReleaseNotes":        localization.ReleaseNotes,
				})
			if result.Error != nil {
				return fmt.Errorf("failed to update localization: %v", result.Error)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit localizations: %v", err)
	}

	return nil
}

// GetOrCreateAppLocalization gets an existing localization or creates a new one
func (s *AppLocalizationService) GetOrCreateAppLocalization(appID uint, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, downloadDescription, shortDescription, longDescription, keywords, releaseNotes string) (*database.AppLocalization, error) {
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

	// Localization doesn't exist, create it
	return s.CreateAppLocalization(appID, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, downloadDescription, shortDescription, longDescription, keywords, releaseNotes)
}

// GetAppSupportedLanguages retrieves all languages supported by an app
func (s *AppLocalizationService) GetAppSupportedLanguages(appID uint) ([]string, error) {
	var localizationCodes []string
	result := s.DB.Model(&database.AppLocalization{}).Where("app_id = ?", appID).Pluck("language_code", &localizationCodes)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve supported languages: %v", result.Error)
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
	appLocalizationServiceInstance = &AppLocalizationService{DB: db}
}

func CreateAppLocalization(appID uint, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, downloadDescription, shortDescription, longDescription, keywords, releaseNotes string) (*database.AppLocalization, error) {
	return appLocalizationServiceInstance.CreateAppLocalization(appID, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, downloadDescription, shortDescription, longDescription, keywords, releaseNotes)
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

func GetOrCreateAppLocalization(appID uint, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, downloadDescription, shortDescription, longDescription, keywords, releaseNotes string) (*database.AppLocalization, error) {
	return appLocalizationServiceInstance.GetOrCreateAppLocalization(appID, languageCode, name, subtitle, privacyURL, marketingURL, supportURL, downloadDescription, shortDescription, longDescription, keywords, releaseNotes)
}

func GetAppSupportedLanguages(appID uint) ([]string, error) {
	return appLocalizationServiceInstance.GetAppSupportedLanguages(appID)
}
