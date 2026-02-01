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
	// Apple App Store Connect locale codes
	"da": {}, "uk": {}, "ru": {}, "hr": {}, "ca": {}, "hu": {}, "hi": {},
	"id": {}, "tr": {}, "he": {}, "el": {}, "de": {}, "it": {}, "no": {}, "cs": {},
	"sk": {}, "ja": {}, "fr": {}, "fr-CA": {}, "pl": {}, "th": {}, "sv": {},
	"zh-Hans": {}, "zh-Hant": {}, "ro": {}, "fi": {}, "en-CA": {}, "en-AU": {}, "en-GB": {},
	"nl": {}, "pt-BR": {}, "pt-PT": {}, "es-MX": {}, "es-ES": {}, "vi": {},
	"ar": {}, "ko": {}, "ms": {},
}

// AppleConnectLanguage is the enum from Apple App Store Connect
type AppleConnectLanguage string

const (
	danish         AppleConnectLanguage = "da"
	ukrainian      AppleConnectLanguage = "uk"
	russian        AppleConnectLanguage = "ru"
	croatian       AppleConnectLanguage = "hr"
	catalan        AppleConnectLanguage = "ca"
	hungarian      AppleConnectLanguage = "hu"
	hindiNorth     AppleConnectLanguage = "hi"
	indonesian     AppleConnectLanguage = "id"
	turkish        AppleConnectLanguage = "tr"
	hebrew         AppleConnectLanguage = "he"
	greek          AppleConnectLanguage = "el"
	german         AppleConnectLanguage = "de"
	italian        AppleConnectLanguage = "it"
	norwegian      AppleConnectLanguage = "no"
	czech          AppleConnectLanguage = "cs"
	slovak         AppleConnectLanguage = "sk"
	japanese       AppleConnectLanguage = "ja"
	french         AppleConnectLanguage = "fr"
	frenchCanada    AppleConnectLanguage = "fr-CA"
	polish         AppleConnectLanguage = "pl"
	thai           AppleConnectLanguage = "th"
	swedish        AppleConnectLanguage = "sv"
	chineseSimplified AppleConnectLanguage = "zh-Hans"
	chineseTraditional AppleConnectLanguage = "zh-Hant"
	romanian       AppleConnectLanguage = "ro"
	finnish        AppleConnectLanguage = "fi"
	englishUS      AppleConnectLanguage = "en-US"
	englishCanada   AppleConnectLanguage = "en-CA"
	englishAustralia AppleConnectLanguage = "en-AU"
	englishUK      AppleConnectLanguage = "en-GB"
	dutch          AppleConnectLanguage = "nl"
	portugueseBrazil AppleConnectLanguage = "pt-BR"
	portuguesePortugal AppleConnectLanguage = "pt-PT"
	spanishMexico   AppleConnectLanguage = "es-MX"
	spanishSpain    AppleConnectLanguage = "es-ES"
	vietnamese     AppleConnectLanguage = "vi"
	arabic         AppleConnectLanguage = "ar"
	korean         AppleConnectLanguage = "ko"
	malay          AppleConnectLanguage = "ms"
)

// AppleConnectLanguages is the list of all Apple Store Connect supported languages
var AppleConnectLanguages = []AppleConnectLanguage{
	danish, ukrainian, russian, croatian, catalan, hungarian, hindiNorth, indonesian,
	turkish, hebrew, greek, german, italian, norwegian, czech, slovak, japanese,
	french, frenchCanada, polish, thai, swedish, chineseSimplified, chineseTraditional,
	romanian, finnish, englishUS, englishCanada, englishAustralia, englishUK, dutch,
	portugueseBrazil, portuguesePortugal, spanishMexico, spanishSpain, vietnamese,
	arabic, korean, malay,
}

// GetAppleConnectLanguages returns the list of Apple Connect languages with metadata
func GetAppleConnectLanguages() []LanguageMetadata {
	languages := make([]LanguageMetadata, 0, len(AppleConnectLanguages))
	for _, langCode := range AppleConnectLanguages {
		if meta, ok := LanguageMetadataMap[string(langCode)]; ok {
			languages = append(languages, meta)
		}
	}
	return languages
}

