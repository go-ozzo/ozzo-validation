// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	s1 := "123"
	s2 := ""
	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", 123, ""},
		{"t2", "", "cannot be blank"},
		{"t3", &s1, ""},
		{"t4", &s2, "cannot be blank"},
	}

	for _, test := range tests {
		r := Required
		err := r.Validate(test.value, nil)
		assertError(t, test.err, err, test.tag)
	}
}

func Test_requiredRule_Error(t *testing.T) {
	r := Required
	assert.Equal(t, "cannot be blank", r.message)
	r2 := r.Error("123")
	assert.Equal(t, "cannot be blank", r.message)
	assert.Equal(t, "123", r2.message)
}
