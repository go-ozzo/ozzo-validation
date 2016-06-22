package validation_test

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Customer struct {
	FirstName string
	LastName  string
	Email     string
	SSN       string
}

func (c Customer) Validate(attrs ...string) error {
	return validation.StructRules{}.
		Add("FirstName", validation.NotEmpty, validation.Length(0, 50)).
		Add("LastName", validation.NotEmpty, validation.Length(0, 50)).
		Add("Email", validation.NotEmpty, is.Email).
		Add("SSN", validation.NotEmpty, is.SSN).
		Validate(c, attrs...)
}

func Example() {
	c := Customer{
		FirstName: "Qiang",
		Email:     "q",
		SSN:       "123456",
	}

	err := c.Validate() // alternatively, err := validation.Validate(c)
	fmt.Println(err)
	// Output:
	// Email: must be a valid email address; LastName: cannot be blank; SSN: must be a valid social security number.
}
