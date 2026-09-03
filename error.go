// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/template"
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

// Error returns the error message.
func (e ErrorObject) Error() string {
	if len(e.params) == 0 {
		return e.message
	}

	res := bytes.Buffer{}
	_ = template.Must(template.New("err").Parse(e.message)).Execute(&res, e.params)

	return res.String()
}

// Error returns the error string of Errors.
func (es Errors) Error() string {
	if len(es) == 0 {
		return ""
	}

	keys := es.sortedKeys()

	var s strings.Builder
	for i, key := range keys {
		if es[key] == nil {
			continue
		}
		if i > 0 {
			s.WriteString("; ")
		}
		if errs, ok := es[key].(Errors); ok {
			_, _ = fmt.Fprintf(&s, "%v: (%v)", key, errs)
		} else {
			_, _ = fmt.Fprintf(&s, "%v: %v", key, es[key].Error())
		}
	}
	s.WriteString(".")
	return s.String()
}

// Is implements the interface relied upon by the standard library's errors.Is.
// It reports whether any of the errors contained in es matches target according
// to errors.Is. This makes a validation.Errors value (such as the one returned
// by Validate or ValidateStruct) usable with errors.Is against a sentinel error,
// including errors nested in the Errors values produced for maps, slices and
// nested structs.
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

// As implements the interface relied upon by the standard library's errors.As.
// It looks for the first error contained in es that matches target according to
// errors.As; if one is found, target is set to that error value and As reports
// true. Like Is, it recurses into the Errors values produced for maps, slices
// and nested structs, examines fields in sorted key order and skips nil field
// errors.
func (es Errors) As(target interface{}) bool {
	for _, key := range es.sortedKeys() {
		if err := es[key]; err != nil && errors.As(err, target) {
			return true
		}
	}
	return false
}

// sortedKeys returns the keys of es sorted in increasing order. It gives the
// Error, Is and As methods a deterministic iteration order over the fields.
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

// Assert that our ErrorObject implements the Error interface.
var _ Error = ErrorObject{}

// Assert that Errors implements the hook interfaces used by errors.Is and errors.As.
var (
	_ interface{ Is(error) bool }       = Errors{}
	_ interface{ As(interface{}) bool } = Errors{}
)
