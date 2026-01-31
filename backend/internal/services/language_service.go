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
	"en": {}, "fr": {}, "pt": {}, "es": {}, "ja": {}, "tr": {}, "ru": {},
	"ar": {}, "ko": {}, "th": {}, "it": {}, "de": {}, "vi": {}, "ms": {}, "id": {},
	"tl": {}, "hi": {}, "pl": {}, "cs": {}, "nl": {}, "km": {}, "my": {},
	"fa": {}, "gu": {}, "ur": {}, "te": {}, "mr": {}, "he": {}, "bn": {}, "ta": {},
	"uk": {}, "bo": {}, "kk": {}, "mn": {}, "ug": {}, "yue": {},
	// Apple App Store Connect locale codes (includes all basic language variants)
	"en-US": {}, "en-GB": {}, "en-AU": {}, "en-CA": {}, "en-IN": {},
	"zh-Hans": {}, "zh-Hant": {}, "zh-CN": {}, "zh-TW": {}, "zh-HK": {},
	"pt-BR": {}, "pt-PT": {},
	"es-ES": {}, "es-MX": {}, "es-AR": {}, "es-CL": {}, "es-CO": {},
	"fr-FR": {}, "fr-CA": {},
	"de-DE": {}, "de-AT": {},
	"it-IT": {},
	"ja-JP": {},
	"ko-KR": {},
	"ru-RU": {},
	"ar-SA": {},
	"th-TH": {},
	"vi-VN": {},
	"id-ID": {},
	"ms-MY": {},
	"hi-IN": {},
	"pl-PL": {},
	"cs-CZ": {},
	"nl-NL": {},
	"tr-TR": {},
	"uk-UA": {},
	"he-IL": {},
	"fa-IR": {},
	"bn-IN": {},
	"ta-IN": {},
	"te-IN": {},
	"mr-IN": {},
	"gu-IN": {},
	"ur-PK": {},
	"kk-KZ": {},
	"mn-MN": {},
	"bo-CN": {},
	"ug-CN": {},
	"yue-HK": {},
	"km-KH": {},
	"my-MM": {},
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
	// Apple App Store Connect locale codes
	"en-US":   {Code: "en-US", Name: "English (US)", NativeName: "English (United States)", Region: "Americas", Direction: "ltr", Emoji: "🇺🇸"},
	"en-GB":   {Code: "en-GB", Name: "English (UK)", NativeName: "English (United Kingdom)", Region: "Europe", Direction: "ltr", Emoji: "🇬🇧"},
	"en-AU":   {Code: "en-AU", Name: "English (Australia)", NativeName: "English (Australia)", Region: "Oceania", Direction: "ltr", Emoji: "🇦🇺"},
	"en-CA":   {Code: "en-CA", Name: "English (Canada)", NativeName: "English (Canada)", Region: "Americas", Direction: "ltr", Emoji: "🇨🇦"},
	"en-IN":   {Code: "en-IN", Name: "English (India)", NativeName: "English (India)", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"zh-Hans": {Code: "zh-Hans", Name: "Chinese (Simplified)", NativeName: "简体中文", Region: "Asia", Direction: "ltr", Emoji: "🇨🇳"},
	"zh-CN":   {Code: "zh-CN", Name: "Chinese (Simplified)", NativeName: "简体中文 (中国)", Region: "Asia", Direction: "ltr", Emoji: "🇨🇳"},
	"zh-TW":   {Code: "zh-TW", Name: "Chinese (Taiwan)", NativeName: "繁體中文 (台灣)", Region: "Asia", Direction: "ltr", Emoji: "🇹🇼"},
	"zh-HK":   {Code: "zh-HK", Name: "Chinese (Hong Kong)", NativeName: "繁體中文 (香港)", Region: "Asia", Direction: "ltr", Emoji: "🇭🇰"},
	"pt-BR":   {Code: "pt-BR", Name: "Portuguese (Brazil)", NativeName: "Português (Brasil)", Region: "Americas", Direction: "ltr", Emoji: "🇧🇷"},
	"pt-PT":   {Code: "pt-PT", Name: "Portuguese (Portugal)", NativeName: "Português (Portugal)", Region: "Europe", Direction: "ltr", Emoji: "🇵🇹"},
	"es-ES":   {Code: "es-ES", Name: "Spanish (Spain)", NativeName: "Español (España)", Region: "Europe", Direction: "ltr", Emoji: "🇪🇸"},
	"es-MX":   {Code: "es-MX", Name: "Spanish (Mexico)", NativeName: "Español (México)", Region: "Americas", Direction: "ltr", Emoji: "🇲🇽"},
	"es-AR":   {Code: "es-AR", Name: "Spanish (Argentina)", NativeName: "Español (Argentina)", Region: "Americas", Direction: "ltr", Emoji: "🇦🇷"},
	"es-CL":   {Code: "es-CL", Name: "Spanish (Chile)", NativeName: "Español (Chile)", Region: "Americas", Direction: "ltr", Emoji: "🇨🇱"},
	"es-CO":   {Code: "es-CO", Name: "Spanish (Colombia)", NativeName: "Español (Colombia)", Region: "Americas", Direction: "ltr", Emoji: "🇨🇴"},
	"fr-FR":   {Code: "fr-FR", Name: "French (France)", NativeName: "Français (France)", Region: "Europe", Direction: "ltr", Emoji: "🇫🇷"},
	"fr-CA":   {Code: "fr-CA", Name: "French (Canada)", NativeName: "Français (Canada)", Region: "Americas", Direction: "ltr", Emoji: "🇨🇦"},
	"de-DE":   {Code: "de-DE", Name: "German (Germany)", NativeName: "Deutsch (Deutschland)", Region: "Europe", Direction: "ltr", Emoji: "🇩🇪"},
	"de-AT":   {Code: "de-AT", Name: "German (Austria)", NativeName: "Deutsch (Österreich)", Region: "Europe", Direction: "ltr", Emoji: "🇦🇹"},
	"it-IT":   {Code: "it-IT", Name: "Italian (Italy)", NativeName: "Italiano (Italia)", Region: "Europe", Direction: "ltr", Emoji: "🇮🇹"},
	"ja-JP":   {Code: "ja-JP", Name: "Japanese (Japan)", NativeName: "日本語 (日本)", Region: "Asia", Direction: "ltr", Emoji: "🇯🇵"},
	"ko-KR":   {Code: "ko-KR", Name: "Korean (South Korea)", NativeName: "한국어 (대한민국)", Region: "Asia", Direction: "ltr", Emoji: "🇰🇷"},
	"ru-RU":   {Code: "ru-RU", Name: "Russian (Russia)", NativeName: "Русский (Россия)", Region: "Europe", Direction: "ltr", Emoji: "🇷🇺"},
	"ar-SA":   {Code: "ar-SA", Name: "Arabic (Saudi Arabia)", NativeName: "العربية (السعودية)", Region: "Middle East", Direction: "rtl", Emoji: "🇸🇦"},
	"th-TH":   {Code: "th-TH", Name: "Thai (Thailand)", NativeName: "ไทย (ไทย)", Region: "Asia", Direction: "ltr", Emoji: "🇹🇭"},
	"vi-VN":   {Code: "vi-VN", Name: "Vietnamese (Vietnam)", NativeName: "Tiếng Việt (Việt Nam)", Region: "Asia", Direction: "ltr", Emoji: "🇻🇳"},
	"id-ID":   {Code: "id-ID", Name: "Indonesian (Indonesia)", NativeName: "Bahasa Indonesia (Indonesia)", Region: "Asia", Direction: "ltr", Emoji: "🇮🇩"},
	"ms-MY":   {Code: "ms-MY", Name: "Malay (Malaysia)", NativeName: "Bahasa Melayu (Malaysia)", Region: "Asia", Direction: "ltr", Emoji: "🇲🇾"},
	"hi-IN":   {Code: "hi-IN", Name: "Hindi (India)", NativeName: "हिन्दी (भारत)", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"pl-PL":   {Code: "pl-PL", Name: "Polish (Poland)", NativeName: "Polski (Polska)", Region: "Europe", Direction: "ltr", Emoji: "🇵🇱"},
	"cs-CZ":   {Code: "cs-CZ", Name: "Czech (Czech Republic)", NativeName: "Čeština (Česko)", Region: "Europe", Direction: "ltr", Emoji: "🇨🇿"},
	"nl-NL":   {Code: "nl-NL", Name: "Dutch (Netherlands)", NativeName: "Nederlands (Nederland)", Region: "Europe", Direction: "ltr", Emoji: "🇳🇱"},
	"tr-TR":   {Code: "tr-TR", Name: "Turkish (Turkey)", NativeName: "Türkçe (Türkiye)", Region: "Asia", Direction: "ltr", Emoji: "🇹🇷"},
	"uk-UA":   {Code: "uk-UA", Name: "Ukrainian (Ukraine)", NativeName: "Українська (Україна)", Region: "Europe", Direction: "ltr", Emoji: "🇺🇦"},
	"he-IL":   {Code: "he-IL", Name: "Hebrew (Israel)", NativeName: "עברית (ישראל)", Region: "Middle East", Direction: "rtl", Emoji: "🇮🇱"},
	"fa-IR":   {Code: "fa-IR", Name: "Persian (Iran)", NativeName: "فارسی (ایران)", Region: "Middle East", Direction: "rtl", Emoji: "🇮🇷"},
	"bn-IN":   {Code: "bn-IN", Name: "Bengali (India)", NativeName: "বাংলা (ভারত)", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"ta-IN":   {Code: "ta-IN", Name: "Tamil (India)", NativeName: "தமிழ் (இந்தியா)", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"te-IN":   {Code: "te-IN", Name: "Telugu (India)", NativeName: "తెలుగు (భారతదేశం)", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"mr-IN":   {Code: "mr-IN", Name: "Marathi (India)", NativeName: "मराठी (भारत)", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"gu-IN":   {Code: "gu-IN", Name: "Gujarati (India)", NativeName: "ગુજરાતી (ભારત)", Region: "Asia", Direction: "ltr", Emoji: "🇮🇳"},
	"ur-PK":   {Code: "ur-PK", Name: "Urdu (Pakistan)", NativeName: "اردو (پاکستان)", Region: "Asia", Direction: "rtl", Emoji: "🇵🇰"},
	"kk-KZ":   {Code: "kk-KZ", Name: "Kazakh (Kazakhstan)", NativeName: "Қазақша (Қазақстан)", Region: "Asia", Direction: "ltr", Emoji: "🇰🇿"},
	"mn-MN":   {Code: "mn-MN", Name: "Mongolian (Mongolia)", NativeName: "Монгол (Монгол)", Region: "Asia", Direction: "ltr", Emoji: "🇲🇳"},
	"bo-CN":   {Code: "bo-CN", Name: "Tibetan (China)", NativeName: "བོད་སྐད (ཀྲུང་གོ)", Region: "Asia", Direction: "ltr", Emoji: "🇨🇳"},
	"ug-CN":   {Code: "ug-CN", Name: "Uyghur (China)", NativeName: "ئۇيغۇرچە (جۇڭگو)", Region: "Asia", Direction: "rtl", Emoji: "🇨🇳"},
	"yue-HK":  {Code: "yue-HK", Name: "Cantonese (Hong Kong)", NativeName: "粵語 (香港)", Region: "Asia", Direction: "ltr", Emoji: "🇭🇰"},
	"km-KH":   {Code: "km-KH", Name: "Khmer (Cambodia)", NativeName: "ខ្មែរ (កម្ពុជា)", Region: "Asia", Direction: "ltr", Emoji: "🇰🇭"},
	"my-MM":   {Code: "my-MM", Name: "Burmese (Myanmar)", NativeName: "မြန်မာ (မြန်မာ)", Region: "Asia", Direction: "ltr", Emoji: "🇲🇲"},
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
