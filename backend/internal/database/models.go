package database

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Username       string `gorm:"uniqueIndex;not null" json:"username"`
	Email          string `gorm:"uniqueIndex;not null" json:"email"`
	Password       string `gorm:"not null" json:"password"` // hashed password
	IsActive       bool   `gorm:"default:true" json:"is_active"`
	IsActivated    bool   `gorm:"default:false" json:"is_activated"`
	ActivationCode string `json:"activation_code,omitempty"`

	// Role and subscription fields
	Role string `gorm:"type:varchar(20);default:'user'" json:"role"` // admin, user

	IsSubscribed     bool       `json:"is_subscribed"`
	SubscriptionType string     `json:"subscription_type"` // free, basic, premium
	SubscriptionEnd  *time.Time `json:"subscription_end,omitempty"`

	// Usage limits based on subscription
	MaxApps         int `json:"max_apps"`          // Max number of apps allowed
	MaxTranslations int `json:"max_translations"`  // Max number of translations per month
	CurrentUsage    int `json:"current_usage"`     // Current monthly usage
	CurrentAppCount int `json:"current_app_count"` // Current number of apps
}

// Project represents an xcstrings project
type Project struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	UserID      uint   `gorm:"not null" json:"user_id"` // Foreign key to User
	User        User   `json:"user"`

	// Store the original xcstrings file content
	FileContent string `gorm:"type:text" json:"file_content"`
	FileName    string `json:"file_name"`

	// Source language for this project
	SourceLanguage string `json:"source_language"`

	// JSON field to store the parsed content structure
	ContentStructure map[string]interface{} `gorm:"type:jsonb" json:"content_structure"`

	// Store project settings
	Settings map[string]interface{} `gorm:"type:jsonb" json:"settings"`
}

// Translation represents a translated string
type Translation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign key to Project
	ProjectID uint    `gorm:"not null" json:"project_id"`
	Project   Project `json:"project"`

	// The original key in the xcstrings file
	Key string `gorm:"not null" json:"key"`

	// The source text
	SourceText string `gorm:"not null" json:"source_text"`

	// The translated text
	TargetText string `json:"target_text"`

	// Target language code (e.g., "zh-Hans", "ja", etc.)
	TargetLanguage string `gorm:"not null" json:"target_language"`

	// State of the translation (e.g., "translated", "needs_review", etc.)
	State string `json:"state"`

	// Optional: translation provider used
	TranslationProvider string `json:"translation_provider"`
}

