// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"testing"

	"strings"

	"reflect"

	"github.com/stretchr/testify/assert"
)

func TestStructRules_Add(t *testing.T) {
	se := StructRules{}
	assert.Equal(t, 0, len(se))
	se1 := se.Add("B", &validateAbc{}, &validateXyz{})
	if assert.Equal(t, 1, len(se1)) {
		assert.Equal(t, 2, len(se1[0].Rules))
	}
}

func TestNewFieldRules(t *testing.T) {
	r := NewFieldRules("name", &validateAbc{}, &validateAbc{})
	assert.Equal(t, "name", r.Field)
	assert.Equal(t, 2, len(r.Rules))
}

func TestRules_shouldSkip(t *testing.T) {
	rules := Rules{}
	assert.False(t, rules.shouldSkip())

	rules = Rules{&validateAbc{}, &validateXyz{}}
	assert.False(t, rules.shouldSkip())

	rules = Rules{Skip}
	assert.True(t, rules.shouldSkip())

	rules = Rules{&validateAbc{}, Skip, &validateXyz{}}
	assert.True(t, rules.shouldSkip())
}

func TestRules_Validate(t *testing.T) {
	rules := Rules{}
	assert.Nil(t, rules.Validate(nil, nil))
	assert.Nil(t, rules.Validate("abc", nil))

	rules = Rules{
		&validateAbc{},
		&validateXyz{},
	}
	err := rules.Validate("123", nil)
	if assert.NotNil(t, err) {
		assert.Equal(t, "error abc", err.Error())
	}
	err = rules.Validate("abc", nil)
	if assert.NotNil(t, err) {
		assert.Equal(t, "error xyz", err.Error())
	}
	assert.Nil(t, rules.Validate("abcxyz", nil))

	rules = Rules{
		&validateAbc{},
		Skip,
		&validateXyz{},
	}
	err = rules.Validate("123", nil)
	if assert.NotNil(t, err) {
		assert.Equal(t, "error abc", err.Error())
	}
	assert.Nil(t, rules.Validate("abc", nil))
}

func TestStructRules_Validate(t *testing.T) {
	tests := []struct {
		tag   string
		model interface{}
		rules StructRules
		attrs []string
		err   string
	}{
		// empty rules
		{"t1", Model1{A: "abc", B: "xyz"}, StructRules{}, []string{}, ""},
		{"t2", Model1{A: "abc", B: "xyz"}, StructRules{}.Add("A").Add("B"), []string{}, ""},
		// normal rules
		{"t3", Model1{A: "abc", B: "xyz"}, StructRules{}.Add("A", &validateAbc{}).Add("B", &validateXyz{}), []string{}, ""},
		{"t4", Model1{A: "xyz", B: "abc"}, StructRules{}.Add("A", &validateAbc{}).Add("B", &validateXyz{}), []string{}, "A: error abc; B: error xyz."},
		// model pointer
		{"t5", &Model1{A: "xyz", B: "abc"}, StructRules{}.Add("A", &validateAbc{}).Add("B", &validateXyz{}), []string{}, "A: error abc; B: error xyz."},
		// private properties
		{"t6", &Model1{A: "xyz", c: "abc"}, StructRules{}.Add("c", &validateAbc{}).Add("A", &validateAbc{}), []string{}, "A: error abc; c: cannot validate private field c in Model1."},
		// property not found
		{"t7", &Model1{A: "xyz", c: "abc"}, StructRules{}.Add("d", &validateAbc{}).Add("A", &validateAbc{}), []string{}, "A: error abc; d: cannot find a field named d in Model1."},
		// validating selected properties
		{"t8", Model1{A: "xyz", B: "abc"}, StructRules{}.Add("A", &validateAbc{}).Add("B", &validateXyz{}), []string{"B"}, "B: error xyz."},
		{"t9", Model1{A: "xyz", B: "abc"}, StructRules{}.Add("A", &validateAbc{}), []string{"B"}, ""},
		// nil pointer
		{"t10", (*Model1)(nil), StructRules{}.Add("A", &validateAbc{}), nil, ""},
		// non-struct
		{"t11", 123, StructRules{}.Add("A", &validateAbc{}), nil, "_: only struct or pointer to struct can be validated."},
		// struct tag
		{"t12", Model1{G: "xyz"}, StructRules{}.Add("G", &validateAbc{}), nil, "g: error abc."},
	}
	for _, test := range tests {
		err := test.rules.Validate(test.model, test.attrs...)
		assertError(t, test.err, err, test.tag)
	}
}

