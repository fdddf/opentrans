package services

import "fmt"

// LanguageMetadata contains information about a supported language
type LanguageMetadata struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	NativeName  string `json:"native_name"`
	Region      string `json:"region,omitempty"`
	Direction   string `json:"direction"` // "ltr" or "rtl"
}

// SupportedLanguages is a curated list of allowed language codes
var SupportedLanguages = map[string]struct{}{
	"en": {}, "en-US": {}, "en-GB": {}, "zh-Hans": {}, "zh-Hant": {}, "ja": {},
	"fr": {}, "de": {}, "es": {}, "ko": {}, "ru": {}, "pt": {}, "pt-BR": {}, "it": {},
	"ar": {}, "pt-PT": {}, "nl": {}, "sv": {}, "da": {}, "fi": {}, "no": {}, "pl": {},
	"tr": {}, "th": {}, "vi": {}, "hi": {}, "id": {}, "cs": {}, "el": {}, "he": {},
	"hu": {}, "ro": {}, "sk": {}, "uk": {},
}

// LanguageMetadataMap contains detailed information about each language
var LanguageMetadataMap = map[string]LanguageMetadata{
	"en":    {Code: "en", Name: "English", NativeName: "English", Region: "World", Direction: "ltr"},
	"en-US": {Code: "en-US", Name: "English (US)", NativeName: "English (US)", Region: "Americas", Direction: "ltr"},
	"en-GB": {Code: "en-GB", Name: "English (UK)", NativeName: "English (UK)", Region: "Europe", Direction: "ltr"},
	"zh-Hans": {Code: "zh-Hans", Name: "Chinese (Simplified)", NativeName: "简体中文", Region: "Asia", Direction: "ltr"},
	"zh-Hant": {Code: "zh-Hant", Name: "Chinese (Traditional)", NativeName: "繁體中文", Region: "Asia", Direction: "ltr"},
	"ja":    {Code: "ja", Name: "Japanese", NativeName: "日本語", Region: "Asia", Direction: "ltr"},
	"ko":    {Code: "ko", Name: "Korean", NativeName: "한국어", Region: "Asia", Direction: "ltr"},
	"fr":    {Code: "fr", Name: "French", NativeName: "Français", Region: "Europe", Direction: "ltr"},
	"de":    {Code: "de", Name: "German", NativeName: "Deutsch", Region: "Europe", Direction: "ltr"},
	"es":    {Code: "es", Name: "Spanish", NativeName: "Español", Region: "Europe", Direction: "ltr"},
	"ru":    {Code: "ru", Name: "Russian", NativeName: "Русский", Region: "Europe", Direction: "ltr"},
	"pt":    {Code: "pt", Name: "Portuguese", NativeName: "Português", Region: "Europe", Direction: "ltr"},
	"pt-BR": {Code: "pt-BR", Name: "Portuguese (Brazil)", NativeName: "Português (Brasil)", Region: "Americas", Direction: "ltr"},
	"pt-PT": {Code: "pt-PT", Name: "Portuguese (Portugal)", NativeName: "Português (Portugal)", Region: "Europe", Direction: "ltr"},
	"it":    {Code: "it", Name: "Italian", NativeName: "Italiano", Region: "Europe", Direction: "ltr"},
	"ar":    {Code: "ar", Name: "Arabic", NativeName: "العربية", Region: "Middle East", Direction: "rtl"},
	"nl":    {Code: "nl", Name: "Dutch", NativeName: "Nederlands", Region: "Europe", Direction: "ltr"},
	"sv":    {Code: "sv", Name: "Swedish", NativeName: "Svenska", Region: "Europe", Direction: "ltr"},
	"da":    {Code: "da", Name: "Danish", NativeName: "Dansk", Region: "Europe", Direction: "ltr"},
	"fi":    {Code: "fi", Name: "Finnish", NativeName: "Suomi", Region: "Europe", Direction: "ltr"},
	"no":    {Code: "no", Name: "Norwegian", NativeName: "Norsk", Region: "Europe", Direction: "ltr"},
	"pl":    {Code: "pl", Name: "Polish", NativeName: "Polski", Region: "Europe", Direction: "ltr"},
	"tr":    {Code: "tr", Name: "Turkish", NativeName: "Türkçe", Region: "Asia", Direction: "ltr"},
	"th":    {Code: "th", Name: "Thai", NativeName: "ไทย", Region: "Asia", Direction: "ltr"},
	"vi":    {Code: "vi", Name: "Vietnamese", NativeName: "Tiếng Việt", Region: "Asia", Direction: "ltr"},
	"hi":    {Code: "hi", Name: "Hindi", NativeName: "हिन्दी", Region: "Asia", Direction: "ltr"},
	"id":    {Code: "id", Name: "Indonesian", NativeName: "Bahasa Indonesia", Region: "Asia", Direction: "ltr"},
	"cs":    {Code: "cs", Name: "Czech", NativeName: "Čeština", Region: "Europe", Direction: "ltr"},
	"el":    {Code: "el", Name: "Greek", NativeName: "Ελληνικά", Region: "Europe", Direction: "ltr"},
	"he":    {Code: "he", Name: "Hebrew", NativeName: "עברית", Region: "Middle East", Direction: "rtl"},
	"hu":    {Code: "hu", Name: "Hungarian", NativeName: "Magyar", Region: "Europe", Direction: "ltr"},
	"ro":    {Code: "ro", Name: "Romanian", NativeName: "Română", Region: "Europe", Direction: "ltr"},
	"sk":    {Code: "sk", Name: "Slovak", NativeName: "Slovenčina", Region: "Europe", Direction: "ltr"},
	"uk":    {Code: "uk", Name: "Ukrainian", NativeName: "Українська", Region: "Europe", Direction: "ltr"},
}

// IsSupportedLanguage checks if a language code is allowed
func IsSupportedLanguage(code string) bool {
	_, ok := SupportedLanguages[code]
	return ok
}

// ValidateLanguages validates a list of language codes
func ValidateLanguages(codes []string) error {
	for _, code := range codes {
		if !IsSupportedLanguage(code) {
			return fmt.Errorf("unsupported language: %s", code)
		}
	}
	return nil
}

// GetSupportedLanguagesList returns a list of all supported languages with metadata
func GetSupportedLanguagesList() []LanguageMetadata {
	languages := make([]LanguageMetadata, 0, len(LanguageMetadataMap))
	for _, lang := range LanguageMetadataMap {
		languages = append(languages, lang)
	}
	return languages
}

// GetLanguageMetadata returns metadata for a specific language code
func GetLanguageMetadata(code string) (LanguageMetadata, bool) {
	lang, ok := LanguageMetadataMap[code]
	return lang, ok
}
