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

	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"` // hashed password
	IsActive bool   `gorm:"default:true" json:"is_active"`
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
