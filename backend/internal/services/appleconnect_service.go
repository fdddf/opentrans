package services

import (
	"fmt"
	"time"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/pkg/appleconnect"
)

// AppleConnectService wraps Apple Connect sync operations.
type AppleConnectService struct {
	AppService            *AppService
	AppLocalizationService *AppLocalizationService
}

// SyncApps syncs apps from Apple Connect into the DB.
func (s *AppleConnectService) SyncApps(userID uint, issuerID, keyID, privateKeyPath, privateKey string) ([]database.App, error) {
	client, err := appleconnect.NewAppleConnectClient(issuerID, keyID, privateKeyPath, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Apple Connect client: %v", err)
	}

	appsResponse, err := client.GetApps()
	if err != nil {
		return nil, fmt.Errorf("failed to get apps from Apple Connect: %v", err)
	}

	var syncedApps []database.App
	for _, appData := range appsResponse.Data {
		existingApp, err := s.AppService.GetAppByBundleID(appData.Attributes.BundleID)
		if err == nil && existingApp != nil {
			err = s.AppService.UpdateApp(existingApp.ID, map[string]interface{}{
				"Name":          appData.Attributes.Name,
				"AppleID":       appData.ID,
				"PrimaryLocale": appData.Attributes.PrimaryLocale,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update app %s: %v", appData.Attributes.BundleID, err)
			}
			syncedApps = append(syncedApps, *existingApp)
			continue
		}

		newApp, err := s.AppService.CreateApp(userID, appData.Attributes.Name, "", appData.Attributes.BundleID, appData.ID, appData.Attributes.PrimaryLocale)
		if err != nil {
			return nil, fmt.Errorf("failed to create app %s: %v", appData.Attributes.BundleID, err)
		}
		syncedApps = append(syncedApps, *newApp)
	}

	return syncedApps, nil
}

// SyncLocalizations syncs localizations from Apple Connect into the DB.
func (s *AppleConnectService) SyncLocalizations(appID uint, issuerID, keyID, privateKeyPath, privateKey string) ([]database.AppLocalization, error) {
	app, err := s.AppService.GetApp(appID)
	if err != nil {
		return nil, err
	}

	client, err := appleconnect.NewAppleConnectClient(issuerID, keyID, privateKeyPath, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Apple Connect client: %v", err)
	}

	localizationsResponse, err := client.GetAppLocalizations(app.AppleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get localizations from Apple Connect: %v", err)
	}

	var syncedLocalizations []database.AppLocalization
	for _, locData := range localizationsResponse.Data {
		existing, err := s.AppLocalizationService.GetAppLocalization(appID, locData.Attributes.Locale)
		if err == nil && existing != nil {
			err = s.AppLocalizationService.UpdateAppLocalization(appID, locData.Attributes.Locale, map[string]interface{}{
				"Name":                locData.Attributes.Name,
				"Subtitle":            locData.Attributes.Subtitle,
				"PrivacyURL":          locData.Attributes.PrivacyURL,
				"MarketingURL":        locData.Attributes.MarketingURL,
				"SupportURL":          locData.Attributes.SupportURL,
				"DownloadDescription": locData.Attributes.DownloadDescription,
				"ShortDescription":    locData.Attributes.ShortDescription,
				"LongDescription":     locData.Attributes.LongDescription,
				"Keywords":            locData.Attributes.Keywords,
				"ReleaseNotes":        locData.Attributes.ReleaseNotes,
				"Source":              "apple",
				"SyncStatus":          "synced",
				"SyncedAt":            time.Now(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update localization %s: %v", locData.Attributes.Locale, err)
			}
			syncedLocalizations = append(syncedLocalizations, *existing)
			continue
		}

		newLoc, err := s.AppLocalizationService.CreateAppLocalization(
			appID,
			locData.Attributes.Locale,
			locData.Attributes.Name,
			locData.Attributes.Subtitle,
			locData.Attributes.PrivacyURL,
			locData.Attributes.MarketingURL,
			locData.Attributes.SupportURL,
			locData.Attributes.DownloadDescription,
			locData.Attributes.ShortDescription,
			locData.Attributes.LongDescription,
			locData.Attributes.Keywords,
			locData.Attributes.ReleaseNotes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create localization %s: %v", locData.Attributes.Locale, err)
		}

		_ = s.AppLocalizationService.UpdateAppLocalization(appID, locData.Attributes.Locale, map[string]interface{}{
			"Source":     "apple",
			"SyncStatus": "synced",
			"SyncedAt":   time.Now(),
		})
		syncedLocalizations = append(syncedLocalizations, *newLoc)
	}

	return syncedLocalizations, nil
}
