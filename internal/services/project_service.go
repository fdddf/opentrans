package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/model"
	"gorm.io/gorm"
)

// CreateProject creates a new project with the provided xcstrings content
func CreateProject(userID uint, name, description, fileName, fileContent, sourceLanguage string) (*database.Project, error) {
	// Parse the xcstrings content to extract structure
	var xcstrings model.XCStrings
	if err := json.Unmarshal([]byte(fileContent), &xcstrings); err != nil {
		return nil, fmt.Errorf("failed to parse xcstrings content: %v", err)
	}

	// Convert the model to the structure format we'll store
	contentStructure := map[string]interface{}{
		"sourceLanguage": xcstrings.SourceLanguage,
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

	result := database.DB.Create(project)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create project: %v", result.Error)
	}

	return project, nil
}

// GetProject retrieves a project by ID
func GetProject(projectID uint) (*database.Project, error) {
	var project database.Project
	result := database.DB.Preload("User").First(&project, projectID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("project not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &project, nil
}

// GetProjectsByUser retrieves all projects for a user
func GetProjectsByUser(userID uint) ([]database.Project, error) {
	var projects []database.Project
	result := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&projects)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve projects: %v", result.Error)
	}

	return projects, nil
}

// UpdateProject updates an existing project
func UpdateProject(projectID uint, updates map[string]interface{}) error {
	result := database.DB.Model(&database.Project{}).Where("id = ?", projectID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update project: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}

// DeleteProject soft deletes a project
func DeleteProject(projectID uint) error {
	result := database.DB.Delete(&database.Project{}, projectID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}

// ParseProjectContent parses the stored content structure back to XCStrings model
func ParseProjectContent(project *database.Project) (*model.XCStrings, error) {
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
func SaveProjectContent(projectID uint, xcstrings *model.XCStrings) error {
	contentStructure := map[string]interface{}{
		"sourceLanguage": xcstrings.SourceLanguage,
		"strings":        xcstrings.Strings,
		"version":        xcstrings.Version,
	}

	return UpdateProject(projectID, map[string]interface{}{
		"ContentStructure": contentStructure,
		"SourceLanguage":   xcstrings.SourceLanguage,
	})
}
