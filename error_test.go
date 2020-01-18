// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInternalError(t *testing.T) {
	err := NewInternalError(errors.New("abc"))
	if assert.NotNil(t, err.InternalError()) {
		assert.Equal(t, "abc", err.InternalError().Error())
	}
}

func TestErrors_Error(t *testing.T) {
	errs := Errors{
		"B": errors.New("B1"),
		"C": errors.New("C1"),
		"A": errors.New("A1"),
	}
	assert.Equal(t, "A: A1; B: B1; C: C1.", errs.Error())

	errs = Errors{
		"B": errors.New("B1"),
	}
	assert.Equal(t, "B: B1.", errs.Error())

	errs = Errors{}
	assert.Equal(t, "", errs.Error())
}

func TestErrors_MarshalMessage(t *testing.T) {
	errs := Errors{
		"A": errors.New("A1"),
		"B": Errors{
			"2": errors.New("B1"),
		},
	}
	errsJSON, err := errs.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "{\"A\":\"A1\",\"B\":{\"2\":\"B1\"}}", string(errsJSON))
}

func TestErrors_Filter(t *testing.T) {
	errs := Errors{
		"B": errors.New("B1"),
		"C": nil,
		"A": errors.New("A1"),
	}
	err := errs.Filter()
	assert.Equal(t, 2, len(errs))
	if assert.NotNil(t, err) {
		assert.Equal(t, "A: A1; B: B1.", err.Error())
	}

	errs = Errors{}
	assert.Nil(t, errs.Filter())

	errs = Errors{
		"B": nil,
		"C": nil,
	}

	assert.Nil(t, errs.Filter())
}

func TestNewValidationErr(t *testing.T) {
	err := newErrMessage("translationKey", "custom_msg")
	assert.Equal(t, err.TranslationKey, "translationKey")
	assert.Equal(t, err.Message, "custom_msg")
	assert.Equal(t, err.DefaultMessage, "")
	assert.Equal(t, err.Params, []interface{}(nil))

	err = err.Default("abc")
	assert.Equal(t, err.DefaultMessage, "abc")

	params := []interface{}{1, 2, 3}
	err = err.SetParams(params)

	assert.Equal(t, err.Params, params)
}

func TestGetValidationErrorWithDefault(t *testing.T) {
	params := []interface{}{1, 2, "3"}
	err := ErrMessageWithDefault("translationKey", "translationKey", params)

	assert.Equal(t, err.TranslationKey, "translationKey")
	assert.Equal(t, err.Message, "")
	assert.Equal(t, err.DefaultMessage, "translationKey")
	assert.Equal(t, err.Params, params)

	err = ErrMessageWithDefault("translationKey", "", nil)

	assert.Equal(t, err.TranslationKey, "translationKey")
	assert.Equal(t, err.Message, "")
	assert.Equal(t, err.DefaultMessage, "")
	assert.Equal(t, err.Params, []interface{}(nil))

	params = []interface{}{1, 2}
	err = err.CustomMessage("custom_message", params)
	assert.Equal(t, err.Message, "custom_message")
	assert.Equal(t, err.Params, params)

}

func TestTranslateInOtherLanguage(t *testing.T) {
	Fa := "fa"
	defer func() {
		delete(TranslationMap, Fa)
		Lang = EnLang
	}()

	AddTranslation(Fa, "required", "نیازه")
	err := newErrMessage("required", "")
	assert.Equal(t, err.Error(), "cannot be blank")

	err2 := err.ToLang(Fa)
	assert.Equal(t, err.Error(), "cannot be blank")
	assert.Equal(t, err2.Error(), "نیازه")
}

func TestChangeErrorsLang(t *testing.T) {
	Fa := "fa"
	defer func() {
		delete(TranslationMap, Fa)
		Lang = EnLang
	}()

	AddTranslation(Fa, "required", "نیازه")
	err := newErrMessage("required", "")
	err2 := err.ToLang(Fa)

	errGroup := Errors{
		"field_1":       err,
		"field_2":       err2,
		"field_group":   Errors{"sub_1": err},
		"unknown_field": errors.New("unknown_err"),
	}

	result := Errors{
		"field_1":       err.ToLang(Fa),
		"field_2":       err2.ToLang(Fa),
		"field_group":   Errors{"sub_1": err.ToLang(Fa)},
		"unknown_field": errors.New("unknown_err"),
	}

	errGroup = errGroup.ToLang(Fa)

	assert.Equal(t, errGroup, result)
}
