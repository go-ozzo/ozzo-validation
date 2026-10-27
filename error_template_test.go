//go:build !ozzo_validation_no_template
// +build !ozzo_validation_no_template

package validation

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestErrorObject_Error_Template(t *testing.T) {
	// "{{. min }}" is not a placeholder: text/template rejects it ("function "min" not
	// defined"), so no message can use it and the no-template build does not support it.
	for _, message := range []string{"{{. min }}", "{{. threshold }}"} {
		_, err := template.New("err").Parse(message)
		assert.Error(t, err, message)
	}

	// Where the builds differ, as documented in the README: the template build prints
	// "<no value>" for a missing key and dereferences pointers.
	n := 5
	e := NewError("code", "{{.missing}} {{.p}}").SetParams(map[string]interface{}{"p": &n})
	assert.Equal(t, "<no value> 5", e.Error())
}
