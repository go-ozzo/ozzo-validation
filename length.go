// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"fmt"
)

// Length returns a validation rule that checks if a value's length is within the specified range.
// If max is 0, it means there is no upper bound for the length.
// This rule should only be used for validating strings, slices, maps, and arrays.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func Length(min, max int) *lengthRule {
	message := "the value must be empty"
	if min == 0 && max > 0 {
		message = fmt.Sprintf("the length must be no more than %v", max)
	} else if min > 0 && max == 0 {
		message = fmt.Sprintf("the length must be no less than %v", min)
	} else if min > 0 && max > 0 {
		message = fmt.Sprintf("the length must be between %v and %v", min, max)
	}
	return &lengthRule{
		min:     min,
		max:     max,
		message: message,
	}
}

type lengthRule struct {
	min, max int
	message  string
}

// Validate checks if the given value is valid or not.
func (v *lengthRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	l, err := LengthOfValue(value)
	if err != nil {
		return err
	}

	if v.min > 0 && l < v.min || v.max > 0 && l > v.max {
		return errors.New(v.message)
	}
	return nil
}

// Error sets the error message for the rule.
func (v *lengthRule) Error(message string) *lengthRule {
	v.message = message
	return v
}
