package cmd

import (
	"strings"
)

// normalizeLanguageCode normalizes language codes and names to standard formats
func normalizeLanguageCode(language string) string {
	if strings.TrimSpace(language) == "" {
		return language
	}

	normalized := strings.ToLower(strings.TrimSpace(language))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	normalized = strings.Join(strings.Fields(normalized), " ")
	normalized = strings.ReplaceAll(normalized, "(", "")
	normalized = strings.ReplaceAll(normalized, ")", "")

	codeByToken := map[string]string{
		"ar":      "ar",
		"cs":      "cs",
		"da":      "da",
		"de":      "de",
		"en":      "en",
		"es":      "es",
		"fi":      "fi",
		"fr":      "fr",
		"he":      "he",
		"hi":      "hi",
		"id":      "id",
		"it":      "it",
		"ja":      "ja",
		"ko":      "ko",
		"nl":      "nl",
		"no":      "no",
		"pl":      "pl",
		"pt":      "pt",
		"ru":      "ru",
		"sv":      "sv",
		"th":      "th",
		"tr":      "tr",
		"vi":      "vi",
		"zh-hans": "zh-Hans",
		"zh-hant": "zh-Hant",
	}
	if code, ok := codeByToken[normalized]; ok {
		return code
	}

	nameByToken := map[string]string{
		"arabic":              "ar",
		"chinese simplified":  "zh-Hans",
		"chinese traditional": "zh-Hant",
		"czech":               "cs",
		"danish":              "da",
		"dutch":               "nl",
		"english":             "en",
		"finnish":             "fi",
		"french":              "fr",
		"german":              "de",
		"hebrew":              "he",
		"hindi":               "hi",
		"indonesian":          "id",
		"italian":             "it",
		"japanese":            "ja",
		"korean":              "ko",
		"norwegian":           "no",
		"polish":              "pl",
		"portuguese":          "pt",
		"russian":             "ru",
		"simplified chinese":  "zh-Hans",
		"spanish":             "es",
		"swedish":             "sv",
		"thai":                "th",
		"traditional chinese": "zh-Hant",
		"turkish":             "tr",
		"vietnamese":          "vi",
	}
	if code, ok := nameByToken[normalized]; ok {
		return code
	}

	return language
}
