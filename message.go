// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

// Msg check if the custom message is not empty, it returns your custom
// message otherwise tries to find the translated message for that rule.
func Msg(traslationKey string, customMsg string) string {
	return msgInLang(Lang, traslationKey, "", customMsg) // TODO: check this line, we should not use Lang variable here.
}

// MsgWithDefault return rules translated message, if
//not found return default message.
func MsgWithDefault(translationKey string, defaultMsg string, customMsg string) string {
	return msgInLang(Lang, translationKey, defaultMsg, customMsg) // TODO: check Lang, we should not use it here.
}

// MsgInLang checks if you pass the custom message, it returns your
// custom message otherwise try to find your rule translation in
// specified lang if not found search in "en" language, finally
// if the English lang does no having any translation for that
// rule, returns your default message.
func MsgInLang(lang string, translationKey string, defaultMsg string, customMsg string) string {
	return msgInLang(lang, translationKey, defaultMsg, customMsg)
}

// msgInLang checks if you pass the custom message, it returns your
// custom message otherwise try to find your rule translation in
// specified lang if not found search in "en" language, finally
// if the English lang does no having any translation for that
// rule, returns your default message.
func msgInLang(lang string, translationKey string, defaultMsg string, customMsg string) string {
	if customMsg != "" {
		return customMsg
	}

	if lang == "" {
		lang = Lang
	}

	if translation := getTranslation(lang, translationKey); translation != "" {
		return translation
	}

	if translation := getTranslation(EnLang, translationKey); translation != "" {
		return translation
	}

	return defaultMsg
}

// getTranslation return rule's translated message.
func getTranslation(lang, translationKey string) string {
	langMap, ok := TranslationMap[lang]

	if !ok {
		return ""
	}

	if translatedRule, ok := langMap[translationKey]; ok {
		return translatedRule
	}

	return ""
}
