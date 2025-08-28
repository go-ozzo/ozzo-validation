package ozzo_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/daetabased/ozzo"
	"github.com/daetabased/ozzo/is"
)

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type Customer struct {
	Name    string
	Gender  string
	Email   string
	Address Address
}

func (a Address) Validate() error {
	return ozzo.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		ozzo.Field(&a.Street, ozzo.Required, ozzo.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		ozzo.Field(&a.City, ozzo.Required, ozzo.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		ozzo.Field(&a.State, ozzo.Required, ozzo.Match(regexp.MustCompile("^[A-Z]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		ozzo.Field(&a.Zip, ozzo.Required, ozzo.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}

func (c Customer) Validate() error {
	return ozzo.ValidateStruct(&c,
		// Name cannot be empty, and the length must be between 5 and 20.
		ozzo.Field(&c.Name, ozzo.Required, ozzo.Length(5, 20)),
		// Gender is optional, and should be either "Female" or "Male".
		ozzo.Field(&c.Gender, ozzo.In("Female", "Male")),
		// Email cannot be empty and should be in a valid email format.
		ozzo.Field(&c.Email, ozzo.Required, is.Email),
		// Validate Address using its own validation rules
		ozzo.Field(&c.Address),
	)
}

func Example() {
	c := Customer{
		Name:  "Qiang Xue",
		Email: "q",
		Address: Address{
			Street: "123 Main Street",
			City:   "Unknown",
			State:  "Virginia",
			Zip:    "12345",
		},
	}

	err := c.Validate()
	fmt.Println(err)
	// Output:
	// Address: (State: must be in a valid format.); Email: must be a valid email address.
}

func Example_second() {
	data := "example"
	err := ozzo.Validate(data,
		ozzo.Required,       // not empty
		ozzo.Length(5, 100), // length between 5 and 100
		is.URL,              // is a valid URL
	)
	fmt.Println(err)
	// Output:
	// must be a valid URL
}

func Example_third() {
	addresses := []Address{
		{State: "MD", Zip: "12345"},
		{Street: "123 Main St", City: "Vienna", State: "VA", Zip: "12345"},
		{City: "Unknown", State: "NC", Zip: "123"},
	}
	err := ozzo.Validate(addresses)
	fmt.Println(err)
	// Output:
	// 0: (City: cannot be blank; Street: cannot be blank.); 2: (Street: cannot be blank; Zip: must be in a valid format.).
}

func Example_four() {
	c := Customer{
		Name:  "Qiang Xue",
		Email: "q",
		Address: Address{
			State: "Virginia",
		},
	}

	err := ozzo.Errors{
		"name":  ozzo.Validate(c.Name, ozzo.Required, ozzo.Length(5, 20)),
		"email": ozzo.Validate(c.Name, ozzo.Required, is.Email),
		"zip":   ozzo.Validate(c.Address.Zip, ozzo.Required, ozzo.Match(regexp.MustCompile("^[0-9]{5}$"))),
	}.Filter()
	fmt.Println(err)
	// Output:
	// email: must be a valid email address; zip: cannot be blank.
}

func Example_five() {
	type Employee struct {
		Name string
	}

	type Manager struct {
		Employee
		Level int
	}

	m := Manager{}
	err := ozzo.ValidateStruct(&m,
		ozzo.Field(&m.Name, ozzo.Required),
		ozzo.Field(&m.Level, ozzo.Required),
	)
	fmt.Println(err)
	// Output:
	// Level: cannot be blank; Name: cannot be blank.
}

type contextKey int

func Example_six() {
	key := contextKey(1)
	rule := ozzo.WithContext(func(ctx context.Context, value interface{}) error {
		s, _ := value.(string)
		if ctx.Value(key) == s {
			return nil
		}
		return errors.New("unexpected value")
	})
	ctx := context.WithValue(context.Background(), key, "good sample")

	err1 := ozzo.ValidateWithContext(ctx, "bad sample", rule)
	fmt.Println(err1)

	err2 := ozzo.ValidateWithContext(ctx, "good sample", rule)
	fmt.Println(err2)

	// Output:
	// unexpected value
	// <nil>
}

func Example_seven() {
	c := map[string]interface{}{
		"Name":  "Qiang Xue",
		"Email": "q",
		"Address": map[string]interface{}{
			"Street": "123",
			"City":   "Unknown",
			"State":  "Virginia",
			"Zip":    "12345",
		},
	}

	err := ozzo.Validate(c,
		ozzo.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			ozzo.Key("Name", ozzo.Required, ozzo.Length(5, 20)),
			// Email cannot be empty and should be in a valid email format.
			ozzo.Key("Email", ozzo.Required, is.Email),
			// Validate Address using its own validation rules
			ozzo.Key("Address", ozzo.Map(
				// Street cannot be empty, and the length must between 5 and 50
				ozzo.Key("Street", ozzo.Required, ozzo.Length(5, 50)),
				// City cannot be empty, and the length must between 5 and 50
				ozzo.Key("City", ozzo.Required, ozzo.Length(5, 50)),
				// State cannot be empty, and must be a string consisting of two letters in upper case
				ozzo.Key("State", ozzo.Required, ozzo.Match(regexp.MustCompile("^[A-Z]{2}$"))),
				// State cannot be empty, and must be a string consisting of five digits
				ozzo.Key("Zip", ozzo.Required, ozzo.Match(regexp.MustCompile("^[0-9]{5}$"))),
			)),
		),
	)
	fmt.Println(err)
	// Output:
	// Address: (State: must be in a valid format; Street: the length must be between 5 and 50.); Email: must be a valid email address.
}
