package validator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// stringValidator returns a validator for type `string`
func stringValidator(rule string) validator[string] {
	// Parse rule
	params, err := parseParams(rule)
	if err != nil {
		return errorValidateFunc[string](err)
	}

	// Parse parameters
	min := -1
	if v, ok := params["min"]; ok && v != "" {
		min, err = strconv.Atoi(v)
		if err != nil {
			return errorValidateFunc[string](fmt.Errorf("parameter 'min' is invalid: failed to cast to int: %v", err))
		}
		if min < 1 {
			return errorValidateFunc[string](errors.New("parameter 'min' must be greater than 0"))
		}
	}
	max := -1
	if v, ok := params["max"]; ok && v != "" {
		max, err = strconv.Atoi(v)
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
	asciiOnly := false
	if _, ok := params["asciionly"]; ok {
		// Boolean option, with no value
		asciiOnly = true
	}
	unorm := norm.NFC
	if unormParam, ok := params["unorm"]; ok {
		switch strings.ToLower(unormParam) {
		case "nfc":
			unorm = norm.NFC
		case "nfd":
			unorm = norm.NFD
		case "nfkc":
			unorm = norm.NFKC
		case "nfkd":
			unorm = norm.NFKD
		default:
			return errorValidateFunc[string](errors.New("parameter 'unorm' is invalid"))
		}
	}

	return func(val string) (res string, err error) {
		// Unicode normalization
		val = unorm.String(val)

		// Trim whitespaces from each end (Unicode-aware)
		// Note that this also trims newlines from both ends, regardless of preserveNewLines
		val = strings.TrimSpace(val)

		// Clean the string
		val = cleanStringInternal(val, cleanStringOpts{
			preserveNewlines:   preserveNewlines,
			replaceWhitespaces: replaceWhitespaces,
			preserveWhitespace: preserveWhitespace,
			asciiOnly:          asciiOnly,
		})

		// Trim whitespaces from each end again
		val = strings.TrimSpace(val)

		// Check if we have length rules
		if min > 0 && len(val) < min {
			return "", fmt.Errorf("value is shorter than %d", min)
		}
		if max > 0 && len(val) > max {
			return "", fmt.Errorf("value is longer than %d", max)
		}

		return val, nil
	}
}

type cleanStringOpts struct {
	preserveNewlines   bool
	replaceWhitespaces bool
	preserveWhitespace bool
	asciiOnly          bool
}

// Iterate through the string to strip control characters
// If needed, also collapse whitespaces and/or replace whitespaces
func cleanStringInternal(val string, opts cleanStringOpts) string {
	var (
		r         rune
		n         int
		a         int
		lastSpace bool
	)
	out := make([]byte, len(val))
	for i, w := 0, 0; i < len(val); i += w {
		if opts.asciiOnly {
			// Go byte-by-byte
			r = rune(val[i])
			w = 1

			// Remove characters that are > 127
			if r > 127 {
				continue
			}
		} else {
			// Get the UTF-8 rune
			r, w = utf8.DecodeRuneInString(val[i:])
		}

		// Remove control characters, but preserve these characters:
		// - tabs (0x09) (which are replaced to regular spaces if preserve-whitespace is not present)
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
		if opts.preserveNewlines && r == '\n' {
			lastSpace = true
			out[n] = '\n'
			n++
		}

		// Collapse consecutive whitespaces
		if lastSpace && !opts.preserveWhitespace {
			continue
		}

		lastSpace = true
		if opts.replaceWhitespaces {
			// Replace with an underscore
			out[n] = '_'
			n++
		} else if !opts.preserveWhitespace {
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
