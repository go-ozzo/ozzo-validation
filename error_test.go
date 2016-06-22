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
