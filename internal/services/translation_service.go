package services

import (
	"errors"
	"fmt"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/model"
	"gorm.io/gorm"
)

// CreateTranslations creates multiple translations in batch
func CreateTranslations(translations []database.Translation) error {
	if len(translations) == 0 {
		return nil
	}

	// Use a transaction to ensure all translations are created together
	tx := database.DB.Begin()
	defer tx.Rollback()

	for _, translation := range translations {
		// Check if translation already exists for this key and language in the project
		var existing database.Translation
		result := tx.Where("project_id = ? AND key = ? AND target_language = ?",
			translation.ProjectID, translation.Key, translation.TargetLanguage).First(&existing)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new translation if it doesn't exist
			result = tx.Create(&translation)
			if result.Error != nil {
				return fmt.Errorf("failed to create translation: %v", result.Error)
			}
		} else if result.Error != nil {
			return fmt.Errorf("database error: %v", result.Error)
		} else {
			// Update existing translation
			result = tx.Model(&database.Translation{}).
				Where("project_id = ? AND key = ? AND target_language = ?",
					translation.ProjectID, translation.Key, translation.TargetLanguage).
				Updates(map[string]interface{}{
					"TargetText":          translation.TargetText,
					"State":               translation.State,
					"TranslationProvider": translation.TranslationProvider,
				})
			if result.Error != nil {
				return fmt.Errorf("failed to update translation: %v", result.Error)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit translations: %v", err)
	}

	return nil
}

// GetTranslationsByProject retrieves all translations for a project
func GetTranslationsByProject(projectID uint) ([]database.Translation, error) {
	var translations []database.Translation
	result := database.DB.Where("project_id = ?", projectID).Find(&translations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve translations: %v", result.Error)
	}

	return translations, nil
}

// GetTranslationsByProjectAndLanguage retrieves translations for a project filtered by target language
func GetTranslationsByProjectAndLanguage(projectID uint, language string) ([]database.Translation, error) {
	var translations []database.Translation
	result := database.DB.Where("project_id = ? AND target_language = ?", projectID, language).Find(&translations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve translations: %v", result.Error)
	}

	return translations, nil
}

// GetTranslationsByProjectAndKeys retrieves translations for specific keys in a project
func GetTranslationsByProjectAndKeys(projectID uint, keys []string) ([]database.Translation, error) {
	var translations []database.Translation
	result := database.DB.Where("project_id = ? AND key IN ?", projectID, keys).Find(&translations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve translations: %v", result.Error)
	}

	return translations, nil
}

// UpdateTranslationsFromXCStrings updates translations in the database based on XCStrings content
func UpdateTranslationsFromXCStrings(projectID uint, xcstrings *model.XCStrings, translationProvider string) error {
	var translations []database.Translation

	for key, entry := range xcstrings.Strings {
		for lang, localization := range entry.Localizations {
			if lang != xcstrings.SourceLanguage { // Only create translations for target languages
				translation := database.Translation{
					ProjectID:           projectID,
					Key:                 key,
					SourceText:          getSourceText(key, entry, xcstrings.SourceLanguage),
					TargetText:          localization.StringUnit.Value,
					TargetLanguage:      lang,
					State:               localization.StringUnit.State,
					TranslationProvider: translationProvider,
				}
				translations = append(translations, translation)
			}
		}
	}

	return CreateTranslations(translations)
}

// getSourceText helper function to extract source text from an entry
func getSourceText(key string, entry model.StringEntry, sourceLanguage string) string {
	if sourceLoc, exists := entry.Localizations[sourceLanguage]; exists {
		return sourceLoc.StringUnit.Value
	}
	return key // fallback to key if source text is not available
}

// ApplyTranslationsToXCStrings applies stored translations to XCStrings model
func ApplyTranslationsToXCStrings(projectID uint, xcstrings *model.XCStrings) error {
	translations, err := GetTranslationsByProject(projectID)
	if err != nil {
		return err
	}

	for _, translation := range translations {
		// Get or create the string entry
		entry, exists := xcstrings.Strings[translation.Key]
		if !exists {
			entry = model.StringEntry{
				Localizations: make(map[string]model.Localization),
			}
		}

		// Update the localization for this language
		localization := model.Localization{
			StringUnit: model.StringUnit{
				Value: translation.TargetText,
				State: translation.State,
			},
		}

		entry.Localizations[translation.TargetLanguage] = localization

		// If this is a source language, update the source language of the entire file
		if translation.TargetLanguage != xcstrings.SourceLanguage {
			// Only set if it's a target language, not source
			xcstrings.Strings[translation.Key] = entry
		} else {
			// If this was a source language translation, we might want to update the source
			// but for now we'll only update target languages
			if _, sourceExists := entry.Localizations[xcstrings.SourceLanguage]; !sourceExists {
				// Add to source if it doesn't exist
				entry.Localizations[xcstrings.SourceLanguage] = localization
				xcstrings.Strings[translation.Key] = entry
			}
		}

		// Update the entry in the string map
		xcstrings.Strings[translation.Key] = entry
	}

	return nil
}

// GetMissingTranslations retrieves keys that are missing translations for specified target languages
func GetMissingTranslations(projectID uint, targetLanguages []string) ([]string, error) {
	// Get all unique keys from the project's content
	project, err := GetProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %v", err)
	}

	xcstrings, err := ParseProjectContent(project)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project content: %v", err)
	}

	var allKeys []string
	for key := range xcstrings.Strings {
		allKeys = append(allKeys, key)
	}

	var missingKeys []string

	for _, key := range allKeys {
		for _, lang := range targetLanguages {
			// Check if translation exists for this key and language
			var count int64
			database.DB.Model(&database.Translation{}).
				Where("project_id = ? AND key = ? AND target_language = ?", projectID, key, lang).
				Count(&count)

			if count == 0 {
				// Check if the key doesn't have a translation for this language in the content structure too
				if entry, exists := xcstrings.Strings[key]; exists {
					if _, hasTranslation := entry.Localizations[lang]; !hasTranslation {
						if !containsString(missingKeys, key) {
							missingKeys = append(missingKeys, key)
						}
					}
				} else {
					if !containsString(missingKeys, key) {
						missingKeys = append(missingKeys, key)
					}
				}
			}
		}
	}

	return missingKeys, nil
}

// containsString helper function
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// BulkCreateTranslationsFromEntries creates translations from XCStrings entries
func BulkCreateTranslationsFromEntries(projectID uint, entries map[string]model.StringEntry, sourceLanguage string, translationProvider string) error {
	var translations []database.Translation

	for key, entry := range entries {
		for lang, localization := range entry.Localizations {
			// Only create translations for non-source languages
			if lang != sourceLanguage {
				translation := database.Translation{
					ProjectID:           projectID,
					Key:                 key,
					SourceText:          getSourceText(key, entry, sourceLanguage),
					TargetText:          localization.StringUnit.Value,
					TargetLanguage:      lang,
					State:               localization.StringUnit.State,
					TranslationProvider: translationProvider,
				}
				translations = append(translations, translation)
			}
		}
	}

	return CreateTranslations(translations)
}

// UpdateProjectSourceLanguage updates the source language for a project and associated translations
func UpdateProjectSourceLanguage(projectID uint, newSourceLanguage string) error {
	// Update the project
	if err := UpdateProject(projectID, map[string]interface{}{
		"SourceLanguage": newSourceLanguage,
	}); err != nil {
		return fmt.Errorf("failed to update project source language: %v", err)
	}

	return nil
}
