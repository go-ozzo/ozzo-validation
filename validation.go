// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package validation provides configurable and extensible rules for validating data of various types.
package validation

import (
	"context"
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

	// ValidatableWithContext is the interface indicating the type implementing it supports context aware data validation.
	ValidatableWithContext interface {
		// ValidateWithContext validates the data and returns an error if validation fails.
		ValidateWithContext(ctx context.Context) error
	}

	// Rule represents a validation rule.
	Rule interface {
		// Validate validates a value and returns a value if validation fails.
		Validate(value interface{}) error
	}

	// RuleWithContext represents a context aware validation rule.
	RuleWithContext interface {
		// ValidateWithContext validates a value and returns a value if validation fails.
		ValidateWithContext(ctx context.Context, value interface{}) error
	}

	// RuleFunc represents a validator function.
	// You may wrap it as a Rule by calling By().
	RuleFunc func(value interface{}) error

	// RuleFuncWithContext represents a context aware validator function.
	// You may wrap it as a RuleWithContext by calling ByWithContext().
	RuleFuncWithContext func(ctx context.Context, value interface{}) error
)

var (
	// ErrorTag is the struct tag name used to customize the error field name for a struct field.
	ErrorTag = "json"

	// Skip is a special validation rule that indicates all rules following it should be skipped.
	Skip = &skipRule{}

	validatableType = reflect.TypeOf((*Validatable)(nil)).Elem()

	validatableWithContextType = reflect.TypeOf((*ValidatableWithContext)(nil)).Elem()
)

// Validate validates the given value and returns the validation error, if any.
//
// Validate performs validation using the following steps:
// - validate the value against the rules passed in as parameters
// - if the value is a map and the map values implement `Validatable`, call `Validate` of every map value
// - if the value is a slice or array whose values implement `Validatable`, call `Validate` of every element
func Validate(value interface{}, rules ...Rule) error {
	for _, rule := range rules {
		if _, ok := rule.(*skipRule); ok {
			return nil
		}
		if err := rule.Validate(value); err != nil {
			return err
		}
	}

	rv := reflect.ValueOf(value)
	if (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface) && rv.IsNil() {
		return nil
	}

	if v, ok := value.(Validatable); ok {
		return v.Validate()
	}

	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Elem().Implements(validatableType) {
			return validateMap(context.Background(), rv)
		}
	case reflect.Slice, reflect.Array:
		if rv.Type().Elem().Implements(validatableType) {
			return validateSlice(context.Background(), rv)
		}
	case reflect.Ptr, reflect.Interface:
		return Validate(rv.Elem().Interface())
	}

	return nil
}

// ValidateWithContext validates the given value and returns the validation error, if any.
//
// ValidateWithContext performs validation using the following steps:
// - validate the value against the rules passed in as parameters
// - if the value is a map and the map values implement `ValidatableWithContext`, call `ValidateWithContext` of every map value
// - if the value is a map and the map values implement `Validatable`, call `Validate` of every map value
// - if the value is a slice or array whose values implement `ValidatableWithContext`, call `ValidateWithContext` of every element
// - if the value is a slice or array whose values implement `Validatable`, call `Validate` of every element
func ValidateWithContext(ctx context.Context, value interface{}, rules ...RuleWithContext) error {
	for _, rule := range rules {
		if _, ok := rule.(*skipRule); ok {
			return nil
		}
		if ctx != nil {
			if err := rule.ValidateWithContext(ctx, value); err != nil {
				return err
			}
		}
	}

	rv := reflect.ValueOf(value)
	if (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface) && rv.IsNil() {
		return nil
	}

	if v, ok := value.(ValidatableWithContext); ok {
		return v.ValidateWithContext(ctx)
	}

	if v, ok := value.(Validatable); ok {
		return v.Validate()
	}

	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Elem().Implements(validatableType) || rv.Type().Elem().Implements(validatableWithContextType) {
			return validateMap(ctx, rv)
		}
	case reflect.Slice, reflect.Array:
		if rv.Type().Elem().Implements(validatableType) || rv.Type().Elem().Implements(validatableWithContextType) {
			return validateSlice(ctx, rv)
		}
	case reflect.Ptr, reflect.Interface:
		return ValidateWithContext(ctx, rv.Elem().Interface())
	}

	return nil
}

// validateMap validates a map of validatable elements
func validateMap(ctx context.Context, rv reflect.Value) error {
	errs := Errors{}
	for _, key := range rv.MapKeys() {
		if mv := rv.MapIndex(key).Interface(); mv != nil {
			if err := validateValueWithContext(ctx, mv); err != nil {
				errs[fmt.Sprintf("%v", key.Interface())] = err
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// validateMap validates a slice/array of validatable elements
func validateSlice(ctx context.Context, rv reflect.Value) error {
	errs := Errors{}
	l := rv.Len()
	for i := 0; i < l; i++ {
		if ev := rv.Index(i).Interface(); ev != nil {
			if err := validateValueWithContext(ctx, ev); err != nil {
				errs[strconv.Itoa(i)] = err
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// validateValueWithContext performs context aware validation on a value if supported.
func validateValueWithContext(ctx context.Context, value interface{}) error {
	if vc, ok := value.(ValidatableWithContext); ok {
		if err := vc.ValidateWithContext(ctx); err != nil {
			return err
		}
	}
	if v, ok := value.(Validatable); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type skipRule struct{}

func (r *skipRule) Validate(interface{}) error {
	return nil
}

func (r *skipRule) ValidateWithContext(context.Context, interface{}) error {
	return nil
}

type inlineRule struct {
	f RuleFunc
	g RuleFuncWithContext
}

func (r *inlineRule) Validate(value interface{}) error {
	return r.f(value)
}

func (r *inlineRule) ValidateWithContext(ctx context.Context, value interface{}) error {
	return r.g(ctx, value)
}

// By wraps a RuleFunc into a Rule.
func By(f RuleFunc) Rule {
	return &inlineRule{f: f}
}

// ByWithContext wraps a RuleFuncWithContext into a RuleWithContext.
func ByWithContext(g RuleFuncWithContext) RuleWithContext {
	return &inlineRule{g: g}
}
