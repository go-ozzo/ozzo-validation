package validation

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Issue #234: behavior of v4.3.0 that v4.4.0 changed. Every case here passes on
// v4.3.0 and fails on v4.4.0/v4.4.1/master.

type issue234Plain struct {
	A int
	B string
}

type issue234Body struct {
	Action string
}

func (b issue234Body) Validate() error {
	return ValidateStruct(&b, Field(&b.Action, Required))
}

func TestIssue234StructIsNotEmpty(t *testing.T) {
	// Only a zero time.Time counts as an empty struct; any other struct is not empty.
	assert.False(t, IsEmpty(issue234Plain{}), "zero struct")
	assert.False(t, IsEmpty(&issue234Plain{}), "pointer to zero struct")
	assert.True(t, IsEmpty(time.Time{}), "zero time.Time")

	assert.NoError(t, Validate(issue234Plain{}, Required), "Required on zero struct")
	assert.NoError(t, Validate(&issue234Plain{}, NilOrNotEmpty), "NilOrNotEmpty on pointer to zero struct")
	assert.EqualError(t, Validate(issue234Plain{}, In(issue234Plain{A: 1})), "must be a valid value", "In on zero struct")
}

func TestIssue234RequiredZeroValidatableStruct(t *testing.T) {
	// Required on a zero Validatable struct used to pass and let the struct report its own
	// field errors (AkamaiOPEN-edgegrid-golang: "Body": validation.Validate(r.Body, validation.Required)).
	err := Validate(issue234Body{}, Required)
	assert.EqualError(t, err, "Action: cannot be blank.")

	req := struct{ Body issue234Body }{}
	err = ValidateStruct(&req, Field(&req.Body, Required))
	if assert.Error(t, err) {
		es, ok := err.(Errors)
		if assert.True(t, ok) {
			assert.EqualError(t, es["Body"], "Action: cannot be blank.")
		}
	}
}

func TestIssue234MinMaxSkipEmpty(t *testing.T) {
	// A zero number or a zero time.Time is "not set" and skipped by Min/Max
	// (edgegrid: Size int with Min(10), 0 = query parameter omitted).
	zero := 0
	tests := []struct {
		tag   string
		rule  Rule
		value interface{}
	}{
		{"Min(10) int 0", Min(10), 0},
		{"Max(-10) int 0", Max(-10), 0},
		{"Min(10) ptr to 0", Min(10), &zero},
		{"Min(0.001) float64 0", Min(0.001), 0.0},
		{"Min(0).Exclusive() int 0", Min(0).Exclusive(), 0},
		{"Min(uint(1)) uint 0", Min(uint(1)), uint(0)},
		{"Min(time) zero time", Min(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)), time.Time{}},
		{"Min(10) empty string", Min(10), ""},
	}
	for _, test := range tests {
		assert.NoError(t, Validate(test.value, test.rule), test.tag)
	}

	req := struct {
		Size     int
		Interval int
		Timeout  float64
	}{}
	assert.NoError(t, ValidateStruct(&req,
		Field(&req.Size, Min(10)),
		Field(&req.Interval, Min(10), Max(3600)),
		Field(&req.Timeout, Min(0.001), Max(60.0)),
	), "unset optional numeric fields")

	// Non-zero values are still checked.
	assert.EqualError(t, Validate(5, Min(10)), "must be no less than 10")
	assert.Error(t, Validate(errors.New("x"), Min(10)), "unsupported type")
}
