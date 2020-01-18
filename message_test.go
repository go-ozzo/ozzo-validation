// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMessage(t *testing.T) {
	assert.Equal(t, "", Msg("unknown_rule", ""))
	assert.Equal(t, "abc", Msg("required", "abc"))
	assert.Equal(t, "cannot be blank", Msg("required", ""))

	assert.Equal(t, "test", MsgWithDefault("unknown_rule", "test", ""))
	assert.Equal(t, "abc", MsgWithDefault("unknown_rule", "test", "abc"))

	assert.Equal(t, "test", MsgInLang(EnLang, "unknown_rule", "test", ""))
	assert.Equal(t, "abc", MsgInLang(EnLang, "unknown_rule", "test", "abc"))
}

func TestGetMessageInUnknownLanguage(t *testing.T) {
	Fa := "fa"

	defer func() {
		Lang = EnLang
		delete(TranslationMap, Fa)
	}()

	// Change current translation language to "fa"
	Lang = Fa

	assert.Equal(t, "abc", Msg("required", "abc"))
	assert.Equal(t, "cannot be blank", Msg("required", ""))
	assert.Equal(t, "cannot be blank", MsgInLang(EnLang, "required", "", ""))

	assert.Equal(t, "", getTranslation(Lang, "required"))
}

func TestGetTranslatedMessageInOtherLanguage(t *testing.T) {
	Fa := "fa"

	defer func() {
		Lang = EnLang
		delete(TranslationMap, Fa)
	}()

	// Change current translation language to "fa"
	Lang = Fa

	langMap := LangMap{
		"required":  "فیلد ضروری می باشد.",
		"not_empty": "فیلد نمیتواند خالی باشد.",
	}

	AddLang(Fa, langMap)
	assert.Equal(t, langMap["required"], Msg("required", ""))
	assert.Equal(t, langMap["not_empty"], Msg("not_empty", ""))

	assert.Equal(t, "cannot be blank", MsgInLang(EnLang, "required", "", ""))
	assert.Equal(t, "abc", MsgInLang(EnLang, "unknown_rule", "abc", ""))
	assert.Equal(t, "custom message", MsgInLang(EnLang, "unknown_rule", "abc", "custom message"))

}
