package validation

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrNotMap is the error that the value being validated is not a map.
	ErrNotMap = errors.New("only a map can be validated")
)

type (
	// ErrKeyWrongType is the error that a key is the wrong type.
	ErrKeyWrongType string

	// ErrKeyNotFound is the error that a key cannot be found in the map.
	ErrKeyNotFound string

	// KeyRules represents a rule set associated with a map key.
	KeyRules struct {
		key   interface{}
		rules []Rule
	}
)

// Error returns the error string of ErrKeyWrongType.
func (e ErrKeyWrongType) Error() string {
	return fmt.Sprintf("key %q is the wrong type", string(e))
}

// Error returns the error string of ErrKeyNotFound.
func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf("key %q cannot be found in the map", string(e))
}

// ValidateMap validates a map by checking the specified map keys against the corresponding validation rules.
// Note that if the value is nil, it is considered valid.
// Use Key() to specify map keys that need to be validated. Each Key() call specifies a single key. A key can
// be associated with multiple rules.
// For example,
//
//    value := map[string]string{
//        "Name":  "name",
//        "Value": "demo",
//    }
//    err := validation.ValidateMap(value,
//         validation.Key("Name", validation.Required),
//         validation.Key("Value", validation.Required, validation.Length(5, 10)),
//    )
//    fmt.Println(err)
//    // Value: the length must be between 5 and 10.
//
// An error will be returned if validation fails.
func ValidateMap(m interface{}, keys ...*KeyRules) error {
	return ValidateMapWithContext(nil, m, keys...)
}

// ValidateMapWithContext validates a map with the given context.
// The only difference between ValidateMapWithContext and ValidateMap is that the former will
// validate map keys with the provided context.
// Please refer to ValidateMap for the detailed instructions on how to use this function.
func ValidateMapWithContext(ctx context.Context, m interface{}, keys ...*KeyRules) error {
	value := reflect.ValueOf(m)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Map {
		// must be a map
		return NewInternalError(ErrNotMap)
	}
	if value.IsNil() {
		// treat a nil map as valid
		return nil
	}

	errs := Errors{}
	kt := value.Type().Key()

	for _, kr := range keys {
		kv := reflect.ValueOf(kr.key)
		if !kt.AssignableTo(kv.Type()) {
			return NewInternalError(ErrKeyWrongType(getErrorKeyName(kr.key)))
		}
		vv := value.MapIndex(kv)
		if !vv.IsValid() {
			return NewInternalError(ErrKeyNotFound(getErrorKeyName(kr.key)))
		}
		var err error
		if ctx == nil {
			err = Validate(vv.Interface(), kr.rules...)
		} else {
			err = ValidateWithContext(ctx, vv.Interface(), kr.rules...)
		}
		if err != nil {
			if ie, ok := err.(InternalError); ok && ie.InternalError() != nil {
				return err
			}
			errs[getErrorKeyName(kr.key)] = err
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Key specifies a map key and the corresponding validation rules.
func Key(key interface{}, rules ...Rule) *KeyRules {
	return &KeyRules{
		key:   key,
		rules: rules,
	}
}

// getErrorKeyName returns the name that should be used to represent the validation error of a map key.
func getErrorKeyName(key interface{}) string {
	return fmt.Sprintf("%v", key)
}

type MapRule struct {
	keys []*KeyRules
}

func (r *MapRule) Validate(value interface{}) error {
	return ValidateMap(value, r.keys...)
}

func (r *MapRule) ValidateWithContext(ctx context.Context, value interface{}) error {
	return ValidateMapWithContext(ctx, value, r.keys...)
}

func Map(keys ...*KeyRules) *MapRule {
	return &MapRule{keys}
}
