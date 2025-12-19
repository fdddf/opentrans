package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fdddf/xcstrings-translator/internal/auth"
	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/model"
	"github.com/fdddf/xcstrings-translator/internal/services"
	"github.com/fdddf/xcstrings-translator/internal/translator"
	"github.com/fdddf/xcstrings-translator/webui"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

// Payload represents the data returned to the UI.
type Payload struct {
	FileName           string           `json:"fileName"`
	SourceLanguage     string           `json:"sourceLanguage"`
	AvailableLanguages []string         `json:"availableLanguages"`
	TotalStrings       int              `json:"totalStrings"`
	Entries            []UILocalization `json:"entries"`
	Warning            string           `json:"warning,omitempty"`
}

// UILocalization is a flattened view for the table UI.
type UILocalization struct {
	Key          string            `json:"key"`
	Source       string            `json:"source"`
	State        string            `json:"state"`
	Translations map[string]string `json:"translations"`
	Missing      []string          `json:"missing"`
}

// TranslateRequest describes the batch translate payload from the UI.
type TranslateRequest struct {
	Provider        string         `json:"provider"`
	TargetLanguages []string       `json:"targetLanguages"`
	SourceLanguage  string         `json:"sourceLanguage"`
	Concurrency     int            `json:"concurrency"`
	TimeoutSeconds  int            `json:"timeoutSeconds"`
	Config          ProviderConfig `json:"config"`
}

