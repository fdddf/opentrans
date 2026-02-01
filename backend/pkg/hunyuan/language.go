package hunyuan

import (
	"strings"
)

// Language 表示混元 AI 支持的语言
type Language struct {
	ShortCode   string // 短语言代码，如 "zh", "en"
	Identifier  string // 混元模型识别的标识符，如 "Chinese", "English"
	DisplayName string // 显示名称
	Emoji       string // 国旗 emoji
}

// hunyuanLanguages 包含混元 AI 支持的所有语言映射
var hunyuanLanguages = []Language{
	{"zh", "Chinese", "Simplified Chinese", "🇨🇳"},
	{"en", "English", "English", "🇺🇸"},
	{"fr", "French", "French", "🇫🇷"},
	{"pt", "Portuguese", "Portuguese", "🇵🇹"},
	{"es", "Spanish", "Spanish", "🇪🇸"},
	{"ja", "Japanese", "Japanese", "🇯🇵"},
	{"tr", "Turkish", "Turkish", "🇹🇷"},
	{"ru", "Russian", "Russian", "🇷🇺"},
	{"ar", "Arabic", "Arabic", "🇸🇦"},
	{"ko", "Korean", "Korean", "🇰🇷"},
	{"th", "Thai", "Thai", "🇹🇭"},
	{"it", "Italian", "Italian", "🇮🇹"},
	{"de", "German", "German", "🇩🇪"},
	{"vi", "Vietnamese", "Vietnamese", "🇻🇳"},
	{"ms", "Malay", "Malay", "🇲🇾"},
	{"id", "Indonesian", "Indonesian", "🇮🇩"},
	{"tl", "Filipino", "Filipino", "🇵🇭"},
	{"hi", "Hindi", "Hindi", "🇮🇳"},
	{"zh-Hant", "Traditional Chinese", "Traditional Chinese", "🌐"},
	{"pl", "Polish", "Polish", "🇵🇱"},
	{"cs", "Czech", "Czech", "🇨🇿"},
	{"nl", "Dutch", "Dutch", "🇳🇱"},
	{"km", "Khmer", "Khmer", "🇰🇭"},
	{"my", "Burmese", "Burmese", "🇲🇲"},
	{"fa", "Persian", "Persian", "🇮🇷"},
	{"gu", "Gujarati", "Gujarati", "🇮🇳"},
	{"ur", "Urdu", "Urdu", "🇵🇰"},
	{"te", "Telugu", "Telugu", "🇮🇳"},
	{"mr", "Marathi", "Marathi", "🇮🇳"},
	{"he", "Hebrew", "Hebrew", "🇮🇱"},
	{"bn", "Bengali", "Bengali", "🇧🇩"},
	{"ta", "Tamil", "Tamil", "🇮🇳"},
	{"uk", "Ukrainian", "Ukrainian", "🇺🇦"},
	{"bo", "Tibetan", "Tibetan", "🇨🇳"},
	{"kk", "Kazakh", "Kazakh", "🇰🇿"},
	{"mn", "Mongolian", "Mongolian", "🇲🇳"},
	{"ug", "Uyghur", "Uyghur", "🇨🇳"},
	{"yue", "Cantonese", "Cantonese", "🇭🇰"},
}

// appStoreLanguageMapping 将 Apple App Store Connect 的 locale 映射到短语言代码
var appStoreLanguageMapping = map[string]string{
	// 简体中文
	"zh-Hans": "zh",
	"zh-CN":   "zh",
	// 繁体中文
	"zh-Hant": "zh-Hant",
	"zh-TW":   "zh-Hant",
	"zh-HK":   "zh-Hant",
	// 英语
	"en":    "en",
	"en-US": "en",
	"en-GB": "en",
	"en-CA": "en",
	"en-AU": "en",
	// 法语
	"fr":    "fr",
	"fr-FR": "fr",
	"fr-CA": "fr",
	// 葡萄牙语
	"pt":    "pt",
	"pt-PT": "pt",
	"pt-BR": "pt",
	// 西班牙语
	"es":    "es",
	"es-ES": "es",
	"es-MX": "es",
	// 日语
	"ja":    "ja",
	"ja-JP": "ja",
	// 土耳其语
	"tr":    "tr",
	"tr-TR": "tr",
	// 俄语
	"ru":    "ru",
	"ru-RU": "ru",
	// 阿拉伯语
	"ar":    "ar",
	"ar-SA": "ar",
	// 韩语
	"ko":    "ko",
	"ko-KR": "ko",
	// 泰语
	"th":    "th",
	"th-TH": "th",
	// 意大利语
	"it":    "it",
	"it-IT": "it",
	// 德语
	"de":    "de",
	"de-DE": "de",
	// 越南语
	"vi":    "vi",
	"vi-VN": "vi",
	// 马来语
	"ms":    "ms",
	"ms-MY": "ms",
	// 印尼语
	"id":    "id",
	"id-ID": "id",
	// 印地语
	"hi":    "hi",
	"hi-IN": "hi",
	// 波兰语
	"pl":    "pl",
	"pl-PL": "pl",
	// 捷克语
	"cs":    "cs",
	"cs-CZ": "cs",
	// 荷兰语
	"nl":    "nl",
	"nl-NL": "nl",
	// 希伯来语
	"he":    "he",
	"he-IL": "he",
	// 乌克兰语
	"uk":    "uk",
	"uk-UA": "uk",
	// 菲律宾语
	"fil":    "tl",
	"fil-PH": "tl",
	"tl":     "tl",
	// 孟加拉语
	"bn":    "bn",
	"bn-IN": "bn",
	// 泰米尔语
	"ta":    "ta",
	"ta-IN": "ta",
	// 高棉语
	"km":    "km",
	"km-KH": "km",
	// 缅甸语
	"my":    "my",
	"my-MM": "my",
	// 波斯语
	"fa":    "fa",
	"fa-IR": "fa",
	// 古吉拉特语
	"gu":    "gu",
	"gu-IN": "gu",
	// 乌尔都语
	"ur":    "ur",
	"ur-PK": "ur",
	// 泰卢固语
	"te":    "te",
	"te-IN": "te",
	// 马拉地语
	"mr":    "mr",
	"mr-IN": "mr",
	// 藏语
	"bo":    "bo",
	"bo-CN": "bo",
	// 哈萨克语
	"kk":    "kk",
	"kk-KZ": "kk",
	// 蒙古语
	"mn":    "mn",
	"mn-MN": "mn",
	// 维吾尔语
	"ug":    "ug",
	"ug-CN": "ug",
	// 粤语
	"yue":    "yue",
	"yue-HK": "yue",
	// "zh-HK":   "yue",
}

