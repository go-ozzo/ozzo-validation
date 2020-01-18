// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"unicode/utf8"
)

// Length returns a validation rule that checks if a value's length is within the specified range.
// If max is 0, it means there is no upper bound for the length.
// This rule should only be used for validating strings, slices, maps, and arrays.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func Length(min, max int) LengthRule {
	r := LengthRule{min: min, max: max}

	r.translationKey, r.errParams = r.detectLengthTranslation()
	return r
}

// RuneLength returns a validation rule that checks if a string's rune length is within the specified range.
// If max is 0, it means there is no upper bound for the length.
// This rule should only be used for validating strings, slices, maps, and arrays.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
// If the value being validated is not a string, the rule works the same as Length.
func RuneLength(min, max int) LengthRule {
	r := Length(min, max)
	r.rune = true

	r.translationKey, r.errParams = r.detectLengthTranslation()

	return r
}

// LengthRule is a validation rule that checks if a value's length is within the specified range.
type LengthRule struct {
	message        string
	translationKey string
	errParams      []interface{}

	min, max int
	rune     bool
}

// Validate checks if the given value is valid or not.
func (v LengthRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	var (
		l   int
		err error
	)
	if s, ok := value.(string); ok && v.rune {
		l = utf8.RuneCountInString(s)
	} else if l, err = LengthOfValue(value); err != nil {
		return err
	}

	if v.min > 0 && l < v.min || v.max > 0 && l > v.max {
		return newErrMessage(v.translationKey, v.message).SetParams(v.errParams)
	}

	return nil
}

func (v LengthRule) detectLengthTranslation() (string, []interface{}) {

	if v.min == 0 && v.max > 0 {
		return "length_more_than", []interface{}{v.max}
	} else if v.min > 0 && v.max == 0 {
		return "length_no_less_than", []interface{}{v.min}
	} else if v.min > 0 && v.max > 0 {
		if v.min == v.max {
			return "length_exactly", []interface{}{v.min}
		}
		return "length_between", []interface{}{v.min, v.max}
	}

	return "length_empty", nil
}

// Error sets the error message for the rule.
func (v LengthRule) Error(message string) LengthRule {
	v.message = message
	return v
}
