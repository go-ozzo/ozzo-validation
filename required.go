// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"reflect"
)

// Required is a validation rule that checks if a value is not nil.
// Required only handles types including interface, pointer, slice, and map.
// All other types are considered valid.
var Required = &requiredRule{message: "is required"}

type requiredRule struct {
	message string
}

// Validate checks if the given value is valid or not.
func (r *requiredRule) Validate(value interface{}, context interface{}) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map:
		if v.IsNil() {
			return errors.New(r.message)
		}
	case reflect.Invalid:
		return errors.New(r.message)
	}
	return nil
}

// Error sets the error message for the rule.
func (r *requiredRule) Error(message string) *requiredRule {
	return &requiredRule{
		message: message,
	}
}
