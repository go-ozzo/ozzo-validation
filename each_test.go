package validation

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestEach(t *testing.T) {
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
		{"t4", map[string]string{"key1": "", "key2": "value2", "key3": ""}, "key1: cannot be blank; key3: cannot be blank."},
		{"t5", map[string]map[string]string{"key1": {"key1.1": "value1"}, "key2": {"key2.1": "value1"}}, ""},
		{"t6", map[string]map[string]string{"": nil}, ": cannot be blank."},
		{"t7", map[interface{}]interface{}{}, ""},
		{"t8", map[interface{}]interface{}{"key1": struct{ foo string }{"foo"}}, ""},
		{"t9", map[interface{}]interface{}{nil: "", "": "", "key1": nil}, ": cannot be blank; key1: cannot be blank."},
		{"t10", []string{"value1", "value2", "value3"}, ""},
		{"t11", []string{"", "value2", ""}, "0: cannot be blank; 2: cannot be blank."},
		{"t12", []interface{}{struct{ foo string }{"foo"}}, ""},
		{"t13", []interface{}{nil, a}, "0: cannot be blank; 1: cannot be blank."},
		{"t14", []interface{}{c0, c1, f}, "0: cannot be blank."},
	}

	for _, test := range tests {
		r := Each(Required)
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestEachWithContext(t *testing.T) {
	rule := Each(WithContext(func(ctx context.Context, value interface{}) error {
		if !strings.Contains(value.(string), ctx.Value(contains).(string)) {
			return errors.New("unexpected value")
		}
		return nil
	}))
	ctx1 := context.WithValue(context.Background(), contains, "abc")
	ctx2 := context.WithValue(context.Background(), contains, "xyz")

	tests := []struct {
		tag   string
		value interface{}
		ctx   context.Context
		err   string
	}{
		{"t1.1", map[string]string{"key": "abc"}, ctx1, ""},
		{"t1.2", map[string]string{"key": "abc"}, ctx2, "key: unexpected value."},
		{"t1.3", map[string]string{"key": "xyz"}, ctx1, "key: unexpected value."},
		{"t1.4", map[string]string{"key": "xyz"}, ctx2, ""},
		{"t1.5", []string{"abc"}, ctx1, ""},
		{"t1.6", []string{"abc"}, ctx2, "0: unexpected value."},
		{"t1.7", []string{"xyz"}, ctx1, "0: unexpected value."},
		{"t1.8", []string{"xyz"}, ctx2, ""},
	}

	for _, test := range tests {
		err := ValidateWithContext(test.ctx, test.value, rule)
		assertError(t, test.err, err, test.tag)
	}
}

func TestEachAndBy(t *testing.T) {
	var byAddr bool
	var s string
	Each(By(func(v interface{}) error {
		_, byAddr = v.(*string)
		return nil
	})).Validate([]*string{&s})

	if !byAddr {
		t.Fatal("slice of pointers does not get passed to `By` function by ref")
	}
}
