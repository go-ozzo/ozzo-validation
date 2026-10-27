//go:build ozzo_validation_no_template
// +build ozzo_validation_no_template

package validation

import (
	"fmt"
	"strings"
)

// Error returns the error message.
//
// This build does not use text/template: it replaces the placeholders {{.key}},
// {{ .key }}, {{.key }} and {{ .key}} with the params, formatted with fmt.Sprint.
// Anything else in the message, such as pipelines, conditionals, functions or a
// placeholder whose key is not in the params, is left as written.
func (e ErrorObject) Error() string {
	if len(e.params) == 0 {
		return e.message
	}

	// One pass over the message: a param value that contains a placeholder is not
	// replaced again, so the result does not depend on the map order. Two placeholders
	// cannot match at the same position, because each one ends with "}}".
	replacements := make([]string, 0, len(e.params)*8)
	for k, v := range e.params {
		s := fmt.Sprint(v)
		replacements = append(replacements,
			"{{."+k+"}}", s,
			"{{ ."+k+" }}", s,
			"{{."+k+" }}", s,
			"{{ ."+k+"}}", s,
		)
	}

	return strings.NewReplacer(replacements...).Replace(e.message)
}
