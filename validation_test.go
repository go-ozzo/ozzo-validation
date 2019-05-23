// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	slice := []String123{String123("abc"), String123("123"), String123("xyz")}
	ctxSlice := []Model4{{A: "abc"}, {A: "def"}}
	mp := map[string]String123{"c": String123("abc"), "b": String123("123"), "a": String123("xyz")}
	var ptr *string
	tests := []struct {
		tag   string
		value interface{}
		err   string
		errWithContext   string
	}{
		{"t1", 123, "", ""},
		{"t2", String123("123"), "", ""},
		{"t3", String123("abc"), "error 123", "error 123"},
		{"t4", []String123{}, "", ""},
		{"t5", slice, "0: error 123; 2: error 123.", "0: error 123; 2: error 123."},
		{"t6", &slice, "0: error 123; 2: error 123.", "0: error 123; 2: error 123."},
		{"t7", ctxSlice, "", "1: (A: error abc.)."},
		{"t8", mp, "a: error 123; c: error 123.", "a: error 123; c: error 123."},
		{"t9", &mp, "a: error 123; c: error 123.", "a: error 123; c: error 123."},
		{"t10", map[string]String123{}, "", ""},
		{"t11", ptr, "", ""},
	}
	for _, test := range tests {
		err := Validate(test.value)
		assertError(t, test.err, err, test.tag)
		err = ValidateWithContext(context.Background(), test.value)
		assertError(t, test.errWithContext, err, test.tag)
	}

	// with rules
	err := Validate("123", &validateAbc{}, &validateXyz{})
	if assert.NotNil(t, err) {
		assert.Equal(t, "error abc", err.Error())
	}
	err = Validate("abc", &validateAbc{}, &validateXyz{})
	if assert.NotNil(t, err) {
		assert.Equal(t, "error xyz", err.Error())
	}
	err = Validate("abcxyz", &validateAbc{}, &validateXyz{})
	assert.Nil(t, err)

	err = Validate("123", &validateAbc{}, Skip, &validateXyz{})
	if assert.NotNil(t, err) {
		assert.Equal(t, "error abc", err.Error())
	}
	err = Validate("abc", &validateAbc{}, Skip, &validateXyz{})
	assert.Nil(t, err)
}

func TestBy(t *testing.T) {
	abcRule := By(func(value interface{}) error {
		s, _ := value.(string)
		if s != "abc" {
			return errors.New("must be abc")
		}
		return nil
	})
	assert.Nil(t, Validate("abc", abcRule))
	err := Validate("xyz", abcRule)
	if assert.NotNil(t, err) {
		assert.Equal(t, "must be abc", err.Error())
	}
}

func TestByWithContext(t *testing.T) {
	abcRule := ByWithContext(func(ctx context.Context, value interface{}) error {
		s, _ := value.(string)
		if s != "abc" {
			return errors.New("must be abc")
		}
		return nil
	})
	assert.Nil(t, ValidateWithContext(context.Background(), "abc", abcRule))
	err := ValidateWithContext(context.Background(), "xyz", abcRule)
	if assert.NotNil(t, err) {
		assert.Equal(t, "must be abc", err.Error())
	}
}

func Test_skipRule_Validate(t *testing.T) {
	assert.Nil(t, Skip.Validate(100))
}

func Test_skipRule_ValidateWithContext(t *testing.T) {
	assert.Nil(t, Skip.ValidateWithContext(context.Background(), 100))
}

func assertError(t *testing.T, expected string, err error, tag string) {
	if expected == "" {
		assert.Nil(t, err, tag)
	} else if assert.NotNil(t, err, tag) {
		assert.Equal(t, expected, err.Error(), tag)
	}
}

type validateAbc struct{}

func (v *validateAbc) Validate(obj interface{}) error {
	if !strings.Contains(obj.(string), "abc") {
		return errors.New("error abc")
	}
	return nil
}

type validateContextAbc struct{}

func (v *validateContextAbc) ValidateWithContext(ctx context.Context, obj interface{}) error {
	if !strings.Contains(obj.(string), "abc") {
		return errors.New("error abc")
	}
	return nil
}

type validateXyz struct{}

func (v *validateXyz) Validate(obj interface{}) error {
	if !strings.Contains(obj.(string), "xyz") {
		return errors.New("error xyz")
	}
	return nil
}

type validateContextXyz struct{}

func (v *validateContextXyz) ValidateWithContext(ctx context.Context, obj interface{}) error {
	if !strings.Contains(obj.(string), "xyz") {
		return errors.New("error xyz")
	}
	return nil
}

type validateInternalError struct{}

func (v *validateInternalError) Validate(obj interface{}) error {
	if strings.Contains(obj.(string), "internal") {
		return NewInternalError(errors.New("error internal"))
	}
	return nil
}

type validateContextInternalError struct{}

func (v *validateContextInternalError) ValidateWithContext(ctx context.Context, obj interface{}) error {
	if strings.Contains(obj.(string), "internal") {
		return NewInternalError(errors.New("error internal"))
	}
	return nil
}

type Model1 struct {
	A string
	B string
	c string
	D *string
	E String123
	F *String123
	G string `json:"g"`
}

type String123 string

func (s String123) Validate() error {
	if !strings.Contains(string(s), "123") {
		return errors.New("error 123")
	}
	return nil
}

type Model2 struct {
	Model3
	M3 Model3
	B  string
}

type Model3 struct {
	A string
}

func (m Model3) Validate() error {
	return ValidateStruct(&m,
		Field(&m.A, &validateAbc{}),
	)
}

type Model4 struct {
	A string
}

func (m Model4) ValidateWithContext(ctx context.Context) error {
	return ValidateStructWithContext(ctx, &m,
		FieldWithContext(&m.A, &validateContextAbc{}),
	)
}

type Model5 struct {
	Model4
	M4 Model4
	B  string
}
