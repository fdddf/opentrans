package services

import (
	"errors"
	"fmt"

	"github.com/fdddf/opentrans/internal/dao/query"
	"github.com/fdddf/opentrans/internal/database"
	"github.com/fdddf/opentrans/internal/model"
	"gorm.io/gorm"
)

// TranslationService handles translation-related operations
type TranslationService struct {
	DB    *database.Database
	Query *query.Query
}

// CreateTranslations creates multiple translations in batch
func (s *TranslationService) CreateTranslations(translations []database.Translation) error {
	if len(translations) == 0 {
		return nil
	}

	// Use a transaction to ensure all translations are created together
	return s.Query.Transaction(func(tx *query.Query) error {
		for _, translation := range translations {
			// Check if translation already exists for this key and language in the project
			existing, err := tx.Translation.Where(
				tx.Translation.ProjectID.Eq(translation.ProjectID),
				tx.Translation.Key.Eq(translation.Key),
				tx.Translation.TargetLanguage.Eq(translation.TargetLanguage),
			).First()

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new translation if it doesn't exist
				if err := tx.Translation.Create(&translation); err != nil {
					return fmt.Errorf("failed to create translation: %v", err)
				}
			} else if err != nil {
				return fmt.Errorf("database error: %v", err)
			} else {
				// Update existing translation
				_, err = tx.Translation.Where(tx.Translation.ID.Eq(existing.ID)).Updates(map[string]interface{}{
					"TargetText":          translation.TargetText,
					"State":               translation.State,
					"TranslationProvider": translation.TranslationProvider,
				})
				if err != nil {
					return fmt.Errorf("failed to update translation: %v", err)
				}
			}
		}
		return nil
	})
}

// GetTranslationsByProject retrieves all translations for a project
func (s *TranslationService) GetTranslationsByProject(projectID uint) ([]database.Translation, error) {
	translations, err := s.Query.Translation.Where(
		s.Query.Translation.ProjectID.Eq(projectID),
	).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve translations: %v", err)
	}

	// Convert slice of pointers to slice of values
	result := make([]database.Translation, len(translations))
	for i, t := range translations {
		result[i] = *t
	}

	return result, nil
}

// GetTranslationsByProjectAndLanguage retrieves translations for a project filtered by target language
func (s *TranslationService) GetTranslationsByProjectAndLanguage(projectID uint, language string) ([]database.Translation, error) {
	translations, err := s.Query.Translation.Where(
		s.Query.Translation.ProjectID.Eq(projectID),
		s.Query.Translation.TargetLanguage.Eq(language),
	).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve translations: %v", err)
	}

	// Convert slice of pointers to slice of values
	result := make([]database.Translation, len(translations))
	for i, t := range translations {
		result[i] = *t
	}

	return result, nil
}

// GetTranslationsByProjectAndKeys retrieves translations for specific keys in a project
func (s *TranslationService) GetTranslationsByProjectAndKeys(projectID uint, keys []string) ([]database.Translation, error) {
	translations, err := s.Query.Translation.Where(
		s.Query.Translation.ProjectID.Eq(projectID),
		s.Query.Translation.Key.In(keys...),
	).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve translations: %v", err)
	}

	// Convert slice of pointers to slice of values
	result := make([]database.Translation, len(translations))
	for i, t := range translations {
		result[i] = *t
	}

	return result, nil
}

// ApplyTranslationsToXCStrings applies stored translations to XCStrings model
func (s *TranslationService) ApplyTranslationsToXCStrings(projectID uint, xcstrings *model.XCStrings) error {
	translations, err := s.GetTranslationsByProject(projectID)
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

// getSourceText helper function to extract source text from an entry
func getSourceText(key string, entry model.StringEntry, sourceLanguage string) string {
	if sourceLoc, exists := entry.Localizations[sourceLanguage]; exists {
		return sourceLoc.StringUnit.Value
	}
	return key // fallback to key if source text is not available
}

// BulkCreateTranslationsFromEntries creates translations from XCStrings entries
func (s *TranslationService) BulkCreateTranslationsFromEntries(projectID uint, entries map[string]model.StringEntry, sourceLanguage string, translationProvider string) error {
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

	return s.CreateTranslations(translations)
}

// UpdateProjectSourceLanguage updates the source language for a project and associated translations
func (s *TranslationService) UpdateProjectSourceLanguage(projectID uint, newSourceLanguage string) error {
	// Update the project - this would need to use project service
	// For now using global function that will work with legacy compatibility
	if err := UpdateProject(projectID, map[string]interface{}{
		"SourceLanguage": newSourceLanguage,
	}); err != nil {
		return fmt.Errorf("failed to update project source language: %v", err)
	}

	return nil
}

// Global functions for backward compatibility
var translationServiceInstance *TranslationService

func SetTranslationService(db *database.Database) {
	translationServiceInstance = &TranslationService{
		DB:    db,
		Query: query.Use(db.DB),
	}
}

func CreateTranslations(translations []database.Translation) error {
	return translationServiceInstance.CreateTranslations(translations)
}

func GetTranslationsByProject(projectID uint) ([]database.Translation, error) {
	return translationServiceInstance.GetTranslationsByProject(projectID)
}

func GetTranslationsByProjectAndLanguage(projectID uint, language string) ([]database.Translation, error) {
	return translationServiceInstance.GetTranslationsByProjectAndLanguage(projectID, language)
}

func GetTranslationsByProjectAndKeys(projectID uint, keys []string) ([]database.Translation, error) {
	return translationServiceInstance.GetTranslationsByProjectAndKeys(projectID, keys)
}

func ApplyTranslationsToXCStrings(projectID uint, xcstrings *model.XCStrings) error {
	return translationServiceInstance.ApplyTranslationsToXCStrings(projectID, xcstrings)
}

func BulkCreateTranslationsFromEntries(projectID uint, entries map[string]model.StringEntry, sourceLanguage string, translationProvider string) error {
	return translationServiceInstance.BulkCreateTranslationsFromEntries(projectID, entries, sourceLanguage, translationProvider)
}

func UpdateProjectSourceLanguage(projectID uint, newSourceLanguage string) error {
	return translationServiceInstance.UpdateProjectSourceLanguage(projectID, newSourceLanguage)
}