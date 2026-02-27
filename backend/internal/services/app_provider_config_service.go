package services

import (
	"errors"
	"fmt"

	"github.com/fdddf/opentrans/internal/dao/query"
	"github.com/fdddf/opentrans/internal/database"
	"gorm.io/gorm"
)

// AppProviderConfigService manages app-provider configuration bindings
type AppProviderConfigService struct {
	DB    *database.Database
	Query *query.Query
}

// BindAppProviderConfig creates or updates a binding between an app and a provider config.
func (s *AppProviderConfigService) BindAppProviderConfig(appID, configID uint, providerType string, isDefault bool) (*database.AppProviderConfig, error) {
	if isDefault {
		if err := s.ClearDefaultForApp(appID, providerType); err != nil {
			return nil, err
		}
	}

	binding, err := s.Query.AppProviderConfig.Where(
		s.Query.AppProviderConfig.AppID.Eq(appID),
		s.Query.AppProviderConfig.ProviderConfigID.Eq(configID),
	).First()
	
	if err == nil && binding != nil {
		_, err = s.Query.AppProviderConfig.Where(s.Query.AppProviderConfig.ID.Eq(binding.ID)).Updates(map[string]interface{}{
			"provider_type": providerType,
			"is_default":    isDefault,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update app provider config: %v", err)
		}
		binding.ProviderType = providerType
		binding.IsDefault = isDefault
		return binding, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check app provider config: %v", err)
	}

	binding = &database.AppProviderConfig{
		AppID:           appID,
		ProviderConfigID: configID,
		ProviderType:    providerType,
		IsDefault:       isDefault,
	}
	if err := s.Query.AppProviderConfig.Create(binding); err != nil {
		return nil, fmt.Errorf("failed to create app provider config: %v", err)
	}

	return binding, nil
}

// GetAppProviderConfig fetches a binding for an app/config pair.
func (s *AppProviderConfigService) GetAppProviderConfig(appID, configID uint) (*database.AppProviderConfig, error) {
	binding, err := s.Query.AppProviderConfig.Where(
		s.Query.AppProviderConfig.AppID.Eq(appID),
		s.Query.AppProviderConfig.ProviderConfigID.Eq(configID),
	).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("app provider config binding not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app provider config: %v", err)
	}

	return binding, nil
}

// GetDefaultAppProviderConfig fetches default binding by app+provider type.
func (s *AppProviderConfigService) GetDefaultAppProviderConfig(appID uint, providerType string) (*database.AppProviderConfig, error) {
	binding, err := s.Query.AppProviderConfig.Where(
		s.Query.AppProviderConfig.AppID.Eq(appID),
		s.Query.AppProviderConfig.ProviderType.Eq(providerType),
		s.Query.AppProviderConfig.IsDefault.Is(true),
	).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("no default app provider config binding found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch default app provider config: %v", err)
	}

	return binding, nil
}

// ClearDefaultForApp removes default flag for an app/provider type.
func (s *AppProviderConfigService) ClearDefaultForApp(appID uint, providerType string) error {
	_, err := s.Query.AppProviderConfig.Where(
		s.Query.AppProviderConfig.AppID.Eq(appID),
		s.Query.AppProviderConfig.ProviderType.Eq(providerType),
		s.Query.AppProviderConfig.IsDefault.Is(true),
	).Update(s.Query.AppProviderConfig.IsDefault, false)
	if err != nil {
		return fmt.Errorf("failed to clear default app provider config: %v", err)
	}
	return nil
}

// Global instance helpers
var appProviderConfigServiceInstance *AppProviderConfigService

func SetAppProviderConfigService(db *database.Database) {
	appProviderConfigServiceInstance = &AppProviderConfigService{
		DB:    db,
		Query: query.Use(db.DB),
	}
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