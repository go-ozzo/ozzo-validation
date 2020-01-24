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
	err2 := err.SetCode("B")

	assert.Equal(t, err.code, "A")
	assert.Equal(t, err.Code(), "A")

	assert.Equal(t, err2.code, "B")
	assert.Equal(t, err2.Code(), "B")
}

func TestError_Message(t *testing.T) {
	err := NewError("code", "A")
	err2 := err.SetMessage("B")

	assert.Equal(t, err.message, "A")
	assert.Equal(t, err.Message(), "A")

	assert.Equal(t, err2.message, "B")
	assert.Equal(t, err2.Message(), "B")
}

func TestError_Params(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "AA": "val2"}
	p2 := map[string]interface{}{"B": "val1", "BB": "val2"}

	err := NewError("code", "A").SetParams(p)
	err2 := err.SetMessage("B").SetParams(p2)

	assert.Equal(t, err.params, p)
	assert.Equal(t, err.Params(), p)

	assert.Equal(t, err2.params, p2)
	assert.Equal(t, err2.Params(), p2)
}

func TestValidationError(t *testing.T) {
	params := map[string]interface{}{
		"A": "B",
	}

	err := NewError("code", "msg").SetParams(params)
	assert.Equal(t, err.code, "code")
	assert.Equal(t, err.message, "msg")
	assert.Equal(t, err.params, params)

	params = map[string]interface{}{"min": 1}
	err = err.SetParams(params)

	assert.Equal(t, err.params, params)
}
