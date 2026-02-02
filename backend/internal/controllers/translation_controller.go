package controllers

import (
	"github.com/fdddf/xcstrings-translator/internal/context"
	"github.com/fdddf/xcstrings-translator/internal/services"
	"github.com/gofiber/fiber/v2"
)

// TranslationController handles translation queue requests
type TranslationController struct{}

// NewTranslationController creates a new TranslationController
func NewTranslationController() *TranslationController {
	return &TranslationController{}
}

// QueueTranslationJobRequest represents the queue translation job request
type QueueTranslationJobRequest struct {
	JobType         string                 `json:"jobType"`
	ProjectID       *uint                  `json:"projectId"`
	AppID           *uint                  `json:"appId"`
	ProviderType    string                 `json:"providerType"`
	SourceLanguage  string                 `json:"sourceLanguage"`
	TargetLanguages []string               `json:"targetLanguages"`
	ConfigData      map[string]interface{} `json:"configData"`
}

// QueueTranslationJob queues a translation job
func (ctrl *TranslationController) QueueTranslationJob(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req QueueTranslationJobRequest
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

// GetQueueJobs gets all queue jobs for the current user
func (ctrl *TranslationController) GetQueueJobs(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	queueService := services.GetQueueService()
	jobs, err := queueService.GetQueueJobsByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"jobs":    jobs,
	})
}

// GetQueueJob gets a specific queue job by ID
func (ctrl *TranslationController) GetQueueJob(c *fiber.Ctx) error {
	userID, ok := context.GetUserIDFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	jobID, err := c.ParamsInt("id", 0)
	if err != nil || jobID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid job ID")
	}

	queueService := services.GetQueueService()
	job, err := queueService.GetQueueJob(uint(jobID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Queue job not found")
	}

	// Verify user has access to this job
	if job.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied to this job")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"job":     job,
	})
}