// GetAppleConnectLanguageMetadata returns metadata for a specific Apple Connect language
func GetAppleConnectLanguageMetadata(code string) (LanguageMetadata, bool) {
	// Map basic language codes to full Apple Connect codes
	appleConnectMap := map[string]string{
		"da": "da", "uk": "uk", "ru": "ru", "hr": "hr", "ca": "ca", "hu": "hu",
		"hi": "hi", "id": "id", "tr": "tr", "he": "he", "el": "el", "de": "de", "it": "it",
		"no": "no", "cs": "cs", "sk": "sk", "ja": "ja", "fr": "fr", "fr-CA": "fr-CA", "pl": "pl",
		"th": "th", "sv": "sv", "zh-Hans": "zh-Hans", "zh-Hant": "zh-Hant", "ro": "ro",
		"fi": "fi", "en-CA": "en-CA", "en-AU": "en-AU", "en-GB": "en-GB", "nl": "nl",
		"pt-BR": "pt-BR", "pt-PT": "pt-PT", "es-MX": "es-MX", "es-ES": "es-ES", "vi": "vi",
		"ar": "ar", "ko": "ko", "ms": "ms",
	}
	
	if mappedCode, ok := appleConnectMap[code]; ok {
		if meta, ok := LanguageMetadataMap[mappedCode]; ok {
			return meta, true
		}
	}
	return LanguageMetadata{}, false
}

