package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type flattenItem struct {
	Name string `json:"name"`
}

func (i flattenItem) Validate() error {
	return ValidateStruct(&i, Field(&i.Name, Required))
}

type flattenOrder struct {
	ID    string            `json:"id"`
	Items []flattenItem     `json:"items"`
	Tags  map[string]string `json:"tags"`
	Ship  flattenItem       `json:"ship"`
}

// paths returns the dotted paths of list, to compare whole lists at a glance.
func paths(list []FieldError) []string {
	s := make([]string, 0, len(list))
	for _, fe := range list {
		s = append(s, fe.Field())
	}
	return s
}

func TestErrors_Flatten_Struct(t *testing.T) {
	o := flattenOrder{
		ID:    "ab",
		Items: []flattenItem{{"x"}, {""}, {"y"}, {""}},
		Tags:  map[string]string{"color": "", "size": "L"},
	}
	err := ValidateStruct(&o,
		Field(&o.ID, Length(3, 10)),
		Field(&o.Items),
		Field(&o.Tags, Each(Required)),
		Field(&o.Ship),
	)
	es, ok := err.(Errors)
	if !assert.True(t, ok, "ValidateStruct returns Errors") {
		return
	}

	list := es.Flatten()
	// keys come from the json tags, and every level is sorted
	assert.Equal(t, []string{"id", "items.1.name", "items.3.name", "ship.name", "tags.color"}, paths(list))

	assert.Equal(t, []string{"id"}, list[0].Path)
	assert.Equal(t, ErrLengthOutOfRange.Code(), list[0].Code)
	// Message is rendered; Error.Message() is the template
	assert.Equal(t, "the length must be between 3 and 10", list[0].Message)
	assert.Equal(t, map[string]interface{}{"min": 3, "max": 10}, list[0].Params)

	assert.Equal(t, []string{"items", "1", "name"}, list[1].Path)
	assert.Equal(t, "/items/1/name", list[1].Pointer())
	for _, fe := range list[1:] {
		assert.Equal(t, ErrRequired.Code(), fe.Code, fe.Field())
		assert.Equal(t, "cannot be blank", fe.Message, fe.Field())
		assert.True(t, errors.Is(fe.Err, ErrRequired), fe.Field())
	}

	// equal inputs give equal output, and es is not modified
	before := es.Error()
	assert.Equal(t, list, es.Flatten())
	assert.Equal(t, before, es.Error())
}

func TestErrors_Flatten_Mixed(t *testing.T) {
	sentinel := errors.New("sentinel")
	es := Errors{
		"b":   nil,
		"a":   sentinel,
		"c":   Errors{"x": nil, "y": Errors{"0": ErrRequired, "1": nil}},
		"d":   fmt.Errorf("wrapped: %w", ErrRequired),
		"e":   fmt.Errorf("profile: %w", Errors{"z": ErrRequired}),
		"f":   customError{code: 7},
		"nil": Errors{"only": nil},
	}
	list := es.Flatten()
	assert.Equal(t, []string{"a", "c.y.0", "d", "e.z", "f"}, paths(list))

	// a plain error: no code, no params, the error itself kept
	assert.Equal(t, FieldError{Path: []string{"a"}, Message: "sentinel", Err: sentinel}, list[0])

	assert.Equal(t, ErrRequired.Code(), list[1].Code)

	// a wrapped validation Error keeps its code, and Message is the wrapping text
	assert.Equal(t, ErrRequired.Code(), list[2].Code)
	assert.Equal(t, "wrapped: cannot be blank", list[2].Message)

	// a wrapped Errors is walked into
	assert.Equal(t, ErrRequired.Code(), list[3].Code)
	assert.Equal(t, "cannot be blank", list[3].Message)

	// a custom error type is recoverable from Err
	var custom customError
	if assert.True(t, errors.As(list[4].Err, &custom)) {
		assert.Equal(t, 7, custom.code)
	}
	assert.Equal(t, "", list[4].Code)
	assert.Equal(t, "custom 7", list[4].Message)
}

func TestErrors_Flatten_WrappedValidatable(t *testing.T) {
	type form struct {
		Profile flattenProfile `json:"profile"`
	}
	f := form{}
	list := Flatten(ValidateStruct(&f, Field(&f.Profile)))
	if assert.Len(t, list, 1) {
		assert.Equal(t, []string{"profile", "name"}, list[0].Path)
		assert.Equal(t, ErrRequired.Code(), list[0].Code)
	}
}

type flattenProfile struct {
	Name string `json:"name"`
}

func (p flattenProfile) Validate() error {
	if err := ValidateStruct(&p, Field(&p.Name, Required)); err != nil {
		return fmt.Errorf("profile: %w", err)
	}
	return nil
}

