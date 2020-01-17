// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// Length returns a validation rule that checks if a value's length is within the specified range.
// If max is 0, it means there is no upper bound for the length.
// This rule should only be used for validating strings, slices, maps, and arrays.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func Length(min, max int) LengthRule {
	return LengthRule{
		min:     min,
		max:     max,
		message: "",
	}
}

// RuneLength returns a validation rule that checks if a string's rune length is within the specified range.
// If max is 0, it means there is no upper bound for the length.
// This rule should only be used for validating strings, slices, maps, and arrays.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
// If the value being validated is not a string, the rule works the same as Length.
func RuneLength(min, max int) LengthRule {
	r := Length(min, max)
	r.rune = true
	return r
}

// LengthRule is a validation rule that checks if a value's length is within the specified range.
type LengthRule struct {
	min, max int
	message  string
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
		return errors.New(v.detectLengthErrMsg())
	}
	return nil
}

func (v LengthRule) detectLengthErrMsg() string {

	if v.min == 0 && v.max > 0 {
		return fmt.Sprintf(Msg("length_more_than", v.message), v.max)
	} else if v.min > 0 && v.max == 0 {
		return fmt.Sprintf(Msg("length_no_less_than", v.message), v.min)
	} else if v.min > 0 && v.max > 0 {
		if v.min == v.max {
			return fmt.Sprintf(Msg("length_exactly", v.message), v.min)
		} else {
			return fmt.Sprintf(Msg("length_between", v.message), v.min, v.max)
		}
	}

	return Msg("length_empty", v.message)
}

// Error sets the error message for the rule.
func (v LengthRule) Error(message string) LengthRule {
	v.message = message
	return v
}
