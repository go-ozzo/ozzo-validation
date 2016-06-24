// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestError_MarshalMessage(t *testing.T) {
	errs := Errors{
		"A": errors.New("A1"),
		"B": SliceErrors{
			2: errors.New("B1"),
		},
	}
	errsJSON, err := errs.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "{\"A\":\"A1\",\"B\":{\"2\":\"B1\"}}", string(errsJSON))
}

func TestSliceErrors_Error(t *testing.T) {
	errs := SliceErrors{
		3: errors.New("B1"),
		0: errors.New("C1"),
		1: errors.New("A1"),
	}
	assert.Equal(t, "0: C1; 1: A1; 3: B1.", errs.Error())

	errs = SliceErrors{
		1: errors.New("B1"),
	}
	assert.Equal(t, "1: B1.", errs.Error())

	errs = SliceErrors{}
	assert.Equal(t, "", errs.Error())
}

func TestSliceError_MarshalMessage(t *testing.T) {
	errs := SliceErrors{
		2: Errors{
			"B": errors.New("B1"),
		},
		0: errors.New("A1"),
	}
	errsJSON, err := errs.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "{\"0\":\"A1\",\"2\":{\"B\":\"B1\"}}", string(errsJSON))
}
