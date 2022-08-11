package validator

import (
	"errors"
)

var ruleSyntaxError = errors.New("invalid rule string: syntax error")

func parseParams(rule string) (params map[string]string, err error) {
	l := len(rule)
	if l == 0 {
		return map[string]string{}, nil
	}

	var start, in, lastClosedParam, end int
	var key, val string
	for i := 0; i < l; i++ {
		// Skip characters that are part of a multi-byte sequence
		if rule[i] > 127 {
			continue
		}

		if rule[i] == '(' {
			if in == 0 {
				start++
			}
			in++
		} else if rule[i] == ')' {
			in--
			if in < 0 {
				return nil, ruleSyntaxError
			}
			lastClosedParam = i
		} else if in == 0 {
			if rule[i] == ',' {
				if params == nil {
					params = map[string]string{}
				}
				end = i
				if lastClosedParam == i-1 {
					end--
				}
				val = rule[start:end]
				// If there's no key, the we only have a key with no value
				if key == "" {
					if start == end {
						return nil, ruleSyntaxError
					}
					params[val] = ""
				} else {
					params[key] = val
					key = ""
				}
				start = i + 1
			} else if rule[i] == '=' {
				end = i
				if lastClosedParam == i-1 {
					end--
				}
				if start == end {
					return nil, ruleSyntaxError
				}
				key = rule[start:end]
				start = i + 1
			}
		}
	}
	if in != 0 {
		return nil, ruleSyntaxError
	}
	if start != l {
		if params == nil {
			params = map[string]string{}
		}
		end = l
		if lastClosedParam == l-1 {
			end--
		}
		if start == end {
			return nil, ruleSyntaxError
		}
		val = rule[start:end]
		if key == "" {
			params[val] = ""
		} else {
			params[key] = val
		}
	}

	return params, nil
}
