// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Struct1 struct {
	Field1 int
	Field2 *int
	Field3 []int
	Field4 [4]int
	field5 int
	Struct2
	S1 *Struct2
	S2 Struct2
}

type Struct2 struct {
	Field21 string
}

func TestFindStructField(t *testing.T) {
	var s1 Struct1
	v1 := reflect.ValueOf(&s1).Elem()
	assert.Nil(t, findStructField(v1, getFieldAddress(s1.Field1)))
	assert.NotNil(t, findStructField(v1, getFieldAddress(&s1.Field1)))
	assert.Nil(t, findStructField(v1, getFieldAddress(s1.Field2)))
	assert.NotNil(t, findStructField(v1, getFieldAddress(&s1.Field2)))
	assert.Nil(t, findStructField(v1, getFieldAddress(s1.Field3)))
	assert.NotNil(t, findStructField(v1, getFieldAddress(&s1.Field3)))
	assert.Nil(t, findStructField(v1, getFieldAddress(s1.Field4)))
	assert.NotNil(t, findStructField(v1, getFieldAddress(&s1.Field4)))
	assert.Nil(t, findStructField(v1, getFieldAddress(s1.field5)))
	assert.NotNil(t, findStructField(v1, getFieldAddress(&s1.field5)))
	assert.Nil(t, findStructField(v1, getFieldAddress(s1.S1)))
	assert.NotNil(t, findStructField(v1, getFieldAddress(&s1.S1)))

	assert.Nil(t, findStructField(v1, getFieldAddress(s1.Field21)))
	assert.NotNil(t, findStructField(v1, getFieldAddress(&s1.Field21)), "field of anonymous struct")
	s2 := reflect.ValueOf(&s1.Struct2).Elem()
	assert.NotNil(t, findStructField(s2, getFieldAddress(&s1.Field21)))
	assert.NotNil(t, findStructField(s2, getFieldAddress(&s1.Struct2.Field21)))
}

func getFieldAddress(f interface{}) uintptr {
	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Ptr {
		return v.Pointer()
	}
	return 0
}

func TestValidateStruct(t *testing.T) {
	var m0 *Model1
	m1 := Model1{A: "abc", B: "xyz", c: "abc", G: "xyz"}
	m2 := Model1{E: String123("xyz")}
	tests := []struct {
		tag   string
		model interface{}
		rules []*FieldRules2
		err   string
	}{
		// empty rules
		{"t1.1", &m1, []*FieldRules2{}, ""},
		{"t1.2", &m1, []*FieldRules2{Field(&m1.A), Field(&m1.B)}, ""},
		// normal rules
		{"t2.1", &m1, []*FieldRules2{Field(&m1.A, &validateAbc{}), Field(&m1.B, &validateXyz{})}, ""},
		{"t2.2", &m1, []*FieldRules2{Field(&m1.A, &validateXyz{}), Field(&m1.B, &validateAbc{})}, "A: error xyz; B: error abc."},
		{"t2.3", &m1, []*FieldRules2{Field(&m1.A, &validateXyz{}), Field(&m1.c, &validateXyz{})}, "A: error xyz; c: error xyz."},
		// non-struct pointer
		{"t3.1", m1, []*FieldRules2{}, StructPointerError.Error()},
		{"t3.2", nil, []*FieldRules2{}, StructPointerError.Error()},
		{"t3.3", m0, []*FieldRules2{}, ""},
		{"t3.4", &m0, []*FieldRules2{}, StructPointerError.Error()},
		// invalid field spec
		{"t4.1", &m1, []*FieldRules2{Field(m1)}, FieldPointerError(0).Error()},
		{"t4.2", &m1, []*FieldRules2{Field(&m1)}, FieldNotFoundError(0).Error()},
		// struct tag
		{"t5.1", &m1, []*FieldRules2{Field(&m1.G, &validateAbc{})}, "g: error abc."},
		// validatable field
		{"t6.1", &m2, []*FieldRules2{Field(&m2.E)}, "E: error 123."},
		{"t6.2", &m2, []*FieldRules2{Field(&m2.E, Skip)}, ""},
		// Required, NotNil
		{"t7.1", &m2, []*FieldRules2{Field(&m2.F, Required)}, "F: cannot be blank."},
		{"t7.2", &m2, []*FieldRules2{Field(&m2.F, NotNil)}, "F: is required."},
		{"t7.3", &m2, []*FieldRules2{Field(&m2.E, Required, Skip)}, ""},
		{"t7.4", &m2, []*FieldRules2{Field(&m2.E, NotNil, Skip)}, ""},
	}
	for _, test := range tests {
		err := ValidateStruct(test.model, test.rules...)
		assertError(t, test.err, err, test.tag)
	}

	a := struct {
		Name  string
		Value string
	}{"name", "demo"}
	err := ValidateStruct(&a,
		Field(&a.Name, Required),
		Field(&a.Value, Required, Length(5, 10)),
	)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Value: the length must be between 5 and 10.", err.Error())
	}
}
