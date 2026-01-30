package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fdddf/xcstrings-translator/internal/auth"
	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/email"
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

// NewAppWithDB is a version of NewApp that accepts a database instance
func NewAppWithDB(db *database.Database) (*fiber.App, error) {
	distFS, err := fs.Sub(webui.EmbeddedFS, "dist")
	if err != nil {
		return nil, fmt.Errorf("embedded UI missing: %w", err)
	}

	state := &ServerState{}

	app := fiber.New()

	// Use middleware to inject database into context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// Public routes (no authentication required)
	api := app.Group("/api")
	api.Post("/upload", state.handleUpload)
	api.Get("/strings", state.handleStrings)
	api.Post("/translate", state.handleTranslate)
	api.Get("/progress", state.handleProgress)
	api.Get("/export", state.handleExport)
	api.Get("/languages", handleGetSupportedLanguages)

	// Authentication routes
	api.Post("/auth/register", handleRegister)
	api.Post("/auth/login", handleLogin)
	api.Post("/auth/logout", handleLogout)
	api.Post("/auth/activate/:code", handleActivateUser) // Add activation endpoint

	// Protected routes (require authentication)
	protected := api.Group("/protected")
	protected.Use(AuthMiddleware)
	adminOnly := protected.Group("/admin")
	adminOnly.Use(AdminOnly)

	// Example admin route to verify wiring (extend as needed)
	adminOnly.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})
	protected.Get("/projects", handleGetProjects)
	protected.Post("/projects", handleCreateProject)
	protected.Get("/projects/:id", handleGetProject)
	protected.Put("/projects/:id", handleUpdateProject)
	protected.Delete("/projects/:id", handleDeleteProject)
	protected.Post("/projects/:id/upload", handleUploadToProject)
	protected.Post("/projects/:id/translate", SubscriptionRequired, handleTranslateProject)
	protected.Get("/projects/:id/export", handleExportProject)
	protected.Get("/projects/:id/translations", handleGetProjectTranslations)
	protected.Get("/projects/:id/missing-translations", handleGetMissingTranslations)
	protected.Get("/projects/:id/translation-status", handleGetTranslationStatus)

	// Provider configuration routes
	protected.Get("/providers", handleGetProviderConfigs)
	protected.Post("/providers", handleCreateProviderConfig)
	protected.Get("/providers/:id", handleGetProviderConfig)
	protected.Put("/providers/:id", handleUpdateProviderConfig)
	protected.Delete("/providers/:id", handleDeleteProviderConfig)
	protected.Get("/providers/:type/default", handleGetDefaultProviderConfig)

	// Activity logging routes
	protected.Get("/activities", handleGetUserActivities)

	// App management routes
	protected.Get("/apps", handleGetApps)
	protected.Post("/apps", SubscriptionRequired, handleCreateApp)
	protected.Get("/apps/:id", handleGetApp)
	protected.Put("/apps/:id", handleUpdateApp)
	protected.Delete("/apps/:id", handleDeleteApp)
	protected.Get("/apps/:id/users", handleGetAppUsers)
	protected.Post("/apps/:id/users", handleAddUserToApp)
	protected.Put("/apps/:id/users/:userId", handleUpdateUserAppRole)
	protected.Delete("/apps/:id/users/:userId", handleRemoveUserFromApp)

	// App localization routes
	protected.Get("/apps/:id/localizations", handleGetAppLocalizations)
	protected.Post("/apps/:id/localizations", handleCreateAppLocalization)
	protected.Get("/apps/:id/localizations/:language", handleGetAppLocalization)
	protected.Put("/apps/:id/localizations/:language", handleUpdateAppLocalization)
	protected.Delete("/apps/:id/localizations/:language", handleDeleteAppLocalization)
	protected.Put("/apps/:id/localizations/bulk", handleBulkUpdateAppLocalizations)
	protected.Get("/apps/:id/languages", handleGetAppLanguages)
	protected.Post("/apps/:id/languages", handleAddAppLanguage)
	protected.Delete("/apps/:id/languages/:language", handleRemoveAppLanguage)

	// Apple Connect API routes
	protected.Post("/apple-connect/sync-apps", handleSyncAppleApps)
	protected.Post("/apple-connect/:appId/sync-localizations", handleSyncAppleAppLocalizations)
	protected.Post("/apps/:id/sync-to-apple", handleSyncAppToApple)
	protected.Get("/apps/:id/stats", handleGetAppStats)

	// Subscription management routes
	protected.Get("/subscription", handleGetUserSubscription)
	protected.Post("/subscription/webhook", handleSubscriptionWebhook)
	protected.Get("/subscription/usage", handleGetUsage)

	// Translation queue routes
	protected.Post("/queue/translate", SubscriptionRequired, handleQueueTranslationJob)

	// Project additional routes
	protected.Get("/projects/:id/stats", handleGetProjectStats)
	protected.Put("/projects/:id/translations/:key/:language", handleUpdateSingleTranslation)
	protected.Put("/projects/:id/translations/bulk", handleBulkUpdateTranslations)
	protected.Get("/projects/:id/languages", handleGetProjectLanguages)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(distFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	return app, nil
}