func TestFieldRules_validate(t *testing.T) {
	v := String123("xyz")
	tests := []struct {
		tag   string
		model interface{}
		rules FieldRules
		err   string
	}{
		{"t1", Model1{A: "abc"}, NewFieldRules("A"), ""},
		{"t2", Model1{A: "abc"}, NewFieldRules("ABC"), "cannot find a field named ABC in Model1"},
		{"t3", Model1{A: "abc"}, NewFieldRules("c"), "cannot validate private field c in Model1"},
		{"t4", Model1{A: "abc"}, NewFieldRules("A", &validateAbc{}), ""},
		{"t5", Model1{A: "xyz"}, NewFieldRules("A", &validateAbc{}), "error abc"},
		{"t6", Model1{E: String123("xyz")}, NewFieldRules("E"), "error 123"},
		{"t7", Model1{E: String123("xyz")}, NewFieldRules("E", Skip), ""},
		{"t8", Model1{F: &v}, NewFieldRules("F"), "error 123"},
		{"t9", Model1{}, NewFieldRules("F"), ""},
		{"t10", Model2{}, NewFieldRules("M3", Skip), ""},
		{"t11", Model2{}, NewFieldRules("M3"), "A: error abc."},
		{"t12", Model2{M3: Model3{A: "abc"}}, NewFieldRules("M3"), ""},
		{"t13", Model2{}, NewFieldRules("Model3", Skip), ""},
		{"t14", Model2{}, NewFieldRules("Model3"), "A: error abc."},
		{"t15", Model2{Model3: Model3{A: "abc"}}, NewFieldRules("Model3"), ""},
	}
	for _, test := range tests {
		err := test.rules.validate(reflect.ValueOf(test.model), test.model, nil)
		assertError(t, test.err, err, test.tag)
	}
}

func TestValidate(t *testing.T) {
	slice := []String123{String123("abc"), String123("123"), String123("xyz")}
	mp := map[string]String123{"c": String123("abc"), "b": String123("123"), "a": String123("xyz")}
	tests := []struct {
		tag   string
		value interface{}
		attrs []string
		err   string
	}{
		{"t1", 123, nil, ""},
		{"t2", String123("123"), nil, ""},
		{"t3", String123("abc"), nil, "error 123"},
		{"t4", []String123{}, nil, ""},
		{"t5", slice, nil, "0: error 123; 2: error 123."},
		{"t6", &slice, nil, "0: error 123; 2: error 123."},
		{"t7", mp, nil, "a: error 123; c: error 123."},
		{"t8", &mp, nil, "a: error 123; c: error 123."},
		{"t9", map[string]String123{}, nil, ""},
	}
	for _, test := range tests {
		err := Validate(test.value, test.attrs...)
		assertError(t, test.err, err, test.tag)
	}
}

func Test_skipRule_Validate(t *testing.T) {
	assert.Nil(t, Skip.Validate(100, nil))
}

func assertError(t *testing.T, expected string, err error, tag string) {
	if expected == "" {
		assert.Nil(t, err, tag)
	} else if assert.NotNil(t, err, tag) {
		assert.Equal(t, expected, err.Error(), tag)
	}
}

type validateAbc struct{}

func (v *validateAbc) Validate(obj interface{}, context interface{}) error {
	if !strings.Contains(obj.(string), "abc") {
		return errors.New("error abc")
	}
	return nil
}

type validateXyz struct{}

func (v *validateXyz) Validate(obj interface{}, context interface{}) error {
	if !strings.Contains(obj.(string), "xyz") {
		return errors.New("error xyz")
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
	G string `validation:"g"`
}

type String123 string

func (s String123) Validate(attrs ...string) error {
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

func (m Model3) Validate(attrs ...string) error {
	return StructRules{}.Add("A", &validateAbc{}).Validate(m, attrs...)
}
