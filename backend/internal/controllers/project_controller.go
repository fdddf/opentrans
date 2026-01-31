package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	appcontext "github.com/fdddf/xcstrings-translator/internal/context"
	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/model"
	"github.com/fdddf/xcstrings-translator/internal/services"
	"github.com/fdddf/xcstrings-translator/internal/translator"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

// ProjectController handles project-related requests
type ProjectController struct{}

// NewProjectController creates a new ProjectController
func NewProjectController() *ProjectController {
	return &ProjectController{}
}

// CreateProjectRequest represents the create project request
type CreateProjectRequest struct {
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	FileName         string                 `json:"fileName"`
	FileContent      string                 `json:"fileContent"`
	SourceLanguage   string                 `json:"sourceLanguage"`
	ContentStructure map[string]interface{} `json:"contentStructure"`
	Settings         map[string]interface{} `json:"settings"`
}

// UpdateProjectRequest represents the update project request
type UpdateProjectRequest struct {
	Name             *string                `json:"name"`
	Description      *string                `json:"description"`
	SourceLanguage   *string                `json:"sourceLanguage"`
	ContentStructure map[string]interface{} `json:"contentStructure"`
	Settings         map[string]interface{} `json:"settings"`
}

// GetProjects retrieves all projects for the authenticated user
func (ctrl *ProjectController) GetProjects(c *fiber.Ctx) error {
	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	projects, err := services.GetProjectsByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"projects": projects,
	})
}

// CreateProject creates a new project
func (ctrl *ProjectController) CreateProject(c *fiber.Ctx) error {
	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.SourceLanguage != "" {
		if err := services.ValidateLanguages([]string{req.SourceLanguage}); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	project, err := services.CreateProject(userID, req.Name, req.Description, req.FileName, req.FileContent, req.SourceLanguage)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"project": project,
	})
}

// GetProject retrieves a single project by ID
func (ctrl *ProjectController) GetProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Verify user owns this project
	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"project": project,
	})
}

// UpdateProject updates an existing project
func (ctrl *ProjectController) UpdateProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["Name"] = *req.Name
	}
	if req.Description != nil {
		updates["Description"] = *req.Description
	}
	if req.SourceLanguage != nil {
		updates["SourceLanguage"] = *req.SourceLanguage
	}
	if req.ContentStructure != nil {
		updates["ContentStructure"] = req.ContentStructure
	}
	if req.Settings != nil {
		updates["Settings"] = req.Settings
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No fields to update")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	err = services.UpdateProject(uint(projectID), updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Project updated successfully",
	})
}

// DeleteProject deletes a project
func (ctrl *ProjectController) DeleteProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	err = services.DeleteProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Project deleted successfully",
	})
}

// UploadToProject uploads a file to a project
func (ctrl *ProjectController) UploadToProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to read file: %v", err))
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to read file: %v", err))
	}

	if len(data) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Empty file")
	}

	var sourceLanguage string
	if c.FormValue("sourceLanguage") != "" {
		sourceLanguage = c.FormValue("sourceLanguage")
	}

	var xcstrings model.XCStrings
	if err := json.Unmarshal(data, &xcstrings); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid xcstrings: %v", err))
	}

	if sourceLanguage == "" {
		sourceLanguage = xcstrings.SourceLanguage
	}

	contentStructure := map[string]interface{}{
		"sourceLanguage": xcstrings.SourceLanguage,
		"strings":        xcstrings.Strings,
		"version":        xcstrings.Version,
	}

	err = services.UpdateProject(uint(projectID), map[string]interface{}{
		"FileContent":      string(data),
		"FileName":         fileHeader.Filename,
		"SourceLanguage":   sourceLanguage,
		"ContentStructure": contentStructure,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"project": project,
	})
}

// TranslateProject translates a project
func (ctrl *ProjectController) TranslateProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	var req TranslateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := services.ValidateLanguages(req.TargetLanguages); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if len(req.TargetLanguages) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Target languages are required")
	}

	// Get the project content
	xcstrings, err := services.ParseProjectContent(project)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to parse project content: %v", err))
	}

	if req.SourceLanguage != "" {
		xcstrings.SourceLanguage = req.SourceLanguage
	}

	requests := translator.CreateTranslationRequests(xcstrings, req.TargetLanguages)
	job := &Job{
		ID:        uuid.NewString(),
		Status:    "running",
		Done:      0,
		Total:     len(requests),
		UpdatedAt: time.Now(),
	}

	// If nothing to do, finish immediately.
	if len(requests) == 0 {
		return c.JSON(fiber.Map{
			"jobId": job.ID,
			"done":  true,
		})
	}

	// Run translation in background
	go ctrl.runProjectTranslation(job, project, xcstrings, req)

	return c.JSON(fiber.Map{
		"jobId": job.ID,
	})
}

