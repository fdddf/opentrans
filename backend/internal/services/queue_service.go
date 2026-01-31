package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/model"
	"github.com/fdddf/xcstrings-translator/internal/translator"
	"gorm.io/gorm"
)

// QueueService handles translation queue operations
type QueueService struct {
	DB                     *database.Database
	AppService             *AppService
	AppLocalizationService *AppLocalizationService
	ProjectService         *ProjectService
	TranslationService     *TranslationService
	SubscriptionService    *SubscriptionService
	ProviderService        *ProviderService
	AppProviderConfigService *AppProviderConfigService

	mu    sync.RWMutex
	queue map[uint]*database.TranslationQueue // In-memory cache of active jobs
}

var queueServiceInstance *QueueService

// GetQueueService returns a singleton instance of QueueService
func GetQueueService() *QueueService {
	if queueServiceInstance == nil {
		queueServiceInstance = &QueueService{
			queue: make(map[uint]*database.TranslationQueue),
		}
	}
	return queueServiceInstance
}

// SetQueueServiceInstance sets the global QueueService instance (used by FX)
func SetQueueServiceInstance(qs *QueueService) {
	queueServiceInstance = qs
}

// SubmitTranslationJob submits a new translation job to the queue
func (qs *QueueService) SubmitTranslationJob(userID uint, projectID *uint, appID *uint, jobType, providerType, sourceLanguage string, targetLanguages []string, configData map[string]interface{}) (*database.TranslationQueue, error) {
	// Check user's subscription and usage limits
	if projectID != nil {
		// This is related to a project, check if user has access
		project, err := qs.ProjectService.GetProject(*projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project: %v", err)
		}

		if project.UserID != userID {
			return nil, errors.New("user does not have access to this project")
		}
	} else if appID != nil {
		// This is related to an app, check if user has access
		hasAccess, _, err := qs.AppService.CheckUserAccessToApp(*appID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check user access to app: %v", err)
		}

		if !hasAccess {
			return nil, errors.New("user does not have access to this app")
		}
	}

	// Create the queue entry
	queueEntry := &database.TranslationQueue{
		UserID:          userID,
		ProjectID:       projectID,
		AppID:           appID,
		Type:            jobType,
		Status:          "pending",
		ProviderType:    providerType,
		SourceLanguage:  sourceLanguage,
		TargetLanguages: targetLanguages,
		ConfigData:      configData,
		Progress:        0,
		Total:           0, // Will be set when processing starts
		Done:            0,
	}

	result := qs.DB.Create(queueEntry)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to submit job to queue: %v", result.Error)
	}

	// Add to in-memory cache
	qs.mu.Lock()
	qs.queue[queueEntry.ID] = queueEntry
	qs.mu.Unlock()

	return queueEntry, nil
}

