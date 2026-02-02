package appleconnect

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AppleConnectClient handles communication with Apple App Store Connect API
type AppleConnectClient struct {
	IssuerID       string
	KeyID          string
	PrivateKeyPath string
	PrivateKey     string
	privateKey     interface{} // Can be *rsa.PrivateKey or *ecdsa.PrivateKey
	baseURL        string
	httpClient     *http.Client
}

// App represents an app in the App Store Connect API
type App struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name          string `json:"name"`
		BundleID      string `json:"bundleId"`
		Sku           string `json:"sku"`
		PrimaryLocale string `json:"primaryLocale"`
	} `json:"attributes"`
}

// AppLocalization represents localization data for an app in App Store Connect
type AppLocalization struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name                string `json:"name"`
		Subtitle            string `json:"subtitle"`
		PrivacyURL          string `json:"privacyUrl"`
		MarketingURL        string `json:"marketingUrl"`
		SupportURL          string `json:"supportUrl"`
		DownloadDescription string `json:"downloadDescription"`
		Description         string `json:"description"`
		ShortDescription    string `json:"shortDescription"`
		Keywords            string `json:"keywords"`
		WhatsNew            string `json:"whatsNew"`
		PromotionalText     string `json:"promotionalText"`
		Locale              string `json:"locale"`
	} `json:"attributes"`
}

// AppInfo represents global app information (name, subtitle)
type AppInfo struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name     string `json:"name"`
		BundleID string `json:"bundleId"`
		Sku      string `json:"sku"`
	} `json:"attributes"`
	Relationships struct {
		AppStoreVersions struct {
			Data []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
		} `json:"appStoreVersions"`
	} `json:"relationships"`
}

// AppStoreVersion represents an app store version in App Store Connect
type AppStoreVersion struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		VersionString string `json:"versionString"`
		Platform      string `json:"platform"`
		State         string `json:"appStoreState"`
		WhatsNew      string `json:"whatsNew"`
	} `json:"attributes"`
}

// AppsResponse represents the response from the apps endpoint
type AppsResponse struct {
	Data []App `json:"data"`
}

// AppLocalizationsResponse represents the response from the app localizations endpoint
type AppLocalizationsResponse struct {
	Data []AppLocalization `json:"data"`
}

// AppLocalizationResponse represents the response for a single app localization
type AppLocalizationResponse struct {
	Data AppLocalization `json:"data"`
}

// AppStoreVersionsResponse represents the response from appStoreVersions endpoint
type AppStoreVersionsResponse struct {
	Data []AppStoreVersion `json:"data"`
}

// NewAppleConnectClient creates a new Apple Connect API client
func NewAppleConnectClient(issuerID, keyID, privateKeyPath, privateKey string) (*AppleConnectClient, error) {
	client := &AppleConnectClient{
		IssuerID:       issuerID,
		KeyID:          keyID,
		PrivateKeyPath: privateKeyPath,
		PrivateKey:     privateKey,
		baseURL:        "https://api.appstoreconnect.apple.com",
		httpClient:     &http.Client{Timeout: 30 * time.Second},
	}

	var privateKeyContent string
	if privateKeyPath != "" {
		content, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
		privateKeyContent = string(content)
	} else if privateKey != "" {
		privateKeyContent = privateKey
	} else {
		return nil, fmt.Errorf("either private_key_path or private_key must be provided")
	}

	parsedKey, err := parsePrivateKey(privateKeyContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	client.privateKey = parsedKey

	return client, nil
}

// parsePrivateKey parses the private key from PEM format
func parsePrivateKey(key string) (interface{}, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
	}

	// Handle both RSA and EC private keys
	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		return k, nil
	case *ecdsa.PrivateKey:
		return k, nil
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", k)
	}
}

// GenerateJWT generates a JWT token for Apple Connect API authentication
func (c *AppleConnectClient) GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": c.IssuerID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(20 * time.Minute).Unix(),
		"aud": "appstoreconnect-v1",
		"sub": c.KeyID,
	})

	token.Header["kid"] = c.KeyID

	return token.SignedString(c.privateKey)
}

