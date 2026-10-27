// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInternalError(t *testing.T) {
	err := NewInternalError(errors.New("abc"))
	if assert.NotNil(t, err.InternalError()) {
		assert.Equal(t, "abc", err.InternalError().Error())
	}
}

func TestErrors_Error(t *testing.T) {
	errs := Errors{
		"B": errors.New("B1"),
		"C": errors.New("C1"),
		"A": errors.New("A1"),
	}
	assert.Equal(t, "A: A1; B: B1; C: C1.", errs.Error())

	errs = Errors{
		"B": errors.New("B1"),
	}
	assert.Equal(t, "B: B1.", errs.Error())

	errs = Errors{}
	assert.Equal(t, "", errs.Error())

	// issue #147: nil entries are skipped, with no stray separator
	nilTests := []struct {
		tag      string
		errs     Errors
		expected string
	}{
		{"nil in the middle", Errors{"A": errors.New("A1"), "B": nil, "C": errors.New("C1")}, "A: A1; C: C1."},
		{"nil first", Errors{"A": nil, "B": errors.New("B1")}, "B: B1."},
		{"nil last", Errors{"A": errors.New("A1"), "B": nil}, "A: A1."},
		{"only nil", Errors{"A": nil}, ""},
		{"all nil", Errors{"A": nil, "B": nil}, ""},
	}
	for _, test := range nilTests {
		assert.Equal(t, test.expected, test.errs.Error(), test.tag)
	}
}

func TestErrors_Is(t *testing.T) {
	sentinel := errors.New("sentinel")

	assert.True(t, errors.Is(Errors{"A": errors.New("a"), "B": sentinel}, sentinel))
	assert.False(t, errors.Is(Errors{"A": errors.New("a"), "B": errors.New("b")}, sentinel))

	// wrapped field errors are matched through the errors chain
	assert.True(t, errors.Is(Errors{"A": fmt.Errorf("wrapped: %w", sentinel)}, sentinel))

	// errors nested in Errors (maps, slices, nested structs) are matched too
	assert.True(t, errors.Is(Errors{"A": Errors{"0": sentinel}}, sentinel))

	// nil field errors and empty Errors do not match, and no target never matches
	assert.False(t, errors.Is(Errors{"A": nil}, sentinel))
	assert.False(t, errors.Is(Errors{}, sentinel))
	assert.False(t, errors.Is(Errors{"A": errors.New("a")}, nil))
	assert.False(t, errors.Is(Errors(nil), sentinel))

	// an Errors is not comparable, so it does not match itself
	es := Errors{"A": errors.New("a")}
	assert.False(t, errors.Is(es, es))
}

func BenchmarkErrorsIs100k(b *testing.B) {
	es := make(Errors, 100000)
	for i := 0; i < 100000; i++ {
		es[strconv.Itoa(i)] = ErrRequired
	}
	sentinel := errors.New("sentinel")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if errors.Is(es, sentinel) {
			b.Fatal("unexpected match")
		}
	}
}

type isTestInner struct {
	Name string
}

func (i isTestInner) Validate() error {
	return ValidateStruct(&i, Field(&i.Name, Required))
}

func TestErrorObject_Is(t *testing.T) {
	// a direct error, also with a custom message, and wrapped
	assert.True(t, errors.Is(Validate("", Required), ErrRequired))
	assert.True(t, errors.Is(Validate("", Required.Error("custom")), ErrRequired))
	assert.True(t, errors.Is(fmt.Errorf("wrapped: %w", Validate("", Required)), ErrRequired))
	// params do not matter, only the code
	assert.True(t, errors.Is(Validate("abc", Length(5, 10)), ErrLengthOutOfRange))

	// through ValidateStruct, a nested struct, Each and Map
	var s struct {
		DB    string
		Inner isTestInner
		Tags  []string
		Attrs map[string]interface{}
	}
	s.Tags = []string{"a", ""}
	s.Attrs = map[string]interface{}{"k": ""}
	err := ValidateStruct(&s,
		Field(&s.Inner),
		Field(&s.Tags, Each(Required)),
		Field(&s.Attrs, Map(Key("k", Required))),
	)
	if assert.Error(t, err) {
		assert.True(t, errors.Is(err, ErrRequired))
		assert.True(t, errors.Is(err.(Errors)["Inner"], ErrRequired))
		assert.True(t, errors.Is(err.(Errors)["Tags"], ErrRequired))
		assert.True(t, errors.Is(err.(Errors)["Attrs"], ErrRequired))
		assert.False(t, errors.Is(err, ErrNilOrNotEmpty))
	}
	assert.True(t, errors.Is(ValidateStruct(&s, Field(&s.DB, Required)), ErrRequired))

	// different codes do not match, and neither do errors without a code
	assert.False(t, errors.Is(Validate("", Required), ErrNilOrNotEmpty))
	assert.False(t, errors.Is(NewError("", "a"), NewError("", "a")))
	assert.False(t, NewError("", "a").(ErrorObject).Is(NewError("", "a")))
	assert.False(t, errors.Is(NewError("code", "a"), errors.New("code")))
	assert.False(t, NewError("code", "a").(ErrorObject).Is(nil))
	// any Error implementation with the same code matches
	assert.True(t, errors.Is(NewError("code", "a"), &ErrorObject{code: "code"}))
	// targets that are not a usable Error report false without panicking
	assert.False(t, errors.Is(NewError("code", "a"), (*ErrorObject)(nil)))
	assert.False(t, errors.Is(NewError("code", "a"), Errors{"A": NewError("code", "a")}))
	assert.False(t, errors.Is(Errors{"A": NewError("code", "a")}, (*ErrorObject)(nil)))

	// sentinels returned by By rules still match, and do not match by code
	sentinel := errors.New("sentinel")
	rule := By(func(interface{}) error { return sentinel })
	err = ValidateStruct(&s, Field(&s.DB, rule))
	assert.True(t, errors.Is(err, sentinel))
	assert.False(t, errors.Is(err, ErrRequired))
}

