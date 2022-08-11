/*
Package validator is used to validate strings, slices, and maps.

# Rules

TODO:

- How to pass rules
- How to pass rules for keys and values (with maps and slices)

# Supported types and rules

These are the supported variable types that can be passed to [Validate] and [ValidateAny], and the rules that are available to them.

## `string`

When passing a value of type `string`, validator performs a set of operations to sanitize the value:

- All leading and trailing whitespace characters are removed, including: spaces, newlines, tabs, and all other characters defined as whitespace by Unicode.
- All whitespace characters–including spaces, newlines, tabs, and all other characters defined as whitespace by Unicode–are replaced with a regular space, and consecutive whitespace characters are collapsed into one. This is the default behavior but can be disabled with the `preserve-whitespace` rule.
- All control characters are removed from the string. This includes almost all characters defined as control characters by Unicode, except tabs and newlines, which are converted to spaces (unless `preserve-whitespace` is set), and the Zero-Width Joiner (ZWJ) character, which is commonly used with emojis.
- The string is normalized to Unicode form NFC (Canonical Composition). [(More info about Unicode normalization)](https://withblue.ink/2019/03/11/why-you-need-to-normalize-unicode-strings.html)

### Optional rules

- **`min=int`**: minimum length–returns an error if the string is shorter.
- **`max=int`**: maximum length–returns an error if the string is longer.
- **`preserve-whitespace`**: boolean flag that preserves all whitespace characters as-is (does not collapse whitespace characters and does not convert Unicode spaces to regular spaces).
- **`preserve-newlines`**: boolean flag that preserves all newlines even when `preserve-whitespace` is not set (note that newlines are still trimmed from the ends of the string).
- **`replace-whitespaces`**: boolean flag that replaces all whitespace characters with an underscore.

## `[]string`

When passing a value of type `[]string` (a slice of strings), validator sanitizes each string value first, using the string validator.

### Optional rules

- **`min=int`**: minimum length–returns an error if the slice's length (number of elements) is smaller than this.
- **`max=int`**: maximum length–returns an error if the slice's length (number of elements) is bigger than this.
- **`sort`**: boolean flag that makes the result sorted alphabetically.
- **`unique`**: boolean flag that removes duplicates in the result (after sorting the values).
- **`value=(rule)`**: rule for validating each value of the slice (see rules for the string validator).

## `map[string]string`

When passing a value of type `map[string]string` (a map where both keys and values are strings), validator sanitizes each key and value first, using the string validator on both.

### Optional rules

These rules apply to the map validator:

- **`min=int`**: minimum length–returns an error if the map's length (number of elements) is smaller than this.
- **`max=int`**: maximum length–returns an error if the map's length (number of elements) is bigger than this.
- **`key=(rule)`**: rule for validating each key of the map (see rules for the string validator).
- **`value=(rule)`**: rule for validating each value of the map (see rules for the string validator).
*/
package validator

// This file does not contain code, and it's only used for docs