// NewApp constructs the Fiber app with embedded assets without starting the server.
// For use without DI, this function remains for backward compatibility
func NewApp() (*fiber.App, error) {
	// Note: Database is initialized by the fx framework
	// This function is for direct server instantiation (e.g., for testing)
	// when not using the fx framework

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
	api.Post("/auth/activate/:code", handleActivateUser) // Add activation endpoint

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

	// Activity logging routes
	protected.Get("/activities", handleGetUserActivities)

	// App management routes
	protected.Get("/apps", handleGetApps)
	protected.Post("/apps", handleCreateApp)
	protected.Get("/apps/:id", handleGetApp)
	protected.Put("/apps/:id", handleUpdateApp)
	protected.Delete("/apps/:id", handleDeleteApp)
	protected.Get("/apps/:id/users", handleGetAppUsers)
	protected.Post("/apps/:id/users", handleAddUserToApp)
	protected.Put("/apps/:id/users/:userId", handleUpdateUserAppRole)
	protected.Delete("/apps/:id/users/:userId", handleRemoveUserFromApp)

	// App localization routes
	protected.Get("/apps/:id/localizations", handleGetAppLocalizations)
	protected.Post("/apps/:id/localizations", handleCreateAppLocalization)
	protected.Get("/apps/:id/localizations/:language", handleGetAppLocalization)
	protected.Put("/apps/:id/localizations/:language", handleUpdateAppLocalization)
	protected.Delete("/apps/:id/localizations/:language", handleDeleteAppLocalization)
	protected.Put("/apps/:id/localizations/bulk", handleBulkUpdateAppLocalizations)

	// Apple Connect API routes
	protected.Post("/apple-connect/sync-apps", handleSyncAppleApps)
	protected.Post("/apple-connect/:appId/sync-localizations", handleSyncAppleAppLocalizations)
	protected.Post("/apps/:id/sync-to-apple", handleSyncAppToApple)
	protected.Get("/apps/:id/stats", handleGetAppStats)

	// Subscription management routes
	protected.Get("/subscription", handleGetUserSubscription)
	protected.Post("/subscription/webhook", handleSubscriptionWebhook)
	protected.Get("/subscription/usage", handleGetUsage)

	// Translation queue routes
	protected.Post("/queue/translate", handleQueueTranslationJob)

	// Project additional routes
	protected.Get("/projects/:id/stats", handleGetProjectStats)
	protected.Put("/projects/:id/translations/:key/:language", handleUpdateSingleTranslation)
	protected.Put("/projects/:id/translations/bulk", handleBulkUpdateTranslations)
	protected.Get("/projects/:id/languages", handleGetProjectLanguages)

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
		if err := services.ValidateLanguages([]string{source}); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
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

	if err := services.ValidateLanguages(req.TargetLanguages); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
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

// handleGetSupportedLanguages returns the supported language codes
func handleGetSupportedLanguages(c *fiber.Ctx) error {
	languages := make([]string, 0, len(services.SupportedLanguages))
	for lang := range services.SupportedLanguages {
		languages = append(languages, lang)
	}
	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
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
		// Return a consistent error response format
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		// Return a consistent error response format
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Username, email, and password are required",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	user, err := auth.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		// Return a consistent error response format
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	// Log registration activity
	logActivity(c, user.ID, "user_registered", fmt.Sprintf("User %s registered", user.Username), "")

	// Send activation email
	go func() {
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			// Default to localhost if not set
			baseURL = "http://localhost:3000"
		}

		err := email.SendActivationEmail(
			user.Email,
			user.Username,
			user.ActivationCode,
			baseURL,
		)
		if err != nil {
			// Log the error but don't fail the registration
			fmt.Printf("Failed to send activation email: %v\n", err)
		}
	}()

	// Return a consistent success response format
	response := struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		User    *database.User `json:"user"`
	}{
		Success: true,
		Message: "Registration successful. Please check your email for activation.",
		User:    user,
	}

	return c.JSON(response)
}