// ProviderConfig is the union of provider-specific options we support.
type ProviderConfig struct {
	APIKey      string  `json:"apiKey"`
	APIBaseURL  string  `json:"apiBaseUrl"`
	Model       string  `json:"model"`
	Glossary    string  `json:"glossary"`
	AppID       string  `json:"appId"`
	AppSecret   string  `json:"appSecret"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
	Formality   string  `json:"formality"`
	IsFree      bool    `json:"isFree"`
}

// ServerState holds the in-memory working copy of the xcstrings data.
type ServerState struct {
	mu              sync.RWMutex
	fileName        string
	xcstrings       *model.XCStrings
	targetLanguages []string
	job             *Job
}

// Job tracks long-running translation progress.
type Job struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"` // running, done, error
	Message   string    `json:"message,omitempty"`
	Done      int       `json:"done"`
	Total     int       `json:"total"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Serve starts the Fiber server using the embedded UI assets.
func Serve(addr string) error {
	app, err := NewApp()
	if err != nil {
		return err
	}
	return app.Listen(addr)
}

// ServeWithListener serves the app using a pre-bound listener (useful for GUI mode).
func ServeWithListener(ln net.Listener) (*fiber.App, <-chan error, error) {
	app, err := NewApp()
	if err != nil {
		return nil, nil, err
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Listener(ln)
	}()
	return app, errCh, nil
}

// NewApp constructs the Fiber app with embedded assets without starting the server.
func NewApp() (*fiber.App, error) {
	// Initialize the database
	err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	distFS, err := fs.Sub(webui.EmbeddedFS, "dist")
	if err != nil {
		return nil, fmt.Errorf("embedded UI missing: %w", err)
	}

	state := &ServerState{}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// Public routes (no authentication required)
	api := app.Group("/api")
	api.Post("/upload", state.handleUpload)
	api.Get("/strings", state.handleStrings)
	api.Post("/translate", state.handleTranslate)
	api.Get("/progress", state.handleProgress)
	api.Get("/export", state.handleExport)

	// Authentication routes
	api.Post("/auth/register", handleRegister)
	api.Post("/auth/login", handleLogin)
	api.Post("/auth/logout", handleLogout)

	// Protected routes (require authentication)
	protected := api.Group("/protected")
	protected.Use(AuthMiddleware)
	protected.Get("/projects", handleGetProjects)
	protected.Post("/projects", handleCreateProject)
	protected.Get("/projects/:id", handleGetProject)
	protected.Put("/projects/:id", handleUpdateProject)
	protected.Delete("/projects/:id", handleDeleteProject)
	protected.Post("/projects/:id/upload", handleUploadToProject)
	protected.Post("/projects/:id/translate", handleTranslateProject)
	protected.Get("/projects/:id/export", handleExportProject)
	protected.Get("/projects/:id/translations", handleGetProjectTranslations)
	protected.Get("/projects/:id/missing-translations", handleGetMissingTranslations)

	// Provider configuration routes
	protected.Get("/providers", handleGetProviderConfigs)
	protected.Post("/providers", handleCreateProviderConfig)
	protected.Get("/providers/:id", handleGetProviderConfig)
	protected.Put("/providers/:id", handleUpdateProviderConfig)
	protected.Delete("/providers/:id", handleDeleteProviderConfig)
	protected.Get("/providers/:type/default", handleGetDefaultProviderConfig)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(distFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	return app, nil
}

func (s *ServerState) handleUpload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to read file: %v", err))
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to read file: %v", err))
	}

	if len(data) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "empty file")
	}

	xcstrings, err := model.ParseXCStrings(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid xcstrings: %v", err))
	}

	if source := c.FormValue("sourceLanguage"); source != "" {
		xcstrings.SourceLanguage = source
	}

	s.mu.Lock()
	s.xcstrings = xcstrings
	s.fileName = fileHeader.Filename
	s.targetLanguages = nil
	s.mu.Unlock()

	payload := s.buildPayload(nil)
	return c.JSON(payload)
}

func (s *ServerState) handleProgress(c *fiber.Ctx) error {
	s.mu.RLock()
	job := s.job
	payload := s.buildPayload(nil)
	s.mu.RUnlock()

	return c.JSON(fiber.Map{
		"job":     job,
		"payload": payload,
	})
}

func (s *ServerState) handleStrings(c *fiber.Ctx) error {
	payload := s.buildPayload(nil)
	if payload == nil {
		return fiber.NewError(fiber.StatusNotFound, "no xcstrings loaded")
	}
	return c.JSON(payload)
}

func (s *ServerState) handleExport(c *fiber.Ctx) error {
	s.mu.RLock()
	xc := s.xcstrings
	name := s.fileName
	s.mu.RUnlock()

	if xc == nil {
		return fiber.NewError(fiber.StatusNotFound, "no xcstrings loaded")
	}

	data, err := model.MarshalXCStrings(xc)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileName := name
	if fileName == "" {
		fileName = "Localizable_translated.xcstrings"
	}
	c.Attachment(fileName)
	c.Set("Content-Type", "application/json")
	return c.Send(data)
}

func (s *ServerState) handleTranslate(c *fiber.Ctx) error {
	var req TranslateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if len(req.TargetLanguages) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "targetLanguages is required")
	}

	s.mu.RLock()
	xc := s.xcstrings
	s.mu.RUnlock()

	if xc == nil {
		return fiber.NewError(fiber.StatusBadRequest, "upload a xcstrings file first")
	}

	if req.SourceLanguage != "" {
		xc.SourceLanguage = req.SourceLanguage
	}

	requests := translator.CreateTranslationRequests(xc, req.TargetLanguages)
	job := s.startJob(len(requests))

	// If nothing to do, finish immediately.
	if len(requests) == 0 {
		s.finishJob("done", "")
		return c.JSON(fiber.Map{"jobId": job.ID})
	}

	go s.runTranslation(job, xc, req)

	return c.JSON(fiber.Map{"jobId": job.ID})
}

func (s *ServerState) buildPayload(targets []string) *Payload {
	s.mu.RLock()
	xc := s.xcstrings
	name := s.fileName
	rememberedTargets := s.targetLanguages
	s.mu.RUnlock()

	if xc == nil {
		return nil
	}

	languages := collectLanguages(xc)

	targetSet := dedupe(targets)
	if len(targetSet) == 0 {
		targetSet = rememberedTargets
	}
	if len(targetSet) == 0 {
		for _, lang := range languages {
			if lang != xc.SourceLanguage {
				targetSet = append(targetSet, lang)
			}
		}
		targetSet = dedupe(targetSet)
	}

	entries := flattenEntries(xc, targetSet)

	return &Payload{
		FileName:           name,
		SourceLanguage:     xc.SourceLanguage,
		AvailableLanguages: languages,
		TotalStrings:       len(xc.Strings),
		Entries:            entries,
	}
}

func flattenEntries(xc *model.XCStrings, targets []string) []UILocalization {
	keys := make([]string, 0, len(xc.Strings))
	for key := range xc.Strings {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	entries := make([]UILocalization, 0, len(keys))
	for _, key := range keys {
		entry := xc.Strings[key]
		translations := make(map[string]string)
		for lang, loc := range entry.Localizations {
			translations[lang] = loc.StringUnit.Value
		}

		sourceText := translations[xc.SourceLanguage]
		if sourceText == "" {
			sourceText = key
		}

		state := ""
		if sourceLoc, ok := entry.Localizations[xc.SourceLanguage]; ok {
			state = sourceLoc.StringUnit.State
		}

		missing := []string{}
		for _, target := range targets {
			if translations[target] == "" {
				missing = append(missing, target)
			}
		}

		entries = append(entries, UILocalization{
			Key:          key,
			Source:       sourceText,
			State:        state,
			Translations: translations,
			Missing:      missing,
		})
	}
	return entries
}

func collectLanguages(xc *model.XCStrings) []string {
	langSet := map[string]struct{}{}
	if xc.SourceLanguage != "" {
		langSet[xc.SourceLanguage] = struct{}{}
	}
	for _, entry := range xc.Strings {
		for lang := range entry.Localizations {
			langSet[lang] = struct{}{}
		}
	}

	langs := make([]string, 0, len(langSet))
	for lang := range langSet {
		langs = append(langs, lang)
	}
	sort.Strings(langs)
	return langs
}

func buildProvider(provider string, cfg ProviderConfig) (model.TranslationProvider, error) {
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

func dedupe(list []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, item := range list {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func (s *ServerState) startJob(total int) *Job {
	job := &Job{
		ID:        uuid.NewString(),
		Status:    "running",
		Done:      0,
		Total:     total,
		UpdatedAt: time.Now(),
	}
	s.mu.Lock()
	s.job = job
	s.mu.Unlock()
	return job
}

func (s *ServerState) incrementJob(delta int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.job == nil {
		return
	}
	s.job.Done += delta
	if s.job.Done > s.job.Total && s.job.Total > 0 {
		s.job.Done = s.job.Total
	}
	s.job.UpdatedAt = time.Now()
}

func (s *ServerState) finishJob(status, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.job == nil {
		return
	}
	s.job.Status = status
	s.job.Message = msg
	s.job.UpdatedAt = time.Now()
}

func (s *ServerState) runTranslation(job *Job, xc *model.XCStrings, req TranslateRequest) {
	provider, err := buildProvider(strings.ToLower(req.Provider), req.Config)
	if err != nil {
		s.finishJob("error", err.Error())
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
				s.applyResponse(resp)
			}
			s.incrementJob(1)
		}
	}

	responses, translateErr := translator.TranslatePerLanguage(ctx, xc, req.TargetLanguages, service, progressBuilder)
	translator.ApplyTranslations(xc, responses)

	if len(req.TargetLanguages) > 0 {
		s.mu.Lock()
		s.targetLanguages = dedupe(req.TargetLanguages)
		s.mu.Unlock()
	}

	if translateErr != nil {
		s.finishJob("error", translateErr.Error())
		return
	}

	s.finishJob("done", "")
}

// applyResponse applies a single successful translation response.
func (s *ServerState) applyResponse(resp model.TranslationResponse) {
	if resp.Error != nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	translator.ApplyTranslations(s.xcstrings, []model.TranslationResponse{resp})
}

// Authentication handlers
func handleRegister(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username, email, and password are required")
	}

	user, err := auth.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

func handleLogin(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, token, err := auth.LoginUser(req.Username, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
		"token":   token,
	})
}

func handleLogout(c *fiber.Ctx) error {
	// In a real application, you might implement token blacklisting
	// For now, just return success
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// Project handlers
func handleGetProjects(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
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

func handleCreateProject(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		Name             string                 `json:"name"`
		Description      string                 `json:"description"`
		FileName         string                 `json:"fileName"`
		FileContent      string                 `json:"fileContent"`
		SourceLanguage   string                 `json:"sourceLanguage"`
		ContentStructure map[string]interface{} `json:"contentStructure"`
		Settings         map[string]interface{} `json:"settings"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
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

func handleGetProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	project, err := services.GetProject(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Verify user owns this project
	userID, ok := GetUserIDFromContext(c)
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

func handleUpdateProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		Name             *string                `json:"name"`
		Description      *string                `json:"description"`
		SourceLanguage   *string                `json:"sourceLanguage"`
		ContentStructure map[string]interface{} `json:"contentStructure"`
		Settings         map[string]interface{} `json:"settings"`
	}

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

func handleDeleteProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := GetUserIDFromContext(c)
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

func handleUploadToProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := GetUserIDFromContext(c)
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

func handleTranslateProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := GetUserIDFromContext(c)
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
	go runProjectTranslation(job, project, xcstrings, req)

	return c.JSON(fiber.Map{
		"jobId": job.ID,
	})
}

