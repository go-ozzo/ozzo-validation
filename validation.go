// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package validation provides configurable and extensible rules for validating data of various types.
package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type (
	// Validatable is the interface indicating the type implementing it supports data validation.
	Validatable interface {
		// Validate validates the data and returns an error if validation fails.
		Validate() error
	}

	// Rule represents a validation rule.
	Rule interface {
		// Validate validates a value and returns a value if validation fails.
		Validate(value interface{}) error
	}

	// Rules represents a list of validation rules.
	Rules []Rule

	// FieldRules represents the validation rules associated with a struct field.
	FieldRules struct {
		Field string // struct field name
		Rules Rules  // rules associated with the field
	}

	// StructRules represents the validation rules for a struct.
	StructRules []FieldRules
)

var (
	// ErrorTag is the struct tag name used to customize the error field name for a struct field.
	ErrorTag = "validation"

	// Skip is a special validation rule that indicates all rules following it should be skipped.
	Skip = &skipRule{}

	validatableType = reflect.TypeOf((*Validatable)(nil)).Elem()
)

// Validate validates the given value and returns the validation error, if any.
//
// Validate only performs validation if
// - the value being validated implements `Validatable`;
// - or the value is a map, slice, or array of elements that implement `Validatable`.
//
// Nil is returned if no validation error or validation is not performed.
//
// If the value is an array, a slice, or a map, and its elements implement Validatable,
// Validate will call Validate of every element and return the validation errors
// in terms of Errors.
func Validate(value interface{}) error {
	if v, ok := value.(Validatable); ok {
		return v.Validate()
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	rt, rk := rv.Type(), rv.Kind()

	if rk == reflect.Map && rt.Elem().Implements(validatableType) {
		errs := Errors{}
		for _, key := range rv.MapKeys() {
			if mv := rv.MapIndex(key).Interface(); mv != nil {
				if err := mv.(Validatable).Validate(); err != nil {
					errs[fmt.Sprintf("%v", key.Interface())] = err
				}
			}
		}
		if len(errs) > 0 {
			return errs
		}
	} else if (rk == reflect.Slice || rk == reflect.Array) && rt.Elem().Implements(validatableType) {
		errs := Errors{}
		l := rv.Len()
		for i := 0; i < l; i++ {
			if ev := rv.Index(i).Interface(); ev != nil {
				if err := ev.(Validatable).Validate(); err != nil {
					errs[strconv.Itoa(i)] = err
				}
			}
		}
		if len(errs) > 0 {
			return errs
		}
	}

	return nil
}

// Validate validates the given value using the validation rules in Rules.
func (rules Rules) Validate(value interface{}) error {
	for _, rule := range rules {
		if _, ok := rule.(*skipRule); ok {
			return nil
		}
		if err := rule.Validate(value); err != nil {
			return err
		}
	}
	return nil
}

func (rules Rules) shouldSkip() bool {
	for _, rule := range rules {
		if _, ok := rule.(*skipRule); ok {
			return true
		}
	}
	return false
}

// Add creates a new FieldRules and adds it to StructRules.
func (r StructRules) Add(name string, rules ...Rule) StructRules {
	return append(r, FieldRules{name, rules})
}

// Validate validates a struct or a pointer to a struct.
// A list of attributes may be provided to specify which fields of the struct should be validated.
// Nil is returned if there is no validation error.
func (r StructRules) Validate(object interface{}, attrs ...string) error {
	// ensure object is a struct
	value := reflect.ValueOf(object)
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		if value.IsNil() {
			// skip nil pointer
			return nil
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return Errors{"_": errors.New("only struct or pointer to struct can be validated")}
	}

	errs := Errors{}

	for _, fieldRules := range r {
		if len(attrs) > 0 {
			found := false
			for _, attr := range attrs {
				if fieldRules.Field == attr {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if err := fieldRules.validate(value); err != nil {
			ft, _ := value.Type().FieldByName(fieldRules.Field)
			if tag := ft.Tag.Get(ErrorTag); tag != "" {
				errs[tag] = err
			} else {
				errs[fieldRules.Field] = err
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// NewFieldRules creates a new FieldRules.
func NewFieldRules(name string, rules ...Rule) FieldRules {
	return FieldRules{name, rules}
}

func (rules FieldRules) validate(object reflect.Value) error {

	fname := rules.Field

	fieldType, present := object.Type().FieldByName(fname)
	if !present {
		return fmt.Errorf("cannot find a field named %v in %v", fname, object.Type().Name())
	}

	if fieldType.PkgPath != "" {
		return fmt.Errorf("cannot validate private field %v in %v", fname, object.Type().Name())
	}

	field := object.FieldByName(fname)
	value := field.Interface()
	if err := rules.Rules.Validate(value); err != nil {
		return err
	}

	// do not dive in validation if the field is a nil pointer or it has a "Skip" rule
	if (field.Kind() == reflect.Interface || field.Kind() == reflect.Ptr) && field.IsNil() || rules.Rules.shouldSkip() {
		return nil
	}

	return Validate(value)
}

type skipRule struct{}

func (r *skipRule) Validate(interface{}) error {
	return nil
}
