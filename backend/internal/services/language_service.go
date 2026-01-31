package services

import "fmt"

// LanguageMetadata contains information about a supported language
type LanguageMetadata struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	NativeName  string `json:"native_name"`
	Region      string `json:"region,omitempty"`
	Direction   string `json:"direction"` // "ltr" or "rtl"
	Emoji       string `json:"emoji"`
}

// SupportedLanguages is a curated list of allowed language codes
var SupportedLanguages = map[string]struct{}{
	"zh": {}, "en": {}, "fr": {}, "pt": {}, "es": {}, "ja": {}, "tr": {}, "ru": {},
	"ar": {}, "ko": {}, "th": {}, "it": {}, "de": {}, "vi": {}, "ms": {}, "id": {},
	"tl": {}, "hi": {}, "zh-Hant": {}, "pl": {}, "cs": {}, "nl": {}, "km": {}, "my": {},
	"fa": {}, "gu": {}, "ur": {}, "te": {}, "mr": {}, "he": {}, "bn": {}, "ta": {},
	"uk": {}, "bo": {}, "kk": {}, "mn": {}, "ug": {}, "yue": {},
}

// LanguageMetadataMap contains detailed information about each language
var LanguageMetadataMap = map[string]LanguageMetadata{
	"zh":      {Code: "zh", Name: "Chinese", NativeName: "简体中文", Region: "Asia", Direction: "ltr", Emoji: "🇨🇳"},
	"en":      {Code: "en", Name: "English", NativeName: "English", Region: "World", Direction: "ltr", Emoji: "🇺🇸"},
	"fr":      {Code: "fr", Name: "French", NativeName: "Français", Region: "Europe", Direction: "ltr", Emoji: "🇫🇷"},
	"pt":      {Code: "pt", Name: "Portuguese", NativeName: "Português", Region: "Europe", Direction: "ltr", Emoji: "🇵🇹"},
	"es":      {Code: "es", Name: "Spanish", NativeName: "Español", Region: "Europe", Direction: "ltr", Emoji: "🇪🇸"},
	"ja":      {Code: "ja", Name: "Japanese", NativeName: "日本語", Region: "Asia", Direction: "ltr", Emoji: "🇯🇵"},
	"tr":      {Code: "tr", Name: "Turkish", NativeName: "Türkçe", Region: "Asia", Direction: "ltr", Emoji: "🇹🇷"},
	"ru":      {Code: "ru", Name: "Russian", NativeName: "Русский", Region: "Europe", Direction: "ltr", Emoji: "🇷🇺"},
	"ar":      {Code: "ar", Name: "Arabic", NativeName: "العربية", Region: "Middle East", Direction: "rtl", Emoji: "🇸🇦"},
	"ko":      {Code: "ko", Name: "Korean", NativeName: "한국어", Region: "Asia", Direction: "ltr", Emoji: "🇰🇷"},
	"th":      {Code: "th", Name: "Thai", NativeName: "ไทย", Region: "Asia", Direction: "ltr", Emoji: "🇹🇭"},
	"it":      {Code: "it", Name: "Italian", NativeName: "Italiano", Region: "Europe", Direction: "ltr", Emoji: "🇮🇹"},
	"de":      {Code: "de", Name: "German", NativeName: "Deutsch", Region: "Europe", Direction: "ltr", Emoji: "🇩🇪"},
	"vi":      {Code: "vi", Name: "Vietnamese", NativeName: "Tiếng Việt", Region: "Asia", Direction: "ltr", Emoji: "🇻🇳"},
	"ms":      {Code: "ms", Name: "Malay", NativeName: "Bahasa Melayu", Region: "Asia", Direction: "ltr", Emoji: "🇲🇾"},
	"id":      {Code: "id", Name: "Indonesian", NativeName: "Bahasa Indonesia", Region: "Asia", Direction: "ltr", Emoji: "🇮🇩"},
	"tl":      {Code: "tl", Name: "Filipino", NativeName: "Filipino", Region: "Asia", Direction: "ltr", Emoji: "🇵🇭"},
	"hi":      {Code: "hi", Name: "Hindi", NativeName: "हिन्दी", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"zh-Hant": {Code: "zh-Hant", Name: "Traditional Chinese", NativeName: "繁體中文", Region: "Asia", Direction: "ltr", Emoji: "🌐"},
	"pl":      {Code: "pl", Name: "Polish", NativeName: "Polski", Region: "Europe", Direction: "ltr", Emoji: "🇵🇱"},
	"cs":      {Code: "cs", Name: "Czech", NativeName: "Čeština", Region: "Europe", Direction: "ltr", Emoji: "🇨🇿"},
	"nl":      {Code: "nl", Name: "Dutch", NativeName: "Nederlands", Region: "Europe", Direction: "ltr", Emoji: "🇳🇱"},
	"km":      {Code: "km", Name: "Khmer", NativeName: "ខ្មែរ", Region: "Asia", Direction: "ltr", Emoji: "🇰🇭"},
	"my":      {Code: "my", Name: "Burmese", NativeName: "မြန်မာ", Region: "Asia", Direction: "ltr", Emoji: "🇲🇲"},
	"fa":      {Code: "fa", Name: "Persian", NativeName: "فارسی", Region: "Middle East", Direction: "rtl", Emoji: "🇮🇷"},
	"gu":      {Code: "gu", Name: "Gujarati", NativeName: "ગુજરાતી", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"ur":      {Code: "ur", Name: "Urdu", NativeName: "اردو", Region: "Asia", Direction: "rtl", Emoji: "🇵🇰"},
	"te":      {Code: "te", Name: "Telugu", NativeName: "తెలుగు", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"mr":      {Code: "mr", Name: "Marathi", NativeName: "मराठी", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"he":      {Code: "he", Name: "Hebrew", NativeName: "עברית", Region: "Middle East", Direction: "rtl", Emoji: "🇮🇱"},
	"bn":      {Code: "bn", Name: "Bengali", NativeName: "বাংলা", Region: "Asia", Direction: "ltr", Emoji: "🇧🇩"},
	"ta":      {Code: "ta", Name: "Tamil", NativeName: "தமிழ்", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"uk":      {Code: "uk", Name: "Ukrainian", NativeName: "Українська", Region: "Europe", Direction: "ltr", Emoji: "🇺🇦"},
	"bo":      {Code: "bo", Name: "Tibetan", NativeName: "བོད་སྐད", Region: "Asia", Direction: "ltr", Emoji: "🇨🇳"},
	"kk":      {Code: "kk", Name: "Kazakh", NativeName: "Қазақша", Region: "Asia", Direction: "ltr", Emoji: "🇰🇿"},
	"mn":      {Code: "mn", Name: "Mongolian", NativeName: "Монгол", Region: "Asia", Direction: "ltr", Emoji: "🇲🇳"},
	"ug":      {Code: "ug", Name: "Uyghur", NativeName: "ئۇيغۇرچە", Region: "Asia", Direction: "rtl", Emoji: "🇨🇳"},
	"yue":     {Code: "yue", Name: "Cantonese", NativeName: "粵語", Region: "Asia", Direction: "ltr", Emoji: "🇭🇰"},
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

// GetSupportedLanguageCodes returns all supported language codes
func GetSupportedLanguageCodes() []string {
	codes := make([]string, 0, len(LanguageMetadataMap))
	for code := range LanguageMetadataMap {
		codes = append(codes, code)
	}
	return codes
}

// GetLanguageEmoji returns the emoji for a specific language code
func GetLanguageEmoji(code string) string {
	if lang, ok := LanguageMetadataMap[code]; ok {
		return lang.Emoji
	}
	return "🌐"
}
