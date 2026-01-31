package config

// Config represents the application configuration
type Config struct {
	Global GlobalConfig `mapstructure:"global"`
	Google GoogleConfig `mapstructure:"google"`
	DeepL  DeepLConfig  `mapstructure:"deepl"`
	Baidu  BaiduConfig  `mapstructure:"baidu"`
	OpenAI OpenAIConfig `mapstructure:"openai"`
	Llama  LlamaConfig  `mapstructure:"llama"`
}

// GlobalConfig contains global configuration settings
type GlobalConfig struct {
	InputFile       string   `mapstructure:"input_file"`
	OutputFile      string   `mapstructure:"output_file"`
	SourceLanguage  string   `mapstructure:"source_language"`
	TargetLanguages []string `mapstructure:"target_languages"`
	Concurrency     int      `mapstructure:"concurrency"`
	Verbose         bool     `mapstructure:"verbose"`
}

// GoogleConfig contains Google Translate configuration
type GoogleConfig struct {
	APIKey   string `mapstructure:"api_key"`
	Model    string `mapstructure:"model"`
	Glossary string `mapstructure:"glossary"`
}

// DeepLConfig contains DeepL configuration
type DeepLConfig struct {
	APIKey    string `mapstructure:"api_key"`
	IsFree    bool   `mapstructure:"is_free"`
	Formality string `mapstructure:"formality"`
}

// BaiduConfig contains Baidu Translate configuration
type BaiduConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
}

// OpenAIConfig contains OpenAI configuration
type OpenAIConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	APIBaseURL  string  `mapstructure:"api_base_url"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

// LlamaConfig contains Llama model configuration
type LlamaConfig struct {
	LibPath         string  `mapstructure:"lib_path"`
	ModelPath       string  `mapstructure:"model_path"`
	Threads         int     `mapstructure:"threads"`
	Seed            int     `mapstructure:"seed"`
	Tokens          int     `mapstructure:"tokens"`
	TopK            int     `mapstructure:"top_k"`
	Tfs             float64 `mapstructure:"tfs"`
	TopP            float64 `mapstructure:"top_p"`
	MinP            float64 `mapstructure:"min_p"`
	TypicalP        float64 `mapstructure:"typical_p"`
	RepeatLastN     int     `mapstructure:"repeat_last_n"`
	RepeatPenalty   float64 `mapstructure:"repeat_penalty"`
	FrequencyPenalty float64 `mapstructure:"frequency_penalty"`
	PresencePenalty float64 `mapstructure:"presence_penalty"`
	Temperature     float64 `mapstructure:"temperature"`
	Verbose         bool    `mapstructure:"verbose"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Global: GlobalConfig{
			InputFile:       "Localizable.xcstrings",
			OutputFile:      "Localizable_translated.xcstrings",
			SourceLanguage:  "en",
			TargetLanguages: []string{"zh-Hans"},
			Concurrency:     5,
			Verbose:         false,
		},
		Google: GoogleConfig{
			Model: "nmt",
		},
		DeepL: DeepLConfig{
			IsFree:    false,
			Formality: "default",
		},
		Baidu: BaiduConfig{},
		OpenAI: OpenAIConfig{
			APIBaseURL:  "https://api.openai.com",
			Model:       "gpt-3.5-turbo",
			Temperature: 0.3,
			MaxTokens:   1024,
		},
	}
}
