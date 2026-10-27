package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// These tests run in both builds, with and without the ozzo_validation_no_template tag,
// so they prove that both render the same text.

func TestErrorObject_Error_BuiltinRules(t *testing.T) {
	threshold := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		tag      string
		rule     Rule
		value    interface{}
		expected string
	}{
		{"length too long", Length(0, 3), "abcd", "the length must be no more than 3"},
		{"length too short", Length(3, 0), "ab", "the length must be no less than 3"},
		{"length invalid", Length(3, 3), "ab", "the length must be exactly 3"},
		{"length out of range", Length(3, 5), "ab", "the length must be between 3 and 5"},
		{"length empty required", Length(0, 0), "ab", "the value must be empty"},
		{"rune length", RuneLength(3, 5), "äö", "the length must be between 3 and 5"},
		{"min int", Min(10), 5, "must be no less than 10"},
		{"min exclusive", Min(10).Exclusive(), 10, "must be greater than 10"},
		{"max float", Max(1.5), 2.5, "must be no greater than 1.5"},
		{"max exclusive", Max(1.5).Exclusive(), 1.5, "must be less than 1.5"},
		{"max large float", Max(1e21), 1e22, "must be no greater than 1e+21"},
		{"max uint", Max(uint(3)), uint(4), "must be no greater than 3"},
		{"max time", Max(threshold), threshold.Add(time.Hour), "must be no greater than 2020-01-02 03:04:05 +0000 UTC"},
		{"multiple of", MultipleOf(3), 10, "must be multiple of 3"},
	}
	for _, test := range tests {
		err := Validate(test.value, test.rule)
		if assert.Error(t, err, test.tag) {
			assert.Equal(t, test.expected, err.Error(), test.tag)
		}
	}
}

func TestErrorObject_Error_Placeholders(t *testing.T) {
	tests := []struct {
		tag      string
		message  string
		params   map[string]interface{}
		expected string
	}{
		{"nil params", "at least {{.min}}", nil, "at least {{.min}}"},
		{"empty params", "at least {{.min}}", map[string]interface{}{}, "at least {{.min}}"},
		{"no placeholder", "plain", map[string]interface{}{"min": 1}, "plain"},
		{"spaced forms", "{{.min}}|{{ .min }}|{{.min }}|{{ .min}}", map[string]interface{}{"min": 3}, "3|3|3|3"},
		{"two params", "{{.min}}-{{.max}}", map[string]interface{}{"min": 1, "max": 2}, "1-2"},
		{"repeated", "{{.min}} {{.min}}", map[string]interface{}{"min": 7}, "7 7"},
		{"prefix keys", "{{.a}} {{.ab}}", map[string]interface{}{"a": 1, "ab": 2}, "1 2"},
		{"string", "got {{.v}}", map[string]interface{}{"v": "x y"}, "got x y"},
		{"float", "got {{.v}}", map[string]interface{}{"v": 0.25}, "got 0.25"},
		{"time", "after {{.v}}", map[string]interface{}{"v": time.Date(2021, 5, 6, 7, 8, 9, 0, time.UTC)},
			"after 2021-05-06 07:08:09 +0000 UTC"},
	}
	for _, test := range tests {
		err := NewError("code", test.message).SetParams(test.params)
		assert.Equal(t, test.expected, err.Error(), test.tag)
	}
}

func TestErrorObject_Error_ParamValueIsNotExpanded(t *testing.T) {
	// A param value that looks like a placeholder is printed as is. The params map is
	// iterated in random order, so repeat to catch an order-dependent result.
	err := NewError("code", "{{.a}} {{.b}}").SetParams(map[string]interface{}{"a": "{{.b}}", "b": "x"})
	for i := 0; i < 100; i++ {
		if !assert.Equal(t, "{{.b}} x", err.Error()) {
			return
		}
	}
}
