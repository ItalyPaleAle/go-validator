package validator

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// keyvalueValidator is a validator for type `map[string]T`
// Supported parameters:
// - `min=int` minimum length
// - `max=int` maximum length
// - `key=rule`: rule for validating map keys
// - `value=rule`: rule for validating map values
func keyvalueValidator[T any](rule string) validator[map[string]T] {
	var zero T
	errFunc := errorValidateFunc[map[string]T]

	// Parse rule
	params, err := parseParams(rule)
	if err != nil {
		return errFunc(err)
	}

	// Rules from parameters
	min := -1
	if v, ok := params["min"]; ok && v != "" {
		min, err = cast.ToIntE(v)
		if err != nil {
			return errFunc(fmt.Errorf("parameter 'min' is invalid: failed to cast to int: %v", err))
		}
		if min < 1 {
			return errFunc(errors.New("parameter 'min' must be greater than 0"))
		}
	}
	max := -1
	if v, ok := params["max"]; ok && v != "" {
		max, err = cast.ToIntE(v)
		if err != nil {
			return errFunc(fmt.Errorf("parameter 'max' is invalid: failed to cast to int: %v", err))
		}
		if max < 1 {
			return errFunc(errors.New("parameter 'max' must be greater than 0"))
		}
	}
	if max > 0 && min > max {
		return errFunc(errors.New("parameter 'max' must not be smaller than parameter 'min'"))
	}

	// Validator function for each key
	keyValidator := stringValidator(params["key"])

	// Validator function for each value
	var valueValidator validator[T]
	var fp reflect.Value

	switch any(zero).(type) {
	case string:
		f := stringValidator(params["value"])
		fp = reflect.ValueOf(&valueValidator).Elem()
		fp.Set(reflect.Indirect(reflect.ValueOf(f)))
	}

	return func(ctx context.Context, val map[string]T) (map[string]T, error) {
		// Check if we have rules
		if min > 0 && len(val) < min {
			return nil, fmt.Errorf("value is shorter than %d", min)
		}
		if max > 0 && len(val) > max {
			return nil, fmt.Errorf("value is longer than %d", max)
		}

		// Validate each item
		res := make(map[string]T, len(val))
		for k, v := range val {
			k, err = keyValidator(ctx, k)
			if err != nil {
				return nil, err
			}
			v, err = valueValidator(ctx, v)
			if err != nil {
				return nil, err
			}
			res[k] = v
		}

		return res, nil
	}
}