// LanguageMetadataMap contains detailed information about each language
var LanguageMetadataMap = map[string]LanguageMetadata{
	// Apple App Store Connect languages (37 languages)
	"da":      {Code: "da", Name: "Danish", NativeName: "Dansk", Region: "Europe", Direction: "ltr", Emoji: "🇩🇰"},
	"uk":      {Code: "uk", Name: "Ukrainian", NativeName: "Українська", Region: "Europe", Direction: "ltr", Emoji: "🇺🇦"},
	"ru":      {Code: "ru", Name: "Russian", NativeName: "Русский", Region: "Europe", Direction: "ltr", Emoji: "🇷🇺"},
	"hr":      {Code: "hr", Name: "Croatian", NativeName: "Hrvatski", Region: "Europe", Direction: "ltr", Emoji: "🇭🇷"},
	"ca":      {Code: "ca", Name: "Catalan", NativeName: "Català", Region: "Europe", Direction: "ltr", Emoji: "🇪🇸"},
	"hu":      {Code: "hu", Name: "Hungarian", NativeName: "Magyar", Region: "Europe", Direction: "ltr", Emoji: "🇭🇺"},
	"hi":      {Code: "hi", Name: "Hindi", NativeName: "हिन्दी", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"id":      {Code: "id", Name: "Indonesian", NativeName: "Bahasa Indonesia", Region: "Asia", Direction: "ltr", Emoji: "🇮🇩"},
	"tr":      {Code: "tr", Name: "Turkish", NativeName: "Türkçe", Region: "Asia", Direction: "ltr", Emoji: "🇹🇷"},
	"he":      {Code: "he", Name: "Hebrew", NativeName: "עברית", Region: "Middle East", Direction: "rtl", Emoji: "🇮🇱"},
	"el":      {Code: "el", Name: "Greek", NativeName: "Ελληνικά", Region: "Europe", Direction: "ltr", Emoji: "🇬🇷"},
	"de":      {Code: "de", Name: "German", NativeName: "Deutsch", Region: "Europe", Direction: "ltr", Emoji: "🇩🇪"},
	"it":      {Code: "it", Name: "Italian", NativeName: "Italiano", Region: "Europe", Direction: "ltr", Emoji: "🇮🇹"},
	"no":      {Code: "no", Name: "Norwegian", NativeName: "Norsk", Region: "Europe", Direction: "ltr", Emoji: "🇳🇴"},
	"cs":      {Code: "cs", Name: "Czech", NativeName: "Čeština", Region: "Europe", Direction: "ltr", Emoji: "🇨🇿"},
	"sk":      {Code: "sk", Name: "Slovak", NativeName: "Slovenčina", Region: "Europe", Direction: "ltr", Emoji: "🇸🇰"},
	"ja":      {Code: "ja", Name: "Japanese", NativeName: "日本語", Region: "Asia", Direction: "ltr", Emoji: "🇯🇵"},
	"fr":      {Code: "fr", Name: "French", NativeName: "Français", Region: "Europe", Direction: "ltr", Emoji: "🇫🇷"},
	"fr-CA":   {Code: "fr-CA", Name: "French (Canada)", NativeName: "Français (Canada)", Region: "Americas", Direction: "ltr", Emoji: "🇨🇦"},
	"pl":      {Code: "pl", Name: "Polish", NativeName: "Polski", Region: "Europe", Direction: "ltr", Emoji: "🇵🇱"},
	"th":      {Code: "th", Name: "Thai", NativeName: "ไทย", Region: "Asia", Direction: "ltr", Emoji: "🇹🇭"},
	"sv":      {Code: "sv", Name: "Swedish", NativeName: "Svenska", Region: "Europe", Direction: "ltr", Emoji: "🇸🇪"},
	"zh-Hans": {Code: "zh-Hans", Name: "Chinese (Simplified)", NativeName: "简体中文", Region: "Asia", Direction: "ltr", Emoji: "🇨🇳"},
	"zh-Hant": {Code: "zh-Hant", Name: "Chinese (Traditional)", NativeName: "繁體中文", Region: "Asia", Direction: "ltr", Emoji: "🇹🇼"},
	"ro":      {Code: "ro", Name: "Romanian", NativeName: "Română", Region: "Europe", Direction: "ltr", Emoji: "🇷🇴"},
	"fi":      {Code: "fi", Name: "Finnish", NativeName: "Suomi", Region: "Europe", Direction: "ltr", Emoji: "🇫🇮"},
	"en-US":   {Code: "en-US", Name: "English (America)", NativeName: "English (America)", Region: "Americas", Direction: "ltr", Emoji: "🇺🇸"},
	"en-CA":   {Code: "en-CA", Name: "English (Canada)", NativeName: "English (Canada)", Region: "Americas", Direction: "ltr", Emoji: "🇨🇦"},
	"en-AU":   {Code: "en-AU", Name: "English (Australia)", NativeName: "English (Australia)", Region: "Oceania", Direction: "ltr", Emoji: "🇦🇺"},
	"en-GB":   {Code: "en-GB", Name: "English (United Kingdom)", NativeName: "English (United Kingdom)", Region: "Europe", Direction: "ltr", Emoji: "🇬🇧"},
	"nl":      {Code: "nl", Name: "Dutch", NativeName: "Nederlands", Region: "Europe", Direction: "ltr", Emoji: "🇳🇱"},
	"pt-BR":   {Code: "pt-BR", Name: "Portuguese (Brazil)", NativeName: "Português (Brasil)", Region: "Americas", Direction: "ltr", Emoji: "🇧🇷"},
	"pt-PT":   {Code: "pt-PT", Name: "Portuguese (Portugal)", NativeName: "Português (Portugal)", Region: "Europe", Direction: "ltr", Emoji: "🇵🇹"},
	"es-MX":   {Code: "es-MX", Name: "Spanish (Mexico)", NativeName: "Español (México)", Region: "Americas", Direction: "ltr", Emoji: "🇲🇽"},
	"es-ES":   {Code: "es-ES", Name: "Spanish (Spain)", NativeName: "Español (España)", Region: "Europe", Direction: "ltr", Emoji: "🇪🇸"},
	"vi":      {Code: "vi", Name: "Vietnamese", NativeName: "Tiếng Việt", Region: "Asia", Direction: "ltr", Emoji: "🇻🇳"},
	"ar":      {Code: "ar", Name: "Arabic", NativeName: "العربية", Region: "Middle East", Direction: "rtl", Emoji: "🇸🇦"},
	"ko":      {Code: "ko", Name: "Korean", NativeName: "한국어", Region: "Asia", Direction: "ltr", Emoji: "🇰🇷"},
	"ms":      {Code: "ms", Name: "Malay", NativeName: "Bahasa Melayu", Region: "Asia", Direction: "ltr", Emoji: "🇲🇾"},
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