// GetApps retrieves all apps from Apple App Store Connect API
func (c *AppleConnectClient) GetApps() (*AppsResponse, error) {
	jwtToken, err := c.GenerateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	req, err := http.NewRequest("GET", c.baseURL+"/v1/apps", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var appsResponse AppsResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &appsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Debug: log the first app to see what we get
	if len(appsResponse.Data) > 0 {
		fmt.Printf("[DEBUG] GetApps returned first app: id=%s, bundleId=%s, name='%s'\n",
			appsResponse.Data[0].ID,
			appsResponse.Data[0].Attributes.BundleID,
			appsResponse.Data[0].Attributes.Name)
	}

	return &appsResponse, nil
}

// GetAppLocalizations retrieves all localizations for a specific app
func (c *AppleConnectClient) GetAppLocalizations(appID string) (*AppLocalizationsResponse, string, string, error) {
	// First get the app's latest version ID, version string, and state
	versionID, versionString, versionState, err := c.getAppLatestVersionID(appID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get app version: %w", err)
	}

	jwtToken, err := c.GenerateJWT()
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/appStoreVersions/%s/appStoreVersionLocalizations", c.baseURL, versionID), nil)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var localizationsResponse AppLocalizationsResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &localizationsResponse)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Debug: log the first localization to see what we get
	if len(localizationsResponse.Data) > 0 {
		fmt.Printf("[DEBUG] Apple API returned first localization: locale=%s, name='%s', subtitle='%s', description='%s'\n",
			localizationsResponse.Data[0].Attributes.Locale,
			localizationsResponse.Data[0].Attributes.Name,
			localizationsResponse.Data[0].Attributes.Subtitle,
			localizationsResponse.Data[0].Attributes.Description)
	}

	return &localizationsResponse, versionString, versionState, nil
}

// GetAppLocalization retrieves a specific localization for an app by locale
func (c *AppleConnectClient) GetAppLocalization(appID, locale string) (*AppLocalization, error) {
	// First get the app's latest version ID (version string and state not needed here)
	versionID, _, _, err := c.getAppLatestVersionID(appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get app version: %w", err)
	}

	// Map locale to Apple's expected format
	locale = mapToAppleLocale(locale)

	jwtToken, err := c.GenerateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/appStoreVersions/%s/appStoreVersionLocalizations?filter[locale]=%s", c.baseURL, versionID, locale), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var localizationsResponse AppLocalizationsResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &localizationsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(localizationsResponse.Data) == 0 {
		return nil, fmt.Errorf("no localization found for app %s and locale %s", appID, locale)
	}

	return &localizationsResponse.Data[0], nil
}

// GetAppVersion retrieves the latest version string and state for an app
func (c *AppleConnectClient) GetAppVersion(appID string) (string, string, error) {
	_, versionString, versionState, err := c.getAppLatestVersionID(appID)
	return versionString, versionState, err
}

// GetAppVersionID retrieves the latest version ID for an app
func (c *AppleConnectClient) GetAppVersionID(appID string) (string, error) {
	versionID, _, _, err := c.getAppLatestVersionID(appID)
	return versionID, err
}

// getAppLatestVersionID retrieves the latest version ID, version string, and state for an app
func (c *AppleConnectClient) getAppLatestVersionID(appID string) (string, string, string, error) {
	jwtToken, err := c.GenerateJWT()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	url := fmt.Sprintf("%s/v1/apps/%s/appStoreVersions?limit=200", c.baseURL, appID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var versionsResponse AppStoreVersionsResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &versionsResponse); err != nil {
		return "", "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(versionsResponse.Data) == 0 {
		return "", "", "", fmt.Errorf("no app store versions found for app %s", appID)
	}

	return versionsResponse.Data[0].ID, versionsResponse.Data[0].Attributes.VersionString, versionsResponse.Data[0].Attributes.State, nil
}

