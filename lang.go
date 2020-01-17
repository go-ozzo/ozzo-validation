package validation

// LangMap use to define language rules translation map.
type LangMap map[string]string

// EnLang is key of translated messages map for english language.
var EnLang = "en"

// Language to use for validation error messages.
// To specify your current language for rules messages
// translation, set this variable.
var Lang = EnLang

// AddRuleTranslation add or replace translation of the rule.
func AddRuleTranslation(lang, ruleName, translation string) {
	setEmptyLangMapIfNotExists(lang)

	TranslationMap[lang][ruleName] = translation
}

// Add language to the TranslationMap, if already,
// exists that language, union two maps.
func AddLang(lang string, langMap LangMap) {
	setEmptyLangMapIfNotExists(lang)

	for k, v := range langMap {
		TranslationMap[lang][k] = v
	}
}

// setEmptyLangMapIfNotExists check if does not exists map for
//specified language, create new zero value map.
func setEmptyLangMapIfNotExists(lang string) {
	if _, ok := TranslationMap[lang]; !ok {
		TranslationMap[lang] = LangMap{}
	}
}