func handleLogin(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		// Return a consistent error response format
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	user, token, err := auth.LoginUser(req.Username, req.Password)
	if err != nil {
		// Return a consistent error response format
		errorResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: err.Error(),
		}
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	// Log login activity
	logActivity(c, user.ID, "user_logged_in", fmt.Sprintf("User %s logged in", user.Username), "")

	// Return a consistent success response format
	response := struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		User    *database.User `json:"user"`
		Token   string         `json:"token"`
	}{
		Success: true,
		Message: "Login successful",
		User:    user,
		Token:   token,
	}

	return c.JSON(response)
}

func handleLogout(c *fiber.Ctx) error {
	// In a real application, you might implement token blacklisting
	// For now, just return success
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "Logged out successfully",
	}

	return c.JSON(response)
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

func handleGetTranslationStatus(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	userID, ok := GetUserIDFromContext(c)
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
		// also allow body parsing fallback
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
		sanitizedConfig := services.SanitizeProviderConfig(&config)
		if sanitizedConfig != nil {
			sanitizedConfigs = append(sanitizedConfigs, *sanitizedConfig)
		}
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

// User activation handler
func handleActivateUser(c *fiber.Ctx) error {
	activationCode := c.Params("code")

	err := auth.ActivateUser(activationCode)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account activated successfully. You can now log in.",
	})
}

// Get user activities handler
func handleGetUserActivities(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Attempt to get database from context or handle properly
	// Since we removed the global DB, we need to pass the database via context
	// This would require updating the server setup to inject DB into fiber context
	db, ok := c.Locals("db").(*database.Database)
	if !ok {
		// If no DB in context, return an error
		// In a proper DI system, the DB would be available in context
		return fiber.NewError(fiber.StatusInternalServerError, "Database not available in context")
	}

	var activities []database.UserActivity
	result := db.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&activities)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"activities": activities,
	})
}

// logActivity logs a user activity
// Note: This function currently needs an injected database instance to work properly
// In a complete DI implementation, this would be passed as parameter or context
func logActivity(c *fiber.Ctx, userID uint, action, details, useragent string) {
	// For now, we'll need to pass db through context or refactor this function
	// Since we removed the global DB, this function needs to be called with proper context
	// This is a temporary fix - in full DI, this would be passed via context or dependency

	// For now, we'll return early since there's no global DB to use
	// In a complete implementation, this would either get DB from context or be refactored
	fmt.Printf("Activity logged (no DB): userID=%d, action=%s, details=%s\n", userID, action, details)

	// Note: This function would need to be refactored to properly handle DI
	// The activity won't be saved to DB since global DB is removed
}