// GetQueueJob retrieves a specific queue job
func (qs *QueueService) GetQueueJob(jobID uint) (*database.TranslationQueue, error) {
	// Check in-memory cache first
	qs.mu.RLock()
	job, exists := qs.queue[jobID]
	qs.mu.RUnlock()

	if exists {
		return job, nil
	}

	// If not in cache, fetch from database
	var queueJob database.TranslationQueue
	result := qs.DB.First(&queueJob, jobID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("queue job not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	// Add to cache
	qs.mu.Lock()
	qs.queue[jobID] = &queueJob
	qs.mu.Unlock()

	return &queueJob, nil
}

// GetQueueJobsByUser retrieves all queue jobs for a user
func (qs *QueueService) GetQueueJobsByUser(userID uint) ([]database.TranslationQueue, error) {
	var queueJobs []database.TranslationQueue
	result := qs.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&queueJobs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve queue jobs: %v", result.Error)
	}

	// Update cache with these jobs
	qs.mu.Lock()
	for _, job := range queueJobs {
		qs.queue[job.ID] = &job
	}
	qs.mu.Unlock()

	return queueJobs, nil
}

// GetPendingJobs retrieves all pending jobs that need processing
func (qs *QueueService) GetPendingJobs() ([]database.TranslationQueue, error) {
	var queueJobs []database.TranslationQueue
	result := qs.DB.Where("status = ?", "pending").Order("priority DESC, created_at ASC").Find(&queueJobs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve pending jobs: %v", result.Error)
	}

	// Update cache with these jobs
	qs.mu.Lock()
	for _, job := range queueJobs {
		qs.queue[job.ID] = &job
	}
	qs.mu.Unlock()

	return queueJobs, nil
}

// UpdateQueueJob updates a queue job
func (qs *QueueService) UpdateQueueJob(jobID uint, updates map[string]interface{}) error {
	result := qs.DB.Model(&database.TranslationQueue{}).Where("id = ?", jobID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update queue job: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("queue job not found")
	}

	// Update in-memory cache
	qs.mu.Lock()
	if job, exists := qs.queue[jobID]; exists {
		for key, value := range updates {
			switch key {
			case "Status":
				job.Status = value.(string)
			case "Progress":
				job.Progress = int(value.(uint))
			case "Total":
				job.Total = int(value.(uint))
			case "Done":
				job.Done = int(value.(uint))
			case "Error":
				job.Error = value.(string)
			case "ResultData":
				job.ResultData = value.(map[string]interface{})
			}
		}
	}
	qs.mu.Unlock()

	return nil
}

// ProcessNextJob processes the next available job in the queue
func (qs *QueueService) ProcessNextJob() error {
	// Get the next pending job
	var nextJob database.TranslationQueue
	result := qs.DB.Where("status = ?", "pending").Order("priority DESC, created_at ASC").First(&nextJob)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// No pending jobs
		return nil
	}

	if result.Error != nil {
		return fmt.Errorf("failed to get next job: %v", result.Error)
	}

	// Update status to processing
	err := qs.UpdateQueueJob(nextJob.ID, map[string]interface{}{
		"Status":   "processing",
		"Progress": 0,
	})
	if err != nil {
		return fmt.Errorf("failed to update job status: %v", err)
	}

	// Process the job based on type
	switch nextJob.Type {
	case "xcstrings":
		err = qs.processXCStringsJob(&nextJob)
	case "app_localization":
		err = qs.processAppLocalizationJob(&nextJob)
	default:
		err = errors.New("unknown job type: " + nextJob.Type)
	}

	if err != nil {
		// Mark as failed
		failureErr := qs.UpdateQueueJob(nextJob.ID, map[string]interface{}{
			"Status":   "failed",
			"Error":    err.Error(),
			"Progress": 100,
		})
		if failureErr != nil {
			return fmt.Errorf("failed to mark job as failed: %v, original error: %v", failureErr, err)
		}
		return err
	} else {
		// Mark as completed
		err = qs.UpdateQueueJob(nextJob.ID, map[string]interface{}{
			"Status":   "completed",
			"Progress": 100,
		})
		if err != nil {
			return fmt.Errorf("failed to mark job as completed: %v", err)
		}
	}

	return nil
}

