# ozzo-validation

[![GoDoc](https://godoc.org/github.com/go-ozzo/ozzo-validation?status.png)](http://godoc.org/github.com/go-ozzo/ozzo-validation)
[![Build Status](https://travis-ci.org/go-ozzo/ozzo-validation.svg?branch=master)](https://travis-ci.org/go-ozzo/ozzo-validation)
[![Coverage](http://gocover.io/_badge/github.com/go-ozzo/ozzo-validation)](http://gocover.io/github.com/go-ozzo/ozzo-validation)
[![Go Report](https://goreportcard.com/badge/github.com/go-ozzo/ozzo-validation)](https://goreportcard.com/report/github.com/go-ozzo/ozzo-validation)

## Description

ozzo-validation is a Go package that provides configurable and extensible validation capabilities for data of
various types. It has the following features:

* rule-based data validation that allows specifying a list of validation rules to validate a data value.
* validation rules are declared via normal code constructs instead of error-prone struct tags.
* support validating selective struct fields.
* can validate data of different types, e.g., strings, byte slices, structs, slices, maps, arrays.
* can validate custom data types provided they implement the `Validatable` interface.
* provide a rich set of validation rules right out of box.
* creating and using a custom configurable validation rule is extremely easy.


## Requirements

Go 1.5 or above.

## Installation

Run the following command to install the package:

```
go get github.com/go-ozzo/ozzo-validation
```

## Getting Started

Create a `main.go` file with the following content:

```go
package main

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

func main() {
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
```

## Validating Simple Values

A validation rule specifies a particular aspect that a data value must satisfy. For example, a rule may require
a data string to be a valid email.

A data value may be validated using a list of validation rules.

For simple data values (e.g. strings, integers), you may declare validation rules by building
a `validation.Rules` slice. For example,

```go
import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

rules := validation.Rules{
	validation.Required,
	is.Email,
}
```

You can then call `Rules.Validate()` to validate the data value:

```go
data := "q"
err := rules.Validate(data, nil)
```

The method `Rules.Validate()` will run through the rules in the order they are declared. If a rule returns an error,
it will return the error and skip the rest of the rules.

## Validating Struct Values

For struct values, validation rules may be declared using `validation.StructRules` which allows you to validate
each individual struct field. For example,

```go
import (
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

rules := validation.StructRules{}.
	Add("Street", validation.NotEmpty).
	Add("City", validation.NotEmpty).
	Add("State", validation.NotEmpty, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))).
	Add("Zip", validation.NotEmpty, validation.Match(regexp.MustCompile("^[0-9]{5}$")))
}
```

You can call `StructRules.Validate()` to validate a struct. Every field listed in `StructRules` will
be validated. You can also specify a subset of the fields to be validated. For example,

```go
address := Customer{}

// validates every field listed in rules
err := rules.Validate(address, "State", "Zip")

// only validates State and Zip
err = rules.Validate(address, "State", "Zip")
```

## Validating Maps, Slices, and Arrays

## Validatable Values

## Processing Validation Errors

## Required vs. Not Empty

## Built-in Validation Rules

## Customizing Error Messages

## Creating Custom Rules

## Credits
