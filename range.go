// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"fmt"
	"reflect"
)

// Range returns a validation rule that checks if a value is within the given range: [min,max].
// Note that the value being checked and the boundary values must be of the same type.
func Range(min interface{}, max interface{}) *rangeRule {
	return &rangeRule{
		min:     min,
		max:     max,
		message: fmt.Sprintf("must be between %v and %v", min, max),
	}
}

type rangeRule struct {
	min     interface{}
	max     interface{}
	message string
}

// Validate checks if the given value is valid or not.
func (r *rangeRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil {
		return nil
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err := ToInt(r.min)
		if err != nil {
			return err
		}
		max, err := ToInt(r.max)
		if err != nil {
			return err
		}
		if min <= rv.Int() && rv.Int() <= max {
			return nil
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		min, err := ToUint(r.min)
		if err != nil {
			return err
		}
		max, err := ToUint(r.max)
		if err != nil {
			return err
		}
		if min <= rv.Uint() && rv.Uint() <= max {
			return nil
		}

	case reflect.Float32, reflect.Float64:
		min, err := ToFloat(r.min)
		if err != nil {
			return err
		}
		max, err := ToFloat(r.max)
		if err != nil {
			return err
		}
		if min <= rv.Float() && rv.Float() <= max {
			return nil
		}

	default:
		r.message = fmt.Sprintf("cannot apply range rule on type %v", rv.Kind())
	}
	return errors.New(r.message)
}

// Error sets the error message for the rule.
func (r *rangeRule) Error(message string) *rangeRule {
	r.message = message
	return r
}