// processXCStringsJob processes an xcstrings translation job
func (qs *QueueService) processXCStringsJob(job *database.TranslationQueue) error {
	if job.ProjectID == nil {
		return errors.New("project ID is required for xcstrings job")
	}

	// Get the project
	project, err := qs.ProjectService.GetProject(*job.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %v", err)
	}

	// Parse the project content
	xcstrings, err := qs.ProjectService.ParseProjectContent(project)
	if err != nil {
		return fmt.Errorf("failed to parse project content: %v", err)
	}

	// Create translation provider
	provider, err := createProviderFromConfig(job.ProviderType, job.ConfigData)
	if err != nil {
		return fmt.Errorf("failed to create provider: %v", err)
	}

	// Get the actual translation provider interface
	translationProvider, ok := provider.(model.TranslationProvider)
	if !ok {
		return errors.New("invalid provider type")
	}

	// Create translation service
	concurrency := 4
	if c, ok := job.ConfigData["concurrency"].(int); ok && c > 0 {
		concurrency = c
	}

	timeout := 300 * time.Second
	if t, ok := job.ConfigData["timeout"].(int); ok && t > 0 {
		timeout = time.Duration(t) * time.Second
	}

	service := translator.NewTranslationService(translationProvider, concurrency, timeout)

	// Create translation requests
	requests := translator.CreateTranslationRequests(xcstrings, job.TargetLanguages)

	// Update job with total count
	err = qs.UpdateQueueJob(job.ID, map[string]interface{}{
		"Total": len(requests),
	})
	if err != nil {
		return fmt.Errorf("failed to update job total: %v", err)
	}

	if len(requests) == 0 {
		// Nothing to translate
		err = qs.UpdateQueueJob(job.ID, map[string]interface{}{
			"Done":     0,
			"Progress": 100,
		})
		if err != nil {
			return fmt.Errorf("failed to update job progress: %v", err)
		}
		return nil
	}

	// Execute translations with progress tracking
	ctx := context.Background()
	doneCount := 0

	progressBuilder := func(target string, total int) translator.ProgressReporter {
		return func(done, total int, resp model.TranslationResponse) {
			if resp.Error == nil {
				// Apply translation to xcstrings
				translator.ApplyTranslations(xcstrings, []model.TranslationResponse{resp})
			}

			// Update progress
			doneCount++
			progress := int(float64(doneCount) / float64(len(requests)) * 100)

			qs.UpdateQueueJob(job.ID, map[string]interface{}{
				"Done":     doneCount,
				"Progress": progress,
			})
		}
	}

	// Execute translation
	_, translateErr := translator.TranslatePerLanguage(ctx, xcstrings, job.TargetLanguages, service, progressBuilder)

	// Save translations to database
	if translateErr == nil {
		err = qs.ProjectService.UpdateTranslationsFromXCStrings(project.ID, xcstrings, job.ProviderType)
		if err != nil {
			return fmt.Errorf("failed to save translations: %v", err)
		}

		// Update project content
		err = qs.ProjectService.SaveProjectContent(project.ID, xcstrings)
		if err != nil {
			return fmt.Errorf("failed to update project content: %v", err)
		}

		// Update user usage
		err = qs.UpdateUserUsage(job.UserID, len(requests))
		if err != nil {
			return fmt.Errorf("failed to update user usage: %v", err)
		}
	}

	// Mark job as completed or failed
	if translateErr != nil {
		return fmt.Errorf("translation failed: %v", translateErr)
	}

	err = qs.UpdateQueueJob(job.ID, map[string]interface{}{
		"Done":     len(requests),
		"Progress": 100,
	})
	if err != nil {
		return fmt.Errorf("failed to update job progress: %v", err)
	}

	return nil
}

