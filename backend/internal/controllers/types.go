package controllers

import "time"

// Job tracks long-running translation progress
type Job struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	Done      int       `json:"done"`
	Total     int       `json:"total"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TranslateRequest represents the translate request
type TranslateRequest struct {
	Provider        string         `json:"provider"`
	TargetLanguages []string       `json:"targetLanguages"`
	SourceLanguage  string         `json:"sourceLanguage"`
	Concurrency     int            `json:"concurrency"`
	TimeoutSeconds  int            `json:"timeoutSeconds"`
	Config          ProviderConfig `json:"config"`
}

// ProviderConfig is the union of provider-specific options
type ProviderConfig struct {
	APIKey      string  `json:"apiKey"`
	APIBaseURL  string  `json:"apiBaseUrl"`
	Model       string  `json:"model"`
	Glossary    string  `json:"glossary"`
	AppID       string  `json:"appId"`
	AppSecret   string  `json:"appSecret"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
	Formality   string  `json:"formality"`
	IsFree      bool    `json:"isFree"`
}