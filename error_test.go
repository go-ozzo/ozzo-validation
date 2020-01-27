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

func TestError_SetCode(t *testing.T) {
	err := NewError("A", "msg")

	assert.Equal(t, err.code, "A")
	assert.Equal(t, err.Code(), "A")

	err.SetCode("B")
	assert.Equal(t, err.code, "B")
}

func TestError_Code(t *testing.T) {
	err := NewError("A", "msg")

	assert.Equal(t, err.code, "A")
	assert.Equal(t, err.Code(), "A")
}

func TestError_Message(t *testing.T) {
	err := NewError("code", "A")

	assert.Equal(t, err.message, "A")
	assert.Equal(t, err.Message(), "A")
}

func TestError_Params(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "AA": "val2"}

	err := NewError("code", "A")
	err.SetParams(p)
	err.SetMessage("B")

	assert.Equal(t, err.params, p)
	assert.Equal(t, err.Params(), p)
}

func TestError_AddParam(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "B": "val2"}

	err := NewError("code", "A")
	err.SetParams(p)
	err.AddParam("C", "val3")

	p["C"] = "val3"

	assert.Equal(t, err.params, p)
	assert.Equal(t, err.Params(), p)
}

func TestValidationError(t *testing.T) {
	params := map[string]interface{}{
		"A": "B",
	}

	err := NewError("code", "msg")
	err.SetParams(params)

	assert.Equal(t, err.code, "code")
	assert.Equal(t, err.message, "msg")
	assert.Equal(t, err.params, params)

	params = map[string]interface{}{"min": 1}
	err.SetParams(params)

	assert.Equal(t, err.params, params)
}