// runProjectTranslation runs the translation process for a project
func runProjectTranslation(job *Job, project *database.Project, xcstrings *model.XCStrings, req TranslateRequest) {
	provider, err := buildProvider(strings.ToLower(req.Provider), req.Config)
	if err != nil {
		// Log error but there's no way to update the job status here
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
				// Update project content with translated response
				translator.ApplyTranslations(xcstrings, []model.TranslationResponse{resp})
			}
		}
	}

	_, translateErr := translator.TranslatePerLanguage(ctx, xcstrings, req.TargetLanguages, service, progressBuilder)
	if translateErr != nil {
		// Log error
		return
	}

	// Update the project in the database with the new translations
	err = services.UpdateTranslationsFromXCStrings(project.ID, xcstrings, req.Provider)
	if err != nil {
		// Log error
		return
	}

	// Update the project content structure
	err = services.SaveProjectContent(project.ID, xcstrings)
	if err != nil {
		// Log error
		return
	}
}

func handleExportProject(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := GetUserIDFromContext(c)
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

func handleGetProjectTranslations(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := GetUserIDFromContext(c)
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

func handleGetMissingTranslations(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := GetUserIDFromContext(c)
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

	missingKeys, err := services.GetMissingTranslations(uint(projectID), req.TargetLanguages)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"missingKeys": missingKeys,
	})
}

