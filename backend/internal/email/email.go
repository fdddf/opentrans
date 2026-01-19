package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Email represents an email to be sent
type Email struct {
	From    string            `json:"from"`
	To      []string          `json:"to"`
	Subject string            `json:"subject"`
	HTML    string            `json:"html"`
	Text    string            `json:"text"`
	Headers map[string]string `json:"headers,omitempty"`
}

// EmailResponse represents the response from the email service
type EmailResponse struct {
	ID string `json:"id"`
}

// SendEmail sends an email using the configured email service
func SendEmail(email Email) (*EmailResponse, error) {
	// Get the email service configuration from environment variables
	serviceType := os.Getenv("EMAIL_SERVICE")

	switch serviceType {
	case "RESENDER":
		return sendWithResender(email)
	default:
		return sendWithSimpleHTTP(email)
	}
}

// sendWithResender sends an email using the Resender API
func sendWithResender(email Email) (*EmailResponse, error) {
	apiKey := os.Getenv("RESENDER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RESENDER_API_KEY environment variable is required")
	}

	url := "https://api.resend.com/emails"

	payload := map[string]interface{}{
		"from":    email.From,
		"to":      email.To,
		"subject": email.Subject,
		"html":    email.HTML,
		"text":    email.Text,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("email service returned status %d", resp.StatusCode)
	}

	var emailResp EmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&emailResp); err != nil {
		return nil, fmt.Errorf("failed to decode email response: %w", err)
	}

	return &emailResp, nil
}

// sendWithSimpleHTTP sends an email using a simple HTTP endpoint
func sendWithSimpleHTTP(email Email) (*EmailResponse, error) {
	// For now, we'll just log the email instead of sending it
	fmt.Printf("Email would be sent:\nFrom: %s\nTo: %v\nSubject: %s\nHTML: %s\n",
		email.From, email.To, email.Subject, email.HTML)

	return &EmailResponse{ID: "test-id"}, nil
}

// SendActivationEmail sends an activation email to the user
func SendActivationEmail(toEmail, username, activationCode, activationURL string) error {
	// Build the activation URL with the activation code
	fullActivationURL := fmt.Sprintf("%s/auth/activate/%s", activationURL, activationCode)

	htmlContent := fmt.Sprintf(`
		<h1>Welcome to XCStrings Translator Studio!</h1>
		<p>Hello %s,</p>
		<p>Thank you for registering. Please click the link below to activate your account:</p>
		<p><a href="%s">Activate Account</a></p>
		<p>Or copy and paste this URL into your browser: %s</p>
		<p>If you didn't register for an account, you can safely ignore this email.</p>
		<p>Best regards,<br>The XCStrings Translator Studio Team</p>
	`, username, fullActivationURL, fullActivationURL)

	textContent := fmt.Sprintf(`
		Welcome to XCStrings Translator Studio!
		
		Hello %s,
		
		Thank you for registering. Please click the link below to activate your account:
		%s
		
		If you didn't register for an account, you can safely ignore this email.
		
		Best regards,
		The XCStrings Translator Studio Team
	`, username, fullActivationURL)

	email := Email{
		From:    "noreply@xcstrings-translator.com",
		To:      []string{toEmail},
		Subject: "Activate Your XCStrings Translator Studio Account",
		HTML:    htmlContent,
		Text:    textContent,
	}

	_, err := SendEmail(email)
	return err
}
