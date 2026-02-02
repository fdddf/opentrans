package server

import (
	"io/fs"
	"net"
	"net/http"

	"github.com/fdddf/xcstrings-translator/internal/context"
	"github.com/fdddf/xcstrings-translator/internal/controllers"
	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/webui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Serve starts the Fiber server using the embedded UI assets.
func Serve(addr string, db *database.Database) error {
	app, err := NewApp(db)
	if err != nil {
		return err
	}
	return app.Listen(addr)
}

// ServeWithListener serves the app using a pre-bound listener (useful for GUI mode).
func ServeWithListener(ln net.Listener, db *database.Database) (*fiber.App, <-chan error, error) {
	app, err := NewApp(db)
	if err != nil {
		return nil, nil, err
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Listener(ln)
	}()
	return app, errCh, nil
}

func NewApp(db *database.Database) (*fiber.App, error) {
	distFS, err := fs.Sub(webui.EmbeddedFS, "dist")
	if err != nil {
		return nil, err
	}

	app := fiber.New()

	// Use middleware to inject database into context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize controllers
	authController := controllers.NewAuthController()
	projectController := controllers.NewProjectController()
	providerController := controllers.NewProviderController()
	appController := controllers.NewAppController()
	appLocalizationController := controllers.NewAppLocalizationController()
	appleConnectController := controllers.NewAppleConnectController()
	subscriptionController := controllers.NewSubscriptionController()
	translationController := controllers.NewTranslationController()
	fileController := controllers.NewFileController()
	userActivityController := controllers.NewUserActivityController()
	userController := controllers.NewUserController()

	// Public routes (no authentication required)
	api := app.Group("/api")
	api.Post("/upload", fileController.HandleUpload)
	api.Get("/strings", fileController.HandleStrings)
	api.Post("/translate", fileController.HandleTranslate)
	api.Get("/progress", fileController.HandleProgress)
	api.Get("/export", fileController.HandleExport)

	// Authentication routes
	api.Post("/auth/register", authController.Register)
	api.Post("/auth/login", authController.Login)
	api.Post("/auth/logout", authController.Logout)
	api.Post("/auth/activate/:code", authController.ActivateUser)

	// Protected routes (require authentication)
	protected := api.Group("/protected")
	protected.Use(context.AuthMiddleware)
	adminOnly := protected.Group("/admin")
	adminOnly.Use(context.AdminOnly)
	protected.Get("/languages", fileController.HandleGetSupportedLanguages)
	protected.Post("/translate/text", fileController.HandleTranslateText)

	// Example admin route to verify wiring (extend as needed)
	adminOnly.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	// Project routes
	protected.Get("/projects", projectController.GetProjects)
	protected.Post("/projects", projectController.CreateProject)
	protected.Get("/projects/:id", projectController.GetProject)
	protected.Put("/projects/:id", projectController.UpdateProject)
	protected.Delete("/projects/:id", projectController.DeleteProject)
	protected.Post("/projects/:id/upload", projectController.UploadToProject)
	protected.Post("/projects/:id/translate", context.SubscriptionRequired, projectController.TranslateProject)
	protected.Get("/projects/:id/export", projectController.ExportProject)
	protected.Get("/projects/:id/translations", projectController.GetProjectTranslations)
	protected.Get("/projects/:id/missing-translations", projectController.GetMissingTranslations)
	protected.Get("/projects/:id/translation-status", projectController.GetTranslationStatus)
	protected.Get("/projects/:id/stats", projectController.GetProjectStats)
	protected.Put("/projects/:id/translations/:key/:language", projectController.UpdateSingleTranslation)
	protected.Put("/projects/:id/translations/bulk", projectController.BulkUpdateTranslations)
	protected.Get("/projects/:id/languages", projectController.GetProjectLanguages)

	// Provider configuration routes
	protected.Get("/providers", providerController.GetProviderConfigs)
	protected.Post("/providers", providerController.CreateProviderConfig)
	protected.Get("/providers/:id", providerController.GetProviderConfig)
	protected.Put("/providers/:id", providerController.UpdateProviderConfig)
	protected.Delete("/providers/:id", providerController.DeleteProviderConfig)
	protected.Get("/providers/:type/default", providerController.GetDefaultProviderConfig)

	// Activity logging routes
	protected.Get("/activities", userActivityController.GetActivities(db))
	adminOnly.Get("/activities", userActivityController.GetAdminActivities(db))

	// App management routes
	protected.Get("/apps", appController.GetApps)
	protected.Post("/apps", context.SubscriptionRequired, appController.CreateApp)
	protected.Get("/apps/:id", appController.GetApp)
	protected.Put("/apps/:id", appController.UpdateApp)
	protected.Delete("/apps/:id", appController.DeleteApp)
	protected.Get("/apps/:id/users", appController.GetAppUsers)
	protected.Post("/apps/:id/users", appController.AddUserToApp)
	protected.Put("/apps/:id/users/:userId", appController.UpdateUserAppRole)
	protected.Delete("/apps/:id/users/:userId", appController.RemoveUserFromApp)
	protected.Get("/apps/:id/stats", appController.GetAppStats)
	protected.Put("/apps/:id/localizations/bulk", appController.BulkUpdateAppLocalizations)
	protected.Post("/apps/:id/sync-to-apple", appController.SyncAppToApple)

	// App localization routes
	protected.Get("/apps/:id/localizations", appLocalizationController.GetAppLocalizations)
	protected.Post("/apps/:id/localizations", appLocalizationController.CreateAppLocalization)
	protected.Get("/apps/:id/localizations/:language", appLocalizationController.GetAppLocalization)
	protected.Put("/apps/:id/localizations/:language", appLocalizationController.UpdateAppLocalization)
	protected.Delete("/apps/:id/localizations/:language", appLocalizationController.DeleteAppLocalization)
	protected.Get("/apps/:id/languages", appLocalizationController.GetAppLanguages)
	protected.Post("/apps/:id/languages", appLocalizationController.AddAppLanguage)
	protected.Delete("/apps/:id/languages/:language", appLocalizationController.RemoveAppLanguage)

	// Apple Connect API routes
	protected.Get("/appleconnections", appleConnectController.GetAppleConnectConfigs)
	protected.Post("/appleconnections", appleConnectController.CreateAppleConnectConfig)
	protected.Get("/appleconnections/:id", appleConnectController.GetAppleConnectConfig)
	protected.Put("/appleconnections/:id", appleConnectController.UpdateAppleConnectConfig)
	protected.Delete("/appleconnections/:id", appleConnectController.DeleteAppleConnectConfig)
	protected.Post("/appleconnections/:id/test", appleConnectController.TestAppleConnectConnection)
	protected.Post("/appleconnections/test", appleConnectController.TestAppleConnectCredentials)
	protected.Post("/apple-connect/sync-apps", appleConnectController.SyncAppleApps)
	protected.Post("/apple-connect/:appId/sync-localizations", appleConnectController.SyncAppleAppLocalizations)

	// Subscription management routes
	protected.Get("/subscription", subscriptionController.GetUserSubscription)
	protected.Post("/subscription/webhook", subscriptionController.SubscriptionWebhook)
	protected.Get("/subscription/usage", subscriptionController.GetUsage)

	// User management routes
	protected.Get("/users/:id", userController.GetUser(db))
	protected.Put("/users/:id", userController.UpdateUser(db))
	protected.Put("/user/profile", userController.UpdateCurrentUser(db))
	protected.Post("/change-password", userController.ChangePassword(db))

	// Translation queue routes
	protected.Post("/queue/translate", context.SubscriptionRequired, translationController.QueueTranslationJob)
	protected.Get("/queue/jobs", translationController.GetQueueJobs)
	protected.Get("/queue/jobs/:id", translationController.GetQueueJob)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(distFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	return app, nil
}
