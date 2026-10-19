//go:build ozzo_validation_no_template
// +build ozzo_validation_no_template

package validation

import (
	"fmt"
	"strings"
)

// Error returns the error message.
func (e ErrorObject) Error() string {
	if len(e.params) == 0 {
		return e.message
	}

	var replacements []string

	for k, v := range e.params {
		replacements = append(replacements,
			"{{."+k+"}}",
			fmt.Sprintf("%v", v),
			"{{. "+k+" }}",
			fmt.Sprintf("%v", v),
		)
	}

	return strings.NewReplacer(replacements...).Replace(e.message)
}
