package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/model"
	"gorm.io/gorm"
)

// ProjectService handles project-related operations
type ProjectService struct {
	DB                     *database.Database
	TranslationService     *TranslationService
	AppLocalizationService *AppLocalizationService
}

// CreateProject creates a new project with the provided xcstrings content
func (s *ProjectService) CreateProject(userID uint, name, description, fileName, fileContent, sourceLanguage string) (*database.Project, error) {
	// Parse the xcstrings content to extract structure
	var xcstrings model.XCStrings
	if err := json.Unmarshal([]byte(fileContent), &xcstrings); err != nil {
		return nil, fmt.Errorf("failed to parse xcstrings content: %v", err)
	}

	// Set source language from xcstrings if not provided
	if sourceLanguage == "" {
		sourceLanguage = xcstrings.SourceLanguage
	}

	// Convert the model to the structure format we'll store
	contentStructure := map[string]interface{}{
		"sourceLanguage": sourceLanguage,
		"strings":        xcstrings.Strings,
		"version":        xcstrings.Version,
	}

	project := &database.Project{
		Name:             name,
		Description:      description,
		UserID:           userID,
		FileContent:      fileContent,
		FileName:         fileName,
		SourceLanguage:   sourceLanguage,
		ContentStructure: contentStructure,
		Settings:         make(map[string]interface{}),
	}

	result := s.DB.Create(project)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create project: %v", result.Error)
	}

	// Create initial translations from the uploaded file
	err := s.TranslationService.BulkCreateTranslationsFromEntries(project.ID, xcstrings.Strings, sourceLanguage, "manual")
	if err != nil {
		// Log error but don't fail the project creation
		fmt.Printf("Warning: failed to create initial translations: %v\n", err)
	}

	return project, nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(projectID uint) (*database.Project, error) {
	var project database.Project
	result := s.DB.Preload("User").First(&project, projectID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("project not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &project, nil
}

// GetProjectsByUser retrieves all projects for a user
func (s *ProjectService) GetProjectsByUser(userID uint) ([]database.Project, error) {
	var projects []database.Project
	result := s.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&projects)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve projects: %v", result.Error)
	}

	return projects, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(projectID uint, updates map[string]interface{}) error {
	result := s.DB.Model(&database.Project{}).Where("id = ?", projectID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update project: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}

// DeleteProject soft deletes a project
func (s *ProjectService) DeleteProject(projectID uint) error {
	result := s.DB.Delete(&database.Project{}, projectID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}

// ParseProjectContent parses the stored content structure back to XCStrings model
func (s *ProjectService) ParseProjectContent(project *database.Project) (*model.XCStrings, error) {
	content, ok := project.ContentStructure["strings"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid content structure")
	}

	// Convert the stored structure back to the model
	stringsMap := make(map[string]model.StringEntry)
	for key, value := range content {
		// Convert interface{} back to StringEntry
		entryBytes, err := json.Marshal(value)
		if err != nil {
			continue
		}

		var entry model.StringEntry
		if err := json.Unmarshal(entryBytes, &entry); err != nil {
			continue
		}

		stringsMap[key] = entry
	}

	xcstrings := &model.XCStrings{
		SourceLanguage: project.SourceLanguage,
		Strings:        stringsMap,
		Version:        "1.0", // Default version
	}

	// Set version if available in content structure
	if version, ok := project.ContentStructure["version"].(string); ok {
		xcstrings.Version = version
	}

	return xcstrings, nil
}

// SaveProjectContent updates the project's content structure
func (s *ProjectService) SaveProjectContent(projectID uint, xcstrings *model.XCStrings) error {
	contentStructure := map[string]interface{}{
		"sourceLanguage": xcstrings.SourceLanguage,
		"strings":        xcstrings.Strings,
		"version":        xcstrings.Version,
	}

	return s.UpdateProject(projectID, map[string]interface{}{
		"ContentStructure": contentStructure,
		"SourceLanguage":   xcstrings.SourceLanguage,
	})
}

// UpdateTranslationsFromXCStrings updates translations from XCStrings model
func (s *ProjectService) UpdateTranslationsFromXCStrings(projectID uint, xcstrings *model.XCStrings, provider string) error {
	// Delete existing translations for this project
	err := s.DB.Where("project_id = ?", projectID).Delete(&database.Translation{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete existing translations: %v", err)
	}

	// Insert new translations
	for key, entry := range xcstrings.Strings {
		for lang, loc := range entry.Localizations {
			if lang == xcstrings.SourceLanguage {
				continue // Skip source language
			}

			translation := &database.Translation{
				ProjectID:           projectID,
				Key:                 key,
				SourceText:          entry.Localizations[xcstrings.SourceLanguage].StringUnit.Value,
				TargetText:          loc.StringUnit.Value,
				TargetLanguage:      lang,
				State:               loc.StringUnit.State,
				TranslationProvider: provider,
			}

			if err := s.DB.Create(translation).Error; err != nil {
				return fmt.Errorf("failed to create translation: %v", err)
			}
		}
	}

	return nil
}

// GetProjectTranslations retrieves all translations for a project
func (s *ProjectService) GetProjectTranslations(projectID uint) ([]database.Translation, error) {
	var translations []database.Translation
	result := s.DB.Where("project_id = ?", projectID).Find(&translations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve translations: %v", result.Error)
	}

	return translations, nil
}

// GetMissingTranslations retrieves missing translations for a project
func (s *ProjectService) GetMissingTranslations(projectID uint, targetLanguages []string) ([]string, error) {
	project, err := s.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	xcstrings, err := s.ParseProjectContent(project)
	if err != nil {
		return nil, err
	}

	var missingKeys []string
	for key, entry := range xcstrings.Strings {
		for _, targetLang := range targetLanguages {
			if _, ok := entry.Localizations[targetLang]; !ok {
				missingKeys = append(missingKeys, key)
				break // Only add once per key
			}
		}
	}

	return missingKeys, nil
}

// GetTranslationStatus returns counts of translated/untranslated per language
func (s *ProjectService) GetTranslationStatus(projectID uint, targetLanguages []string) (map[string]map[string]int, error) {
	project, err := s.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	xcstrings, err := s.ParseProjectContent(project)
	if err != nil {
		return nil, err
	}

	translations, err := s.TranslationService.GetTranslationsByProject(projectID)
	if err != nil {
		return nil, err
	}

	languageStats := make(map[string]map[string]int)

	for _, lang := range targetLanguages {
		languageStats[lang] = map[string]int{"translated": 0, "untranslated": 0}
	}

	for key := range xcstrings.Strings {
		for _, lang := range targetLanguages {
			var translatedText string
			for _, translation := range translations {
				if translation.Key == key && translation.TargetLanguage == lang {
					translatedText = translation.TargetText
					break
				}
			}

			if translatedText != "" {
				languageStats[lang]["translated"]++
			} else {
				languageStats[lang]["untranslated"]++
			}
		}
	}

	return languageStats, nil
}

// GetProjectStats returns statistics for a project
func (s *ProjectService) GetProjectStats(projectID uint) (map[string]interface{}, error) {
	project, err := s.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	xcstrings, err := s.ParseProjectContent(project)
	if err != nil {
		return nil, err
	}

	// Get all translations
	translations, err := s.TranslationService.GetTranslationsByProject(projectID)
	if err != nil {
		return nil, err
	}

	// Count total strings
	totalStrings := len(xcstrings.Strings)

	// Count languages
	languages := make(map[string]bool)
	languages[project.SourceLanguage] = true
	for _, translation := range translations {
		languages[translation.TargetLanguage] = true
	}

	// Per-language stats
	languageStats := make(map[string]map[string]int)

	// Count translated vs untranslated
	translatedCount := 0
	untranslatedCount := 0
	for key := range xcstrings.Strings {
		hasTranslation := false
		for targetLang := range languages {
			if _, ok := languageStats[targetLang]; !ok {
				languageStats[targetLang] = map[string]int{"translated": 0, "untranslated": 0}
			}

			var translatedText string
			for _, translation := range translations {
				if translation.Key == key && translation.TargetLanguage == targetLang {
					translatedText = translation.TargetText
					break
				}
			}

			if translatedText != "" {
				languageStats[targetLang]["translated"]++
				hasTranslation = true
			} else {
				languageStats[targetLang]["untranslated"]++
			}
		}

		if hasTranslation {
			translatedCount++
		} else {
			untranslatedCount++
		}
	}

	stats := map[string]interface{}{
		"totalStrings":      totalStrings,
		"totalLanguages":    len(languages),
		"translatedCount":   translatedCount,
		"untranslatedCount": untranslatedCount,
		"perLanguage":       languageStats,
		"languages":         languages,
		"sourceLanguage":    project.SourceLanguage,
	}

	return stats, nil
}

// UpdateSingleTranslation manually updates a single translation
func (s *ProjectService) UpdateSingleTranslation(projectID uint, key, targetLanguage, targetText, state string) (*database.Translation, error) {
	// Check if translation exists
	var translation database.Translation
	result := s.DB.Where("project_id = ? AND key = ? AND target_language = ?", projectID, key, targetLanguage).First(&translation)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new translation
		project, err := s.GetProject(projectID)
		if err != nil {
			return nil, err
		}

		// Get source text from project content
		xcstrings, err := s.ParseProjectContent(project)
		if err != nil {
			return nil, err
		}

		sourceText := ""
		if entry, exists := xcstrings.Strings[key]; exists {
			if sourceLoc, exists := entry.Localizations[project.SourceLanguage]; exists {
				sourceText = sourceLoc.StringUnit.Value
			}
		}

		if sourceText == "" {
			sourceText = key
		}

		translation = database.Translation{
			ProjectID:           projectID,
			Key:                 key,
			SourceText:          sourceText,
			TargetText:          targetText,
			TargetLanguage:      targetLanguage,
			State:               state,
			TranslationProvider: "manual",
		}

		result = s.DB.Create(&translation)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to create translation: %v", result.Error)
		}
	} else if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	} else {
		// Update existing translation
		result = s.DB.Model(&translation).Updates(map[string]interface{}{
			"TargetText": targetText,
			"State":      state,
		})
		if result.Error != nil {
			return nil, fmt.Errorf("failed to update translation: %v", result.Error)
		}
	}

	// Update project content structure
	project, err := s.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	xcstrings, err := s.ParseProjectContent(project)
	if err != nil {
		return nil, err
	}

	// Update the localization in the xcstrings structure
	if entry, exists := xcstrings.Strings[key]; exists {
		if entry.Localizations == nil {
			entry.Localizations = make(map[string]model.Localization)
		}
		entry.Localizations[targetLanguage] = model.Localization{
			StringUnit: model.StringUnit{
				Value: targetText,
				State: state,
			},
		}
		xcstrings.Strings[key] = entry

		// Save updated content
		err = s.SaveProjectContent(projectID, xcstrings)
		if err != nil {
			return nil, err
		}
	}

	return &translation, nil
}

// BulkUpdateTranslations updates multiple translations at once
func (s *ProjectService) BulkUpdateTranslations(projectID uint, updates []map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	// Use a transaction
	tx := s.DB.Begin()
	defer tx.Rollback()

	for _, update := range updates {
		key, ok := update["key"].(string)
		if !ok {
			continue
		}

		targetLanguage, ok := update["targetLanguage"].(string)
		if !ok {
			continue
		}

		targetText, ok := update["targetText"].(string)
		if !ok {
			continue
		}

		state := "translated"
		if s, ok := update["state"].(string); ok {
			state = s
		}

		// Update or create translation
		var translation database.Translation
		result := tx.Where("project_id = ? AND key = ? AND target_language = ?", projectID, key, targetLanguage).First(&translation)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Get source text
			project, err := s.GetProject(projectID)
			if err != nil {
				return err
			}

			xcstrings, err := s.ParseProjectContent(project)
			if err != nil {
				return err
			}

			sourceText := ""
			if entry, exists := xcstrings.Strings[key]; exists {
				if sourceLoc, exists := entry.Localizations[project.SourceLanguage]; exists {
					sourceText = sourceLoc.StringUnit.Value
				}
			}

			if sourceText == "" {
				sourceText = key
			}

			translation = database.Translation{
				ProjectID:           projectID,
				Key:                 key,
				SourceText:          sourceText,
				TargetText:          targetText,
				TargetLanguage:      targetLanguage,
				State:               state,
				TranslationProvider: "manual",
			}

			result = tx.Create(&translation)
			if result.Error != nil {
				return fmt.Errorf("failed to create translation: %v", result.Error)
			}
		} else if result.Error != nil {
			return fmt.Errorf("database error: %v", result.Error)
		} else {
			// Update existing
			result = tx.Model(&translation).Updates(map[string]interface{}{
				"TargetText": targetText,
				"State":      state,
			})
			if result.Error != nil {
				return fmt.Errorf("failed to update translation: %v", result.Error)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit bulk update: %v", err)
	}

	// Update project content structure
	project, err := s.GetProject(projectID)
	if err != nil {
		return err
	}

	xcstrings, err := s.ParseProjectContent(project)
	if err != nil {
		return err
	}

	// Apply all updates to xcstrings
	for _, update := range updates {
		key, _ := update["key"].(string)
		targetLanguage, _ := update["targetLanguage"].(string)
		targetText, _ := update["targetText"].(string)
		state := "translated"
		if s, ok := update["state"].(string); ok {
			state = s
		}

		if entry, exists := xcstrings.Strings[key]; exists {
			if entry.Localizations == nil {
				entry.Localizations = make(map[string]model.Localization)
			}
			entry.Localizations[targetLanguage] = model.Localization{
				StringUnit: model.StringUnit{
					Value: targetText,
					State: state,
				},
			}
			xcstrings.Strings[key] = entry
		}
	}

	// Save updated content
	err = s.SaveProjectContent(projectID, xcstrings)
	if err != nil {
		return err
	}

	return nil
}

// ExportProjectAsXCStrings exports the project as a complete xcstrings file
func (s *ProjectService) ExportProjectAsXCStrings(projectID uint) ([]byte, error) {
	project, err := s.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	// Parse project content
	xcstrings, err := s.ParseProjectContent(project)
	if err != nil {
		return nil, err
	}

	// Apply all stored translations to the xcstrings
	err = s.TranslationService.ApplyTranslationsToXCStrings(projectID, xcstrings)
	if err != nil {
		return nil, err
	}

	// Marshal to JSON
	data, err := model.MarshalXCStrings(xcstrings)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetProjectLanguages returns all languages available in a project
func (s *ProjectService) GetProjectLanguages(projectID uint) ([]string, error) {
	project, err := s.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	xcstrings, err := s.ParseProjectContent(project)
	if err != nil {
		return nil, err
	}

	// Collect languages from xcstrings
	languageSet := make(map[string]bool)
	languageSet[project.SourceLanguage] = true

	for _, entry := range xcstrings.Strings {
		for lang := range entry.Localizations {
			languageSet[lang] = true
		}
	}

	// Convert to slice
	languages := make([]string, 0, len(languageSet))
	for lang := range languageSet {
		languages = append(languages, lang)
	}

	return languages, nil
}

// Global functions for backward compatibility
var projectServiceInstance *ProjectService

func SetProjectService(db *database.Database, translationService *TranslationService, appLocalizationService *AppLocalizationService) {
	projectServiceInstance = &ProjectService{
		DB:                     db,
		TranslationService:     translationService,
		AppLocalizationService: appLocalizationService,
	}
}

func CreateProject(userID uint, name, description, fileName, fileContent, sourceLanguage string) (*database.Project, error) {
	return projectServiceInstance.CreateProject(userID, name, description, fileName, fileContent, sourceLanguage)
}

func GetProject(projectID uint) (*database.Project, error) {
	return projectServiceInstance.GetProject(projectID)
}

func GetProjectsByUser(userID uint) ([]database.Project, error) {
	return projectServiceInstance.GetProjectsByUser(userID)
}

func UpdateProject(projectID uint, updates map[string]interface{}) error {
	return projectServiceInstance.UpdateProject(projectID, updates)
}

func DeleteProject(projectID uint) error {
	return projectServiceInstance.DeleteProject(projectID)
}

func ParseProjectContent(project *database.Project) (*model.XCStrings, error) {
	return projectServiceInstance.ParseProjectContent(project)
}

func SaveProjectContent(projectID uint, xcstrings *model.XCStrings) error {
	return projectServiceInstance.SaveProjectContent(projectID, xcstrings)
}

func UpdateTranslationsFromXCStrings(projectID uint, xcstrings *model.XCStrings, provider string) error {
	return projectServiceInstance.UpdateTranslationsFromXCStrings(projectID, xcstrings, provider)
}

func GetProjectTranslations(projectID uint) ([]database.Translation, error) {
	return projectServiceInstance.GetProjectTranslations(projectID)
}

func GetMissingTranslations(projectID uint, targetLanguages []string) ([]string, error) {
	return projectServiceInstance.GetMissingTranslations(projectID, targetLanguages)
}

// GetTranslationStatus returns per-language status counts
func GetTranslationStatus(projectID uint, targetLanguages []string) (map[string]map[string]int, error) {
	return projectServiceInstance.GetTranslationStatus(projectID, targetLanguages)
}

func GetProjectStats(projectID uint) (map[string]interface{}, error) {
	return projectServiceInstance.GetProjectStats(projectID)
}

func UpdateSingleTranslation(projectID uint, key, targetLanguage, targetText, state string) (*database.Translation, error) {
	return projectServiceInstance.UpdateSingleTranslation(projectID, key, targetLanguage, targetText, state)
}

func BulkUpdateTranslations(projectID uint, updates []map[string]interface{}) error {
	return projectServiceInstance.BulkUpdateTranslations(projectID, updates)
}

func ExportProjectAsXCStrings(projectID uint) ([]byte, error) {
	return projectServiceInstance.ExportProjectAsXCStrings(projectID)
}

func GetProjectLanguages(projectID uint) ([]string, error) {
	return projectServiceInstance.GetProjectLanguages(projectID)
}
