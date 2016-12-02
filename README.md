# ozzo-validation

[![GoDoc](https://godoc.org/github.com/go-ozzo/ozzo-validation?status.png)](http://godoc.org/github.com/go-ozzo/ozzo-validation)
[![Build Status](https://travis-ci.org/go-ozzo/ozzo-validation.svg?branch=master)](https://travis-ci.org/go-ozzo/ozzo-validation)
[![Coverage Status](https://coveralls.io/repos/github/go-ozzo/ozzo-validation/badge.svg?branch=master)](https://coveralls.io/github/go-ozzo/ozzo-validation?branch=master)
[![Go Report](https://goreportcard.com/badge/github.com/go-ozzo/ozzo-validation)](https://goreportcard.com/report/github.com/go-ozzo/ozzo-validation)

## Description

ozzo-validation is a Go package that provides configurable and extensible data validation capabilities.
It uses programming constructs to specify how data should be validated rather than relying on error-prone struct tags,
which makes your code more flexible and less error prone. ozzo-validation has the following features:

* rule-based data validation that allows validating a data value with multiple rules.
* validation rules are declared via normal programming constructs instead of error-prone struct tags.
* can validate data of different types, e.g., structs, strings, byte slices, slices, maps, arrays.
* can validate custom data types as long as they implement the `Validatable` interface.
* support validating data types that implement the `sql.Valuer` interface (e.g. `sql.NullString`).
* support validating selective struct fields.
* customizable and well-formatted validation errors.
* provide a rich set of validation rules right out of box.
* extremely easy to create and use custom validation rules.


## Requirements

Go 1.5 or above.

## Installation

Run the following command to install the package:

```
go get github.com/go-ozzo/ozzo-validation
```

You may also get specified release of the package by:

```
go get gopkg.in/go-ozzo/ozzo-validation.v2
```


## Validating Structs

Struct validation is perhaps the most common use case for data validation. Typically, validation is needed
after a struct is populated with the client-side data. You can use `validation.StructRules` to specify how
struct fields should be validated, and then call `StructRules.Validate()` to perform the validation.

For example,

```go
package main

import (
	"fmt"
	"regexp"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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
	return validation.StructRules{}.
		// Street cannot be empty, and the length must between 5 and 50
		Add("Street", validation.Required, validation.Length(5, 50)).
		// City cannot be empty, and the length must between 5 and 50
		Add("City", validation.Required, validation.Length(5, 50)).
		// State cannot be empty, and must be a string consisting of two letters in upper case
		Add("State", validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))).
		// Zip cannot be empty, and must be a string consisting of five digits
		Add("Zip", validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))).
		// performs validation
		Validate(a)
}

func (c Customer) Validate() error {
	return validation.StructRules{}.
		// Name cannot be empty, and the length must be between 5 and 20.
		Add("Name", validation.Required, validation.Length(5, 20)).
		// Gender is optional, and should be either "Female" or "Male".
		Add("Gender", validation.In("Female", "Male")).
		// Email cannot be empty and should be in a valid email format.
		Add("Email", validation.Required, is.Email).
		// Validate Address using its own validation rules
		Add("Address").
		// performs validation
		Validate(c)
}

func main() {
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

	err := validation.Validate(c) // or alternatively, err := c.Validate()
	fmt.Println(err)
	// Output:
	// Address: (State: must be in a valid format.); Email: must be a valid email address.
}
```

The method `StructRules.Add()` is used to specify the rules for validating a particular struct field.
A single field can be associated with multiple rules, and a single struct can have rules for multiple fields.

When the validation is performed, the fields are validated in the order they are added to `StructRules`. Similarly, for
each field being validated, the rules are also executed in the order they are associated with the field.
If a rule fails, an error is recorded for that field, and the validation will continue with the next field.

The method `StructRules.Validate()` returns validation errors as `validation.Errors` which is a map of fields
and their corresponding errors. Nil is returned if validation passes.

Only public struct fields can be validated. A validation error will be reported if trying to validate a private
or non-existing struct field.

### Nested Validation

If a struct field implements the `validation.Validatable` interface, besides running through the validation
rules associated with the field, the field's `Validate()` method will also be called when validating the whole struct.
In the above example, `Address` is such a field because the `Address` type implements `validation.Validatable`.

If a struct field is a map, slice, or array, and its elements implement the `validation.Validatable` interface,
the validation will also be carried for each element while validating the field. See "Validating Maps, Slices, and Arrays"
for more details.

Note that in order for a field's `Validate()` to be called, the field must be listed in `StructRules` even if
the field has no associated rules.

Sometimes, you may want to skip the field's `Validate()`. To do so, you may associate a `validation.Skip` rule
with the field.

### Validating Selected Fields of a Struct

By default, `StructRules.Validate()` will validate every field that has rules. You can explicitly specify which
fields should be validated by passing the field names to the method. For example, the following code only
validate the `Name` and `Email` fields even though more fields have associated validation rules:

```go
err := validation.StructRules{}.
	// Name cannot be empty, and the length must be between 5 and 20.
	Add("Name", validation.Required, validation.Length(5, 20)).
	// Gender is optional, and should be either "Female" or "Male".
	Add("Gender", validation.In("Female", "Male")).
	// Email cannot be empty and should be in a valid email format.
	Add("Email", validation.Required, is.Email).
	// Validate Address using its own validation rules
	Add("Address").
	// only validate Name and Email
	Validate(customer, "Name", "Email")
```

## Validating Simple Values

A simple data value (e.g. strings, integers) can be validated by building `validation.Rules`
which represents a list of validation rules and calling `Rules.Validate()`. For example,

```go
rules := validation.Rules{
	validation.Required,          // not empty
	validation.Length(5, 100),    // length between 5 and 100
	is.URL,                       // is a valid URL
}

data := "example"
err := rules.Validate(data)
fmt.Println(err)
// Output:
// must be a valid URL
```

The method `Rules.Validate()` will run through the rules in the order that they are declared. If a rule returns an error,
it will return the error and skip the rest of the rules.


## Validating Maps, Slices, and Arrays

When the elements of a map, a slice, or an array implement the `validation.Validatable` interface, calling
`validation.Validate()` on the map, slice, or array will automatically call the `Validate()` method of each non-nil
element. The validation errors of the elements will be returned as `validation.Errors` which maps the keys of the
invalid elements to their corresponding validation errors. For example,

```go
addresses := []Address{
	Address{State: "MD", Zip: "12345"},
	Address{Street: "123 Main St", City: "Vienna", State: "VA", Zip: "12345"},
	Address{City: "Unknown", State: "NC", Zip: "123"},
}
err := validation.Validate(addresses)
fmt.Println(err)
// Output:
// 0: (City: cannot be blank; Street: cannot be blank.); 2: (Street: cannot be blank; Zip: must be in a valid format.).
```

## Validating Pointers

When a value being validated is a pointer, most validation rules will validate the actual value pointed to by the pointer.
If the pointer is nil, these rules will skip the validation.

An exception is the `validation.Required` and `validation.Required` rules. When a pointer is nil, they
will report a validation error.


## Validating `sql.Valuer`

If a data type implements the `sql.Valuer` interface (e.g. `sql.NullString`), the built-in validation rules will handle
it properly. In particular, when a rule is validating such data, it will call the `Value()` method and validate
the returned value instead.


## Processing Validation Errors

All validation methods return an `error` when validation fails. The `error` may be typecast into a `validation.Errors`
if the value being validated is a struct, a map/slice/array of validatables. `validation.Errors` implements both
`error` and `json.Marshaler` interfaces and can return a well-formatted text or JSON string.

By default, `validation.Errors` uses struct field names as its keys when the validation errors come from a struct.
You can customize the key names using struct tags named `validation`. For example,

```go
type Address struct {
	Street string `validation:"street"`
	City   string `validation:"city"`
	State  string `validation:"state"`
	Zip    string `validation:"zip"`
}
```

This could be useful if you are using snake case or camelCase in your JSON responses.

You may customize the tag name by changing `validation.ErrorTag`. For example, you may set `validation.ErrorTag`
to be `"json"` if you want to reuse the JSON field names as error field names.


## Required vs. Not Nil

When validating input values, there are two different scenarios about checking if input values are provided or not.

In the first scenario, an input value is considered missing if it is not entered or it is entered as a zero value
(e.g. an empty string, a zero integer). You can use the `validation.Required` rule in this case.

In the second scenario, an input value is considered missing only if it is not entered. A pointer field is usually
used in this case so that you can detect if a value is entered or not by checking if the pointer is nil or not.
You can use the `validation.NotNil` rule to ensure a value is entered (even if it is a zero value).

## Built-in Validation Rules

The following rules are provided in the `validation` package:

* `In(...interface{})`: checks if a value can be found in the given list of values.
* `Length(min, max int)`: checks if the length of a value is within the specified range.
  This rule should only be used for validating strings, slices, maps, and arrays.
* `Range(min, max int)`: checks if a value is within the specified range.
  This rule should only be used for validating int, uint and float types.
* `Match(*regexp.Regexp)`: checks if a value matches the specified regular expression.
  This rule should only be used for strings and byte slices.
* `Required`: checks if a value is not empty (neither nil nor zero).
* `NotNil`: checks if a pointer value is not nil. Non-pointer values are considered valid.
* `Skip`: this is a special rule used to indicate that all rules following it should be skipped (including the nested ones).

The `is` sub-package provides a list of commonly used string validation rules that can be used to check if the format
of a value satisfies certain requirements. Note that these rules only handle strings and byte slices and if a string
 or byte slice is empty, it is considered valid. You may use a `Required` rule to ensure a value is not empty.
Below is the whole list of the rules provided by the `is` package:

* `Email`: validates if a string is an email or not
* `URL`: validates if a string is a valid URL
* `RequestURL`: validates if a string is a valid request URL
* `RequestURI`: validates if a string is a valid request URI
* `Alpha`: validates if a string contains English letters only (a-zA-Z)
* `Digit`: validates if a string contains digits only (0-9)
* `Alphanumeric`: validates if a string contains English letters and digits only (a-zA-Z0-9)
* `UTFLetter`: validates if a string contains unicode letters only
* `UTFDigit`: validates if a string contains unicode decimal digits only
* `UTFLetterNumeric`: validates if a string contains unicode letters and numbers only
* `UTFNumeric`: validates if a string contains unicode number characters (category N) only
* `LowerCase`: validates if a string contains lower case unicode letters only
* `UpperCase`: validates if a string contains upper case unicode letters only
* `Hexadecimal`: validates if a string is a valid hexadecimal number
* `HexColor`: validates if a string is a valid hexadecimal color code
* `RGBColor`: validates if a string is a valid RGB color in the form of rgb(R, G, B)
* `Int`: validates if a string is a valid integer number
* `Float`: validates if a string is a floating point number
* `UUIDv3`: validates if a string is a valid version 3 UUID
* `UUIDv4`: validates if a string is a valid version 4 UUID
* `UUIDv5`: validates if a string is a valid version 5 UUID
* `UUID`: validates if a string is a valid UUID
* `CreditCard`: validates if a string is a valid credit card number
* `ISBN10`: validates if a string is an ISBN version 10
* `ISBN13`: validates if a string is an ISBN version 13
* `ISBN`: validates if a string is an ISBN (either version 10 or 13)
* `JSON`: validates if a string is in valid JSON format
* `ASCII`: validates if a string contains ASCII characters only
* `PrintableASCII`: validates if a string contains printable ASCII characters only
* `Multibyte`: validates if a string contains multibyte characters
* `FullWidth`: validates if a string contains full-width characters
* `HalfWidth`: validates if a string contains half-width characters
* `VariableWidth`: validates if a string contains both full-width and half-width characters
* `Base64`: validates if a string is encoded in Base64
* `DataURI`: validates if a string is a valid base64-encoded data URI
* `CountryCode2`: validates if a string is a valid ISO3166 Alpha 2 country code
* `CountryCode3`: validates if a string is a valid ISO3166 Alpha 3 country code
* `DialString`: validates if a string is a valid dial string that can be passed to Dial()
* `MAC`: validates if a string is a MAC address
* `IP`: validates if a string is a valid IP address (either version 4 or 6)
* `IPv4`: validates if a string is a valid version 4 IP address
* `IPv6`: validates if a string is a valid version 6 IP address
* `DNSName`: validates if a string is valid DNS name
* `Host`: validates if a string is a valid IP (both v4 and v6) or a valid DNS name
* `Port`: validates if a string is a valid port number
* `MongoID`: validates if a string is a valid Mongo ID
* `Latitude`: validates if a string is a valid latitude
* `Longitude`: validates if a string is a valid longitude
* `SSN`: validates if a string is a social security number (SSN)
* `Semver`: validates if a string is a valid semantic version

### Customizing Error Messages

All the built-in validation rules allow you to customize error messages. To do so, simply call the `Error()` method
of the rules, e.g.,

```go
rules := validation.Rules{
	validation.Required.Error("is required"),
	validation.Match(regexp.MustCompile("^[0-9]{5}$")).Error("must be a string with five digits"),
}

data := "2123"
err := rules.Validate(data)
fmt.Println(err)
// Output:
// must be a string with five digits
```

## Creating Custom Rules

Creating a custom rule is as simple as implementing the `validation.Rule` interface. The interface contains a single
method as shown below, which should validate the value and return the validation error, if any.

```go
// Validate validates a value and returns an error if validation fails.
Validate(value interface{}) error
```

## Credits

The `is` sub-package wraps the excellent validators provided by the [govalidator](https://github.com/asaskevich/govalidator) package.
