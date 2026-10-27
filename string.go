// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"math"
	"reflect"
)

type stringValidator func(string) bool

// StringRule is a rule that checks a string variable using a specified stringValidator.
type StringRule struct {
	validate stringValidator
	err      Error
	number   func(float64) bool
}

// NewStringRule creates a new validation rule using a function that takes a string value and returns a bool.
// The rule returned will use the function to check if a given string or byte slice is valid or not.
// An empty value is considered to be valid. Please use the Required rule to make sure a value is not empty.
func NewStringRule(validator stringValidator, message string) StringRule {
	return StringRule{
		validate: validator,
		err:      NewError("", message),
	}
}

// NewStringRuleWithError creates a new validation rule using a function that takes a string value and returns a bool.
// The rule returned will use the function to check if a given string or byte slice is valid or not.
// An empty value is considered to be valid. Please use the Required rule to make sure a value is not empty.
func NewStringRuleWithError(validator stringValidator, err Error) StringRule {
	return StringRule{
		validate: validator,
		err:      err,
	}
}

// Error sets the error message for the rule.
func (r StringRule) Error(message string) StringRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r StringRule) ErrorObject(err Error) StringRule {
	r.err = err
	return r
}

// AllowNumber makes the rule accept numbers (int, uint and float types) in addition to strings and byte slices.
// A number is valid if it is finite and the given function returns true for it; NaN and infinities are invalid.
// Strings and byte slices are still checked by the rule's string validator.
func (r StringRule) AllowNumber(valid func(float64) bool) StringRule {
	r.number = valid
	return r
}

// Validate checks if the given value is valid or not.
func (r StringRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	if r.number != nil {
		if f, ok := numberToFloat(value); ok {
			if !math.IsNaN(f) && !math.IsInf(f, 0) && r.number(f) {
				return nil
			}
			return r.err
		}
	}

	str, err := EnsureString(value)
	if err != nil {
		return err
	}

	if r.validate(str) {
		return nil
	}

	return r.err
}

// numberToFloat converts a value of an int, uint or float kind to float64.
// It reports false for any other kind, strings included.
func numberToFloat(value interface{}) (float64, bool) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	}
	return 0, false
}
