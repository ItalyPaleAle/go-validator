package validator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var validators sync.Map

// Validate and sanitize a value
// val is generally a pointer, and it returns a pointer
func Validate(ctx context.Context, val any, rule string) (res any, err error) {
	rule = strings.TrimSpace(rule)

	// Get the type of the value
	t := reflect.TypeOf(val)
	isPtr := false
	if t.Kind() == reflect.Pointer {
		isPtr = true
		v := reflect.ValueOf(val)
		t = t.Elem()
		if v.IsZero() {
			val = nil
		} else {
			val = v.Elem().Interface()
		}
	}

	// Get the validator from the cache
	cacheKey := t.String() + "|" + rule
	f, _ := validators.Load(cacheKey)

	// Work with the type of the value
	var ok bool
	switch t.String() {
	case "string":
		var (
			x  string
			fT validator[string]
		)
		if val != nil {
			x, ok = val.(string)
			if !ok {
				return nil, errors.New("failed to assert value as string")
			}
		}
		if f != nil {
			if fT, ok = f.(validator[string]); ok {
				x, err = fT(ctx, x)
				if err != nil {
					return nil, err
				}
				if isPtr {
					return &x, nil
				} else {
					return x, nil
				}
			}
		}
		fT = stringValidator(rule)
		defer validators.Store(cacheKey, fT)
		x, err = fT(ctx, x)
		if err != nil {
			return nil, err
		}
		if isPtr {
			return &x, nil
		} else {
			return x, nil
		}
	case "[]string":
		var (
			x  []string
			fT validator[[]string]
		)
		if val != nil {
			x, ok = val.([]string)
			if !ok {
				return nil, errors.New("failed to assert value as string")
			}
		}
		if f != nil {
			if fT, ok = f.(validator[[]string]); ok {
				x, err = fT(ctx, x)
				if err != nil {
					return nil, err
				}
				if isPtr {
					return &x, nil
				} else {
					return x, nil
				}
			}
		}
		fT = listValidator[string](rule)
		defer validators.Store(cacheKey, fT)
		x, err = fT(ctx, x)
		if err != nil {
			return nil, err
		}
		if isPtr {
			return &x, nil
		} else {
			return x, nil
		}
	case "map[string]string":
		var (
			x  map[string]string
			fT validator[map[string]string]
		)
		if val != nil {
			x, ok = val.(map[string]string)
			if !ok {
				return nil, errors.New("failed to assert value as string")
			}
		}
		if f != nil {
			if fT, ok = f.(validator[map[string]string]); ok {
				x, err = fT(ctx, x)
				if err != nil {
					return nil, err
				}
				if isPtr {
					return &x, nil
				} else {
					return x, nil
				}
			}
		}
		fT = keyvalueValidator[string](rule)
		defer validators.Store(cacheKey, fT)
		x, err = fT(ctx, x)
		if err != nil {
			return nil, err
		}
		if isPtr {
			return &x, nil
		} else {
			return x, nil
		}
	}

	return nil, fmt.Errorf("cannot find a validator for type %T", val)
}

// ValidateString is a short-hand function to validate a string
func ValidateString(val string, rule string) (string, error) {
	res, err := Validate(context.Background(), val, rule)
	if err != nil {
		return "", err
	}

	resStr, ok := res.(string)
	if !ok {
		return "", errors.New("returned value from Validate is not a string")
	}

	return resStr, nil
}

// validator is the type of a validator function
type validator[T any] func(ctx context.Context, val T) (res T, err error)

// errorValidateFunc returns a validator function that returns an error
func errorValidateFunc[T any](err error) validator[T] {
	return func(ctx context.Context, val T) (T, error) {
		var zero T
		return zero, err
	}
}
