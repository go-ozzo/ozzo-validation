package validation

import (
	"errors"
	"sort"
	"strings"
)

// FieldError is one error of an Errors, with the path of keys that leads to it.
// Flatten and Errors.Flatten return them.
//
// Message is the Error() text of the error. For an InternalError and any other error that
// is not a validation Error, Code is empty: check Err before exposing Message to clients.
//
//	for _, fe := range validation.Flatten(err) {
//		var ie validation.InternalError
//		if errors.As(fe.Err, &ie) {
//			// log it and answer with a server error: its message is not for clients
//			continue
//		}
//		// report fe.Pointer(), fe.Code and fe.Message
//	}
type FieldError struct {
	// Path holds the keys from the outermost Errors to the error: struct field names
	// (as chosen by ErrorTag), map keys and slice indexes, e.g. ["items", "0", "name"].
	// It is nil for an error that is not in an Errors (see Flatten).
	Path []string `json:"path,omitempty"`
	// Code is the code of a validation Error, or "" for any other error.
	Code string `json:"code"`
	// Message is the rendered text of the error, err.Error(), with the params of a
	// validation Error applied. Error.Message() is the template it was rendered from.
	Message string `json:"message"`
	// Params are the params of a validation Error, or nil for any other error. It is the
	// error's own map, not a copy: do not modify it.
	Params map[string]interface{} `json:"params,omitempty"`
	// Err is the error itself, so that callers can recover their own error types and tell
	// an InternalError apart. It is never marshaled.
	Err error `json:"-"`
}

// Field returns the path joined with dots, e.g. "items.0.name". It is ambiguous when a key
// contains a dot; use Pointer or Path when that matters.
func (fe FieldError) Field() string {
	return strings.Join(fe.Path, ".")
}

// pointerEscaper escapes a JSON Pointer reference token (RFC 6901, section 3).
var pointerEscaper = strings.NewReplacer("~", "~0", "/", "~1")

// Pointer returns the path as an RFC 6901 JSON Pointer, e.g. "/items/0/name", with "~"
// escaped as "~0" and "/" as "~1". It returns "" (the whole document) for a nil Path.
func (fe FieldError) Pointer() string {
	var s strings.Builder
	for _, key := range fe.Path {
		s.WriteByte('/')
		s.WriteString(pointerEscaper.Replace(key))
	}
	return s.String()
}

// Flatten returns err as a flat list of field errors. It returns nil for a nil error, and
// Errors.Flatten for an Errors, also a wrapped one. Any other error, such as the single
// error that Validate returns, is one FieldError with a nil Path.
func Flatten(err error) []FieldError {
	if err == nil {
		return nil
	}
	if es, ok := nestedErrors(err); ok {
		return es.Flatten()
	}
	return []FieldError{newFieldError(nil, err)}
}

// Flatten returns the errors in es as a flat list, walking into the Errors nested for
// nested structs, maps and slices. Each FieldError carries the path of keys to its error
// and owns its Path slice.
//
// The keys of each level are sorted, so equal Errors give equal lists. Keys that are
// decimal indexes ("0", or digits without a leading zero) come first, in numeric order,
// so "items.2" comes before "items.10"; the other keys follow in string order. This is
// not the order of Error(), which compares all keys as strings.
//
// Nil entries are skipped, and an empty Errors, nested or not, contributes nothing: an
// Errors without errors returns nil. An entry that is or wraps an Errors, such as
// fmt.Errorf("profile: %w", errs), is walked into; the wrapping text is dropped. An
// Errors that contains itself is not supported, as in Error(). Flatten does not modify es.
//
// Code and Params come from the validation Error found in an entry with errors.As, so a
// wrapped validation Error keeps its code. Other errors, such as one returned by a By
// rule, get an empty Code and nil Params. An entry that is or wraps an InternalError is
// one FieldError with an empty Code and nil Params, whatever it wraps: Validate and
// ValidateStruct return an InternalError instead of an Errors, but Each and the elements
// of a slice or map of Validatable put one into the map. Check Err with errors.As before
// exposing Message to clients.
func (es Errors) Flatten() []FieldError {
	var list []FieldError
	es.flatten(nil, &list)
	return list
}

func (es Errors) flatten(path []string, list *[]FieldError) {
	keys := make([]string, 0, len(es))
	for key, err := range es {
		if err != nil {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keyLess(keys[i], keys[j]) })

	for _, key := range keys {
		// A new slice per entry: appending to the shared path could overwrite a sibling's.
		p := make([]string, len(path)+1)
		copy(p, path)
		p[len(path)] = key

		err := es[key]
		if nested, ok := nestedErrors(err); ok {
			nested.flatten(p, list)
			continue
		}
		*list = append(*list, newFieldError(p, err))
	}
}

// nestedErrors returns the Errors that err is or wraps, unless err is or wraps an
// InternalError: that is reported as one error, never walked into.
func nestedErrors(err error) (Errors, bool) {
	var ie InternalError
	if errors.As(err, &ie) {
		return nil, false
	}
	var es Errors
	ok := errors.As(err, &es)
	return es, ok
}

func newFieldError(path []string, err error) FieldError {
	fe := FieldError{Path: path, Message: err.Error(), Err: err}
	var ie InternalError
	if errors.As(err, &ie) {
		return fe
	}
	var ve Error
	if errors.As(err, &ve) {
		fe.Code = ve.Code()
		fe.Params = ve.Params()
	}
	return fe
}

// keyLess orders the keys of one level of Errors: decimal indexes first, numerically,
// then the other keys as strings.
func keyLess(a, b string) bool {
	ai, bi := isIndexKey(a), isIndexKey(b)
	switch {
	case ai && bi:
		// Without a leading zero, the shorter number is the smaller one; this needs no
		// parsing and does not overflow.
		if len(a) != len(b) {
			return len(a) < len(b)
		}
		return a < b
	case ai != bi:
		return ai
	default:
		return a < b
	}
}

// isIndexKey reports whether key is a decimal index as strconv.Itoa writes it for a
// non-negative int: "0", or digits without a leading zero.
func isIndexKey(key string) bool {
	if key == "" || key[0] == '0' && len(key) > 1 {
		return false
	}
	for i := 0; i < len(key); i++ {
		if key[i] < '0' || key[i] > '9' {
			return false
		}
	}
	return true
}