type customError struct {
	code int
}

func (e customError) Error() string { return fmt.Sprintf("custom %d", e.code) }

// errors.As must not look inside an Errors: callers use errors.As(err, &validation.Error)
// to tell a single-value error from a field map, and a match would lose the field names.
func TestErrors_NotAs(t *testing.T) {
	type config struct {
		DB string
	}
	var c config
	err := ValidateStruct(&c, Field(&c.DB, Required))
	if assert.Error(t, err) {
		var single Error
		assert.False(t, errors.As(err, &single), "a field map must not match validation.Error")
		assert.Nil(t, single)

		var fields Errors
		if assert.True(t, errors.As(err, &fields)) {
			assert.Equal(t, "cannot be blank", fields["DB"].Error())
		}
	}

	var custom customError
	assert.False(t, errors.As(Errors{"A": customError{code: 7}}, &custom))
	assert.False(t, errors.As(Errors{"A": Errors{"0": customError{code: 9}}}, &custom))
	assert.Equal(t, customError{}, custom)

	// a single-value error still matches validation.Error
	var single Error
	if assert.True(t, errors.As(Validate("", Required), &single)) {
		assert.Equal(t, ErrRequired.Code(), single.Code())
	}
}

func TestErrors_MarshalMessage(t *testing.T) {
	errs := Errors{
		"A": errors.New("A1"),
		"B": Errors{
			"2": errors.New("B1"),
		},
	}
	errsJSON, err := errs.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "{\"A\":\"A1\",\"B\":{\"2\":\"B1\"}}", string(errsJSON))
}

func TestErrors_Filter(t *testing.T) {
	errs := Errors{
		"B": errors.New("B1"),
		"C": nil,
		"A": errors.New("A1"),
	}
	err := errs.Filter()
	assert.Equal(t, 2, len(errs))
	if assert.NotNil(t, err) {
		assert.Equal(t, "A: A1; B: B1.", err.Error())
	}

	errs = Errors{}
	assert.Nil(t, errs.Filter())

	errs = Errors{
		"B": nil,
		"C": nil,
	}

	assert.Nil(t, errs.Filter())
}

func TestErrorObject_SetCode(t *testing.T) {
	err := NewError("A", "msg").(ErrorObject)

	assert.Equal(t, err.code, "A")
	assert.Equal(t, err.Code(), "A")

	err = err.SetCode("B").(ErrorObject)
	assert.Equal(t, "B", err.code)
}

func TestErrorObject_Code(t *testing.T) {
	err := NewError("A", "msg").(ErrorObject)

	assert.Equal(t, err.Code(), "A")
}

func TestErrorObject_SetMessage(t *testing.T) {
	err := NewError("code", "A").(ErrorObject)

	assert.Equal(t, err.message, "A")
	assert.Equal(t, err.Message(), "A")

	err = err.SetMessage("abc").(ErrorObject)
	assert.Equal(t, err.message, "abc")
	assert.Equal(t, err.Message(), "abc")
}

func TestErrorObject_Message(t *testing.T) {
	err := NewError("code", "A").(ErrorObject)

	assert.Equal(t, err.message, "A")
	assert.Equal(t, err.Message(), "A")
}

func TestErrorObject_Params(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "AA": "val2"}

	err := NewError("code", "A").(ErrorObject)
	err = err.SetParams(p).(ErrorObject)
	err = err.SetMessage("B").(ErrorObject)

	assert.Equal(t, err.params, p)
	assert.Equal(t, err.Params(), p)
}

func TestErrorObject_AddParam2(t *testing.T) {
	p := map[string]interface{}{"key": "val"}
	err := NewError("code", "A").(ErrorObject)
	err = err.AddParam("key", "val").(ErrorObject)

	assert.Equal(t, err.params, p)
	assert.Equal(t, err.Params(), p)
}

func TestErrorObject_AddParam(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "B": "val2"}

	err := NewError("code", "A").(ErrorObject)
	err = err.SetParams(p).(ErrorObject)
	err = err.AddParam("C", "val3").(ErrorObject)

	p["C"] = "val3"

	assert.Equal(t, err.params, p)
	assert.Equal(t, err.Params(), p)
}

func TestError_Code(t *testing.T) {
	err := NewError("A", "msg")

	assert.Equal(t, err.Code(), "A")
}

func TestError_SetMessage(t *testing.T) {
	err := NewError("code", "A")

	assert.Equal(t, err.Message(), "A")

	err = err.SetMessage("abc")
	assert.Equal(t, err.Message(), "abc")
}

func TestError_Message(t *testing.T) {
	err := NewError("code", "A")

	assert.Equal(t, err.Message(), "A")
}

func TestError_Params(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "AA": "val2"}

	err := NewError("code", "A")
	err = err.SetParams(p)
	err = err.SetMessage("B")

	assert.Equal(t, err.Params(), p)
}

func TestValidationError(t *testing.T) {
	params := map[string]interface{}{
		"A": "B",
	}

	err := NewError("code", "msg")
	err = err.SetParams(params)

	assert.Equal(t, err.Code(), "code")
	assert.Equal(t, err.Message(), "msg")
	assert.Equal(t, err.Params(), params)

	params = map[string]interface{}{"min": 1}
	err = err.SetParams(params)

	assert.Equal(t, err.Params(), params)
}
