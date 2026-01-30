package services

import (
	"errors"
	"fmt"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"gorm.io/gorm"
)

// ProviderService handles provider configuration operations
type ProviderService struct {
	DB *database.Database
}

// CreateProviderConfig creates a new provider configuration for a user
func (s *ProviderService) CreateProviderConfig(userID uint, providerType string, configData map[string]interface{}, isDefault bool) (*database.ProviderConfig, error) {
	// If setting this as default, unset any existing default for this provider type
	if isDefault {
		err := s.UnsetDefaultProviderConfig(userID, providerType)
		if err != nil {
			return nil, fmt.Errorf("failed to unset existing default: %v", err)
		}
	}

	config := &database.ProviderConfig{
		UserID:       userID,
		ProviderType: providerType,
		ConfigData:   configData,
		IsDefault:    isDefault,
	}

	result := s.DB.Create(config)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create provider config: %v", result.Error)
	}

	return config, nil
}

// GetProviderConfig retrieves a specific provider configuration by ID
func (s *ProviderService) GetProviderConfig(configID uint) (*database.ProviderConfig, error) {
	var config database.ProviderConfig
	result := s.DB.First(&config, configID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("provider configuration not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &config, nil
}

// GetProviderConfigsByUser retrieves all provider configurations for a user
func (s *ProviderService) GetProviderConfigsByUser(userID uint) ([]database.ProviderConfig, error) {
	var configs []database.ProviderConfig
	result := s.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&configs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve provider configs: %v", result.Error)
	}

	return configs, nil
}

// GetProviderConfigsByUserAndType retrieves provider configurations for a user filtered by provider type
func (s *ProviderService) GetProviderConfigsByUserAndType(userID uint, providerType string) ([]database.ProviderConfig, error) {
	var configs []database.ProviderConfig
	result := s.DB.Where("user_id = ? AND provider_type = ?", userID, providerType).Order("created_at DESC").Find(&configs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve provider configs: %v", result.Error)
	}

	return configs, nil
}

// GetDefaultProviderConfig retrieves the default provider configuration for a user and provider type
func (s *ProviderService) GetDefaultProviderConfig(userID uint, providerType string) (*database.ProviderConfig, error) {
	var config database.ProviderConfig
	result := s.DB.Where("user_id = ? AND provider_type = ? AND is_default = ?", userID, providerType, true).First(&config)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("no default provider configuration found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &config, nil
}

// UpdateProviderConfig updates an existing provider configuration
func (s *ProviderService) UpdateProviderConfig(configID uint, updates map[string]interface{}) error {
	// Check if the config ID exists
	var existing database.ProviderConfig
	result := s.DB.First(&existing, configID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("provider configuration not found")
	}

	if result.Error != nil {
		return fmt.Errorf("database error: %v", result.Error)
	}

	// If updating IsDefault to true, unset any existing default for this provider type
	if isDefault, ok := updates["IsDefault"].(bool); ok && isDefault {
		err := s.UnsetDefaultProviderConfig(existing.UserID, existing.ProviderType)
		if err != nil {
			return fmt.Errorf("failed to unset existing default: %v", err)
		}
	}

	result = s.DB.Model(&database.ProviderConfig{}).Where("id = ?", configID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update provider config: %v", result.Error)
	}

	return nil
}

// DeleteProviderConfig deletes a provider configuration
func (s *ProviderService) DeleteProviderConfig(configID uint) error {
	result := s.DB.Delete(&database.ProviderConfig{}, configID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete provider config: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("provider configuration not found")
	}

	return nil
}

// UnsetDefaultProviderConfig sets all configurations of a specific type for a user to not be default
func (s *ProviderService) UnsetDefaultProviderConfig(userID uint, providerType string) error {
	result := s.DB.Model(&database.ProviderConfig{}).
		Where("user_id = ? AND provider_type = ? AND is_default = ?", userID, providerType, true).
		Updates(map[string]interface{}{"is_default": false})

	if result.Error != nil {
		return fmt.Errorf("failed to unset default provider config: %v", result.Error)
	}

	return nil
}

// ValidateProviderConfigData validates the structure of provider configuration data
func (s *ProviderService) ValidateProviderConfigData(providerType string, configData map[string]interface{}) error {
	switch providerType {
	case "openai":
		// Required fields: apiKey
		if apiKey, exists := configData["apiKey"]; !exists || apiKey == "" {
			return errors.New("apiKey is required for OpenAI provider")
		}
	case "google":
		// Required fields: apiKey
		if apiKey, exists := configData["apiKey"]; !exists || apiKey == "" {
			return errors.New("apiKey is required for Google provider")
		}
	case "deepl":
		// Required fields: apiKey
		if apiKey, exists := configData["apiKey"]; !exists || apiKey == "" {
			return errors.New("apiKey is required for DeepL provider")
		}
	case "baidu":
		// Required fields: appId, appSecret
		if appID, exists := configData["appId"]; !exists || appID == "" {
			return errors.New("appId is required for Baidu provider")
		}
		if appSecret, exists := configData["appSecret"]; !exists || appSecret == "" {
			return errors.New("appSecret is required for Baidu provider")
		}
	case "appleconnect":
		// Required fields: issuerID, keyID, privateKey
		if issuerID, exists := configData["issuerID"]; !exists || issuerID == "" {
			return errors.New("issuerID is required for Apple Connect provider")
		}
		if keyID, exists := configData["keyID"]; !exists || keyID == "" {
			return errors.New("keyID is required for Apple Connect provider")
		}
		if privateKey, exists := configData["privateKey"]; !exists || privateKey == "" {
			return errors.New("privateKey is required for Apple Connect provider")
		}
	case "llama":
		// Required fields: modelPath, libPath
		if modelPath, exists := configData["modelPath"]; !exists || modelPath == "" {
			return errors.New("modelPath is required for Llama provider")
		}
		if libPath, exists := configData["libPath"]; !exists || libPath == "" {
			return errors.New("libPath is required for Llama provider")
		}
	default:
		return fmt.Errorf("unknown provider type: %s", providerType)
	}

	return nil
}

// SanitizeProviderConfig removes sensitive fields from config data before returning to client
func (s *ProviderService) SanitizeProviderConfig(config *database.ProviderConfig) *database.ProviderConfig {
	sanitized := *config

	// Create a copy of ConfigData to avoid modifying the original
	safeConfigData := make(map[string]interface{})
	for key, value := range config.ConfigData {
		// Don't include sensitive fields in the response
		if isSensitiveField(key) {
			safeConfigData[key] = "***REDACTED***" // Show placeholder for sensitive values
		} else {
			safeConfigData[key] = value
		}
	}

	sanitized.ConfigData = safeConfigData
	return &sanitized
}

// isSensitiveField checks if a configuration field is sensitive
func isSensitiveField(field string) bool {
	sensitiveFields := []string{
		"apiKey", "api_key", "password", "secret", "appSecret", "app_secret",
		"token", "access_token", "refresh_token", "private_key", "client_secret",
	}

	for _, sensitiveField := range sensitiveFields {
		if field == sensitiveField {
			return true
		}
	}

	return false
}

// GetProviderConfigForTranslation prepares configuration data specifically for translation services
func (s *ProviderService) GetProviderConfigForTranslation(config *database.ProviderConfig) map[string]interface{} {
	// Return the config data without sanitization for use in translation services
	// This is safe since it's only used internally in the translation process
	return config.ConfigData
}

// Global functions for backward compatibility
var providerServiceInstance *ProviderService

func SetProviderService(db *database.Database) {
	providerServiceInstance = &ProviderService{DB: db}
}

func CreateProviderConfig(userID uint, providerType string, configData map[string]interface{}, isDefault bool) (*database.ProviderConfig, error) {
	return providerServiceInstance.CreateProviderConfig(userID, providerType, configData, isDefault)
}

func GetProviderConfig(configID uint) (*database.ProviderConfig, error) {
	return providerServiceInstance.GetProviderConfig(configID)
}

func GetProviderConfigsByUser(userID uint) ([]database.ProviderConfig, error) {
	return providerServiceInstance.GetProviderConfigsByUser(userID)
}

func GetProviderConfigsByUserAndType(userID uint, providerType string) ([]database.ProviderConfig, error) {
	return providerServiceInstance.GetProviderConfigsByUserAndType(userID, providerType)
}

func GetDefaultProviderConfig(userID uint, providerType string) (*database.ProviderConfig, error) {
	return providerServiceInstance.GetDefaultProviderConfig(userID, providerType)
}

func UpdateProviderConfig(configID uint, updates map[string]interface{}) error {
	return providerServiceInstance.UpdateProviderConfig(configID, updates)
}

func DeleteProviderConfig(configID uint) error {
	return providerServiceInstance.DeleteProviderConfig(configID)
}

func UnsetDefaultProviderConfig(userID uint, providerType string) error {
	return providerServiceInstance.UnsetDefaultProviderConfig(userID, providerType)
}

func ValidateProviderConfigData(providerType string, configData map[string]interface{}) error {
	return providerServiceInstance.ValidateProviderConfigData(providerType, configData)
}

func SanitizeProviderConfig(config *database.ProviderConfig) *database.ProviderConfig {
	return providerServiceInstance.SanitizeProviderConfig(config)
}

func GetProviderConfigForTranslation(config *database.ProviderConfig) map[string]interface{} {
	return providerServiceInstance.GetProviderConfigForTranslation(config)
}
