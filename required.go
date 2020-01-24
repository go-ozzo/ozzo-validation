// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

// Required is a validation rule that checks if a value is not empty.
// A value is considered not empty if
// - integer, float: not zero
// - bool: true
// - string, array, slice, map: len() > 0
// - interface, pointer: not nil and the referenced value is not empty
// - any other types
var Required = requiredRule{err: NewError("validation_required_is_blank", "cannot be blank"), skipNil: false}

// NilOrNotEmpty checks if a value is a nil pointer or a value that is not empty.
// NilOrNotEmpty differs from Required in that it treats a nil pointer as valid.
var NilOrNotEmpty = requiredRule{err: NewError("validation_nil_or_not_empty_is_blank", "cannot be blank"), skipNil: true}

type requiredRule struct {
	skipNil bool
	err     Error
}

// Validate checks if the given value is valid or not.
func (v requiredRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if v.skipNil && !isNil && IsEmpty(value) || !v.skipNil && (isNil || IsEmpty(value)) {
		return v.err
	}
	return nil
}

// Error sets the error message for the rule.
func (v requiredRule) Error(message string) requiredRule {
	v.err = v.err.SetMessage(message)

	return v
}
