package validator

import (
	"reflect"
	"testing"
)

func Test_sliceValidator(t *testing.T) {
	rules := []string{
		"",
		"min=2",
		"max=3",
		"min=2,max=3",
		"value=(min=2)",
		"sort",
		"unique",

		// Invalid rules
		"min=0",
		"max=-1",
		"min=5,max=1",
	}

	tests := []struct {
		name    string
		rule    string
		value   []string
		wantRes []string
		wantErr bool
	}{
		{name: "empty slice", rule: rules[0], value: []string{}, wantRes: []string{}},
		{name: "slice of strings", rule: rules[0], value: []string{"hello", "world"}, wantRes: []string{"hello", "world"}},
		{name: "trim spaces in strings", rule: rules[0], value: []string{" hi! ", "ciao mondo"}, wantRes: []string{"hi!", "ciao mondo"}},
		{name: "normalize strings", rule: rules[0], value: []string{"e\u0300", "\u00e8", "😀"}, wantRes: []string{"\u00e8", "\u00e8", "😀"}},
		{name: "min length ok", rule: rules[1], value: []string{"a", "b"}, wantRes: []string{"a", "b"}},
		{name: "min length fail", rule: rules[1], value: []string{"a"}, wantErr: true},
		{name: "min length empty fail", rule: rules[1], value: []string{}, wantErr: true},
		{name: "max length ok 1", rule: rules[2], value: []string{"a", "b"}, wantRes: []string{"a", "b"}},
		{name: "max length ok 2", rule: rules[2], value: []string{}, wantRes: []string{}},
		{name: "max length fail", rule: rules[2], value: []string{"a", "b", "c", "d"}, wantErr: true},
		{name: "min and max length ok", rule: rules[3], value: []string{"a", "b"}, wantRes: []string{"a", "b"}},
		{name: "min and max length fail 1", rule: rules[3], value: []string{"a"}, wantErr: true},
		{name: "min and max length fail 2", rule: rules[3], value: []string{"a", "b", "c", "d"}, wantErr: true},
		{name: "value min length ok", rule: rules[4], value: []string{"hello", "world"}, wantRes: []string{"hello", "world"}},
		{name: "value min length ok empty", rule: rules[4], value: []string{}, wantRes: []string{}},
		{name: "value min length fail 1", rule: rules[4], value: []string{"h"}, wantErr: true},
		{name: "value min length fail 2", rule: rules[4], value: []string{"hi", "h"}, wantErr: true},

		{name: "sort result 1", rule: rules[5], value: []string{"b", "a"}, wantRes: []string{"a", "b"}},
		{name: "sort result 2", rule: rules[5], value: []string{"b", "a", "c"}, wantRes: []string{"a", "b", "c"}},
		{name: "unique 1", rule: rules[6], value: []string{"hello", "hello"}, wantRes: []string{"hello"}},
		{name: "unique 2", rule: rules[6], value: []string{"ciao", "hello", "hola", "hello"}, wantRes: []string{"ciao", "hello", "hola"}},
		{name: "unique 3", rule: rules[6], value: []string{"ciao", "hola", "hello"}, wantRes: []string{"ciao", "hello", "hola"}},
		{name: "unique also sorts 1", rule: rules[6], value: []string{"c", "a", "b"}, wantRes: []string{"a", "b", "c"}},
		{name: "unique also sorts 1", rule: rules[6], value: []string{"c", "a", "a", "b", "c"}, wantRes: []string{"a", "b", "c"}},

		{name: "invalid rule: min<1", rule: rules[7], value: []string{}, wantErr: true},
		{name: "invalid rule: max<1", rule: rules[8], value: []string{}, wantErr: true},
		{name: "invalid rule: min>max", rule: rules[9], value: []string{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := sliceValidator[string](tt.rule)
			gotRes, err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("sliceValidator().validator error = %v, wantErr %v (value = %s)", err, tt.wantErr, gotRes)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("sliceValidator().validator = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
