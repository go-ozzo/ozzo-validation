# ozzo-validation

[![GoDoc](https://godoc.org/github.com/go-ozzo/ozzo-validation?status.png)](http://godoc.org/github.com/go-ozzo/ozzo-validation)
[![Build Status](https://travis-ci.org/go-ozzo/ozzo-validation.svg?branch=master)](https://travis-ci.org/go-ozzo/ozzo-validation)
[![Coverage](http://gocover.io/_badge/github.com/go-ozzo/ozzo-validation)](http://gocover.io/github.com/go-ozzo/ozzo-validation)
[![Go Report](https://goreportcard.com/badge/github.com/go-ozzo/ozzo-validation)](https://goreportcard.com/report/github.com/go-ozzo/ozzo-validation)

## Description

ozzo-validation is a Go package that provides configurable and extensible data validation capabilities.
It uses programming constructs to specify how data should be validated rather than relying on struct tags,
which makes your code more flexible and less error prone. ozzo-validation has the following features:

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


## Validating Struct Values

Struct validation is perhaps the most common use case for data validation. Typically, validation is needed
after a struct is populated with the client-side data. You can use `validation.StructRules` to specify how
struct fields should be validated, and then call `StructRules.Validate()` to perform validation when needed.
For example,

```go
package main

import (
	"fmt"
	"regexp"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Customer struct {
	Name  string
	Email string
	Zip   string
}

func main() {
	rules := validation.StructRules{}.
		Add("Name", validation.NotEmpty).
		Add("Email", validation.NotEmpty, is.Email).
		Add("Zip", validation.NotEmpty, validation.Match(regexp.MustCompile("^[0-9]{5}$")))

	customer := Customer{}

	// validates every field listed in rules
	err := rules.Validate(customer)
	fmt.Println(err)

	// only validates Email and Zip
	err = rules.Validate(customer, "Email", "Zip")
	fmt.Println(err)
}
```

The method `StructRules.Add()` can be used to specify the rules for validating a particular struct field.
A single field can be associated with multiple rules, and a single struct can have rules for multiple fields.

When validation is performed, the fields are validated in the order they are added to `StructRules`, and for
each field being validated, the rules are also executed in the order they are associated with the field.
If a rule fails, an error is recorded for that field, and the validation will continue with the next field.

The method `StructRules.Validate()` returns validation errors as `validation.Errors` which is a map of fields
and their corresponding errors. Nil is returned if validation passes.



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


## Validating Maps, Slices, and Arrays

## Validatable Values

## Processing Validation Errors

## Required vs. Not Empty

## Built-in Validation Rules

## Customizing Error Messages

## Creating Custom Rules

## Credits