// App management handlers
func handleGetApps(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	apps, err := services.GetAppsByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"apps":    apps,
	})
}

func handleCreateApp(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		Name             string `json:"name"`
		Description      string `json:"description"`
		BundleID         string `json:"bundleId"`
		AppleID          string `json:"appleId"`
		PrimaryLocale    string `json:"primaryLocale"`
		ShortDescription string `json:"shortDescription"`
		LongDescription  string `json:"longDescription"`
		Keywords         string `json:"keywords"`
		SupportURL       string `json:"supportUrl"`
		MarketingURL     string `json:"marketingUrl"`
		PrivacyURL       string `json:"privacyUrl"`
		Version          string `json:"version"`
		AppCategory      string `json:"appCategory"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	app, err := services.CreateApp(userID, req.Name, req.Description, req.BundleID, req.AppleID, req.PrimaryLocale)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"app":     app,
	})
}

func handleGetApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"app":     app,
	})
}

func handleUpdateApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req struct {
		Name             *string `json:"name"`
		Description      *string `json:"description"`
		BundleID         *string `json:"bundleId"`
		AppleID          *string `json:"appleId"`
		PrimaryLocale    *string `json:"primaryLocale"`
		ShortDescription *string `json:"shortDescription"`
		LongDescription  *string `json:"longDescription"`
		Keywords         *string `json:"keywords"`
		SupportURL       *string `json:"supportUrl"`
		MarketingURL     *string `json:"marketingUrl"`
		PrivacyURL       *string `json:"privacyUrl"`
		Version          *string `json:"version"`
		AppCategory      *string `json:"appCategory"`
		IsReadyForReview *bool   `json:"isReadyForReview"`
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
	if req.BundleID != nil {
		updates["BundleID"] = *req.BundleID
	}
	if req.AppleID != nil {
		updates["AppleID"] = *req.AppleID
	}
	if req.PrimaryLocale != nil {
		updates["PrimaryLocale"] = *req.PrimaryLocale
	}
	if req.ShortDescription != nil {
		updates["ShortDescription"] = *req.ShortDescription
	}
	if req.LongDescription != nil {
		updates["LongDescription"] = *req.LongDescription
	}
	if req.Keywords != nil {
		updates["Keywords"] = *req.Keywords
	}
	if req.SupportURL != nil {
		updates["SupportURL"] = *req.SupportURL
	}
	if req.MarketingURL != nil {
		updates["MarketingURL"] = *req.MarketingURL
	}
	if req.PrivacyURL != nil {
		updates["PrivacyURL"] = *req.PrivacyURL
	}
	if req.Version != nil {
		updates["Version"] = *req.Version
	}
	if req.AppCategory != nil {
		updates["AppCategory"] = *req.AppCategory
	}
	if req.IsReadyForReview != nil {
		updates["IsReadyForReview"] = *req.IsReadyForReview
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No fields to update")
	}

	err = services.UpdateApp(uint(appID), updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App updated successfully",
	})
}

func handleDeleteApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	err = services.DeleteApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App deleted successfully",
	})
}

func handleGetAppUsers(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	users, err := services.GetUsersForApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

func handleAddUserToApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner of the app
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can add users")
	}

	var req struct {
		UserID uint   `json:"userId"`
		Role   string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate role
	validRoles := map[string]bool{
		"owner":  true,
		"admin":  true,
		"editor": true,
		"viewer": true,
	}
	if !validRoles[req.Role] {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role. Valid roles: owner, admin, editor, viewer")
	}

	err = services.AddUserToApp(uint(appID), req.UserID, req.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User added to app successfully",
	})
}

func handleUpdateUserAppRole(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userIDParam, err := c.ParamsInt("userId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner of the app
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can update user roles")
	}

	var req struct {
		Role string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate role
	validRoles := map[string]bool{
		"owner":  true,
		"admin":  true,
		"editor": true,
		"viewer": true,
	}
	if !validRoles[req.Role] {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role. Valid roles: owner, admin, editor, viewer")
	}

	err = services.UpdateUserAppRole(uint(appID), uint(userIDParam), req.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User role updated successfully",
	})
}

func handleRemoveUserFromApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userIDParam, err := c.ParamsInt("userId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user is the owner of the app
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}

	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can remove users")
	}

	// Prevent removing the owner from their own app
	if uint(userIDParam) == app.UserID {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot remove the owner from the app")
	}

	err = services.RemoveUserFromApp(uint(appID), uint(userIDParam))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User removed from app successfully",
	})
}

// App localization handlers
func handleGetAppLocalizations(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	localizations, err := services.GetAppLocalizations(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"localizations": localizations,
	})
}

func handleGetAppLanguages(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	languages, err := services.GetAppSupportedLanguages(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
}

func handleAddAppLanguage(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Only owner may manage languages
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}
	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can manage languages")
	}

	var req struct {
		Language string `json:"language"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := services.ValidateLanguages([]string{req.Language}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = services.GetOrCreateAppLocalization(
		uint(appID),
		req.Language,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	languages, err := services.GetAppSupportedLanguages(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
}

func handleRemoveAppLanguage(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	language := c.Params("language")

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Only owner may manage languages
	app, err := services.GetApp(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}
	if app.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Only app owner can manage languages")
	}

	if err := services.ValidateLanguages([]string{language}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := services.DeleteAppLocalization(uint(appID), language); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	languages, err := services.GetAppSupportedLanguages(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
}

func handleCreateAppLocalization(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	role, _ := GetUserRoleFromContext(c)
	if role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Only admin can create localizations")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req struct {
		LanguageCode        string `json:"languageCode"`
		Name                string `json:"name"`
		Subtitle            string `json:"subtitle"`
		PrivacyURL          string `json:"privacyUrl"`
		MarketingURL        string `json:"marketingUrl"`
		SupportURL          string `json:"supportUrl"`
		DownloadDescription string `json:"downloadDescription"`
		ShortDescription    string `json:"shortDescription"`
		LongDescription     string `json:"longDescription"`
		Keywords            string `json:"keywords"`
		ReleaseNotes        string `json:"releaseNotes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := services.ValidateLanguages([]string{req.LanguageCode}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	localization, err := services.CreateAppLocalization(
		uint(appID),
		req.LanguageCode,
		req.Name,
		req.Subtitle,
		req.PrivacyURL,
		req.MarketingURL,
		req.SupportURL,
		req.DownloadDescription,
		req.ShortDescription,
		req.LongDescription,
		req.Keywords,
		req.ReleaseNotes,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	_ = services.UpdateAppLocalization(uint(appID), req.LanguageCode, map[string]interface{}{
		"SyncStatus": "pending",
		"Source":     "local",
	})

	return c.JSON(fiber.Map{
		"success":      true,
		"localization": localization,
	})
}

func handleGetAppLocalization(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	languageCode := c.Params("language")

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	localization, err := services.GetAppLocalization(uint(appID), languageCode)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"localization": localization,
	})
}

func handleUpdateAppLocalization(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	languageCode := c.Params("language")

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	role, _ := GetUserRoleFromContext(c)
	if role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Only admin can update localizations")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req struct {
		Name                string `json:"name"`
		Subtitle            string `json:"subtitle"`
		PrivacyURL          string `json:"privacyUrl"`
		MarketingURL        string `json:"marketingUrl"`
		SupportURL          string `json:"supportUrl"`
		DownloadDescription string `json:"downloadDescription"`
		ShortDescription    string `json:"shortDescription"`
		LongDescription     string `json:"longDescription"`
		Keywords            string `json:"keywords"`
		ReleaseNotes        string `json:"releaseNotes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	updates := map[string]interface{}{
		"Name":                req.Name,
		"Subtitle":            req.Subtitle,
		"PrivacyURL":          req.PrivacyURL,
		"MarketingURL":        req.MarketingURL,
		"SupportURL":          req.SupportURL,
		"DownloadDescription": req.DownloadDescription,
		"ShortDescription":    req.ShortDescription,
		"LongDescription":     req.LongDescription,
		"Keywords":            req.Keywords,
		"ReleaseNotes":        req.ReleaseNotes,
		"SyncStatus":          "pending",
		"Source":              "local",
	}

	err = services.UpdateAppLocalization(uint(appID), languageCode, updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App localization updated successfully",
	})
}

func handleDeleteAppLocalization(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	languageCode := c.Params("language")

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	role, _ := GetUserRoleFromContext(c)
	if role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Only admin can delete localizations")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	err = services.DeleteAppLocalization(uint(appID), languageCode)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App localization deleted successfully",
	})
}