// GetLanguageByShortCode 根据短语言代码获取语言信息
func GetLanguageByShortCode(shortCode string) *Language {
	for i := range hunyuanLanguages {
		if hunyuanLanguages[i].ShortCode == shortCode {
			return &hunyuanLanguages[i]
		}
	}
	return nil
}

// GetLanguageByIdentifier 根据混元标识符获取语言信息
func GetLanguageByIdentifier(identifier string) *Language {
	for i := range hunyuanLanguages {
		if hunyuanLanguages[i].Identifier == identifier {
			return &hunyuanLanguages[i]
		}
	}
	return nil
}

// MapAppStoreLocaleToHunyuan 将 Apple App Store Connect 的 locale 映射到混元 AI 的语言标识符
func MapAppStoreLocaleToHunyuan(locale string) string {
	// 标准化 locale 格式：将下划线替换为连字符
	normalized := strings.ReplaceAll(locale, "_", "-")
	normalized = strings.TrimSpace(normalized)

	// 1. 直接查找完整的 locale 映射
	if shortCode, ok := appStoreLanguageMapping[normalized]; ok {
		if lang := GetLanguageByShortCode(shortCode); lang != nil {
			return lang.Identifier
		}
	}

	// 2. 尝试按连字符分割，只使用基础语言代码
	parts := strings.Split(normalized, "-")
	if len(parts) > 0 {
		baseCode := parts[0]
		if shortCode, ok := appStoreLanguageMapping[baseCode]; ok {
			if lang := GetLanguageByShortCode(shortCode); lang != nil {
				return lang.Identifier
			}
		}
		// 3. 直接使用基础代码作为短代码查找
		if lang := GetLanguageByShortCode(baseCode); lang != nil {
			return lang.Identifier
		}
	}

	// 4. 如果都找不到，返回原始 locale
	return locale
}

// MapAppStoreLocaleToDisplayName 将 Apple App Store Connect 的 locale 映射到显示名称
func MapAppStoreLocaleToDisplayName(locale string) string {
	identifier := MapAppStoreLocaleToHunyuan(locale)
	if lang := GetLanguageByIdentifier(identifier); lang != nil {
		return lang.DisplayName
	}
	return locale
}

// MapAppStoreLocaleToEmoji 将 Apple App Store Connect 的 locale 映射到 emoji
func MapAppStoreLocaleToEmoji(locale string) string {
	identifier := MapAppStoreLocaleToHunyuan(locale)
	if lang := GetLanguageByIdentifier(identifier); lang != nil {
		return lang.Emoji
	}
	return "🌐"
}

// MapAppStoreLocaleToShortCode 将 Apple App Store Connect 的 locale 映射到短语言代码
func MapAppStoreLocaleToShortCode(locale string) string {
	// 标准化 locale 格式
	normalized := strings.ReplaceAll(locale, "_", "-")
	normalized = strings.TrimSpace(normalized)

	// 直接查找
	if shortCode, ok := appStoreLanguageMapping[normalized]; ok {
		return shortCode
	}

	// 尝试基础语言代码
	parts := strings.Split(normalized, "-")
	if len(parts) > 0 {
		baseCode := parts[0]
		if shortCode, ok := appStoreLanguageMapping[baseCode]; ok {
			return shortCode
		}
	}

	// 返回基础代码
	return parts[0]
}

// GetHunyuanIdentifiers 返回所有混元 AI 支持的语言标识符
func GetHunyuanIdentifiers() []string {
	identifiers := make([]string, len(hunyuanLanguages))
	for i := range hunyuanLanguages {
		identifiers[i] = hunyuanLanguages[i].Identifier
	}
	return identifiers
}

// GetShortCodes 返回所有短语言代码
func GetShortCodes() []string {
	codes := make([]string, len(hunyuanLanguages))
	for i := range hunyuanLanguages {
		codes[i] = hunyuanLanguages[i].ShortCode
	}
	return codes
}
