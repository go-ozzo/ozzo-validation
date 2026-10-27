//go:build ozzo_validation_no_template
// +build ozzo_validation_no_template

package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorObject_Error_NoTemplate(t *testing.T) {
	// Where the builds differ, as documented in the README: anything but a plain or spaced
	// field placeholder is left as written, and a pointer prints as a pointer.
	n := 5
	e := NewError("code", "{{.missing}} {{.min | printf \"%d\"}} {{if .min}}x{{end}} {{  .min  }} {{.p}}").
		SetParams(map[string]interface{}{"min": 3, "p": &n})
	assert.Equal(t,
		fmt.Sprintf("{{.missing}} {{.min | printf \"%%d\"}} {{if .min}}x{{end}} {{  .min  }} %p", &n),
		e.Error())
}
