// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTranslator is a mocked struct that implements Translator.
type TestTranslator struct {
	mock.Mock
}

func (t TestTranslator) TranslateStructFieldErr(field string, err Error) (string, error) {
	m := t.Called(field, err)
	return m.String(0), m.Error(1)
}

func (t TestTranslator) TranslateSingleFieldErr(err Error) (string, error) {
	m := t.Called(err)
	return m.String(0), m.Error(1)
}

var _ Translator = TestTranslator{}

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

func TestError_Code(t *testing.T) {
	err := NewError("A", "msg")
	err2 := err.Code("B")

	assert.Equal(t, err.code, "A")
	assert.Equal(t, err.GetCode(), "A")

	assert.Equal(t, err2.code, "B")
	assert.Equal(t, err2.GetCode(), "B")
}

func TestError_Message(t *testing.T) {
	err := NewError("code", "A")
	err2 := err.Message("B")

	assert.Equal(t, err.message, "A")
	assert.Equal(t, err.GetMessage(), "A")

	assert.Equal(t, err2.message, "B")
	assert.Equal(t, err2.GetMessage(), "B")
}

func TestError_Params(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "AA": "val2"}
	p2 := map[string]interface{}{"B": "val1", "BB": "val2"}

	err := NewError("code", "A").Params(p)
	err2 := err.Message("B").Params(p2)

	assert.Equal(t, err.params, p)
	assert.Equal(t, err.GetParams(), p)

	assert.Equal(t, err2.params, p2)
	assert.Equal(t, err2.GetParams(), p2)
}

func TestValidationError(t *testing.T) {
	params := map[string]interface{}{
		"A": "B",
	}

	err := NewError("code", "msg").Params(params)
	assert.Equal(t, err.code, "code")
	assert.Equal(t, err.message, "msg")
	assert.Equal(t, err.params, params)

	params = map[string]interface{}{"min": 1}
	err = err.Params(params)

	assert.Equal(t, err.params, params)
}

func TestError_Translate(t *testing.T) {
	validationErr := NewError("code", "msg")
	translator := new(TestTranslator)

	translator.On("TranslateSingleFieldErr", validationErr).Return("abc", nil).Once()

	msg, err := validationErr.Translate(translator)

	assert.Equal(t, "abc", msg)
	assert.Nil(t, err)
}

func TestErrorReturnErr_Translate(t *testing.T) {
	validationErr := NewError("code", "msg")
	translator := new(TestTranslator)

	translationErr := errors.New("err")

	translator.On("TranslateSingleFieldErr", validationErr).Return("", translationErr).Once()

	msg, err := validationErr.Translate(translator)

	assert.Equal(t, "", msg)
	assert.Equal(t, translationErr, err)
}

func TestErrors_Translate(t *testing.T) {
	vErrors := Errors{
		"A": NewError("code", "msg"),
		"B": NewError("code2", "msg2"),
	}

	result := map[string]interface{}{"A": "TA", "B": "TB"}

	translator := new(TestTranslator)

	translator.On("TranslateStructFieldErr", "A", NewError("code", "msg")).Return("TA", nil).Once()
	translator.On("TranslateStructFieldErr", "B", NewError("code2", "msg2")).Return("TB", nil).Once()

	translatedMessages, err := vErrors.Translate(translator)

	assert.Equal(t, result, translatedMessages)
	assert.Nil(t, err)
}

func TestErrorsReturnErr_Translate(t *testing.T) {
	vErrors := Errors{
		"A": NewError("code", "msg"),
		"B": NewError("code2", "msg2"),
	}

	translateErr := errors.New("err")
	translator := new(TestTranslator)

	translator.On("TranslateStructFieldErr", "A", NewError("code", "msg")).Return("TA", nil).Once()
	translator.On("TranslateStructFieldErr", "B", NewError("code2", "msg2")).Return("TB", translateErr).Once()

	translatedMessages, err := vErrors.Translate(translator)

	assert.Nil(t, translatedMessages)
	assert.Equal(t, translateErr, err)
}

func TestErrorsWithInnerErrors_Translate(t *testing.T) {
	vErrors := Errors{
		"A": NewError("code", "msg"),
		"B": Errors{"b": NewError("code2", "msg2")},
		"C": errors.New("C"),
	}
	result := map[string]interface{}{
		"A": "TA",
		"B": map[string]interface{}{
			"b": "Tb",
		},
		"C": errors.New("C"),
	}
	translator := new(TestTranslator)

	translator.On("TranslateStructFieldErr", "A", NewError("code", "msg")).Return("TA", nil).Once()
	translator.On("TranslateStructFieldErr", "b", NewError("code2", "msg2")).Return("Tb", nil).Once()

	translatedMessages, err := vErrors.Translate(translator)

	assert.Equal(t, result, translatedMessages)
	assert.Nil(t, err)
}

func TestErrorsWithInnerErrorsReturnErr_Translate(t *testing.T) {
	vErrors := Errors{
		"A": NewError("code", "msg"),
		"B": Errors{"b": NewError("code2", "msg2")},
	}

	translateErr := errors.New("err")
	translator := new(TestTranslator)

	translator.On("TranslateStructFieldErr", "A", NewError("code", "msg")).Return("TA", nil).Once()
	translator.On("TranslateStructFieldErr", "b", NewError("code2", "msg2")).Return("Tb", translateErr).Once()

	translatedMessages, err := vErrors.Translate(translator)

	assert.Nil(t, translatedMessages)
	assert.Equal(t, translateErr, err)
}
