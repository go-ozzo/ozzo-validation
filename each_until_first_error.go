// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"reflect"
	"strconv"
)

// EachUntilFirstError is the same as Each but stops early once the first item with a validation error was encountered.
// Use this instead of Each for array's or maps that may potentially contain ten-thousands of erroneous items and
// you want to avoid returning ten-thousands of validation errors (for memory and cpu reasons).
func EachUntilFirstError(rules ...Rule) EachUntilFirstErrorRule {
	return EachUntilFirstErrorRule{
		rules: rules,
	}
}

// EachUntilFirstErrorRule is the same as EachRule but stops early once the first item with a validation error was encountered.
// Use this instead of EachRule for array's or maps that may potentially contain ten-thousands of erroneous items and
// you want to avoid returning ten-thousands of validation errors (for memory and cpu reasons).
type EachUntilFirstErrorRule struct {
	rules []Rule
}

// Validate loops through the given iterable and calls the Ozzo Validate() method for each value.
func (r EachUntilFirstErrorRule) Validate(value interface{}) error {
	errs := Errors{}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			val := r.getInterface(v.MapIndex(k))
			if err := Validate(val, r.rules...); err != nil {
				errs[r.getString(k)] = err
				break
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			val := r.getInterface(v.Index(i))
			if err := Validate(val, r.rules...); err != nil {
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

func (r EachUntilFirstErrorRule) getInterface(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return nil
		}
		return value.Elem().Interface()
	default:
		return value.Interface()
	}
}

func (r EachUntilFirstErrorRule) getString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return ""
		}
		return value.Elem().String()
	default:
		return value.String()
	}
}
