// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

type (
	// Errors represents the validation errors for a struct object. The keys are the struct field names.
	Errors map[string]error
	// SliceErrors represents the validation errors for a slice. The keys are the indices of the slice elements having errors.
	SliceErrors map[int]error
)

// Error returns the error string of SliceErrors.
func (es SliceErrors) Error() string {
	if len(es) == 0 {
		return ""
	}

	keys := []int{}
	for key := range es {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	s := ""
	for i, key := range keys {
		if i > 0 {
			s += "; "
		}
		s += formatError(key, es[key])
	}
	return s + "."
}

// MarshalJSON converts SliceErrors
// into a valid JSON.
func (es SliceErrors) MarshalJSON() ([]byte, error) {
	errs := map[string]string{}
	for key := range es {
		errs[strconv.Itoa(key)] = es[key].Error()
	}
	return json.Marshal(errs)
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

	s := ""
	for i, key := range keys {
		if i > 0 {
			s += "; "
		}
		s += formatError(key, es[key])
	}
	return s + "."
}

// MarshalJSON converts the Errors
// into a valid JSON.
func (es Errors) MarshalJSON() ([]byte, error) {
	errs := map[string]string{}
	for key := range es {
		errs[key] = es[key].Error()
	}
	return json.Marshal(errs)
}

func formatError(key interface{}, err error) string {
	switch err.(type) {
	case SliceErrors, Errors:
		return fmt.Sprintf("%v: (%v)", key, err.Error())
	default:
		return fmt.Sprintf("%v: %v", key, err.Error())
	}
}
