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
	// Apple App Store Connect locale codes (41 languages)
	"ar-SA": {}, "ca": {}, "zh-Hans": {}, "zh-Hant": {}, "hr": {}, "cs": {}, "da": {},
	"nl-NL": {}, "en-AU": {}, "en-CA": {}, "en-GB": {}, "en-US": {}, "fi": {},
	"fr-CA": {}, "fr-FR": {}, "de-DE": {}, "el": {}, "he": {}, "hi": {}, "hu": {},
	"id": {}, "it": {}, "ja": {}, "ko": {}, "ms": {}, "no": {}, "pl": {},
	"pt-BR": {}, "pt-PT": {}, "ro": {}, "ru": {}, "sk": {}, "es-MX": {}, "es-ES": {},
	"sv": {}, "th": {}, "tr": {}, "uk": {}, "vi": {},
}

// AppleConnectLanguage is the enum from Apple App Store Connect
type AppleConnectLanguage string

const (
	arabic            AppleConnectLanguage = "ar-SA"
	catalan           AppleConnectLanguage = "ca"
	chineseSimplified AppleConnectLanguage = "zh-Hans"
	chineseTraditional AppleConnectLanguage = "zh-Hant"
	croatian          AppleConnectLanguage = "hr"
	czech             AppleConnectLanguage = "cs"
	danish            AppleConnectLanguage = "da"
	dutch             AppleConnectLanguage = "nl-NL"
	englishAustralia  AppleConnectLanguage = "en-AU"
	englishCanada     AppleConnectLanguage = "en-CA"
	englishUK         AppleConnectLanguage = "en-GB"
	englishUS         AppleConnectLanguage = "en-US"
	finnish           AppleConnectLanguage = "fi"
	frenchCanada      AppleConnectLanguage = "fr-CA"
	frenchFrance      AppleConnectLanguage = "fr-FR"
	german            AppleConnectLanguage = "de-DE"
	greek             AppleConnectLanguage = "el"
	hebrew            AppleConnectLanguage = "he"
	hindi             AppleConnectLanguage = "hi"
	hungarian         AppleConnectLanguage = "hu"
	indonesian        AppleConnectLanguage = "id"
	italian           AppleConnectLanguage = "it"
	japanese          AppleConnectLanguage = "ja"
	korean            AppleConnectLanguage = "ko"
	malay             AppleConnectLanguage = "ms"
	norwegian         AppleConnectLanguage = "no"
	polish            AppleConnectLanguage = "pl"
	portugueseBrazil  AppleConnectLanguage = "pt-BR"
	portuguesePortugal AppleConnectLanguage = "pt-PT"
	romanian          AppleConnectLanguage = "ro"
	russian           AppleConnectLanguage = "ru"
	slovak            AppleConnectLanguage = "sk"
	spanishMexico     AppleConnectLanguage = "es-MX"
	spanishSpain      AppleConnectLanguage = "es-ES"
	swedish           AppleConnectLanguage = "sv"
	thai              AppleConnectLanguage = "th"
	turkish           AppleConnectLanguage = "tr"
	ukrainian         AppleConnectLanguage = "uk"
	vietnamese        AppleConnectLanguage = "vi"
)