// processAppLocalizationJob processes an app localization job
func (qs *QueueService) processAppLocalizationJob(job *database.TranslationQueue) error {
	if job.AppID == nil {
		return errors.New("app ID is required for app localization job")
	}

	// Get the app
	app, err := qs.AppService.GetApp(*job.AppID)
	if err != nil {
		return fmt.Errorf("failed to get app: %v", err)
	}

	// Get existing localizations
	localizations, err := qs.AppLocalizationService.GetAppLocalizations(*job.AppID)
	if err != nil {
		return fmt.Errorf("failed to get app localizations: %v", err)
	}

	// Get the primary locale as the source
	sourceLoc, err := qs.AppLocalizationService.GetAppLocalization(*job.AppID, app.PrimaryLocale)
	if err != nil {
		return fmt.Errorf("failed to get source localization: %v", err)
	}

	// Create translation provider
	provider, err := createProviderFromConfig(job.ProviderType, job.ConfigData)
	if err != nil {
		return fmt.Errorf("failed to create provider: %v", err)
	}

	// Get the actual translation provider interface
	translationProvider, ok := provider.(model.TranslationProvider)
	if !ok {
		return errors.New("invalid provider type")
	}

	// Create translation service with concurrency 1 for Llama to avoid thread safety issues
	concurrency := 1
	if job.ProviderType != "llama" {
		concurrency = 4
	}

	timeout := 300 * time.Second
	if t, ok := job.ConfigData["timeout"].(int); ok && t > 0 {
		timeout = time.Duration(t) * time.Second
	}

	service := translator.NewTranslationService(translationProvider, concurrency, timeout)

	// Count total fields to translate
	fieldsToTranslate := []string{
		"Name", "Subtitle", "ShortDescription", "LongDescription", "Keywords",
		"ReleaseNotes", "PromotionalText", "DownloadDescription",
	}
	totalFields := len(job.TargetLanguages) * len(fieldsToTranslate)

	// Update job with total count
	err = qs.UpdateQueueJob(job.ID, map[string]interface{}{
		"Total": totalFields,
	})
	if err != nil {
		return fmt.Errorf("failed to update job total: %v", err)
	}

	// Process each target language
	doneCount := 0
	ctx := context.Background()

	for _, targetLang := range job.TargetLanguages {
		// Check if localization already exists
		var targetLoc *database.AppLocalization
		exists := false
		for _, loc := range localizations {
			if loc.LanguageCode == targetLang {
				exists = true
				targetLoc = &loc
				break
			}
		}

		// If it doesn't exist, create one from source
		if !exists {
			targetLoc = &database.AppLocalization{
				AppID:              *job.AppID,
				LanguageCode:       targetLang,
				Name:               sourceLoc.Name,
				Subtitle:           sourceLoc.Subtitle,
				PrivacyURL:         sourceLoc.PrivacyURL,
				MarketingURL:       sourceLoc.MarketingURL,
				SupportURL:         sourceLoc.SupportURL,
				DownloadDescription: sourceLoc.DownloadDescription,
				ShortDescription:   sourceLoc.ShortDescription,
				LongDescription:    sourceLoc.LongDescription,
				Keywords:           sourceLoc.Keywords,
				ReleaseNotes:       sourceLoc.ReleaseNotes,
				PromotionalText:    sourceLoc.PromotionalText,
			}
		}

		// Translate each field
		for _, field := range fieldsToTranslate {
			sourceText := ""
			switch field {
			case "Name":
				sourceText = sourceLoc.Name
			case "Subtitle":
				sourceText = sourceLoc.Subtitle
			case "ShortDescription":
				sourceText = sourceLoc.ShortDescription
			case "LongDescription":
				sourceText = sourceLoc.LongDescription
			case "Keywords":
				sourceText = sourceLoc.Keywords
			case "ReleaseNotes":
				sourceText = sourceLoc.ReleaseNotes
			case "PromotionalText":
				sourceText = sourceLoc.PromotionalText
			case "DownloadDescription":
				sourceText = sourceLoc.DownloadDescription
			}

			// Skip empty fields
			if sourceText == "" {
				doneCount++
				progress := int(float64(doneCount) / float64(totalFields) * 100)
				qs.UpdateQueueJob(job.ID, map[string]interface{}{
					"Done":     doneCount,
					"Progress": progress,
				})
				continue
			}

			// Create translation request
			req := model.TranslationRequest{
				Key:            fmt.Sprintf("%s.%s", targetLang, field),
				Text:           sourceText,
				SourceLanguage: job.SourceLanguage,
				TargetLanguage: targetLang,
			}

			// Translate using batch method (single item batch)
			responses, err := service.TranslateBatch(ctx, []model.TranslationRequest{req}, nil)
			if err != nil || len(responses) == 0 || responses[0].Error != nil {
				// Log error but continue with other fields
				if err != nil {
					fmt.Printf("Error translating %s: %v\n", field, err)
				} else if len(responses) > 0 && responses[0].Error != nil {
					fmt.Printf("Error translating %s: %v\n", field, responses[0].Error)
				}
				doneCount++
				progress := int(float64(doneCount) / float64(totalFields) * 100)
				qs.UpdateQueueJob(job.ID, map[string]interface{}{
					"Done":     doneCount,
					"Progress": progress,
				})
				continue
			}

			// Update the field with translated text
			switch field {
			case "Name":
				targetLoc.Name = responses[0].TranslatedText
			case "Subtitle":
				targetLoc.Subtitle = responses[0].TranslatedText
			case "ShortDescription":
				targetLoc.ShortDescription = responses[0].TranslatedText
			case "LongDescription":
				targetLoc.LongDescription = responses[0].TranslatedText
			case "Keywords":
				targetLoc.Keywords = responses[0].TranslatedText
			case "ReleaseNotes":
				targetLoc.ReleaseNotes = responses[0].TranslatedText
			case "PromotionalText":
				targetLoc.PromotionalText = responses[0].TranslatedText
			case "DownloadDescription":
				targetLoc.DownloadDescription = responses[0].TranslatedText
			}

			doneCount++
			progress := int(float64(doneCount) / float64(totalFields) * 100)
			qs.UpdateQueueJob(job.ID, map[string]interface{}{
				"Done":     doneCount,
				"Progress": progress,
			})
		}

		// Save or update the localization
		updates := map[string]interface{}{
			"name":                targetLoc.Name,
			"subtitle":            targetLoc.Subtitle,
			"privacy_url":         targetLoc.PrivacyURL,
			"marketing_url":       targetLoc.MarketingURL,
			"support_url":         targetLoc.SupportURL,
			"download_description": targetLoc.DownloadDescription,
			"short_description":   targetLoc.ShortDescription,
			"long_description":    targetLoc.LongDescription,
			"keywords":            targetLoc.Keywords,
			"release_notes":       targetLoc.ReleaseNotes,
			"promotional_text":    targetLoc.PromotionalText,
		}

		if exists {
			err = qs.AppLocalizationService.UpdateAppLocalizationWithValidation(*job.AppID, targetLang, updates)
		} else {
			_, err = qs.AppLocalizationService.CreateAppLocalization(*job.AppID, targetLang,
				targetLoc.Name, targetLoc.Subtitle, targetLoc.PrivacyURL,
				targetLoc.MarketingURL, targetLoc.SupportURL, targetLoc.DownloadDescription,
				targetLoc.ShortDescription, targetLoc.LongDescription, targetLoc.Keywords,
				targetLoc.ReleaseNotes, targetLoc.PromotionalText, "",
			)
		}

		if err != nil {
			return fmt.Errorf("failed to save app localization for %s: %v", targetLang, err)
		}
	}

	// Mark job as completed
	err = qs.UpdateQueueJob(job.ID, map[string]interface{}{
		"Done":     totalFields,
		"Progress": 100,
	})
	if err != nil {
		return fmt.Errorf("failed to update job progress: %v", err)
	}

	return nil
}

