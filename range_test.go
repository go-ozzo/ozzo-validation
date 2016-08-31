// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	var i int
	var u uint
	var f float64

	tests := []struct {
		tag      string
		min, max interface{}
		value    interface{}
		err      string
	}{
		{"t1", 2, 4, "", "cannot apply range rule on type string"},
		{"t2", 2, 4, "abc", "cannot apply range rule on type string"},
		{"t3", 2, 4, []int{1, 2}, "cannot apply range rule on type slice"},
		{"t4", 2, 4, map[string]int{"A": 1}, "cannot apply range rule on type map"},
		{"t5", 0, 2, nil, ""},

		{"t6", 2, 4, 2, ""},
		{"t7", 2, 4, 3, ""},
		{"t8", 2, 4, 4, ""},
		{"t9", 0, 2, &i, ""},
		{"t10", 1, 2, &i, "must be between 1 and 2"},
		{"t11", 2, 4, 5, "must be between 2 and 4"},
		{"t12", uint(2), 4, 3, "cannot convert uint to int64"},
		{"t13", 2, uint(4), 3, "cannot convert uint to int64"},

		{"t14", uint(0), uint(2), &u, ""},
		{"t15", uint(1), uint(2), &u, "must be between 1 and 2"},
		{"t16", uint(0), uint(1), uint(1), ""},
		{"t17", uint(0), uint(1), uint(2), "must be between 0 and 1"},
		{"t18", 0, uint(2), uint(1), "cannot convert int to uint64"},
		{"t19", uint(0), 2, uint(1), "cannot convert int to uint64"},

		{"t20", 0.0, 2.0, &f, ""},
		{"t21", 1.0, 2.0, &f, "must be between 1 and 2"},
		{"t22", 0.0, 0.1, 1.0, "must be between 0 and 0.1"},
		{"t23", 0, 1.0, 0.5, "cannot convert int to float64"},
		{"t24", 0.0, 1, 1.0, "cannot convert int to float64"},
	}

	for _, test := range tests {
		r := Range(test.min, test.max)
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestRangeError(t *testing.T) {
	r := Range(10, 20)
	assert.Equal(t, "must be between 10 and 20", r.message)

	r.Error("123")
	assert.Equal(t, "123", r.message)
}
