package services

import (
	"sync"

	"go.uber.org/fx"

	"github.com/fdddf/xcstrings-translator/internal/auth"
	"github.com/fdddf/xcstrings-translator/internal/database"
)

// Module is the FX module for services
var Module = fx.Module("services",
	fx.Provide(NewAppService),
	fx.Provide(NewAppLocalizationService),
	fx.Provide(NewProjectService),
	fx.Provide(NewProviderService),
	fx.Provide(NewQueueService),
	fx.Provide(NewSubscriptionService),
	fx.Provide(NewTranslationService),
	fx.Invoke(InitializeServices),
)

// ServiceParams holds the common dependencies for services
type ServiceParams struct {
	fx.In

	DB *database.Database
}

// AppServiceDeps holds dependencies for AppService
type AppServiceDeps struct {
	fx.In
	DB                     *database.Database
	AppLocalizationService *AppLocalizationService
}

// ProjectServiceDeps holds dependencies for ProjectService
type ProjectServiceDeps struct {
	fx.In
	DB                     *database.Database
	TranslationService     *TranslationService
	AppLocalizationService *AppLocalizationService
}

// QueueServiceDeps holds dependencies for QueueService
type QueueServiceDeps struct {
	fx.In
	DB                     *database.Database
	AppService             *AppService
	AppLocalizationService *AppLocalizationService
	ProjectService         *ProjectService
	TranslationService     *TranslationService
	SubscriptionService    *SubscriptionService
	ProviderService        *ProviderService
}

// NewAppService creates a new AppService with database dependency
func NewAppService(deps AppServiceDeps) *AppService {
	return &AppService{
		DB:                     deps.DB,
		AppLocalizationService: deps.AppLocalizationService,
	}
}

// NewAppLocalizationService creates a new AppLocalizationService with database dependency
func NewAppLocalizationService(p ServiceParams) *AppLocalizationService {
	return &AppLocalizationService{
		DB: p.DB,
	}
}

// NewProjectService creates a new ProjectService with database dependency
func NewProjectService(deps ProjectServiceDeps) *ProjectService {
	return &ProjectService{
		DB:                     deps.DB,
		TranslationService:     deps.TranslationService,
		AppLocalizationService: deps.AppLocalizationService,
	}
}

// NewProviderService creates a new ProviderService with database dependency
func NewProviderService(p ServiceParams) *ProviderService {
	return &ProviderService{
		DB: p.DB,
	}
}

// NewQueueService creates a new QueueService with database dependency
func NewQueueService(deps QueueServiceDeps) *QueueService {
	return &QueueService{
		DB:                     deps.DB,
		AppService:             deps.AppService,
		AppLocalizationService: deps.AppLocalizationService,
		ProjectService:         deps.ProjectService,
		TranslationService:     deps.TranslationService,
		SubscriptionService:    deps.SubscriptionService,
		ProviderService:        deps.ProviderService,
		mu:                     sync.RWMutex{},
		queue:                  make(map[uint]*database.TranslationQueue),
	}
}

// NewSubscriptionService creates a new SubscriptionService with database dependency
func NewSubscriptionService(p ServiceParams) *SubscriptionService {
	return &SubscriptionService{
		DB: p.DB,
	}
}

// NewTranslationService creates a new TranslationService with database dependency
func NewTranslationService(p ServiceParams) *TranslationService {
	return &TranslationService{
		DB: p.DB,
	}
}

// InitializeServices sets up the global service instances for backward compatibility
func InitializeServices(
	db *database.Database,
	authService *auth.Auth,
	appService *AppService,
	appLocalizationService *AppLocalizationService,
	projectService *ProjectService,
	providerService *ProviderService,
	queueService *QueueService,
	subscriptionService *SubscriptionService,
	translationService *TranslationService,
) {
	// Initialize the auth instance for backward compatibility
	auth.SetAuthInstance(authService)

	// Set all services for backward compatibility
	SetAppLocalizationService(db)
	SetTranslationService(db)
	SetProviderService(db)
	SetSubscriptionService(db)
	SetProjectService(db, translationService, appLocalizationService)
	SetAppService(db, appLocalizationService)
	SetQueueService(db)
}
