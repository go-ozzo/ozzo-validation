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
	mpCtx := map[string]StringValidateContext{"c": StringValidateContext("abc"), "b": StringValidateContext("123"), "a": StringValidateContext("xyz")}
	var (
		ptr     *string
		noCtx   StringValidate        = "abc"
		withCtx StringValidateContext = "xyz"
	)
	tests := []struct {
		tag            string
		value          interface{}
		err            string
		errWithContext string
	}{
		{"t1", 123, "", ""},
		{"t2", String123("123"), "", ""},
		{"t3", String123("abc"), "error 123", "error 123"},
		{"t4", []String123{}, "", ""},
		{"t4.1", []StringValidateContext{}, "", ""},
		{"t4.2", map[string]StringValidateContext{}, "", ""},
		{"t5", slice, "0: error 123; 2: error 123.", "0: error 123; 2: error 123."},
		{"t6", &slice, "0: error 123; 2: error 123.", "0: error 123; 2: error 123."},
		{"t7", ctxSlice, "", "1: (A: error abc.)."},
		{"t8", mp, "a: error 123; c: error 123.", "a: error 123; c: error 123."},
		{"t8.1", mpCtx, "a: must be abc; b: must be abc.", "a: must be abc with context; b: must be abc with context."},
		{"t9", &mp, "a: error 123; c: error 123.", "a: error 123; c: error 123."},
		{"t10", map[string]String123{}, "", ""},
		{"t11", ptr, "", ""},
		{"t12", noCtx, "called validate", "called validate"},
		{"t13", withCtx, "must be abc", "must be abc with context"},
	}
	for _, test := range tests {
		err := Validate(test.value)
		assertError(t, test.err, err, test.tag)
		// rules that are not context-aware should still be applied in context-aware validation
		err = ValidateWithContext(context.Background(), test.value)
		assertError(t, test.errWithContext, err, test.tag)
	}

	// with rules
	err := Validate("123", &validateAbc{}, &validateXyz{})
	assert.EqualError(t, err, "error abc")
	err = Validate("abc", &validateAbc{}, &validateXyz{})
	assert.EqualError(t, err, "error xyz")
	err = Validate("abcxyz", &validateAbc{}, &validateXyz{})
	assert.NoError(t, err)

	err = Validate("123", &validateAbc{}, Skip, &validateXyz{})
	assert.EqualError(t, err, "error abc")
	err = Validate("abc", &validateAbc{}, Skip, &validateXyz{})
	assert.NoError(t, err)

	err = Validate("123", &validateAbc{}, Skip.When(true), &validateXyz{})
	assert.EqualError(t, err, "error abc")
	err = Validate("abc", &validateAbc{}, Skip.When(true), &validateXyz{})
	assert.NoError(t, err)

	err = Validate("123", &validateAbc{}, Skip.When(false), &validateXyz{})
	assert.EqualError(t, err, "error abc")
	err = Validate("abc", &validateAbc{}, Skip.When(false), &validateXyz{})
	assert.EqualError(t, err, "error xyz")
}

func stringEqual(str string) RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		if s != str {
			return errors.New("unexpected string")
		}
		return nil
	}
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

	xyzRule := By(stringEqual("xyz"))
	assert.Nil(t, Validate("xyz", xyzRule))
	assert.NotNil(t, Validate("abc", xyzRule))
	assert.Nil(t, ValidateWithContext(context.Background(), "xyz", xyzRule))
	assert.NotNil(t, ValidateWithContext(context.Background(), "abc", xyzRule))
}

type key int

func TestByWithContext(t *testing.T) {
	k := key(1)
	abcRule := WithContext(func(ctx context.Context, value interface{}) error {
		if ctx.Value(k) != value.(string) {
			return errors.New("must be abc")
		}
		return nil
	})
	ctx := context.WithValue(context.Background(), k, "abc")
	assert.Nil(t, ValidateWithContext(ctx, "abc", abcRule))
	err := ValidateWithContext(ctx, "xyz", abcRule)
	if assert.NotNil(t, err) {
		assert.Equal(t, "must be abc", err.Error())
	}

	assert.NotNil(t, Validate("abc", abcRule))
}

func Test_skipRule_Validate(t *testing.T) {
	assert.Nil(t, Skip.Validate(100))
}

func assertError(t *testing.T, expected string, err error, tag string) {
	if expected == "" {
		assert.NoError(t, err, tag)
	} else {
		assert.EqualError(t, err, expected, tag)
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

func (v *validateContextAbc) Validate(obj interface{}) error {
	return v.ValidateWithContext(context.Background(), obj)
}

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

func (v *validateContextXyz) Validate(obj interface{}) error {
	return v.ValidateWithContext(context.Background(), obj)
}

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

type Model1 struct {
	A string
	B string
	c string
	D *string
	E String123
	F *String123
	G string `json:"g"`
	H []string
	I map[string]string
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
		Field(&m.A, &validateContextAbc{}),
	)
}

type Model5 struct {
	Model4
	M4 Model4
	B  string
}

type StringValidate string

func (s StringValidate) Validate() error {
	return errors.New("called validate")
}

type StringValidateContext string

func (s StringValidateContext) Validate() error {
	if string(s) != "abc" {
		return errors.New("must be abc")
	}
	return nil
}

func (s StringValidateContext) ValidateWithContext(ctx context.Context) error {
	if string(s) != "abc" {
		return errors.New("must be abc with context")
	}
	return nil
}
