package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type validateTypes interface {
	string | map[string]string | []string
}

var validators sync.Map

// Validate and sanitize a value, using generics to define the supported types
// Rule follows the format for the given type
func Validate[T validateTypes](val T, rule string) (res T, err error) {
	var zero T
	if reflect.ValueOf(val).IsZero() {
		return zero, nil
	}

	rule = strings.TrimSpace(rule)

	var ok bool
	switch x := any(val).(type) {
	case string:
		cacheKey := "string|" + rule
		f, _ := validators.Load(cacheKey)
		var fT validator[string]
		if f != nil {
			fT, ok = f.(validator[string])
			if !ok {
				fT = nil
			}
		}
		if fT == nil {
			fT = stringValidator(rule)
			defer validators.Store(cacheKey, fT)
		}
		x, err = fT(x)
		if err != nil {
			return zero, err
		}
		return any(x).(T), nil
	case []string:
		if len(x) == 0 {
			return val, nil
		}
		cacheKey := "[]string|" + rule
		f, _ := validators.Load(cacheKey)
		var fT validator[[]string]
		if f != nil {
			fT, ok = f.(validator[[]string])
			if !ok {
				fT = nil
			}
		}
		if fT == nil {
			fT = sliceValidator[string](rule)
			defer validators.Store(cacheKey, fT)
		}
		x, err = fT(x)
		if err != nil {
			return zero, err
		}
		return any(x).(T), nil
	case map[string]string:
		if len(x) == 0 {
			return val, nil
		}
		cacheKey := "map[string]string|" + rule
		f, _ := validators.Load(cacheKey)
		var fT validator[map[string]string]
		if f != nil {
			fT, ok = f.(validator[map[string]string])
			if !ok {
				fT = nil
			}
		}
		if fT == nil {
			fT = mapValidator[string](rule)
			defer validators.Store(cacheKey, fT)
		}
		x, err = fT(x)
		if err != nil {
			return zero, err
		}
		return any(x).(T), nil
	default:
		return zero, fmt.Errorf("cannot find a validator for type %T", val)
	}
}

// ValidateAny validates and sanitizes a value with type any
// Supported types are: `string`, `map[string]string`, `[]string`, and pointers to those types
// Rule follows the format for the given type
func ValidateAny(val any, rule string) (res any, err error) {
	if val == nil {
		return nil, nil
	}

	// Get the type of the value
	isPtr := false
	if reflect.TypeOf(val).Kind() == reflect.Pointer {
		isPtr = true
		v := reflect.ValueOf(val)
		if v.IsZero() {
			return val, nil
		}
		val = v.Elem().Interface()
	}

	// Switch based on the type of the value
	switch x := val.(type) {
	case string:
		x, err = Validate(x, rule)
		if err != nil {
			return nil, err
		}
		if isPtr {
			return &x, nil
		}
		return x, nil
	case []string:
		x, err = Validate(x, rule)
		if err != nil {
			return nil, err
		}
		if isPtr {
			return &x, nil
		}
		return x, nil
	case map[string]string:
		x, err = Validate(x, rule)
		if err != nil {
			return nil, err
		}
		if isPtr {
			return &x, nil
		}
		return x, nil
	default:
		return nil, fmt.Errorf("cannot find a validator for type %T", val)
	}
}

// validator is the type of a validator function
type validator[T any] func(val T) (res T, err error)

// errorValidateFunc returns a validator function that returns an error
func errorValidateFunc[T any](err error) validator[T] {
	return func(val T) (T, error) {
		var zero T
		return zero, err
	}
}