func TestErrors_Flatten_Empty(t *testing.T) {
	assert.Nil(t, Errors{}.Flatten())
	assert.Nil(t, Errors(nil).Flatten())
	assert.Nil(t, Errors{"a": nil, "b": Errors{}, "c": Errors{"d": nil}}.Flatten())

	// an empty nested Errors contributes nothing next to real errors
	list := Errors{"a": Errors{}, "b": ErrRequired, "c": Errors{"d": nil}}.Flatten()
	assert.Equal(t, []string{"b"}, paths(list))
}

func TestErrors_Flatten_Order(t *testing.T) {
	es := Errors{}
	for i := 11; i >= 0; i-- {
		es[strconv.Itoa(i)] = ErrRequired
	}
	assert.Equal(t, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}, paths(es.Flatten()))

	// a mixed level: decimal indexes first, numerically, then the other keys as strings;
	// "01", "-1", "1.5" and "" are not indexes
	es = Errors{}
	for _, k := range []string{"b", "10", "", "01", "a", "-1", "2", "1.5", "100", "B"} {
		es[k] = ErrRequired
	}
	list := es.Flatten()
	keys := make([]string, 0, len(list))
	for _, fe := range list {
		keys = append(keys, fe.Path[0])
	}
	assert.Equal(t, []string{"2", "10", "100", "", "-1", "01", "1.5", "B", "a", "b"}, keys)

	// a level is sorted independently of the others
	list = Errors{
		"items": Errors{"10": Errors{"name": ErrRequired}, "2": Errors{"name": ErrRequired}},
		"id":    ErrRequired,
	}.Flatten()
	assert.Equal(t, []string{"id", "items.2.name", "items.10.name"}, paths(list))

	// Error() keeps comparing all keys as strings
	assert.Equal(t, "10: cannot be blank; 2: cannot be blank.",
		Errors{"10": ErrRequired, "2": ErrRequired}.Error())
}

func TestErrors_Flatten_PathOwnership(t *testing.T) {
	// Siblings under a prefix with spare capacity: appending to the prefix would make
	// them share one backing array, and every Path would end with the last key.
	prefix := make([]string, 1, 8)
	prefix[0] = "items"
	var list []FieldError
	Errors{"0": ErrRequired, "1": ErrRequired, "2": ErrRequired}.flatten(prefix, &list)
	assert.Equal(t, []string{"items.0", "items.1", "items.2"}, paths(list))
	assert.Equal(t, []string{"items"}, prefix, "the prefix is not modified")

	// the same through nested Errors
	inner := Errors{}
	for _, k := range []string{"a", "b", "c", "d", "e"} {
		inner[k] = Errors{"x": ErrRequired, "y": ErrRequired}
	}
	list = Errors{"p": Errors{"q": inner}}.Flatten()
	assert.Equal(t, []string{
		"p.q.a.x", "p.q.a.y", "p.q.b.x", "p.q.b.y", "p.q.c.x",
		"p.q.c.y", "p.q.d.x", "p.q.d.y", "p.q.e.x", "p.q.e.y",
	}, paths(list))

	// changing one Path leaves the others as they were
	list[0].Path[0] = "changed"
	list[0].Path[3] = "changed"
	assert.Equal(t, []string{"p", "q", "a", "y"}, list[1].Path)
	assert.Equal(t, "p.q.e.y", list[9].Field())
}

func TestErrors_Flatten_InternalError(t *testing.T) {
	internal := errors.New("database is down")
	failing := By(func(interface{}) error { return NewInternalError(internal) })

	// Each puts an element's InternalError into the Errors it returns
	err := Validate([]string{"a", "b"}, Each(failing))
	es, ok := err.(Errors)
	if !assert.True(t, ok) {
		return
	}
	list := es.Flatten()
	assert.Equal(t, []string{"0", "1"}, paths(list))
	for _, fe := range list {
		assert.Equal(t, "", fe.Code)
		var ie InternalError
		if assert.True(t, errors.As(fe.Err, &ie), fe.Field()) {
			assert.Equal(t, internal, ie.InternalError())
		}
	}

	// so does a slice of Validatable elements
	err = Validate([]Validatable{flattenInternal{internal}})
	es, ok = err.(Errors)
	if assert.True(t, ok) {
		list = es.Flatten()
		if assert.Len(t, list, 1) {
			assert.Equal(t, []string{"0"}, list[0].Path)
			var ie InternalError
			assert.True(t, errors.As(list[0].Err, &ie))
		}
	}

	// an InternalError is one entry with no code, whatever it wraps
	list = Errors{
		"a": NewInternalError(ErrRequired),
		"b": NewInternalError(Errors{"c": ErrRequired}),
		"d": fmt.Errorf("wrapped: %w", NewInternalError(ErrRequired)),
	}.Flatten()
	assert.Equal(t, []string{"a", "b", "d"}, paths(list))
	for _, fe := range list {
		assert.Equal(t, "", fe.Code, fe.Field())
		assert.Nil(t, fe.Params, fe.Field())
	}
}