// Apple Connect API handlers
func handleSyncAppleApps(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		ConfigID uint `json:"configId"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if req.ConfigID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "configId is required")
	}

	// Get the user's Apple Connect configuration
	config, err := services.GetProviderConfig(req.ConfigID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Provider configuration not found")
	}
	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}
	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Provider configuration is not Apple Connect")
	}
	_, ok = config.ConfigData["issuerID"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "issuerID is missing from configuration")
	}
	_, ok = config.ConfigData["keyID"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "keyID is missing from configuration")
	}
	_, ok = config.ConfigData["privateKey"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "privateKey is missing from configuration")
	}

	queueService := services.GetQueueService()
	if queueService.AppService == nil || queueService.AppLocalizationService == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Apple Connect service unavailable")
	}
	apps, err := services.NewAppleConnectService(services.AppleConnectServiceDeps{
		AppService:             queueService.AppService,
		AppLocalizationService: queueService.AppLocalizationService,
	}).SyncApps(userID,
		config.ConfigData["issuerID"].(string),
		config.ConfigData["keyID"].(string),
		fmt.Sprintf("%v", config.ConfigData["privateKeyPath"]),
		fmt.Sprintf("%v", config.ConfigData["privateKey"]),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, syncedApp := range apps {
		_, _ = services.BindAppProviderConfig(syncedApp.ID, config.ID, "appleconnect", false)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Synced %d apps from Apple Connect", len(apps)),
		"count":   len(apps),
	})
}

func handleSyncAppleAppLocalizations(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("appId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		ConfigID uint `json:"configId"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if req.ConfigID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "configId is required")
	}

	// Check if user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	config, err := services.GetProviderConfig(req.ConfigID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Provider configuration not found")
	}
	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}
	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Provider configuration is not Apple Connect")
	}

	_, err = services.BindAppProviderConfig(uint(appID), config.ID, "appleconnect", true)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	binding, err := services.GetAppProviderConfig(uint(appID), config.ID)
	if err != nil || binding.ProviderConfigID != config.ID {
		return fiber.NewError(fiber.StatusForbidden, "Apple Connect configuration not bound to this app")
	}

	queueService := services.GetQueueService()
	if queueService.AppService == nil || queueService.AppLocalizationService == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Apple Connect service unavailable")
	}
	localizations, err := services.NewAppleConnectService(services.AppleConnectServiceDeps{
		AppService:             queueService.AppService,
		AppLocalizationService: queueService.AppLocalizationService,
	}).SyncLocalizations(uint(appID),
		config.ConfigData["issuerID"].(string),
		config.ConfigData["keyID"].(string),
		fmt.Sprintf("%v", config.ConfigData["privateKeyPath"]),
		fmt.Sprintf("%v", config.ConfigData["privateKey"]),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Synced %d localizations from Apple Connect", len(localizations)),
		"count":   len(localizations),
	})
}

