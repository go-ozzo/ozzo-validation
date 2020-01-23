// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/template"
)

type (
	// Error represents an validation error
	Error struct {
		code    string
		message string
		params  map[string]interface{}
	}

	// Errors represents the validation errors that are indexed by struct field names, map or slice keys.
	// values are Error or Errors (for map,slice and array error value is Errors).
	Errors map[string]error

	// InternalError represents an error that should NOT be treated as a validation error.
	InternalError interface {
		error
		InternalError() error
	}

	internalError struct {
		error
	}

	// Translator is the interface that need to be implemented if we
	// need to error translation feature.
	Translator interface {
		TranslateStructFieldErr(field string, err Error) (string, error)
		TranslateSingleFieldErr(err Error) (string, error)
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

// Code set the error's translation code.
func (e Error) Code(code string) Error {
	e.code = code

	return e
}

// GetCode get the error's translation code.
func (e Error) GetCode() string {
	return e.code
}

// Params set the error's params.
func (e Error) Params(params map[string]interface{}) Error {
	e.params = params
	return e
}

// GetParams returns the error's params.
func (e Error) GetParams() map[string]interface{} {
	return e.params
}

// Message set the error's message.
func (e Error) Message(message string) Error {
	e.message = message
	return e
}

// GetMessage return the error's message.
func (e Error) GetMessage() string {
	return e.message
}

// Error returns the error message.
func (e Error) Error() string {
	if len(e.params) == 0 {
		return e.message
	}

	res := bytes.Buffer{}
	_ = template.Must(template.New("err").Parse(e.message)).Execute(&res, e.params)

	return res.String()
}

// Translate get a translator that must implemented Translator and
// return translated errors.
func (e Error) Translate(t Translator) (string, error) {
	return t.TranslateSingleFieldErr(e)
}

// Translate method recursively change language of all Error errors.
func (es Errors) Translate(t Translator) (map[string]interface{}, error) {
	var errCollection = make(map[string]interface{})
	var err error

	for k, e := range es {
		if errors, ok := e.(Errors); ok {
			errCollection[k], err = errors.Translate(t)

			if err != nil {
				return nil, err
			}

		} else if ve, ok := e.(Error); ok {
			errCollection[k], err = t.TranslateStructFieldErr(k, ve)

			if err != nil {
				return nil, err
			}

		} else {
			errCollection[k] = e
		}
	}

	return errCollection, nil
}

// Error returns the error string of Errors.
func (es Errors) Error() string {
	if len(es) == 0 {
		return ""
	}

	keys := []string{}
	for key := range es {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var s strings.Builder
	for i, key := range keys {
		if i > 0 {
			s.WriteString("; ")
		}
		if errs, ok := es[key].(Errors); ok {
			fmt.Fprintf(&s, "%v: (%v)", key, errs)
		} else {
			fmt.Fprintf(&s, "%v: %v", key, es[key].Error())
		}
	}
	s.WriteString(".")
	return s.String()
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
	return Error{
		code:    code,
		message: message,
	}
}