// runProjectTranslation runs the translation process for a project
func (ctrl *ProjectController) runProjectTranslation(job *Job, project *database.Project, xcstrings *model.XCStrings, req TranslateRequest) {
	provider, err := ctrl.buildProvider(strings.ToLower(req.Provider), req.Config)
	if err != nil {
		return
	}

	concurrency := req.Concurrency
	if concurrency <= 0 {
		concurrency = 4
	}

	timeout := time.Duration(req.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 300 * time.Second
	}

	service := translator.NewTranslationService(provider, concurrency, timeout)
	ctx := context.Background()

	progressBuilder := func(target string, total int) translator.ProgressReporter {
		return func(done, total int, resp model.TranslationResponse) {
			if resp.Error == nil {
				translator.ApplyTranslations(xcstrings, []model.TranslationResponse{resp})
			}
		}
	}

	_, translateErr := translator.TranslatePerLanguage(ctx, xcstrings, req.TargetLanguages, service, progressBuilder)
	if translateErr != nil {
		return
	}

	// Update the project in the database with the new translations
	err = services.UpdateTranslationsFromXCStrings(project.ID, xcstrings, req.Provider)
	if err != nil {
		return
	}

	// Update the project content structure
	err = services.SaveProjectContent(project.ID, xcstrings)
	if err != nil {
		return
	}
}

// ExportProject exports a project as xcstrings file
func (ctrl *ProjectController) ExportProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	// Get project content and apply stored translations
	xcstrings, err := services.ParseProjectContent(project)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to parse project content: %v", err))
	}

	// Apply stored translations to the XCStrings model
	err = services.ApplyTranslationsToXCStrings(project.ID, xcstrings)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to apply translations: %v", err))
	}

	data, err := model.MarshalXCStrings(xcstrings)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileName := project.FileName
	if fileName == "" {
		fileName = "Localizable_translated.xcstrings"
	}
	c.Attachment(fileName)
	c.Set("Content-Type", "application/json")
	return c.Send(data)
}

// GetProjectTranslations retrieves translations for a project
func (ctrl *ProjectController) GetProjectTranslations(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	translations, err := services.GetTranslationsByProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"translations": translations,
	})
}

// GetMissingTranslations retrieves missing translations for a project
func (ctrl *ProjectController) GetMissingTranslations(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	var req struct {
		TargetLanguages []string `json:"targetLanguages"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if len(req.TargetLanguages) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Target languages are required")
	}
	if err := services.ValidateLanguages(req.TargetLanguages); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	missingKeys, err := services.GetMissingTranslations(uint(projectID), req.TargetLanguages)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"missingKeys": missingKeys,
	})
}

// GetTranslationStatus retrieves translation status for a project
func (ctrl *ProjectController) GetTranslationStatus(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	var req struct {
		TargetLanguages []string `json:"targetLanguages"`
	}

	if err := c.QueryParser(&req); err != nil {
		_ = c.BodyParser(&req)
	}

	if len(req.TargetLanguages) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Target languages are required")
	}
	if err := services.ValidateLanguages(req.TargetLanguages); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	stats, err := services.GetTranslationStatus(uint(projectID), req.TargetLanguages)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"stats":   stats,
	})
}

// GetProjectStats retrieves statistics for a project
func (ctrl *ProjectController) GetProjectStats(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	stats, err := services.GetProjectStats(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"stats":   stats,
	})
}

// UpdateSingleTranslation updates a single translation
func (ctrl *ProjectController) UpdateSingleTranslation(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	key := c.Params("key")
	language := c.Params("language")

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	var req struct {
		TargetText string `json:"targetText"`
		State      string `json:"state"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.TargetText == "" {
		return fiber.NewError(fiber.StatusBadRequest, "targetText is required")
	}

	if req.State == "" {
		req.State = "translated"
	}

	translation, err := services.UpdateSingleTranslation(uint(projectID), key, language, req.TargetText, req.State)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"translation": translation,
	})
}

// BulkUpdateTranslations bulk updates translations
func (ctrl *ProjectController) BulkUpdateTranslations(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	var req struct {
		Updates []map[string]interface{} `json:"updates"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if len(req.Updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "updates are required")
	}

	err = services.BulkUpdateTranslations(uint(projectID), req.Updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Translations updated successfully",
	})
}

// GetProjectLanguages retrieves languages for a project
func (ctrl *ProjectController) GetProjectLanguages(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := appcontext.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user owns this project
	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	if project.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
	}

	languages, err := services.GetProjectLanguages(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
}

// buildProvider builds a translation provider
func (ctrl *ProjectController) buildProvider(provider string, cfg ProviderConfig) (model.TranslationProvider, error) {
	switch provider {
	case "google":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("apiKey required for Google provider")
		}
		return translator.NewGoogleTranslator(cfg.APIKey), nil
	case "deepl":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("apiKey required for DeepL provider")
		}
		return translator.NewDeepLTranslator(cfg.APIKey, cfg.IsFree), nil
	case "baidu":
		if cfg.AppID == "" || cfg.AppSecret == "" {
			return nil, fmt.Errorf("appId and appSecret are required for Baidu provider")
		}
		return translator.NewBaiduTranslator(cfg.AppID, cfg.AppSecret), nil
	default:
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("apiKey required for OpenAI provider")
		}
		return translator.NewOpenAITranslator(cfg.APIKey, cfg.APIBaseURL, cfg.Model, cfg.Temperature, cfg.MaxTokens), nil
	}
}