// CreateTranslationRequestsForLanguages creates translation requests for specific languages
func CreateTranslationRequestsForLanguages(xcstrings *model.XCStrings, targetLanguages []string) []interface{} { // Using interface{} as a placeholder
	// This function would create specific translation requests based on missing keys
	// For now, returning an empty slice - actual implementation would analyze the xcstrings
	// content and create requests for missing translations
	return make([]interface{}, 0)
}

// createProviderFromConfig creates a provider based on configuration data
func createProviderFromConfig(providerType string, configData map[string]interface{}) (model.TranslationProvider, error) {
	providerType = strings.ToLower(providerType)

	switch providerType {
	case "google":
		apiKey, ok := configData["apiKey"].(string)
		if !ok || apiKey == "" {
			return nil, errors.New("apiKey required for Google provider")
		}
		return translator.NewGoogleTranslator(apiKey), nil

	case "deepl":
		apiKey, ok := configData["apiKey"].(string)
		if !ok || apiKey == "" {
			return nil, errors.New("apiKey required for DeepL provider")
		}
		isFree := false
		if free, ok := configData["isFree"].(bool); ok {
			isFree = free
		}
		return translator.NewDeepLTranslator(apiKey, isFree), nil

	case "baidu":
		appID, ok := configData["appId"].(string)
		if !ok || appID == "" {
			return nil, errors.New("appId required for Baidu provider")
		}
		appSecret, ok := configData["appSecret"].(string)
		if !ok || appSecret == "" {
			return nil, errors.New("appSecret required for Baidu provider")
		}
		return translator.NewBaiduTranslator(appID, appSecret), nil

	case "openai":
		apiKey, ok := configData["apiKey"].(string)
		if !ok || apiKey == "" {
			return nil, errors.New("apiKey required for OpenAI provider")
		}
		apiBaseURL := "https://api.openai.com"
		if url, ok := configData["apiBaseUrl"].(string); ok && url != "" {
			apiBaseURL = url
		}
		model := "gpt-3.5-turbo"
		if m, ok := configData["model"].(string); ok && m != "" {
			model = m
		}
		temperature := 0.3
		if t, ok := configData["temperature"].(float64); ok {
			temperature = t
		}
		maxTokens := 1024
		if mt, ok := configData["maxTokens"].(int); ok {
			maxTokens = mt
		}
		return translator.NewOpenAITranslator(apiKey, apiBaseURL, model, temperature, maxTokens), nil

	case "llama":
		modelPath, ok := configData["modelPath"].(string)
		if !ok || modelPath == "" {
			return nil, errors.New("modelPath required for Llama provider")
		}
		libPath, ok := configData["libPath"].(string)
		if !ok || libPath == "" {
			return nil, errors.New("libPath required for Llama provider")
		}

		// Build Llama options from config data
		options := translator.LlamaOptions{
			ModelPath:   modelPath,
			GrammarPath: "",
			Threads:     4,
			Seed:        -1,
			Tokens:      4096,
			TopK:        20,
			Tfs:         1.0,
			TopP:        0.6,
			MinP:        0.1,
			TypicalP:    1.0,
			RepeatLastN: 64,
			RepeatPenalty: 1.05,
			FrequencyPenalty: 0.0,
			PresencePenalty: 0.0,
			Temperature: 0.7,
			Verbose:     false,
		}

		// Override with config values if present
		if t, ok := configData["threads"].(int); ok && t > 0 {
			options.Threads = t
		}
		if s, ok := configData["seed"].(int); ok {
			options.Seed = s
		}
		if tok, ok := configData["tokens"].(int); ok && tok > 0 {
			options.Tokens = tok
		}
		if tk, ok := configData["topK"].(int); ok && tk > 0 {
			options.TopK = tk
		}
		if tf, ok := configData["tfs"].(float64); ok {
			options.Tfs = tf
		}
		if tp, ok := configData["topP"].(float64); ok {
			options.TopP = tp
		}
		if mp, ok := configData["minP"].(float64); ok {
			options.MinP = mp
		}
		if typ, ok := configData["typicalP"].(float64); ok {
			options.TypicalP = typ
		}
		if rln, ok := configData["repeatLastN"].(int); ok && rln > 0 {
			options.RepeatLastN = rln
		}
		if rp, ok := configData["repeatPenalty"].(float64); ok {
			options.RepeatPenalty = rp
		}
		if fp, ok := configData["frequencyPenalty"].(float64); ok {
			options.FrequencyPenalty = fp
		}
		if pp, ok := configData["presencePenalty"].(float64); ok {
			options.PresencePenalty = pp
		}
		if temp, ok := configData["temperature"].(float64); ok {
			options.Temperature = temp
		}
		if v, ok := configData["verbose"].(bool); ok {
			options.Verbose = v
		}

		// Initialize llama library before creating translator
		if err := translator.InitLlamaLibrary(libPath); err != nil {
			return nil, fmt.Errorf("failed to initialize llama library: %v", err)
		}

		return translator.NewLlamaTranslator(options)

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

// PollForNextJob runs continuously to process jobs in the queue
func (qs *QueueService) PollForNextJob(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := qs.ProcessNextJob()
			if err != nil {
				// Log error but continue processing
				fmt.Printf("Error processing queue job: %v\n", err)
			}
		}
	}
}

