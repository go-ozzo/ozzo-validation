// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"context"
	"errors"
)

// Nil is a validation rule that checks if a value is nil.
// All other types are considered valid.
var Nil = &nilRule{message: "must be empty"}

type nilRule struct {
	message string
}

// Validate checks if the given value is valid or not.
func (r *nilRule) Validate(value interface{}) error {
	if value == nil {
		return nil
	}
	return errors.New(r.message)
}

// ValidateWithContext checks if the given value is valid or not.
func (r *nilRule) ValidateWithContext(ctx context.Context, value interface{}) error {
	return r.Validate(value)
}

// Error sets the error message for the rule.
func (r *nilRule) Error(message string) *nilRule {
	return &nilRule{
		message: message,
	}
}
