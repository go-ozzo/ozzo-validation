// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func validateMe(s string) bool {
	return s == "me"
}

func TestNewStringRule(t *testing.T) {
	v := NewStringRule(validateMe, "abc")
	assert.NotNil(t, v.validate)
	assert.Equal(t, "abc", v.message)
}

func TestStringValidator_Error(t *testing.T) {
	v := NewStringRule(validateMe, "abc")
	assert.Equal(t, "abc", v.message)
	v2 := v.Error("correct")
	assert.Equal(t, "correct", v2.message)
	assert.Equal(t, "abc", v.message)
}

func TestStringValidator_Validate(t *testing.T) {
	v := NewStringRule(validateMe, "wrong")

	err := v.Validate("me", nil)
	assert.Nil(t, err)

	err = v.Validate("", nil)
	assert.Nil(t, err)

	err = v.Validate("not me", nil)
	if assert.NotNil(t, err) {
		assert.Equal(t, "wrong", err.Error())
	}

	err = v.Validate(100, nil)
	if assert.NotNil(t, err) {
		assert.NotEqual(t, "wrong", err.Error())
	}

	v2 := v.Error("Wrong!")
	err = v2.Validate("not me", nil)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Wrong!", err.Error())
	}
}
