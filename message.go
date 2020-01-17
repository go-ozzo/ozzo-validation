package validation

// LangMap use to define language rules translation map.
type LangMap map[string]string

// EnLang is key of translated messages map for english language.
var EnLang = "en"

// Language to use for validation error messages.
// To specify your current language for rules messages
// translation, set this variable.
var Lang = EnLang

// Msg check if custom message is not empty, it return
// your custom messgae, otherwise try to find translated
// message for that rule.
func Msg(ruleName string, customMsg string) string {
	return msgInLang(Lang, ruleName, "", customMsg)
}

// MsgWithDefault return rules translated message , if not found return
// default message.
func MsgWithDefault(ruleName string, defaultMsg string, customMsg string) string {
	return msgInLang(Lang, ruleName, defaultMsg, customMsg)
}

// MsgInLang check if you pass custom message, it return your custom message.
// otherwise try to find your rule translation in specified lang, if not found,
// search in "en" lang , finally if not found again, return your efault message.
func MsgInLang(lang string, ruleName string, defaultMsg string, customMsg string) string {
	return msgInLang(lang, ruleName, defaultMsg, customMsg)
}

// msgInLang check if you pass custom message, it return your custom message.
// otherwise try to find your rule translation in specified lang, if not found,
// search in "en" lang , finally if not found again, return your efault message.
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

// getRuleTranslation return rule's translated message.
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
