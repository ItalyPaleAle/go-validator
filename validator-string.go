package validator

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cast"
	"golang.org/x/text/unicode/norm"
)

// stringValidator is a validator for type `string`
// Supported parameters:
// - `min=int` minimum length
// - `max=int` maximum length
// - `preserve-whitespace` boolean flag that preserves all whitespaces (does not collapse whitespaces and does not convert Unicode spaces to regular spaces)
// - `preserve-newlines` boolean flag that preserves all newlines even when `preserve-whitespace` is not set (note that newlines are still trimmed from the ends of the string)
// - `replace-whitespaces` boolean flag that replaces all whitespaces with an underscore
func stringValidator(rule string) validator[string] {
	// Parse rule
	params, err := parseParams(rule)
	if err != nil {
		return errorValidateFunc[string](err)
	}

	// Parse parameters
	min := -1
	if v, ok := params["min"]; ok && v != "" {
		min, err = cast.ToIntE(v)
		if err != nil {
			return errorValidateFunc[string](fmt.Errorf("parameter 'min' is invalid: failed to cast to int: %v", err))
		}
		if min < 1 {
			return errorValidateFunc[string](errors.New("parameter 'min' must be greater than 0"))
		}
	}
	max := -1
	if v, ok := params["max"]; ok && v != "" {
		max, err = cast.ToIntE(v)
		if err != nil {
			return errorValidateFunc[string](fmt.Errorf("parameter 'max' is invalid: failed to cast to int: %v", err))
		}
		if max < 1 {
			return errorValidateFunc[string](errors.New("parameter 'max' must be greater than 0"))
		}
	}
	if max > 0 && min > max {
		return errorValidateFunc[string](errors.New("parameter 'max' must not be smaller than parameter 'min'"))
	}
	preserveWhitespace := false
	if _, ok := params["preserve-whitespace"]; ok {
		// Boolean option, with no value
		preserveWhitespace = true
	}
	preserveNewlines := false
	if _, ok := params["preserve-newlines"]; ok {
		// Boolean option, with no value
		preserveNewlines = true
	}
	replaceWhitespaces := false
	if _, ok := params["replace-whitespaces"]; ok {
		// Boolean option, with no value
		replaceWhitespaces = true
	}

	return func(ctx context.Context, val string) (res string, err error) {
		// Normalize to form NFC
		val = norm.NFC.String(val)

		// Trim whitespaces from each end (Unicode-aware)
		// Note that this also trims newlines from both ends, regardless of preserveNewLines
		val = strings.TrimSpace(val)

		// Clean the string
		val = cleanStringInternal(val, preserveNewlines, replaceWhitespaces, preserveWhitespace)

		// Check if we have rules
		if min > 0 && len(val) < min {
			return "", fmt.Errorf("value is shorter than %d", min)
		}
		if max > 0 && len(val) > max {
			return "", fmt.Errorf("value is longer than %d", max)
		}

		return val, nil
	}
}

// Iterate through the string to strip control characters
// If needed, also collapse whitespaces and/or replace whitespaces
func cleanStringInternal(val string, preserveNewlines, replaceWhitespaces, preserveWhitespace bool) string {
	var (
		r         rune
		n         int
		a         int
		lastSpace bool
	)
	out := make([]byte, len(val))
	for i, w := 0, 0; i < len(val); i += w {
		// Get the rune
		r, w = utf8.DecodeRuneInString(val[i:])

		// Remove control characters, but preserve these characters:
		// - tabs (0x09)
		// - newlines (0x0A)
		// - Zero-Width Joiner (ZWJ), which is used by emojis (U+200D)
		if r != 0x09 && r != 0x0A && r != 0x200D && unicode.Is(unicode.C, r) {
			continue
		}

		// Add runes that are not spaces right away
		if !unicode.IsSpace(r) {
			lastSpace = false
			a = utf8.EncodeRune(out[n:], r)
			n += a
			continue
		}
		// If preserving newlines, keep those too
		if preserveNewlines && r == '\n' {
			lastSpace = true
			out[n] = '\n'
			n++
		}

		// Collapse consecutive whitespaces
		if lastSpace && !preserveWhitespace {
			continue
		}

		lastSpace = true
		if replaceWhitespaces {
			// Replace with an underscore
			out[n] = '_'
			n++
		} else if !preserveWhitespace {
			// Replace with a regular space
			out[n] = ' '
			n++
		} else {
			// Add the rune as-is
			a = utf8.EncodeRune(out[n:], r)
			n += a
		}
	}

	return string(out[:n])
}