type flattenInternal struct {
	err error
}

func (f flattenInternal) Validate() error {
	return NewInternalError(f.err)
}

func TestFlatten(t *testing.T) {
	assert.Nil(t, Flatten(nil))
	assert.Nil(t, Flatten(Errors{}))
	assert.Nil(t, Flatten(Errors{"a": nil}))

	// a single validation error: one entry with a nil Path
	list := Flatten(Validate("", Required))
	if assert.Len(t, list, 1) {
		assert.Nil(t, list[0].Path)
		assert.Equal(t, "", list[0].Field())
		assert.Equal(t, "", list[0].Pointer())
		assert.Equal(t, ErrRequired.Code(), list[0].Code)
		assert.Equal(t, "cannot be blank", list[0].Message)
	}

	// an Errors, also wrapped, is flattened
	assert.Equal(t, []string{"a", "b.c"}, paths(Flatten(Errors{"a": ErrRequired, "b": Errors{"c": ErrRequired}})))
	assert.Equal(t, []string{"a"}, paths(Flatten(fmt.Errorf("x: %w", Errors{"a": ErrRequired}))))

	// an InternalError, such as the one ValidateStruct returns for a non-pointer
	list = Flatten(ValidateStruct(flattenItem{}))
	if assert.Len(t, list, 1) {
		assert.Nil(t, list[0].Path)
		assert.Equal(t, "", list[0].Code)
		var ie InternalError
		assert.True(t, errors.As(list[0].Err, &ie))
		assert.True(t, errors.Is(list[0].Err, ErrStructPointer))
	}
	list = Flatten(NewInternalError(Errors{"a": ErrRequired}))
	if assert.Len(t, list, 1) {
		assert.Nil(t, list[0].Path)
		assert.Equal(t, "", list[0].Code)
	}

	// any other error
	sentinel := errors.New("sentinel")
	assert.Equal(t, []FieldError{{Message: "sentinel", Err: sentinel}}, Flatten(sentinel))
}

func TestFieldError_FieldAndPointer(t *testing.T) {
	tests := []struct {
		tag     string
		path    []string
		field   string
		pointer string
	}{
		{"no path", nil, "", ""},
		{"empty key", []string{""}, "", "/"},
		{"empty keys", []string{"", ""}, ".", "//"},
		{"plain", []string{"items", "0", "name"}, "items.0.name", "/items/0/name"},
		{"dot key", []string{"."}, ".", "/."},
		{"dot", []string{"a.b", "c"}, "a.b.c", "/a.b/c"},
		{"slash key", []string{"/"}, "/", "/~1"},
		{"slash", []string{"a/b"}, "a/b", "/a~1b"},
		{"tilde key", []string{"~"}, "~", "/~0"},
		{"tilde", []string{"m~n"}, "m~n", "/m~0n"},
		// "~1" in a key is "~" followed by "1": it must not decode back to "/"
		{"tilde one", []string{"~1"}, "~1", "/~01"},
		{"both", []string{"~/", "/~"}, "~/./~", "/~0~1/~1~0"},
	}
	for _, test := range tests {
		fe := FieldError{Path: test.path}
		assert.Equal(t, test.field, fe.Field(), test.tag)
		assert.Equal(t, test.pointer, fe.Pointer(), test.tag)
	}

	list := Errors{"a.b": Errors{"c/d": Errors{"e~f": ErrRequired}}}.Flatten()
	if assert.Len(t, list, 1) {
		assert.Equal(t, []string{"a.b", "c/d", "e~f"}, list[0].Path)
		assert.Equal(t, "a.b.c/d.e~f", list[0].Field())
		assert.Equal(t, "/a.b/c~1d/e~0f", list[0].Pointer())
	}
}

func TestFieldError_JSON(t *testing.T) {
	list := Errors{
		"id":   Validate("ab", Length(3, 10)),
		"name": NewInternalError(errors.New("secret")),
	}.Flatten()
	b, err := json.Marshal(list)
	assert.NoError(t, err)
	assert.Equal(t, `[`+
		`{"path":["id"],"code":"validation_length_out_of_range","message":"the length must be between 3 and 10","params":{"max":10,"min":3}},`+
		`{"path":["name"],"code":"","message":"secret"}`+
		`]`, string(b))

	// no path for a single error, and never the error itself
	b, err = json.Marshal(Flatten(Validate("", Required)))
	assert.NoError(t, err)
	assert.Equal(t, `[{"code":"validation_required","message":"cannot be blank"}]`, string(b))
}
