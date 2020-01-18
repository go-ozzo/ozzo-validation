// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type (
	// ErrMessage return by each rule or validate method that want to return validation errors.
	ErrMessage struct {
		Lang           string
		TranslationKey string
		DefaultMessage string
		Message        string
		Params         []interface{}
	}

	// Errors represents the validation errors that are indexed by struct field names, map or slice keys.
	// values are ErrMessage or Errors (for map,slice and array error value is Errors).
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

// SetParams set the message params that using to format the message.
func (e ErrMessage) SetParams(params []interface{}) ErrMessage {
	e.Params = params
	return e
}

// CustomMessage set our custom message, so ErrMessage ignore
// the translated message or default message and just return
// custom message on "Error" method.
func (e ErrMessage) CustomMessage(message string, params []interface{}) ErrMessage {
	e.Message = message
	e.Params = params
	return e
}

// Default set default message for ErrMessage. ErrMessage use
// this message as its message if does not exists any translation
// the key for it in the translation map.
func (e ErrMessage) Default(defaultMessage string) ErrMessage {
	e.DefaultMessage = defaultMessage
	return e
}

// ToLang set error message language, ErrMessage use
// this language to detect the translated message.
func (e ErrMessage) ToLang(lang string) ErrMessage {
	e.Lang = lang
	return e
}

// Error returns the error message in the specified language.
// Priority to returning message is :
// - The custom message
// - The translated message in ErrMessage language
// - The translated message in Lang language
// - The translated message in the English language
// - The default message
// - Empty string
func (e ErrMessage) Error() string {
	return fmt.Sprintf(msgInLang(e.Lang, e.TranslationKey, e.DefaultMessage, e.Message), e.Params...)
}

// ToLang method recursively change language of all ErrMessage errors.
func (es Errors) ToLang(lang string) Errors {
	newErrors := Errors{}

	for k, e := range es {
		if errors, ok := e.(Errors); ok {
			newErrors[k] = errors.ToLang(lang)
		} else if ve, ok := e.(ErrMessage); ok {
			newErrors[k] = ve.ToLang(lang)
		} else {
			newErrors[k] = e
		}
	}

	return newErrors
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

// newErr generate new error.we use this function just for our validation errors, for users defined rules
// its better to use ErrWithDefault function, it get default message instead of custom message.
func newErrMessage(translationKey, customMessage string) ErrMessage {
	return ErrMessage{
		TranslationKey: translationKey,
		Message:        customMessage,
	}
}

// ErrMessageWithDefault generate new validation error
// with specified default message.
func ErrMessageWithDefault(translationKey, defaultMessage string, params []interface{}) ErrMessage {
	return ErrMessage{
		TranslationKey: translationKey,
		DefaultMessage: defaultMessage,
		Params:         params,
	}
}
