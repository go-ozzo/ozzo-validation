//go:build !ozzo_validation_no_template
// +build !ozzo_validation_no_template

package validation

import (
	"bytes"
	"text/template"
)

// Error returns the error message.
func (e ErrorObject) Error() string {
	if len(e.params) == 0 {
		return e.message
	}

	res := bytes.Buffer{}
	_ = template.Must(template.New("err").Parse(e.message)).Execute(&res, e.params)

	return res.String()
}
