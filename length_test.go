// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	var v *string
	tests := []struct {
		tag      string
		min, max int
		value    interface{}
		err      string
	}{
		{"t1", 2, 4, "abc", ""},
		{"t2", 2, 4, "", ""},
		{"t3", 2, 4, "abcdf", "the length must be between 2 and 4"},
		{"t4", 0, 4, "ab", ""},
		{"t5", 0, 4, "abcde", "the length must be no more than 4"},
		{"t6", 2, 0, "ab", ""},
		{"t7", 2, 0, "a", "the length must be no less than 2"},
		{"t8", 2, 0, v, ""},
		{"t9", 2, 0, 123, "cannot get the length of int"},
	}

	for _, test := range tests {
		r := Length(test.min, test.max)
		err := r.Validate(test.value, nil)
		assertError(t, test.err, err, test.tag)
	}
}

func Test_lengthRule_Error(t *testing.T) {
	r := Length(10, 20)
	assert.Equal(t, "the length must be between 10 and 20", r.message)

	r = Length(0, 20)
	assert.Equal(t, "the length must be no more than 20", r.message)

	r = Length(10, 0)
	assert.Equal(t, "the length must be no less than 10", r.message)

	r.Error("123")
	assert.Equal(t, "123", r.message)
}
