// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"sort"
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
	for key, _ := range es {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	s := ""
	for i, key := range keys {
		if i > 0 {
			s += "; "
		}
		s += fmt.Sprintf("%v: %v", key, es[key].Error())
	}
	return s + "."
}

// Error returns the error string of Errors.
func (es Errors) Error() string {
	if len(es) == 0 {
		return ""
	}

	keys := []string{}
	for key, _ := range es {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	s := ""
	for i, key := range keys {
		if i > 0 {
			s += "; "
		}
		s += fmt.Sprintf("%v: %v", key, es[key].Error())
	}
	return s + "."
}
