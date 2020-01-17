package validation

// LangMap use to define language rules translation map.
type LangMap map[string]string

var EnLang = "en"

// Language to use for validation error messages.
var Lang = EnLang

func Msg(ruleName string, customMsg string) string {
	return msgInLang(Lang, ruleName, "", customMsg)
}

func MsgWithDefault(ruleName string, defaultMsg string, customMsg string) string {
	return msgInLang(Lang, ruleName, defaultMsg, customMsg)
}

func MsgInLang(lang string, ruleName string, defaultMsg string, customMsg string) string {
	return msgInLang(lang, ruleName, defaultMsg, customMsg)
}

// msgInLang check if rule exists in the language, return it, otherwise
// check that rule in en language, and if does not exists, finally return
// default message.
func msgInLang(lang string, ruleName string, defaultMsg string, customMsg string) string {
	if customMsg != "" {
		return customMsg
	}

	if translation := getRuleTranslation(lang, ruleName); translation != "" {
		return translation
	}

	if translation := getRuleTranslation(EnLang, ruleName); translation != "" {
		return translation
	}

	return defaultMsg
}

func getRuleTranslation(lang, ruleName string) string {
	langMap, ok := TranslationMap[lang]

	if !ok {
		return ""
	}

	if translatedRule, ok := langMap[ruleName]; ok {
		return translatedRule
	}

	return ""
}
