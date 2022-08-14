# go-validator

[![Go Reference](https://pkg.go.dev/badge/github.com/italypaleale/go-validator.svg)](https://pkg.go.dev/github.com/italypaleale/go-validator) [![Continuous Integration](https://github.com/ItalyPaleAle/go-validator/actions/workflows/ci.yaml/badge.svg)](https://github.com/ItalyPaleAle/go-validator/actions/workflows/ci.yaml)

This package is a Go library for sanitizing and validating strings, maps, and slices.

Features:

- Can be used as a standalone library or with GraphQL directives.
- Supports values of type `string`, `[]string` and `map[string]string`.
- Many rules to control the sanitizer's behavior and add validation rules.
- Designed to work with Unicode.

# Using validator

## Import the package

Add the package:

```sh
go get -u github.com/italypaleale/go-validator
```

Then import it into your Go files:

```go
import (
	validator "github.com/italypaleale/go-validator"
)
```

## Validating objects

Validator currently support 3 types of variables:

- `string`
- `map[string]string`
- `[]string`

If your variable is already of one of those known types, you can use the [`Validate`](https://pkg.go.dev/github.com/italypaleale/go-validator#Validate) method:

```go
// Validate(val T, rule string) (res T, err error)
cleanedVal, err := validator.Validate(myVal, rules)
```

Otherwise, you can pass a variable of type `any` (i.e. `interface{}`) to the [`ValidateAny`](https://pkg.go.dev/github.com/italypaleale/go-validator#ValidateAny) method:

```go
// ValidateAny(val any, rule string) (res any, err error)
cleanedAny, err := validator.ValidateAny(myAny, rules)
```

## Using with GraphQL directives

Validator has been designed to work with GraphQL directives too. It's currently tested with [`99designs/gqlgen`](https://github.com/99designs/gqlgen).

> [Official documentation](https://gqlgen.com/reference/directives/) for directives in `99designs/gqlgen`.

For example, you can create a GraphQL mutation that accepts strings, slices, or maps (through custom scalar types) and validate them. Example schema:

```graphql
directive @validate(
  rule: String!
) on INPUT_FIELD_DEFINITION | ARGUMENT_DEFINITION

type Mutation {
  createObject(
    name: String! @validate(rule: "max=200")
    tags: [String!] @validate(rule: "value=(min=1,max=60,asciionly),unique")
  ): Object
}
```

Enable the directive by a method like the one below to your gqlgen server's `DirectiveRoot`:

```go
import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	validator "github.com/italypaleale/go-validator"
)

func main() {
	resolvers := // ...

	c := generated.Config{
		Resolvers: resolvers,
		Directives: generated.DirectiveRoot{
			Validate: ValidateDirective,
		},
	}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	// ...
}

// ValidateDirective is the handler for the @validate directive, which validates and sanitizes an input or value
func ValidateDirective(ctx context.Context, obj interface{}, next graphql.Resolver, rule string) (res interface{}, err error) {
	// Get the value
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}

	// Validate the value
	// This uses a cache for validator functions
	val, err = validator.ValidateAny(val, rule)
	if err != nil {
		return nil, err
	}

	return val, nil
}
```

# Rules

Rules allow configuring the behavior of the validator functions. They are passed as a string as the last parameter for [`Validate`](https://pkg.go.dev/github.com/italypaleale/go-validator#Validate) and [`ValidateAny`](https://pkg.go.dev/github.com/italypaleale/go-validator#ValidateAny).

Example of rule (for the string validator):

```go
cleanedVal, err := validator.Validate(myVal, "min=3,preserve-newlines")
```

> You can find all supported rules in the [Supported types and rules](#supported-types-and-rules) section below.

Rules are comma-separated. They can have a value (such as `min=3`, indicating the minimum length required), or they can be boolean (such as `preserve-newlines`), whose presence alone is enough to enable the rule.

With slices and maps rules can be used to control how values (and in the case of maps, keys too) are validated too. For example, while validating a `[]string`:

```text
value=(min=3,preserve-newlines),min=2
```

The rule above requires all values to comply with `min=3,preserve-newlines`. It additionally requires the slice itself to have at least 2 elements.

# Supported types and rules

These are the supported variable types that can be passed to [`Validate`](https://pkg.go.dev/github.com/italypaleale/go-validator#Validate) and [`ValidateAny`](https://pkg.go.dev/github.com/italypaleale/go-validator#ValidateAny), and the rules that are available to them.

## `string`

When passing a value of type `string`, validator performs a set of operations to sanitize the value:

- All leading and trailing whitespace characters are removed, including: spaces, newlines, tabs, and all other characters defined as whitespace by Unicode.
- All whitespace characters–including spaces, newlines, tabs, and all other characters defined as whitespace by Unicode–are replaced with a regular space, and consecutive whitespace characters are collapsed into one. This is the default behavior but can be disabled with the `preserve-whitespace` rule.
- All control characters are removed from the string. This includes almost all characters defined as control characters by Unicode, except tabs and newlines, which are converted to spaces (unless `preserve-whitespace` is set), and the Zero-Width Joiner (ZWJ) character, which is commonly used with emojis.
- The string is normalized to Unicode form NFC (Canonical Composition); other forms can be selected with the `unorm` option. (More info about [Unicode normalization](https://withblue.ink/2019/03/11/why-you-need-to-normalize-unicode-strings.html))

### Optional rules

- **`min=int`**: minimum length–returns an error if the string is shorter.
- **`max=int`**: maximum length–returns an error if the string is longer.
- **`preserve-whitespace`**: boolean flag that preserves all whitespace characters as-is (does not collapse whitespace characters and does not convert Unicode spaces to regular spaces).
- **`preserve-newlines`**: boolean flag that preserves all newlines even when `preserve-whitespace` is not set (note that newlines are still trimmed from the ends of the string).
- **`replace-whitespaces`**: boolean flag that replaces all whitespace characters with an underscore.
- **`asciionly`**: boolean flag that removes all non-ASCII characters from the string. Note: this is executed after normalizing the string.
- **`unorm=string`**: Unicode normalization form to use. Possible values: `nfc` (default), `nfd`, `nfkc`, `nfkd`.

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