// CreateAppLocalization creates a new localization for an app
// Note: Only fields supported by appStoreVersionLocalizations are included:
// - description (long description)
// - keywords
// - marketingUrl
// - promotionalText
// - supportUrl
// - whatsNew (What's new in this version)
// Other fields like name, subtitle, privacyUrl need to be updated on the app resource
func (c *AppleConnectClient) CreateAppLocalization(appID, locale, marketingURL, supportURL, longDescription, keywords, whatsNew, promotionalText string) (*AppLocalization, error) {
	jwtToken, err := c.GenerateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// First get the app's version ID to associate the localization with
	versionID, _, _, err := c.getAppLatestVersionID(appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get app version: %w", err)
	}

	// Map locale to Apple's expected format
	locale = mapToAppleLocale(locale)

	// Clean text fields to remove invalid characters
	longDescription = cleanAppleText(longDescription)
	whatsNew = cleanAppleText(whatsNew)
	promotionalText = cleanAppleText(promotionalText)

	// Debug logging
	fmt.Printf("[DEBUG] CreateAppLocalization: locale=%s, promotionalText='%s'\n", locale, promotionalText)

	// Truncate fields to Apple's limits before sending
	marketingURL = truncateString(marketingURL, 255)
	supportURL = truncateString(supportURL, 255)
	keywords = truncateString(keywords, 100)
	promotionalText = truncateString(promotionalText, 170)

	// Build attributes map with only supported fields
	attributes := map[string]interface{}{
		"locale": locale,
	}

	if longDescription != "" {
		attributes["description"] = longDescription
	}

	if keywords != "" {
		attributes["keywords"] = keywords
	}

	if marketingURL != "" {
		attributes["marketingUrl"] = marketingURL
	}

	if promotionalText != "" {
		attributes["promotionalText"] = promotionalText
	}

	if supportURL != "" {
		attributes["supportUrl"] = supportURL
	}

	if whatsNew != "" {
		attributes["whatsNew"] = whatsNew
	}

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "appStoreVersionLocalizations",
			"attributes": attributes,
			"relationships": map[string]interface{}{
				"appStoreVersion": map[string]interface{}{
					"data": map[string]interface{}{
						"type": "appStoreVersions",
						"id":   versionID,
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/v1/appStoreVersionLocalizations", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	// Create an io.Reader from payloadBytes
	req.Body = io.NopCloser(bytes.NewReader(payloadBytes))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdLocalization AppLocalizationResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &createdLocalization)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &createdLocalization.Data, nil
}

// UpdateAppLocalization updates an existing localization in App Store Connect
// Note: Only fields supported by appStoreVersionLocalizations are included:
// - description (long description)
// - keywords
// - marketingUrl
// - promotionalText
// - supportUrl
// - whatsNew (What's new in this version)
// Other fields like name, subtitle, privacyUrl need to be updated on the app resource
func (c *AppleConnectClient) UpdateAppLocalization(localizationID, marketingURL, supportURL, longDescription, keywords, whatsNew, promotionalText string) (*AppLocalization, error) {

	jwtToken, err := c.GenerateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Clean text fields to remove invalid characters
	longDescription = cleanAppleText(longDescription)
	whatsNew = cleanAppleText(whatsNew)
	promotionalText = cleanAppleText(promotionalText)

	// Debug logging
	fmt.Printf("[DEBUG] UpdateAppLocalization: localizationID=%s, promotionalText='%s'\n", localizationID, promotionalText)

	// Truncate fields to Apple's limits before sending
	keywords = truncateString(keywords, 100)
	marketingURL = truncateString(marketingURL, 255)
	supportURL = truncateString(supportURL, 255)
	promotionalText = truncateString(promotionalText, 170)

	// Build attributes map with only supported fields
	attributes := map[string]interface{}{}
	if longDescription != "" {
		attributes["description"] = longDescription
	}

	if keywords != "" {
		attributes["keywords"] = keywords
	}

	if marketingURL != "" {
		attributes["marketingUrl"] = marketingURL
	}

	if promotionalText != "" {
		attributes["promotionalText"] = promotionalText
	}

	if supportURL != "" {
		attributes["supportUrl"] = supportURL
	}

	if whatsNew != "" {
		attributes["whatsNew"] = whatsNew
	}

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type":       "appStoreVersionLocalizations",
			"id":         localizationID,
			"attributes": attributes,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v1/appStoreVersionLocalizations/%s", c.baseURL, localizationID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewReader(payloadBytes))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var updated AppLocalizationResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &updated); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &updated.Data, nil
}

// DeleteAppLocalization deletes an existing localization in App Store Connect
func (c *AppleConnectClient) DeleteAppLocalization(localizationID string) error {
	jwtToken, err := c.GenerateJWT()
	if err != nil {
		return fmt.Errorf("failed to generate JWT: %w", err)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/appStoreVersionLocalizations/%s", c.baseURL, localizationID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// truncateString truncates a string to the specified maximum length
// It tries to avoid cutting in the middle of a word when possible
// Handles UTF-8 characters correctly by using runes instead of bytes
func truncateString(s string, maxLen int) string {
	// Convert to runes to handle multi-byte characters
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	// Truncate to max rune length
	truncatedRunes := runes[:maxLen]
	truncated := string(truncatedRunes)

	// Try to find last space or comma to avoid cutting words
	lastSpace := strings.LastIndexAny(truncated, " ,")
	if lastSpace > maxLen/2 {
		truncated = truncated[:lastSpace]
	}

	return strings.TrimSpace(truncated)
}

// mapToAppleLocale converts a language code to Apple's expected locale format
func mapToAppleLocale(locale string) string {
	// Normalize locale format: replace underscore with hyphen
	normalized := strings.ReplaceAll(locale, "_", "-")
	normalized = strings.TrimSpace(normalized)

	// Apple App Store Connect locale mapping (short code to full locale)
	appleLocaleMap := map[string]string{
		"ar":    "ar-SA",
		"ca":    "ca",
		// Chinese
		"zh":         "zh-Hans",
		"zh-CN":      "zh-Hans",
		"zh-Hans":    "zh-Hans",
		"zh-Hans-CN": "zh-Hans",
		"zh-TW":      "zh-Hant",
		"zh-HK":      "zh-Hant",
		"zh-Hant":    "zh-Hant",
		"zh-Hant-TW": "zh-Hant",
		"zh-Hant-HK": "zh-Hant",
		"hr":    "hr",   // Croatian
		"cs":    "cs",
		"da":    "da",
		"nl":    "nl-NL",
		"en":    "en-US",
		"en-AU": "en-AU",
		"en-CA": "en-CA",
		"en-GB": "en-GB",
		"en-US": "en-US",
		"fi":    "fi",
		"fr":    "fr-FR",
		"fr-CA": "fr-CA",
		"de":    "de-DE",
		"el":    "el",
		"he":    "he",
		"hi":    "hi",
		"hu":    "hu",   // Hungarian
		"id":    "id",
		"it":    "it",
		"ja":    "ja",
		"ko":    "ko",
		"ms":    "ms",
		"no":    "no",
		"pl":    "pl",
		"pt":    "pt-PT",
		"pt-BR": "pt-BR",
		"pt-PT": "pt-PT",
		"ro":    "ro",   // Romanian
		"ru":    "ru",
		"sk":    "sk",
		"es":    "es-ES",
		"es-MX": "es-MX",
		"es-ES": "es-ES",
		"sv":    "sv",
		"th":    "th",
		"tr":    "tr",
		"uk":    "uk",
		"vi":    "vi",
	}

	// Try exact match first
	if mapped, ok := appleLocaleMap[normalized]; ok {
		return mapped
	}

	// Try base language code
	parts := strings.Split(normalized, "-")
	if len(parts) > 0 {
		baseCode := parts[0]
		if mapped, ok := appleLocaleMap[baseCode]; ok {
			fmt.Printf("[DEBUG] mapToAppleLocale: %s -> %s (via base code %s)\n", locale, mapped, baseCode)
			return mapped
		}
	}

	// Return normalized version if no mapping found
	fmt.Printf("[DEBUG] mapToAppleLocale: %s -> %s (no mapping found, returning normalized)\n", locale, normalized)
	return normalized
}

// cleanAppleText removes characters that are not allowed by Apple's API
// This includes characters like em dash (⸻) and other special characters
// Preserves newlines and normalizes spaces within each line
func cleanAppleText(text string) string {
	if text == "" {
		return text
	}

	// Characters that Apple doesn't allow in descriptions
	forbiddenChars := []string{
		"⸻", // em dash (U+2E3B)
		"—",  // em dash (U+2014)
		"–",  // en dash (U+2013)
		"�",  // replacement character (U+FFFD) - indicates invalid UTF-8
		"\x00", // null character
		"\x01", // start of heading
		"\x02", // start of text
		"\x03", // end of text
		"\x04", // end of transmission
		"\x05", // enquiry
		"\x06", // acknowledge
		"\x07", // bell
		"\x08", // backspace
		"\x0B", // vertical tab
		"\x0C", // form feed
		"\x0E", // shift out
		"\x0F", // shift in
		"\x10", // data link escape
		"\x11", // device control 1
		"\x12", // device control 2
		"\x13", // device control 3
		"\x14", // device control 4
		"\x15", // negative acknowledge
		"\x16", // synchronous idle
		"\x17", // end of transmission block
		"\x18", // cancel
		"\x19", // end of medium
		"\x1A", // substitute
		"\x1B", // escape
		"\x1C", // file separator
		"\x1D", // group separator
		"\x1E", // record separator
		"\x1F", // unit separator
		"\u2028", // line separator
		"\u2029", // paragraph separator
	}

	result := text
	for _, fc := range forbiddenChars {
		result = strings.ReplaceAll(result, fc, " ")
	}

	// Normalize spaces within each line while preserving line breaks
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		// Remove multiple consecutive spaces within each line
		trimmed := strings.Join(strings.Fields(line), " ")
		lines[i] = trimmed
	}
	result = strings.Join(lines, "\n")

	// Normalize consecutive newlines (collapse 3+ newlines to 2)
	result = strings.ReplaceAll(result, "\n\n\n", "\n\n")

	return result
}
