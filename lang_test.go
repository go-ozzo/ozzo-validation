// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddLang(t *testing.T) {
	Fa := "fa"

	defer func() {
		delete(TranslationMap, Fa)
	}()

	langMap := LangMap{
		"required":  "فیلد ضروری می باشد.",
		"not_empty": "فیلد نمیتواند خالی باشد.",
	}

	AddLang(Fa, langMap)
	realLangMap, ok := TranslationMap[Fa]

	assert.True(t, ok)
	assert.Equal(t, langMap, realLangMap)
}

func TestAddRuleTranslation(t *testing.T) {
	Fa := "fa"

	defer func() {
		delete(TranslationMap, Fa)
	}()

	dateRage := "بازه زمانی اشتباه می باشد."
	match := "فرمت داده نامعتبر می باشد."

	AddTranslation(Fa, "date_range", dateRage)
	AddTranslation(Fa, "match", match)

	assert.Equal(t, dateRage, TranslationMap[Fa]["date_range"])
	assert.Equal(t, match, TranslationMap[Fa]["match"])
}
