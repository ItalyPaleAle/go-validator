package validator

import (
	"reflect"
	"testing"
)

func Test_mapValidator(t *testing.T) {
	rules := []string{
		"",
		"min=2",
		"max=3",
		"min=2,max=3",
		"value=(min=2)",
		"key=(min=2),value=(max=5)",
		"key=(asciionly)",
	}
	invalidRules := []string{
		"min=0",
		"max=-1",
		"min=5,max=1",
	}

	tests := []struct {
		name    string
		rule    string
		value   map[string]string
		wantRes map[string]string
		wantErr bool
	}{
		{name: "empty map", rule: rules[0], value: map[string]string{}, wantRes: map[string]string{}},
		{name: "one element", rule: rules[0], value: map[string]string{"hello": "world"}, wantRes: map[string]string{"hello": "world"}},
		{name: "two elements", rule: rules[0], value: map[string]string{"hello": "world", "foo": "bar"}, wantRes: map[string]string{"hello": "world", "foo": "bar"}},
		{name: "trim spaces in strings", rule: rules[0], value: map[string]string{"1": " hi! ", "  2  ": " ciao "}, wantRes: map[string]string{"1": "hi!", "2": "ciao"}},
		{name: "normalize strings", rule: rules[0], value: map[string]string{"e\u0300": "e\u0300", "foo2": "\u00e8", "foo3": "😀"}, wantRes: map[string]string{"\u00e8": "\u00e8", "foo2": "\u00e8", "foo3": "😀"}},
		{name: "min length ok", rule: rules[1], value: map[string]string{"a": "", "b": ""}, wantRes: map[string]string{"a": "", "b": ""}},
		{name: "min length fail", rule: rules[1], value: map[string]string{"a": "1"}, wantErr: true},
		{name: "min length empty fail", rule: rules[1], value: map[string]string{}, wantErr: true},
		{name: "max length ok 1", rule: rules[2], value: map[string]string{"a": "", "b": ""}, wantRes: map[string]string{"a": "", "b": ""}},
		{name: "max length ok 2", rule: rules[2], value: map[string]string{}, wantRes: map[string]string{}},
		{name: "max length fail", rule: rules[2], value: map[string]string{"a": "", "b": "", "c": "", "d": ""}, wantErr: true},
		{name: "min and max length ok", rule: rules[3], value: map[string]string{"a": "", "b": ""}, wantRes: map[string]string{"a": "", "b": ""}},
		{name: "min and max length fail 1", rule: rules[3], value: map[string]string{"a": ""}, wantErr: true},
		{name: "min and max length fail 2", rule: rules[3], value: map[string]string{"a": "", "b": "", "c": "", "d": ""}, wantErr: true},
		{name: "value min length ok", rule: rules[4], value: map[string]string{"1": "hello", "2": "world"}, wantRes: map[string]string{"1": "hello", "2": "world"}},
		{name: "value min length ok empty", rule: rules[4], value: map[string]string{}, wantRes: map[string]string{}},
		{name: "value min length fail 1", rule: rules[4], value: map[string]string{"1": "h"}, wantErr: true},
		{name: "value min length fail 2", rule: rules[4], value: map[string]string{"1": "hi", "2": "h"}, wantErr: true},
		{name: "key min length ok", rule: rules[5], value: map[string]string{"foo": "hello", "bar": "world"}, wantRes: map[string]string{"foo": "hello", "bar": "world"}},
		{name: "key min length ok empty", rule: rules[5], value: map[string]string{}, wantRes: map[string]string{}},
		{name: "key min length fail 1", rule: rules[5], value: map[string]string{"1": "h"}, wantErr: true},
		{name: "key min length fail 2", rule: rules[5], value: map[string]string{"foo": "hi", "2": "h"}, wantErr: true},
		{name: "value max length ok", rule: rules[5], value: map[string]string{"foo": "hello", "bar": "world!"}, wantErr: true},
		{name: "asciionly for key", rule: rules[6], value: map[string]string{"hi🤣": "😕"}, wantRes: map[string]string{"hi": "😕"}},

		{name: "invalid rule: min<1", rule: invalidRules[0], value: map[string]string{}, wantErr: true},
		{name: "invalid rule: max<1", rule: invalidRules[1], value: map[string]string{}, wantErr: true},
		{name: "invalid rule: min>max", rule: invalidRules[2], value: map[string]string{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := mapValidator[string](tt.rule)
			gotRes, err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapValidator().validator error = %v, wantErr %v (value = %s)", err, tt.wantErr, gotRes)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("mapValidator().validator = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
