package validation

import (
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
