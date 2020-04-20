// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNil(t *testing.T) {
	s1 := "123"
	s2 := ""
	var time1 time.Time
	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", 123, "must be blank"},
		{"t2", "", "must be blank"},
		{"t3", &s1, "must be blank"},
		{"t4", &s2, "must be blank"},
		{"t5", nil, ""},
		{"t6", time1, "must be blank"},
	}

	for _, test := range tests {
		r := Nil
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestEmpty(t *testing.T) {
	s1 := "123"
	s2 := ""
	time1 := time.Now()
	var time2 time.Time
	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", 123, "must be blank"},
		{"t2", "", ""},
		{"t3", &s1, "must be blank"},
		{"t4", &s2, ""},
		{"t5", nil, ""},
		{"t6", time1, "must be blank"},
		{"t7", time2, ""},
	}

	for _, test := range tests {
		r := Empty
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestAbsentRule_When(t *testing.T) {
	r := Nil.When(false)
	err := Validate(42, r)
	assert.Nil(t, err)

	r = Nil.When(true)
	err = Validate(42, r)
	assert.Equal(t, ErrNil, err)
}

func Test_absentRule_Error(t *testing.T) {
	r := Nil
	assert.Equal(t, "must be blank", r.Validate("42").Error())
	assert.False(t, r.skipNil)
	r2 := r.Error("123")
	assert.Equal(t, "must be blank", r.Validate("42").Error())
	assert.False(t, r.skipNil)
	assert.Equal(t, "123", r2.err.Message())
	assert.False(t, r2.skipNil)

	r = Empty
	assert.Equal(t, "must be blank", r.Validate("42").Error())
	assert.True(t, r.skipNil)
	r2 = r.Error("123")
	assert.Equal(t, "must be blank", r.Validate("42").Error())
	assert.True(t, r.skipNil)
	assert.Equal(t, "123", r2.err.Message())
	assert.True(t, r2.skipNil)
}

func TestAbsentRule_Error(t *testing.T) {
	r := Nil

	err := NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.err)
	assert.Equal(t, err.Code(), r.err.Code())
	assert.Equal(t, err.Message(), r.err.Message())
	assert.NotEqual(t, err, Nil.err)
}