// Subscription management handlers
func handleGetUserSubscription(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	user, err := services.GetUserSubscriptionInfo(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

func handleSubscriptionWebhook(c *fiber.Ctx) error {
	// For now, just log the webhook payload - in a real application, you'd process specific Stripe events
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	fmt.Printf("Received subscription webhook: %+v\n", req)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Webhook received",
	})
}

func handleGetUsage(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	overLimit, usage, limit, err := services.CheckUserUsage(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"overLimit":  overLimit,
		"usage":      usage,
		"limit":      limit,
		"percentage": float64(usage) / float64(limit) * 100,
	})
}

// Translation queue handlers
func handleQueueTranslationJob(c *fiber.Ctx) error {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		JobType         string                 `json:"jobType"`
		ProjectID       *uint                  `json:"projectId"`
		AppID           *uint                  `json:"appId"`
		ProviderType    string                 `json:"providerType"`
		SourceLanguage  string                 `json:"sourceLanguage"`
		TargetLanguages []string               `json:"targetLanguages"`
		ConfigData      map[string]interface{} `json:"configData"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := services.ValidateLanguages(req.TargetLanguages); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verify user has access to project or app if provided
	if req.ProjectID != nil {
		project, err := services.GetProject(*req.ProjectID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Project not found")
		}
		if project.UserID != userID {
			return fiber.NewError(fiber.StatusForbidden, "Access denied to this project")
		}
	}

	if req.AppID != nil {
		hasAccess, _, err := services.CheckUserAccessToApp(*req.AppID, userID)
		if err != nil || !hasAccess {
			return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
		}
	}

	queueService := services.GetQueueService()
	job, err := queueService.SubmitTranslationJob(
		userID,
		req.ProjectID,
		req.AppID,
		req.JobType,
		req.ProviderType,
		req.SourceLanguage,
		req.TargetLanguages,
		req.ConfigData,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"job":     job,
	})
}

// Project additional handlers
func handleGetProjectStats(c *fiber.Ctx) error {
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

	stats, err := services.GetProjectStats(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"stats":   stats,
	})
}

func handleUpdateSingleTranslation(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}

	key := c.Params("key")
	language := c.Params("language")

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

func handleBulkUpdateTranslations(c *fiber.Ctx) error {
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

func handleGetProjectLanguages(c *fiber.Ctx) error {
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

	languages, err := services.GetProjectLanguages(uint(projectID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"languages": languages,
	})
}

// App additional handlers
func handleGetAppStats(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	stats, err := services.GetAppStats(uint(appID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"stats":   stats,
	})
}

func handleBulkUpdateAppLocalizations(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
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

	err = services.BulkUpdateAppLocalizations(uint(appID), req.Updates)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App localizations updated successfully",
	})
}

func handleSyncAppToApple(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid app ID")
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Verify user has access to this app
	hasAccess, _, err := services.CheckUserAccessToApp(uint(appID), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !hasAccess {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this app")
	}

	var req struct {
		ConfigID uint `json:"configId"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.ConfigID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "configId is required")
	}

	config, err := services.GetProviderConfig(req.ConfigID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Provider configuration not found")
	}
	if config.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this configuration")
	}
	if config.ProviderType != "appleconnect" {
		return fiber.NewError(fiber.StatusBadRequest, "Provider configuration is not Apple Connect")
	}
	_, ok = config.ConfigData["issuerID"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "issuerID is missing from configuration")
	}
	_, ok = config.ConfigData["keyID"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "keyID is missing from configuration")
	}
	_, ok = config.ConfigData["privateKey"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "privateKey is missing from configuration")
	}

	err = services.SyncAppToAppleConnect(uint(appID),
		config.ConfigData["issuerID"].(string),
		config.ConfigData["keyID"].(string),
		config.ConfigData["privateKey"].(string),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "App synced to Apple Connect successfully",
	})
}
