package validation_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Street, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.City, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&a.State, validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}

func (c Customer) Validate() error {
	return validation.ValidateStruct(&c,
		// Name cannot be empty, and the length must be between 5 and 20.
		validation.Field(&c.Name, validation.Required, validation.Length(5, 20)),
		// Gender is optional, and should be either "Female" or "Male".
		validation.Field(&c.Gender, validation.In("Female", "Male")),
		// Email cannot be empty and should be in a valid email format.
		validation.Field(&c.Email, validation.Required, is.Email),
		// Validate Address using its own validation rules
		validation.Field(&c.Address),
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
	err := validation.Validate(data,
		validation.Required,       // not empty
		validation.Length(5, 100), // length between 5 and 100
		is.URL,                    // is a valid URL
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
	err := validation.Validate(addresses)
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

	err := validation.Errors{
		"name":  validation.Validate(c.Name, validation.Required, validation.Length(5, 20)),
		"email": validation.Validate(c.Name, validation.Required, is.Email),
		"zip":   validation.Validate(c.Address.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
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
	err := validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Level, validation.Required),
	)
	fmt.Println(err)
	// Output:
	// Level: cannot be blank; Name: cannot be blank.
}

type contextKey int

func Example_six() {
	key := contextKey(1)
	rule := validation.WithContext(func(ctx context.Context, value interface{}) error {
		s, _ := value.(string)
		if ctx.Value(key) == s {
			return nil
		}
		return errors.New("unexpected value")
	})
	ctx := context.WithValue(context.Background(), key, "good sample")

	err1 := validation.ValidateWithContext(ctx, "bad sample", rule)
	fmt.Println(err1)

	err2 := validation.ValidateWithContext(ctx, "good sample", rule)
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

	err := validation.Validate(c,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("Name", validation.Required, validation.Length(5, 20)),
			// Email cannot be empty and should be in a valid email format.
			validation.Key("Email", validation.Required, is.Email),
			// Validate Address using its own validation rules
			validation.Key("Address", validation.Map(
				// Street cannot be empty, and the length must between 5 and 50
				validation.Key("Street", validation.Required, validation.Length(5, 50)),
				// City cannot be empty, and the length must between 5 and 50
				validation.Key("City", validation.Required, validation.Length(5, 50)),
				// State cannot be empty, and must be a string consisting of two letters in upper case
				validation.Key("State", validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
				// State cannot be empty, and must be a string consisting of five digits
				validation.Key("Zip", validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
			)),
		),
	)
	fmt.Println(err)
	// Output:
	// Address: (State: must be in a valid format; Street: the length must be between 5 and 50.); Email: must be a valid email address.
}

func ExampleThresholdRule_ValidateEmpty() {
	// Min and Max skip empty values: 0 counts as "not set".
	fmt.Println(validation.Validate(0, validation.Min(0).Exclusive()))
	// ValidateEmpty checks 0 against the threshold as well.
	fmt.Println(validation.Validate(0, validation.Min(0).Exclusive().ValidateEmpty()))
	// Output:
	// <nil>
	// must be greater than 0
}

type Location struct {
	Latitude  string
	Longitude string
}

// IsEmpty reports a Location with no coordinates as empty.
func (l Location) IsEmpty() bool { return l == Location{} }

type Item struct {
	Name string `json:"name"`
}

func (i Item) Validate() error {
	return validation.ValidateStruct(&i, validation.Field(&i.Name, validation.Required))
}

type Order struct {
	ID    string `json:"id"`
	Items []Item `json:"items"`
}

func ExampleFlatten() {
	o := Order{ID: "a1", Items: make([]Item, 11)}
	for i := range o.Items {
		o.Items[i].Name = "pen"
	}
	o.Items[2].Name, o.Items[10].Name = "", ""
	err := validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Length(3, 10)),
		validation.Field(&o.Items),
	)
	for _, fe := range validation.Flatten(err) {
		var ie validation.InternalError
		if errors.As(fe.Err, &ie) {
			continue // log it; its message is not for clients
		}
		fmt.Printf("%s %s %s: %s\n", fe.Field(), fe.Pointer(), fe.Code, fe.Message)
	}

	// a single error is one entry with no path
	for _, fe := range validation.Flatten(validation.Validate("", validation.Required)) {
		fmt.Printf("%q %s: %s\n", fe.Pointer(), fe.Code, fe.Message)
	}
	// Output:
	// id /id validation_length_out_of_range: the length must be between 3 and 10
	// items.2.name /items/2/name validation_required: cannot be blank
	// items.10.name /items/10/name validation_required: cannot be blank
	// "" validation_required: cannot be blank
}

func ExampleErrors_Flatten() {
	o := Order{ID: "a1", Items: []Item{{Name: "pen"}, {}, {}}}
	err := validation.ValidateStruct(&o, validation.Field(&o.Items))

	var errs validation.Errors
	if errors.As(err, &errs) {
		for _, fe := range errs.Flatten() {
			fmt.Println(fe.Pointer(), fe.Code)
		}
	}
	// Output:
	// /items/1/name validation_required
	// /items/2/name validation_required
}

func ExampleEmptyChecker() {
	// Location implements validation.EmptyChecker, so Required can tell when it is not set.
	fmt.Println(validation.Validate(Location{}, validation.Required))
	fmt.Println(validation.Validate(Location{Latitude: "52.52", Longitude: "13.40"}, validation.Required))
	// A struct without IsEmpty is never empty.
	fmt.Println(validation.Validate(struct{ Name string }{}, validation.Required))
	// Output:
	// cannot be blank
	// <nil>
	// <nil>
}
