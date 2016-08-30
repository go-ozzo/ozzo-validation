// Copyright (C) 2016 Etix Labs - All Rights Reserved.
// All information contained herein is, and remains the property of Etix Labs and its suppliers,
// if any. The intellectual and technical concepts contained herein are proprietary to Etix Labs
// Dissemination of this information or reproduction of this material is strictly forbidden unless
// prior written permission is obtained from Etix Labs.

package validation

import "testing"

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

		{"t5", 2, 4, 2, ""},
		{"t6", 2, 4, 3, ""},
		{"t7", 2, 4, 4, ""},
		{"t8", 0, 2, &i, ""},
		{"t9", 1, 2, &i, "must be between 1 and 2"},
		{"t10", 2, 4, 5, "must be between 2 and 4"},
		{"t11", uint(2), 4, 3, "cannot convert uint to int64"},
		{"t12", 2, uint(4), 3, "cannot convert uint to int64"},

		{"t13", uint(0), uint(2), &u, ""},
		{"t14", uint(1), uint(2), &u, "must be between 1 and 2"},
		{"t15", uint(0), uint(1), uint(1), ""},
		{"t16", uint(0), uint(1), uint(2), "must be between 0 and 1"},
		{"t17", 0, uint(2), uint(1), "cannot convert int to uint64"},
		{"t18", uint(0), 2, uint(1), "cannot convert int to uint64"},

		{"t19", 0.0, 2.0, &f, ""},
		{"t20", 1.0, 2.0, &f, "must be between 1 and 2"},
		{"t21", 0.0, 0.1, 1.0, "must be between 0 and 0.1"},
		{"t22", 0, 1.0, 0.5, "cannot convert int to float64"},
		{"t23", 0.0, 1, 1.0, "cannot convert int to float64"},
	}

	for _, test := range tests {
		r := Range(test.min, test.max)
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}