// Provider configuration handlers
func handleGetProviderConfigs(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	configs, err := services.GetProviderConfigsByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Sanitize configs before returning
	var sanitizedConfigs []database.ProviderConfig
	for _, config := range configs {
		sanitizedConfigs = append(sanitizedConfigs, *services.SanitizeProviderConfig(&config))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"configs": sanitizedConfigs,
	})
}

func handleCreateProviderConfig(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		ProviderType string                 `json:"providerType"`
		ConfigData   map[string]interface{} `json:"configData"`
		IsDefault    bool                   `json:"isDefault"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate the config data
	err := services.ValidateProviderConfigData(req.ProviderType, req.ConfigData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	config, err := services.CreateProviderConfig(userID, req.ProviderType, req.ConfigData, req.IsDefault)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"config":  services.SanitizeProviderConfig(config),
	})
}

func handleGetProviderConfig(c *fiber.Ctx) error {
	configID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid config ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	config, err := services.GetProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"config":  services.SanitizeProviderConfig(config),
	})
}

func handleUpdateProviderConfig(c *fiber.Ctx) error {
	configID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid config ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		ProviderType *string                `json:"providerType"`
		ConfigData   map[string]interface{} `json:"configData"`
		IsDefault    *bool                  `json:"isDefault"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing config to check ownership
	config, err := services.GetProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Configuration not found")
	}

	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}

	updates := make(map[string]interface{})
	if req.ProviderType != nil {
		updates["ProviderType"] = *req.ProviderType
	}
	if req.ConfigData != nil {
		// Validate the config data
		err := services.ValidateProviderConfigData(*req.ProviderType, req.ConfigData)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		updates["ConfigData"] = req.ConfigData
	}
	if req.IsDefault != nil {
		updates["IsDefault"] = *req.IsDefault
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No fields to update")
	}

	err = services.UpdateProviderConfig(uint(configID), updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Get the updated config to return
	updatedConfig, err := services.GetProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"config":  services.SanitizeProviderConfig(updatedConfig),
	})
}

func handleDeleteProviderConfig(c *fiber.Ctx) error {
	configID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid config ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get existing config to check ownership
	config, err := services.GetProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Configuration not found")
	}

	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}

	err = services.DeleteProviderConfig(uint(configID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Configuration deleted successfully",
	})
}

func handleGetDefaultProviderConfig(c *fiber.Ctx) error {
	providerType := c.Params("type")

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	config, err := services.GetDefaultProviderConfig(userID, providerType)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"config":  services.SanitizeProviderConfig(config),
	})
}
