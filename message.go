// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

// Msg check if the custom message is not empty, it returns your custom
// message otherwise tries to find the translated message for that rule.
func Msg(ruleName string, customMsg string) string {
	return msgInLang(Lang, ruleName, "", customMsg)
}

// MsgWithDefault return rules translated message, if
//not found return default message.
func MsgWithDefault(ruleName string, defaultMsg string, customMsg string) string {
	return msgInLang(Lang, ruleName, defaultMsg, customMsg)
}

// MsgInLang checks if you pass the custom message, it returns your
// custom message otherwise try to find your rule translation in
// specified lang if not found search in "en" language, finally
// if the English lang does no having any translation for that
// rule, returns your default message.
func MsgInLang(lang string, ruleName string, defaultMsg string, customMsg string) string {
	return msgInLang(lang, ruleName, defaultMsg, customMsg)
}

// msgInLang checks if you pass the custom message, it returns your
// custom message otherwise try to find your rule translation in
// specified lang if not found search in "en" language, finally
// if the English lang does no having any translation for that
// rule, returns your default message.
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
