// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

type stringValidator func(string) bool

// StringRule is a rule that checks a string variable using a specified stringValidator.
type StringRule struct {
	validate stringValidator

	translationKey    string
	defaultMessage    string
	message           string
	translationParams []interface{}
}

// NewStringRule creates a new validation rule using a function that takes a string value and returns a bool.
// The rule returned will use the function to check if a given string or byte slice is valid or not.
// An empty value is considered to be valid. Please use the Required rule to make sure a value is not empty.
func NewStringRule(validator stringValidator, translationKey string) StringRule {
	return StringRule{
		validate:       validator,
		translationKey: translationKey,
	}
}

// Error sets the error message for the rule.
func (v StringRule) Error(message string) StringRule {
	v.message = message
	return v
}

// TranslationParams sets the rule's translation parameters.
func (v StringRule) TranslationParams(params []interface{}) StringRule {
	v.translationParams = params
	return v
}

// TranslationKey sets the rule's translation key.
func (v StringRule) TranslationKey(key string) StringRule {
	v.translationKey = key
	return v
}

// DefaultMessage sets the default message for the rule.
func (v StringRule) DefaultMessage(message string) StringRule {
	v.defaultMessage = message
	return v
}

// Validate checks if the given value is valid or not.
func (v StringRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	str, err := EnsureString(value)
	if err != nil {
		return err
	}

	if v.validate(str) {
		return nil
	}

	return newErrMessage(v.translationKey, v.message).
		Default(v.defaultMessage).
		SetParams(v.translationParams)
}
