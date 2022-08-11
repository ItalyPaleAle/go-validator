package validator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// sliceValidator returns a validator for type `[]T`
func sliceValidator[T any](rule string) validator[[]T] {
	var zero T
	errFunc := errorValidateFunc[[]T]

	// Parse rule
	params, err := parseParams(rule)
	if err != nil {
		return errFunc(err)
	}

	// Parse parameters
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
	sortFlag := false
	if _, ok := params["sort"]; ok {
		// Boolean option, with no value
		sortFlag = true
	}
	uniqueFlag := false
	if _, ok := params["unique"]; ok {
		// Boolean option, with no value
		uniqueFlag = true
	}

	// Validator function for each value, as well as sort and unique functions
	var (
		valueValidator        validator[T]
		valueSorter           func([]T)     = nil
		valueDuplicateRemover func([]T) []T = nil
		fp                    reflect.Value
	)
	switch any(zero).(type) {
	case string:
		f := stringValidator(params["value"])
		fp = reflect.ValueOf(&valueValidator).Elem()
		fp.Set(reflect.Indirect(reflect.ValueOf(f)))

		if sortFlag || uniqueFlag {
			fp = reflect.ValueOf(&valueSorter).Elem()
			fp.Set(reflect.Indirect(reflect.ValueOf(SortSlice[string])))
		}
		if uniqueFlag {
			fp = reflect.ValueOf(&valueDuplicateRemover).Elem()
			fp.Set(reflect.Indirect(reflect.ValueOf(RemoveDuplicatesInSortedSlice[string])))
		}
	default:
		return errFunc(fmt.Errorf("type of value '%T' is not supported", zero))
	}

	return func(list []T) (res []T, err error) {
		// Check if we have rules
		if min > 0 && len(list) < min {
			return nil, fmt.Errorf("value is shorter than %d", min)
		}
		if max > 0 && len(list) > max {
			return nil, fmt.Errorf("value is longer than %d", max)
		}

		// Validate each item
		for i := 0; i < len(list); i++ {
			list[i], err = valueValidator(list[i])
			if err != nil {
				return nil, err
			}
		}

		// Sort if needed
		if valueSorter != nil {
			valueSorter(list)
		}

		// Unique values if needed
		if valueDuplicateRemover != nil {
			list = valueDuplicateRemover(list)
		}

		return list, nil
	}
}
