// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// StructPointerError is the error that a struct being validated is not specified as a pointer.
	StructPointerError = errors.New("only a pointer to a struct can be validated")
)

type (
	// FieldPointerError is the error that a field is not specified as a pointer.
	FieldPointerError int

	// FieldNotFoundError is the error that a field cannot be found in the struct.
	FieldNotFoundError int

	FieldRules struct {
		fieldPtr interface{}
		rules    []Rule
	}
)

// Error returns the error string of FieldPointerError.
func (e FieldPointerError) Error() string {
	return fmt.Sprintf("field #%v must be specified as a pointer", int(e))
}

// Error returns the error string of FieldNotFoundError.
func (e FieldNotFoundError) Error() string {
	return fmt.Sprintf("field #%v cannot be found in the struct", int(e))
}

// ValidateStruct validates a struct by checking the specified struct fields against the corresponding validation rules.
// Note that the struct being validated must be specified as a pointer to it. If the pointer is nil, it is considered valid.
// Use Field() to specify struct fields that need to be validated. Each Field() call specifies a single field which
// should be specified as a pointer to the field. A field can be associated with multiple rules.
// For example,
//
//    value := struct {
//        Name  string
//        Value string
//    }{"name", "demo"}
//    err := validation.ValidateStruct(&value,
//         validation.Field(&a.Name, validation.Required),
//         validation.Field(&a.Value, validation.Required, validation.Length(5, 10)),
//    )
//    fmt.Println(err)
//    // Value: the length must be between 5 and 10.
//
// An error will be returned if validation fails.
func ValidateStruct(structPtr interface{}, fields ...*FieldRules) error {
	value := reflect.ValueOf(structPtr)
	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		// must be a pointer to a struct
		return StructPointerError
	}
	if value.IsNil() {
		// treat a nil struct pointer as valid
		return nil
	}
	value = value.Elem()

	errs := Errors{}

	for i, fr := range fields {
		fv := reflect.ValueOf(fr.fieldPtr)
		if fv.Kind() != reflect.Ptr {
			return FieldPointerError(i)
		}
		ft := findStructField(value, fv)
		if ft == nil {
			return FieldNotFoundError(i)
		}
		if err := Validate(fv.Elem().Interface(), fr.rules...); err != nil {
			errs[getErrorFieldName(ft)] = err
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Field specifies a struct field and the corresponding validation rules.
// The struct field must be specified as a pointer to it.
func Field(fieldPtr interface{}, rules ...Rule) *FieldRules {
	return &FieldRules{
		fieldPtr: fieldPtr,
		rules:    rules,
	}
}

// findStructField looks for a field in the given struct.
// The field being looked for should be a pointer to the actual struct field.
// If found, the field info will be returned. Otherwise, nil will be returned.
func findStructField(structValue reflect.Value, fieldValue reflect.Value) *reflect.StructField {
	ptr := fieldValue.Pointer()
	for i := structValue.NumField() - 1; i >= 0; i-- {
		if ptr == structValue.Field(i).UnsafeAddr() {
			f := structValue.Type().Field(i)
			// do additional type comparison because it's possible that the address of
			// an embedded struct is the same as the first field of the embedded struct
			if f.Type == fieldValue.Elem().Type() {
				return &f
			}
		}
	}
	return nil
}

// getErrorFieldName returns the name that should be used to represent the validation error of a struct field.
func getErrorFieldName(f *reflect.StructField) string {
	if tag := f.Tag.Get(ErrorTag); tag != "" {
		return tag
	}
	return f.Name
}