// AppleConnectLanguages is the list of all Apple Store Connect supported languages (41 languages)
var AppleConnectLanguages = []AppleConnectLanguage{
	arabic, catalan, chineseSimplified, chineseTraditional, croatian, czech, danish,
	dutch, englishAustralia, englishCanada, englishUK, englishUS, finnish,
	frenchCanada, frenchFrance, german, greek, hebrew, hindi, hungarian,
	indonesian, italian, japanese, korean, malay, norwegian, polish,
	portugueseBrazil, portuguesePortugal, romanian, russian, slovak, spanishMexico, spanishSpain,
	swedish, thai, turkish, ukrainian, vietnamese,
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
	// Map alternative language codes to Apple Connect codes
	appleConnectMap := map[string]string{
		"ar": "ar-SA", "ar-SA": "ar-SA",
		"ca": "ca",
		"zh-CN": "zh-Hans", "zh-Hans": "zh-Hans", "zh-Hans-CN": "zh-Hans",
		"zh-TW": "zh-Hant", "zh-Hant": "zh-Hant", "zh-Hant-TW": "zh-Hant",
		"hr": "hr", "cs": "cs", "da": "da",
		"nl": "nl-NL", "nl-NL": "nl-NL",
		"en-AU": "en-AU", "en-CA": "en-CA", "en-GB": "en-GB", "en-US": "en-US",
		"fi": "fi", "fr-CA": "fr-CA", "fr-FR": "fr-FR", "fr": "fr-FR",
		"de-DE": "de-DE", "de": "de-DE",
		"el": "el", "he": "he", "hi": "hi", "hu": "hu", "id": "id",
		"it": "it", "ja": "ja", "ko": "ko", "ms": "ms", "no": "no",
		"pl": "pl", "pt-BR": "pt-BR", "pt-PT": "pt-PT", "ro": "ro",
		"ru": "ru", "sk": "sk", "es-MX": "es-MX", "es-ES": "es-ES", "es": "es-ES",
		"sv": "sv", "th": "th", "tr": "tr", "uk": "uk", "vi": "vi",
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
	// Apple App Store Connect languages (41 languages)
	"ar-SA":   {Code: "ar-SA", Name: "Arabic", NativeName: "العربية", Region: "Middle East", Direction: "rtl", Emoji: "🇸🇦"},
	"ca":      {Code: "ca", Name: "Catalan", NativeName: "Català", Region: "Europe", Direction: "ltr", Emoji: "🇪🇸"},
	"zh-Hans": {Code: "zh-Hans", Name: "Chinese (Simplified)", NativeName: "简体中文", Region: "Asia", Direction: "ltr", Emoji: "🇨🇳"},
	"zh-Hant": {Code: "zh-Hant", Name: "Chinese (Traditional)", NativeName: "繁體中文", Region: "Asia", Direction: "ltr", Emoji: "🇹🇼"},
	"hr":      {Code: "hr", Name: "Croatian", NativeName: "Hrvatski", Region: "Europe", Direction: "ltr", Emoji: "🇭🇷"},
	"cs":      {Code: "cs", Name: "Czech", NativeName: "Čeština", Region: "Europe", Direction: "ltr", Emoji: "🇨🇿"},
	"da":      {Code: "da", Name: "Danish", NativeName: "Dansk", Region: "Europe", Direction: "ltr", Emoji: "🇩🇰"},
	"nl-NL":   {Code: "nl-NL", Name: "Dutch", NativeName: "Nederlands", Region: "Europe", Direction: "ltr", Emoji: "🇳🇱"},
	"en-AU":   {Code: "en-AU", Name: "English (Australia)", NativeName: "English (Australia)", Region: "Oceania", Direction: "ltr", Emoji: "🇦🇺"},
	"en-CA":   {Code: "en-CA", Name: "English (Canada)", NativeName: "English (Canada)", Region: "Americas", Direction: "ltr", Emoji: "🇨🇦"},
	"en-GB":   {Code: "en-GB", Name: "English (U.K.)", NativeName: "English (U.K.)", Region: "Europe", Direction: "ltr", Emoji: "🇬🇧"},
	"en-US":   {Code: "en-US", Name: "English (U.S.)", NativeName: "English (U.S.)", Region: "Americas", Direction: "ltr", Emoji: "🇺🇸"},
	"fi":      {Code: "fi", Name: "Finnish", NativeName: "Suomi", Region: "Europe", Direction: "ltr", Emoji: "🇫🇮"},
	"fr-CA":   {Code: "fr-CA", Name: "French (Canada)", NativeName: "Français (Canada)", Region: "Americas", Direction: "ltr", Emoji: "🇨🇦"},
	"fr-FR":   {Code: "fr-FR", Name: "French (France)", NativeName: "Français (France)", Region: "Europe", Direction: "ltr", Emoji: "🇫🇷"},
	"de-DE":   {Code: "de-DE", Name: "German", NativeName: "Deutsch", Region: "Europe", Direction: "ltr", Emoji: "🇩🇪"},
	"el":      {Code: "el", Name: "Greek", NativeName: "Ελληνικά", Region: "Europe", Direction: "ltr", Emoji: "🇬🇷"},
	"he":      {Code: "he", Name: "Hebrew", NativeName: "עברית", Region: "Middle East", Direction: "rtl", Emoji: "🇮🇱"},
	"hi":      {Code: "hi", Name: "Hindi", NativeName: "हिन्दी", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"hu":      {Code: "hu", Name: "Hungarian", NativeName: "Magyar", Region: "Europe", Direction: "ltr", Emoji: "🇭🇺"},
	"id":      {Code: "id", Name: "Indonesian", NativeName: "Bahasa Indonesia", Region: "Asia", Direction: "ltr", Emoji: "🇮🇩"},
	"it":      {Code: "it", Name: "Italian", NativeName: "Italiano", Region: "Europe", Direction: "ltr", Emoji: "🇮🇹"},
	"ja":      {Code: "ja", Name: "Japanese", NativeName: "日本語", Region: "Asia", Direction: "ltr", Emoji: "🇯🇵"},
	"ko":      {Code: "ko", Name: "Korean", NativeName: "한국어", Region: "Asia", Direction: "ltr", Emoji: "🇰🇷"},
	"ms":      {Code: "ms", Name: "Malay", NativeName: "Bahasa Melayu", Region: "Asia", Direction: "ltr", Emoji: "🇲🇾"},
	"no":      {Code: "no", Name: "Norwegian", NativeName: "Norsk", Region: "Europe", Direction: "ltr", Emoji: "🇳🇴"},
	"pl":      {Code: "pl", Name: "Polish", NativeName: "Polski", Region: "Europe", Direction: "ltr", Emoji: "🇵🇱"},
	"pt-BR":   {Code: "pt-BR", Name: "Portuguese (Brazil)", NativeName: "Português (Brasil)", Region: "Americas", Direction: "ltr", Emoji: "🇧🇷"},
	"pt-PT":   {Code: "pt-PT", Name: "Portuguese (Portugal)", NativeName: "Português (Portugal)", Region: "Europe", Direction: "ltr", Emoji: "🇵🇹"},
	"ro":      {Code: "ro", Name: "Romanian", NativeName: "Română", Region: "Europe", Direction: "ltr", Emoji: "🇷🇴"},
	"ru":      {Code: "ru", Name: "Russian", NativeName: "Русский", Region: "Europe", Direction: "ltr", Emoji: "🇷🇺"},
	"sk":      {Code: "sk", Name: "Slovak", NativeName: "Slovenčina", Region: "Europe", Direction: "ltr", Emoji: "🇸🇰"},
	"es-MX":   {Code: "es-MX", Name: "Spanish (Mexico)", NativeName: "Español (México)", Region: "Americas", Direction: "ltr", Emoji: "🇲🇽"},
	"es-ES":   {Code: "es-ES", Name: "Spanish (Spain)", NativeName: "Español (España)", Region: "Europe", Direction: "ltr", Emoji: "🇪🇸"},
	"sv":      {Code: "sv", Name: "Swedish", NativeName: "Svenska", Region: "Europe", Direction: "ltr", Emoji: "🇸🇪"},
	"th":      {Code: "th", Name: "Thai", NativeName: "ไทย", Region: "Asia", Direction: "ltr", Emoji: "🇹🇭"},
	"tr":      {Code: "tr", Name: "Turkish", NativeName: "Türkçe", Region: "Asia", Direction: "ltr", Emoji: "🇹🇷"},
	"uk":      {Code: "uk", Name: "Ukrainian", NativeName: "Українська", Region: "Europe", Direction: "ltr", Emoji: "🇺🇦"},
	"vi":      {Code: "vi", Name: "Vietnamese", NativeName: "Tiếng Việt", Region: "Asia", Direction: "ltr", Emoji: "🇻🇳"},
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