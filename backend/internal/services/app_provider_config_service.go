package services

import (
	"errors"
	"fmt"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"gorm.io/gorm"
)

// AppProviderConfigService manages app-provider configuration bindings
type AppProviderConfigService struct {
	DB *database.Database
}

// BindAppProviderConfig creates or updates a binding between an app and a provider config.
func (s *AppProviderConfigService) BindAppProviderConfig(appID, configID uint, providerType string, isDefault bool) (*database.AppProviderConfig, error) {
	if isDefault {
		if err := s.ClearDefaultForApp(appID, providerType); err != nil {
			return nil, err
		}
	}

	var binding database.AppProviderConfig
	result := s.DB.Where("app_id = ? AND provider_config_id = ?", appID, configID).First(&binding)
	if result.Error == nil {
		updates := map[string]interface{}{
			"provider_type": providerType,
			"is_default":    isDefault,
		}
		if err := s.DB.Model(&database.AppProviderConfig{}).Where("id = ?", binding.ID).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update app provider config: %v", err)
		}
		binding.ProviderType = providerType
		binding.IsDefault = isDefault
		return &binding, nil
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check app provider config: %v", result.Error)
	}

	binding = database.AppProviderConfig{
		AppID:           appID,
		ProviderConfigID: configID,
		ProviderType:    providerType,
		IsDefault:       isDefault,
	}
	if err := s.DB.Create(&binding).Error; err != nil {
		return nil, fmt.Errorf("failed to create app provider config: %v", err)
	}

	return &binding, nil
}

// GetAppProviderConfig fetches a binding for an app/config pair.
func (s *AppProviderConfigService) GetAppProviderConfig(appID, configID uint) (*database.AppProviderConfig, error) {
	var binding database.AppProviderConfig
	result := s.DB.Where("app_id = ? AND provider_config_id = ?", appID, configID).First(&binding)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("app provider config binding not found")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch app provider config: %v", result.Error)
	}

	return &binding, nil
}

// GetDefaultAppProviderConfig fetches default binding by app+provider type.
func (s *AppProviderConfigService) GetDefaultAppProviderConfig(appID uint, providerType string) (*database.AppProviderConfig, error) {
	var binding database.AppProviderConfig
	result := s.DB.Where("app_id = ? AND provider_type = ? AND is_default = ?", appID, providerType, true).First(&binding)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("no default app provider config binding found")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch default app provider config: %v", result.Error)
	}

	return &binding, nil
}

// ClearDefaultForApp removes default flag for an app/provider type.
func (s *AppProviderConfigService) ClearDefaultForApp(appID uint, providerType string) error {
	result := s.DB.Model(&database.AppProviderConfig{}).
		Where("app_id = ? AND provider_type = ? AND is_default = ?", appID, providerType, true).
		Updates(map[string]interface{}{"is_default": false})
	if result.Error != nil {
		return fmt.Errorf("failed to clear default app provider config: %v", result.Error)
	}
	return nil
}

// Global instance helpers
var appProviderConfigServiceInstance *AppProviderConfigService

func SetAppProviderConfigService(db *database.Database) {
	appProviderConfigServiceInstance = &AppProviderConfigService{DB: db}
}

func BindAppProviderConfig(appID, configID uint, providerType string, isDefault bool) (*database.AppProviderConfig, error) {
	return appProviderConfigServiceInstance.BindAppProviderConfig(appID, configID, providerType, isDefault)
}

func GetAppProviderConfig(appID, configID uint) (*database.AppProviderConfig, error) {
	return appProviderConfigServiceInstance.GetAppProviderConfig(appID, configID)
}

func GetDefaultAppProviderConfig(appID uint, providerType string) (*database.AppProviderConfig, error) {
	return appProviderConfigServiceInstance.GetDefaultAppProviderConfig(appID, providerType)
}

func ClearDefaultForAppProvider(appID uint, providerType string) error {
	return appProviderConfigServiceInstance.ClearDefaultForApp(appID, providerType)
}
