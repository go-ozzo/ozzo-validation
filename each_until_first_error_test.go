package validation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEachUntilFirstError(t *testing.T) {
	var a *int
	var f = func(v string) string { return v }
	var c0 chan int
	c1 := make(chan int)

	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", nil, "must be an iterable (map, slice or array)"},
		{"t2", map[string]string{}, ""},
		{"t3", map[string]string{"key1": "value1", "key2": "value2"}, ""},
		{"t5", map[string]map[string]string{"key1": {"key1.1": "value1"}, "key2": {"key2.1": "value1"}}, ""},
		{"t6", map[string]map[string]string{"": nil}, ": cannot be blank."},
		{"t7", map[interface{}]interface{}{}, ""},
		{"t8", map[interface{}]interface{}{"key1": struct{ foo string }{"foo"}}, ""},
		{"t10", []string{"value1", "value2", "value3"}, ""},
		{"t11", []string{"", "value2", ""}, "0: cannot be blank."},
		{"t12", []interface{}{struct{ foo string }{"foo"}}, ""},
		{"t13", []interface{}{nil, a}, "0: cannot be blank."},
		{"t14", []interface{}{c0, c1, f}, "0: cannot be blank."},
	}

	for _, test := range tests {
		r := EachUntilFirstError(Required)
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestEachUntilFirstError_StopsAtFirst(t *testing.T) {
	// For maps with multiple errors, verify exactly one error is returned
	// (map iteration order is non-deterministic, so we can't predict which key)
	r := EachUntilFirstError(Required)
	err := r.Validate(map[string]string{"a": "", "b": "", "c": ""})
	if assert.NotNil(t, err) {
		errs := err.(Errors)
		assert.Equal(t, 1, len(errs), "should stop at first error")
	}
}

func TestEachUntilFirstErrorWithContext(t *testing.T) {
	ctx := context.Background()

	// Slice: stops at first error
	r := EachUntilFirstError(Required)
	err := r.ValidateWithContext(ctx, []string{"", "value", ""})
	if assert.NotNil(t, err) {
		errs := err.(Errors)
		assert.Equal(t, 1, len(errs), "should stop at first error")
		assert.NotNil(t, errs["0"], "first element should fail")
	}

	// Valid slice
	err = r.ValidateWithContext(ctx, []string{"a", "b"})
	assert.Nil(t, err)

	// Empty map
	err = r.ValidateWithContext(ctx, map[string]string{})
	assert.Nil(t, err)
}