// UserActivity represents user activity logs
type UserActivity struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID    uint   `gorm:"not null" json:"user_id"` // Foreign key to User
	User      User   `json:"user"`
	Action    string `gorm:"not null" json:"action"`
	Details   string `json:"details"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

// ProviderConfig represents user-specific translation provider configurations
type ProviderConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID uint `gorm:"not null" json:"user_id"` // Foreign key to User
	User   User `json:"user"`

	// Provider type (e.g., "google", "deepl", "baidu", "openai")
	ProviderType string `gorm:"not null" json:"provider_type"`

	// Configuration settings (stored as JSON for flexibility)
	ConfigData map[string]interface{} `gorm:"type:jsonb" json:"config_data"`

	// Whether this configuration is the default one
	IsDefault bool `json:"is_default"`
}

// App represents an iOS application in the system
type App struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	UserID      uint   `gorm:"not null" json:"user_id"` // Foreign key to User
	User        User   `json:"user"`

	BundleID          string `gorm:"not null;uniqueIndex" json:"bundle_id"`           // App's bundle identifier
	AppleID           string `json:"apple_id"`                                        // App's Apple ID
	PrimaryLocale     string `json:"primary_locale"`                                  // Primary language of the app
	AppleConnectToken string `json:"apple_connect_token"`                             // Token for connecting to App Store Connect
	Origin            string `gorm:"type:varchar(20);default:'manual'" json:"origin"` // manual, synced

	// App metadata
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	Keywords         string `json:"keywords"` // Comma-separated keywords
	SupportURL       string `json:"support_url"`
	MarketingURL     string `json:"marketing_url"`
	PrivacyURL       string `json:"privacy_url"`

	// App Store status
	Version          string `json:"version"`
	AppCategory      string `json:"app_category"`
	IsReadyForReview bool   `json:"is_ready_for_review"`
}

// AppLocalization represents localization data for an app in App Store Connect
type AppLocalization struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	AppID uint `gorm:"not null" json:"app_id"` // Foreign key to App
	App   App  `json:"app"`

	LanguageCode string `gorm:"not null" json:"language_code"` // Language code (e.g., "en-US", "zh-Hans")

	// App Store Connect localization fields
	Name                string `json:"name"`                 // App name in this language
	Subtitle            string `json:"subtitle"`             // App subtitle in this language
	PrivacyURL          string `json:"privacy_url"`          // Privacy URL in this language
	MarketingURL        string `json:"marketing_url"`        // Marketing URL in this language
	SupportURL          string `json:"support_url"`          // Support URL in this language
	DownloadDescription string `json:"download_description"` // Download description in this language
	ShortDescription    string `json:"short_description"`    // Short description in this language
	LongDescription     string `json:"long_description"`     // Long description in this language
	Keywords            string `json:"keywords"`             // Keywords in this language (comma-separated)
	ReleaseNotes        string `json:"release_notes"`        // Release notes in this language
	SyncedAt            *time.Time `json:"synced_at"`
	Source              string     `gorm:"type:varchar(20);default:'local'" json:"source"`
	SyncStatus          string     `gorm:"type:varchar(20);default:'pending'" json:"sync_status"`
}

// Subscription represents user subscription information
type Subscription struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID uint `gorm:"not null" json:"user_id"` // Foreign key to User
	User   User `json:"user"`

	// Subscription details
	StripeSubscriptionID string     `json:"stripe_subscription_id"` // Stripe subscription ID
	StripeCustomerID     string     `json:"stripe_customer_id"`     // Stripe customer ID
	SubscriptionType     string     `json:"subscription_type"`      // free, basic, premium
	SubscriptionStatus   string     `json:"subscription_status"`    // active, canceled, past_due, unpaid
	CurrentPeriodStart   time.Time  `json:"current_period_start"`   // Current billing period start
	CurrentPeriodEnd     time.Time  `json:"current_period_end"`     // Current billing period end
	TrialEnd             *time.Time `json:"trial_end,omitempty"`    // When trial period ends (if any)
	CancelAtPeriodEnd    bool       `json:"cancel_at_period_end"`   // Whether to cancel at period end
}

// AppUser represents the many-to-many relationship between users and apps for collaboration
type AppUser struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	AppID  uint `gorm:"not null" json:"app_id"`
	UserID uint `gorm:"not null" json:"user_id"`

	// Role in the app: owner, admin, editor, viewer
	Role string `gorm:"default:'viewer'" json:"role"`

	App  App  `json:"app"`
	User User `json:"user"`
}

// AppProviderConfig represents the binding between an app and a provider configuration
type AppProviderConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	AppID           uint          `gorm:"not null" json:"app_id"`
	ProviderConfigID uint          `gorm:"not null" json:"provider_config_id"`
	ProviderType    string        `gorm:"not null" json:"provider_type"`
	IsDefault       bool          `json:"is_default"`
	App             App           `json:"app"`
	ProviderConfig  ProviderConfig `json:"provider_config"`
}

// TranslationQueue represents items in the translation queue for batch processing
type TranslationQueue struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign keys
	UserID    uint  `gorm:"not null" json:"user_id"` // User who submitted the job
	ProjectID *uint `json:"project_id,omitempty"`    // Optional: if related to a project
	AppID     *uint `json:"app_id,omitempty"`        // Optional: if related to an app

	// Queue job details
	Type     string `gorm:"not null" json:"type"`      // "xcstrings", "app_localization"
	Status   string `gorm:"not null" json:"status"`    // "pending", "processing", "completed", "failed"
	Priority int    `gorm:"default:0" json:"priority"` // Higher number = higher priority
	Progress int    `json:"progress"`                  // Percentage complete
	Total    int    `json:"total"`                     // Total items to process
	Done     int    `json:"done"`                      // Items completed
	Error    string `json:"error,omitempty"`           // Error message if failed

	// Configuration
	ProviderType    string                 `gorm:"serializer:json" json:"provider_type"`    // "google", "deepl", etc.
	SourceLanguage  string                 `gorm:"serializer:json" json:"source_language"`  // Source language code
	TargetLanguages []string               `gorm:"serializer:json" json:"target_languages"` // Target language codes
	ConfigData      map[string]interface{} `gorm:"serializer:json" json:"config_data"`      // Provider-specific config

	// Result data
	ResultData map[string]interface{} `gorm:"type:jsonb" json:"result_data"` // Result of translation job
}
