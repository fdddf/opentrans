package services

import (
	"fmt"
	"strings"
)

// LocalizationValidationError represents a validation error for a localization field
type LocalizationValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// LocalizationValidationResult represents the result of localization validation
type LocalizationValidationResult struct {
	IsValid bool                           `json:"is_valid"`
	Errors  []LocalizationValidationError `json:"errors"`
}

// AppLocalizationValidator handles validation rules for app localizations
type AppLocalizationValidator struct{}

// NewAppLocalizationValidator creates a new validator
func NewAppLocalizationValidator() *AppLocalizationValidator {
	return &AppLocalizationValidator{}
}

// Validation rules based on Apple App Store Connect requirements
const (
	MaxAppNameLength        = 30
	MaxSubtitleLength       = 30
	MaxShortDescriptionLength = 80
	MaxKeywordsLength       = 100
	MaxPrivacyURLLength     = 255
	MaxMarketingURLLength   = 255
	MaxSupportURLLength     = 255
	MaxPromotionalTextLength = 170
	MaxReleaseNotesLength   = 4000
)

// ValidateAppLocalization validates an app localization based on Apple App Store Connect requirements
func (v *AppLocalizationValidator) ValidateAppLocalization(languageCode, name, subtitle, shortDescription, keywords, privacyURL, marketingURL, supportURL, promotionalText, releaseNotes string) *LocalizationValidationResult {
	result := &LocalizationValidationResult{
		IsValid: true,
		Errors:  []LocalizationValidationError{},
	}

	// Validate name (required for all languages except primary locale might be optional)
	if strings.TrimSpace(name) == "" {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "name",
			Message: "App name is required",
		})
		result.IsValid = false
	} else if len(name) > MaxAppNameLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "name",
			Message: fmt.Sprintf("App name must not exceed %d characters (current: %d)", MaxAppNameLength, len(name)),
		})
		result.IsValid = false
	}

	// Validate subtitle
	if len(subtitle) > MaxSubtitleLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "subtitle",
			Message: fmt.Sprintf("Subtitle must not exceed %d characters (current: %d)", MaxSubtitleLength, len(subtitle)),
		})
		result.IsValid = false
	}

	// Validate short description
	if len(shortDescription) > MaxShortDescriptionLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "shortDescription",
			Message: fmt.Sprintf("Short description must not exceed %d characters (current: %d)", MaxShortDescriptionLength, len(shortDescription)),
		})
		result.IsValid = false
	}

	// Validate keywords
	if len(keywords) > MaxKeywordsLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "keywords",
			Message: fmt.Sprintf("Keywords must not exceed %d characters (current: %d)", MaxKeywordsLength, len(keywords)),
		})
		result.IsValid = false
	}

	// Validate URLs
	if len(privacyURL) > 0 && len(privacyURL) > MaxPrivacyURLLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "privacyURL",
			Message: fmt.Sprintf("Privacy URL must not exceed %d characters (current: %d)", MaxPrivacyURLLength, len(privacyURL)),
		})
		result.IsValid = false
	}

	if len(marketingURL) > 0 && len(marketingURL) > MaxMarketingURLLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "marketingURL",
			Message: fmt.Sprintf("Marketing URL must not exceed %d characters (current: %d)", MaxMarketingURLLength, len(marketingURL)),
		})
		result.IsValid = false
	}

	if len(supportURL) > 0 && len(supportURL) > MaxSupportURLLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "supportURL",
			Message: fmt.Sprintf("Support URL must not exceed %d characters (current: %d)", MaxSupportURLLength, len(supportURL)),
		})
		result.IsValid = false
	}

	// Validate promotional text
	if len(promotionalText) > MaxPromotionalTextLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "promotionalText",
			Message: fmt.Sprintf("Promotional text must not exceed %d characters (current: %d)", MaxPromotionalTextLength, len(promotionalText)),
		})
		result.IsValid = false
	}

	// Validate release notes
	if len(releaseNotes) > MaxReleaseNotesLength {
		result.Errors = append(result.Errors, LocalizationValidationError{
			Field:   "releaseNotes",
			Message: fmt.Sprintf("Release notes must not exceed %d characters (current: %d)", MaxReleaseNotesLength, len(releaseNotes)),
		})
		result.IsValid = false
	}

	return result
}

// ValidateURLFormat validates if a string is a valid URL format
func (v *AppLocalizationValidator) ValidateURLFormat(url string) bool {
	if strings.TrimSpace(url) == "" {
		return true // Empty URLs are allowed
	}

	// Basic URL validation - must start with http:// or https://
	url = strings.TrimSpace(url)
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// Global validator instance
var localizationValidator *AppLocalizationValidator

func SetLocalizationValidator() {
	localizationValidator = NewAppLocalizationValidator()
}

func ValidateAppLocalization(languageCode, name, subtitle, shortDescription, keywords, privacyURL, marketingURL, supportURL, promotionalText, releaseNotes string) *LocalizationValidationResult {
	if localizationValidator == nil {
		SetLocalizationValidator()
	}
	return localizationValidator.ValidateAppLocalization(languageCode, name, subtitle, shortDescription, keywords, privacyURL, marketingURL, supportURL, promotionalText, releaseNotes)
}

func ValidateURLFormat(url string) bool {
	if localizationValidator == nil {
		SetLocalizationValidator()
	}
	return localizationValidator.ValidateURLFormat(url)
}