// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type (
	// Error interface represents an validation error
	Error interface {
		Error() string
		Code() string
		Message() string
		SetMessage(string) Error
		Params() map[string]interface{}
		SetParams(map[string]interface{}) Error
	}

	// ErrorObject is the default validation error
	// that implements the Error interface.
	ErrorObject struct {
		code    string
		message string
		params  map[string]interface{}
	}

	// Errors represents the validation errors that are indexed by struct field names, map or slice keys.
	// values are Error or Errors (for map, slice and array error value is Errors).
	Errors map[string]error

	// InternalError represents an error that should NOT be treated as a validation error.
	InternalError interface {
		error
		InternalError() error
	}

	internalError struct {
		error
	}
)

// NewInternalError wraps a given error into an InternalError.
func NewInternalError(err error) InternalError {
	return internalError{error: err}
}

// InternalError returns the actual error that it wraps around.
func (e internalError) InternalError() error {
	return e.error
}

// SetCode set the error's translation code.
func (e ErrorObject) SetCode(code string) Error {
	e.code = code
	return e
}

// Code get the error's translation code.
func (e ErrorObject) Code() string {
	return e.code
}

// SetParams set the error's params.
func (e ErrorObject) SetParams(params map[string]interface{}) Error {
	e.params = params
	return e
}

// AddParam add parameter to the error's parameters.
func (e ErrorObject) AddParam(name string, value interface{}) Error {
	if e.params == nil {
		e.params = make(map[string]interface{})
	}

	e.params[name] = value
	return e
}

// Params returns the error's params.
func (e ErrorObject) Params() map[string]interface{} {
	return e.params
}

// SetMessage set the error's message.
func (e ErrorObject) SetMessage(message string) Error {
	e.message = message
	return e
}

// Message return the error's message.
func (e ErrorObject) Message() string {
	return e.message
}

// Is reports whether target is a validation Error with the same non-empty code. It makes
// errors.Is(err, ErrRequired) and the like work: an ErrorObject holds a params map, so it
// cannot be compared with ==, and a rule with a custom message (Required.Error("...")) keeps
// the code of the error it was built from. Params are not compared.
//
// An error with an empty code never matches, not even itself, so give the sentinel errors
// of your own rules a code: NewError("my_code", "...").
func (e ErrorObject) Is(target error) bool {
	if e.code == "" {
		return false
	}
	if p, ok := target.(*ErrorObject); ok && p == nil {
		// Code() on a nil *ErrorObject would panic.
		return false
	}
	t, ok := target.(Error)
	return ok && e.code == t.Code()
}

// Error returns the error string of Errors.
func (es Errors) Error() string {
	if len(es) == 0 {
		return ""
	}

	keys := es.sortedKeys()

	var s strings.Builder
	for _, key := range keys {
		if es[key] == nil {
			continue
		}
		// Decide the separator by what was written, not by the key index: skipped nil entries
		// must not leave a leading "; ".
		if s.Len() > 0 {
			s.WriteString("; ")
		}
		if errs, ok := es[key].(Errors); ok {
			_, _ = fmt.Fprintf(&s, "%v: (%v)", key, errs)
		} else {
			_, _ = fmt.Fprintf(&s, "%v: %v", key, es[key].Error())
		}
	}
	if s.Len() == 0 {
		// Only nil entries: report the same as an empty Errors.
		return ""
	}
	s.WriteString(".")
	return s.String()
}

// Is implements the interface relied upon by the standard library's errors.Is.
// It reports whether any entry of es matches target according to errors.Is,
// including entries of the Errors nested for maps, slices and nested structs.
// What an entry matches is up to the entry: an error returned by a By rule
// matches its sentinel, and an ErrorObject matches a validation Error with the
// same code (see ErrorObject.Is), so errors.Is(err, ErrRequired) reports whether
// any field failed Required.
//
// Fields are examined in sorted key order and nil field errors are skipped.
func (es Errors) Is(target error) bool {
	if target == nil {
		return false
	}
	for _, key := range es.sortedKeys() {
		if err := es[key]; err != nil && errors.Is(err, target) {
			return true
		}
	}
	return false
}

// Errors deliberately has no As method (and no Unwrap() []error). Code that tells
// a single-value error from a field map relies on errors.As(err, &validation.Error)
// reporting false for an Errors; an As that looks inside the map would make it
// return the first field's error and lose the field names.

// sortedKeys returns the keys of es sorted in increasing order. It gives the
// Error and Is methods a deterministic iteration order over the fields.
func (es Errors) sortedKeys() []string {
	keys := make([]string, 0, len(es))
	for key := range es {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// MarshalJSON converts the Errors into a valid JSON.
func (es Errors) MarshalJSON() ([]byte, error) {
	errs := map[string]interface{}{}
	for key, err := range es {
		if ms, ok := err.(json.Marshaler); ok {
			errs[key] = ms
		} else {
			errs[key] = err.Error()
		}
	}
	return json.Marshal(errs)
}

// Filter removes all nils from Errors and returns back the updated Errors as an error.
// If the length of Errors becomes 0, it will return nil.
func (es Errors) Filter() error {
	for key, value := range es {
		if value == nil {
			delete(es, key)
		}
	}
	if len(es) == 0 {
		return nil
	}
	return es
}

// NewError create new validation error.
func NewError(code, message string) Error {
	return ErrorObject{
		code:    code,
		message: message,
	}
}

// Assert that our ErrorObject implements the Error interface and the hook used by errors.Is.
var (
	_ Error                       = ErrorObject{}
	_ interface{ Is(error) bool } = ErrorObject{}
)

// Assert that Errors implements the hook interface used by errors.Is.
var _ interface{ Is(error) bool } = Errors{}
