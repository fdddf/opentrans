package services

import (
	"fmt"
	"log"
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
	Warnings []LocalizationValidationError `json:"warnings"`
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
// Note: This function now only adds warnings for length violations, allowing database storage
// The actual truncation happens during sync with Apple
// Note: name is no longer required during validation as it can be inherited from primary locale
func (v *AppLocalizationValidator) ValidateAppLocalization(languageCode, name, subtitle, shortDescription, keywords, privacyURL, marketingURL, supportURL, promotionalText, releaseNotes string) *LocalizationValidationResult {
	result := &LocalizationValidationResult{
		IsValid:  true,
		Errors:   []LocalizationValidationError{},
		Warnings: []LocalizationValidationError{},
	}

	// Validate name length (name is optional, can be inherited from primary locale)
	if len(name) > MaxAppNameLength {
		warning := fmt.Sprintf("App name exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxAppNameLength, len(name))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "name",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	// Validate subtitle
	if len(subtitle) > MaxSubtitleLength {
		warning := fmt.Sprintf("Subtitle exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxSubtitleLength, len(subtitle))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "subtitle",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	// Validate short description
	if len(shortDescription) > MaxShortDescriptionLength {
		warning := fmt.Sprintf("Short description exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxShortDescriptionLength, len(shortDescription))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "shortDescription",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	// Validate keywords
	if len(keywords) > MaxKeywordsLength {
		warning := fmt.Sprintf("Keywords exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxKeywordsLength, len(keywords))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "keywords",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	// Validate URLs
	if len(privacyURL) > 0 && len(privacyURL) > MaxPrivacyURLLength {
		warning := fmt.Sprintf("Privacy URL exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxPrivacyURLLength, len(privacyURL))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "privacyURL",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	if len(marketingURL) > 0 && len(marketingURL) > MaxMarketingURLLength {
		warning := fmt.Sprintf("Marketing URL exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxMarketingURLLength, len(marketingURL))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "marketingURL",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	if len(supportURL) > 0 && len(supportURL) > MaxSupportURLLength {
		warning := fmt.Sprintf("Support URL exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxSupportURLLength, len(supportURL))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "supportURL",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	// Validate promotional text
	if len(promotionalText) > MaxPromotionalTextLength {
		warning := fmt.Sprintf("Promotional text exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxPromotionalTextLength, len(promotionalText))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "promotionalText",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	// Validate release notes
	if len(releaseNotes) > MaxReleaseNotesLength {
		warning := fmt.Sprintf("Release notes exceeds %d characters (current: %d), will be truncated when syncing to Apple", MaxReleaseNotesLength, len(releaseNotes))
		result.Warnings = append(result.Warnings, LocalizationValidationError{
			Field:   "releaseNotes",
			Message: warning,
		})
		log.Printf("[Validation Warning] [%s] %s", languageCode, warning)
	}

	return result
}

// TruncateField truncates a field to the specified maximum length
// Used when syncing to Apple App Store Connect
func TruncateField(field string, maxLength int) string {
	if len(field) <= maxLength {
		return field
	}
	// Truncate at maxLength, trying to avoid cutting in the middle of a word
	truncated := field[:maxLength]
	// Try to find last space or comma to avoid cutting words
	lastSpace := strings.LastIndexAny(truncated, " ,")
	if lastSpace > maxLength/2 {
		truncated = truncated[:lastSpace]
	}
	return strings.TrimSpace(truncated)
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