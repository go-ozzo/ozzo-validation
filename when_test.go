// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"
)

func abcValidation(val string) bool {
	return val == "abc"
}

func TestWhen(t *testing.T) {
	abcRule := NewStringRule(abcValidation, "wrong_abc")
	validateMeRule := NewStringRule(validateMe, "wrong_me")

	tests := []struct {
		tag       string
		condition bool
		value     interface{}
		rules     []Rule
		err       string
	}{
		// True condition
		{"t1.1", true, nil, []Rule{}, ""},
		{"t1.2", true, "", []Rule{}, ""},
		{"t1.3", true, "", []Rule{abcRule}, ""},
		{"t1.4", true, 12, []Rule{Required}, ""},
		{"t1.5", true, nil, []Rule{Required}, "cannot be blank"},
		{"t1.6", true, "123", []Rule{abcRule}, "wrong_abc"},
		{"t1.7", true, "abc", []Rule{abcRule}, ""},
		{"t1.8", true, "abc", []Rule{abcRule, abcRule}, ""},
		{"t1.9", true, "abc", []Rule{abcRule, validateMeRule}, "wrong_me"},
		{"t1.10", true, "me", []Rule{abcRule, validateMeRule}, "wrong_abc"},

		// False condition
		{"t2.1", false, "", []Rule{}, ""},
		{"t2.2", false, "", []Rule{abcRule}, ""},
		{"t2.3", false, "abc", []Rule{abcRule}, ""},
		{"t2.4", false, "abc", []Rule{abcRule, abcRule}, ""},
		{"t2.5", false, "abc", []Rule{abcRule, validateMeRule}, ""},
		{"t2.6", false, "me", []Rule{abcRule, validateMeRule}, ""},
		{"t2.7", false, "", []Rule{abcRule, validateMeRule}, ""},
	}

	for _, test := range tests {
		err := Validate(test.value, When(test.condition, test.rules...))
		assertError(t, test.err, err, test.tag)
	}
}
