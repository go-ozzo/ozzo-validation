// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

// LangMap contain translation of rules in languages
type LangMap map[string]string

// EnLang is the key to translated messages map for the English language.
var EnLang = "en"

// Lang is the language to use for validation error messages.
// To specify your current language for rules messages
// translation, set this variable.
var Lang = EnLang

// AddRuleTranslation add or replace translation of the rule.
func AddRuleTranslation(lang, ruleName, translation string) {
	setEmptyLangMapIfNotExists(lang)

	TranslationMap[lang][ruleName] = translation
}

// AddLang add language to the TranslationMap, if already
// exists that language, union two maps.
func AddLang(lang string, langMap LangMap) {
	setEmptyLangMapIfNotExists(lang)

	for k, v := range langMap {
		TranslationMap[lang][k] = v
	}
}

// setEmptyLangMapIfNotExists check if does not exists the translation
// map of the language, create a new zero value map.
func setEmptyLangMapIfNotExists(lang string) {
	if _, ok := TranslationMap[lang]; !ok {
		TranslationMap[lang] = LangMap{}
	}
}