// UpdateUserUsage updates the user's monthly translation usage
func (qs *QueueService) UpdateUserUsage(userID uint, count int) error {
	// Get the user
	var user database.User
	result := qs.DB.First(&user, userID)
	if result.Error != nil {
		return fmt.Errorf("failed to get user: %v", result.Error)
	}

	// Check if user has reached the limit
	if user.MaxTranslations > 0 && user.CurrentUsage+count > user.MaxTranslations {
		return errors.New("translation limit reached for this month")
	}

	// Update usage
	result = qs.DB.Model(&database.User{}).Where("id = ?", userID).UpdateColumn("current_usage", gorm.Expr("current_usage + ?", count))
	if result.Error != nil {
		return fmt.Errorf("failed to update user usage: %v", result.Error)
	}

	return nil
}

// Global functions for backward compatibility
func SetQueueService(db *database.Database) {
	queueServiceInstance = &QueueService{
		DB:                     db,
		queue:                    make(map[uint]*database.TranslationQueue),
		AppService:               appServiceInstance,
		AppLocalizationService:   appLocalizationServiceInstance,
		ProjectService:           projectServiceInstance,
		TranslationService:       translationServiceInstance,
		SubscriptionService:      subscriptionServiceInstance,
		ProviderService:          providerServiceInstance,
		AppProviderConfigService: appProviderConfigServiceInstance,
	}
}
