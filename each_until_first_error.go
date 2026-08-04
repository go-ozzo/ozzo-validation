// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"context"
	"errors"
	"reflect"
	"strconv"
)

// EachUntilFirstError returns a validation rule that loops through the given iterable (map, slice, or array)
// and validates each item with the given rules. It stops at the first item that has a validation error.
// Use this instead of Each for collections that may contain many erroneous items and you want to avoid
// generating a large number of validation errors.
func EachUntilFirstError(rules ...Rule) EachUntilFirstErrorRule {
	return EachUntilFirstErrorRule{
		rules: rules,
	}
}

// EachUntilFirstErrorRule is a validation rule that validates each item in an iterable
// and stops at the first error. See EachUntilFirstError().
type EachUntilFirstErrorRule struct {
	rules []Rule
}

// Validate loops through the given iterable and validates each value, stopping at the first error.
func (r EachUntilFirstErrorRule) Validate(value interface{}) error {
	return r.ValidateWithContext(nil, value)
}

// ValidateWithContext loops through the given iterable and validates each value with context,
// stopping at the first error.
func (r EachUntilFirstErrorRule) ValidateWithContext(ctx context.Context, value interface{}) error {
	errs := Errors{}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			val := getIterableInterface(v.MapIndex(k))
			var err error
			if ctx == nil {
				err = Validate(val, r.rules...)
			} else {
				err = ValidateWithContext(ctx, val, r.rules...)
			}
			if err != nil {
				errs[getIterableString(k)] = err
				break
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			val := getIterableInterface(v.Index(i))
			var err error
			if ctx == nil {
				err = Validate(val, r.rules...)
			} else {
				err = ValidateWithContext(ctx, val, r.rules...)
			}
			if err != nil {
				errs[strconv.Itoa(i)] = err
				break
			}
		}
	default:
		return errors.New("must be an iterable (map, slice or array)")
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
