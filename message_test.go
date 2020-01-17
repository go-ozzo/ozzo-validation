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

func TestGetMessageInOtherLanguage(t *testing.T) {
	defer func() {
		Lang = EnLang
	}()

	Lang = "fa"
	assert.Equal(t, "abc", Msg("required", "abc"))
	assert.Equal(t, "cannot be blank", Msg("required", ""))
	assert.Equal(t, "cannot be blank", MsgInLang(EnLang, "required", "", ""))

	assert.Equal(t, "", getRuleTranslation(Lang, "required"))
}
