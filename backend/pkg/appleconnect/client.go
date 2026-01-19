package appleconnect

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
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

// AppsResponse represents the response from the apps endpoint
type AppsResponse struct {
	Data []App `json:"data"`
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

	fmt.Printf("Using private key from %s", privateKeyPath)
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

	return &appsResponse, nil
}
