package services

import "fmt"

// SupportedLanguages is a curated list of allowed language codes
var SupportedLanguages = map[string]struct{}{
	"en": {}, "en-US": {}, "en-GB": {}, "zh-Hans": {}, "zh-Hant": {}, "ja": {},
	"fr": {}, "de": {}, "es": {}, "ko": {}, "ru": {}, "pt": {}, "pt-BR": {}, "it": {},
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
