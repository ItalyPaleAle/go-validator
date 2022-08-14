package validator

import (
	"reflect"
	"testing"
)

func Test_stringValidator(t *testing.T) {
	rules := []string{
		"",
		"min=3",
		"max=5",
		"min=2,max=5",
		"preserve-whitespace",
		"preserve-newlines",
		"replace-whitespaces",
		"replace-whitespaces,preserve-newlines",
		"asciionly",
		"unorm=nfd",
		"unorm=nfkc",
		"asciionly,unorm=nfd",
	}
	rulesInvalid := []string{
		"min=0",
		"max=-1",
		"min=5,max=1",
		"unorm=invalid",
	}

	tests := []struct {
		name    string
		rule    string
		value   string
		wantRes string
		wantErr bool
	}{
		{name: "empty string", rule: rules[0], value: "", wantRes: ""},
		{name: "ascii string", rule: rules[0], value: "hi!", wantRes: "hi!"},
		{name: "trim spaces", rule: rules[0], value: " hi! ", wantRes: "hi!"},
		{name: "trim unicode spaces", rule: rules[0], value: "  hi! \n ", wantRes: "hi!"},
		{name: "only spaces", rule: rules[0], value: "   ", wantRes: ""},
		{name: "normalize string 1", rule: rules[0], value: "\u00e8", wantRes: "\u00e8"},
		{name: "normalize string 2", rule: rules[0], value: "e\u0300", wantRes: "\u00e8"},
		{name: "normalize string 3", rule: rules[0], value: "\u00df", wantRes: "\u00df"},
		{name: "normalize string 4", rule: rules[0], value: "😀", wantRes: "😀"},
		{name: "normalize string 5", rule: rules[0], value: "\U0001f600", wantRes: "😀"},
		{name: "emojis with modifiers and ZWJ", rule: rules[0], value: "😀 👨‍👩‍👧‍👦 1️⃣ 💁🏽‍♂️ 🧑🏻‍🍼", wantRes: "😀 👨‍👩‍👧‍👦 1️⃣ 💁🏽‍♂️ 🧑🏻‍🍼"},

		{name: "min length ok", rule: rules[1], value: " hi! ", wantRes: "hi!"},
		{name: "min length fail", rule: rules[1], value: " hi ", wantErr: true},
		{name: "min length empty string 1", rule: rules[1], value: "", wantErr: true},
		{name: "min length empty string 2", rule: rules[1], value: "  ", wantErr: true},
		{name: "max length ok", rule: rules[2], value: "  hi!  ", wantRes: "hi!"},
		{name: "max length fail", rule: rules[2], value: "hello world", wantErr: true},
		{name: "min and max length ok 1", rule: rules[3], value: "hello", wantRes: "hello"},
		{name: "min and max length ok 2", rule: rules[3], value: "hello   ", wantRes: "hello"},
		{name: "min and max length ok 3", rule: rules[3], value: "hi -!   ", wantRes: "hi -!"},
		{name: "min and max length fail 1", rule: rules[3], value: "hello!   ", wantErr: true},
		{name: "min and max length fail 2", rule: rules[3], value: "hi - !!   ", wantErr: true},

		{name: "collapse whitespaces", rule: rules[0], value: "hi   !", wantRes: "hi !"},
		{name: "replace Unicode spaces", rule: rules[0], value: "hi !", wantRes: "hi !"},
		{name: "replace and collapse Unicode spaces", rule: rules[0], value: "h   i !", wantRes: "h i !"},
		{name: "tabs are replaced with spaces", rule: rules[0], value: "hi\t !", wantRes: "hi !"},

		{name: "preserve whitespaces", rule: rules[4], value: "hi   !", wantRes: "hi   !"},
		{name: "preserve whitespaces keeps newlines", rule: rules[4], value: "hi   !\n \n\n hi", wantRes: "hi   !\n \n\n hi"},
		{name: "preserve whitespaces keeps tabs", rule: rules[4], value: "hi\t   !\n \n\n hi", wantRes: "hi\t   !\n \n\n hi"},
		{name: "preserve Unicode spaces", rule: rules[4], value: "hi !", wantRes: "hi !"},
		{name: "preserve and don't collapse Unicode spaces", rule: rules[4], value: "h   i !", wantRes: "h   i !"},

		{name: "remove control chars", rule: rules[0], value: "he\x07l\u001el\uFEFFo\u2064", wantRes: "hello"},
		{name: "newlines are removed", rule: rules[0], value: "hello\nworld", wantRes: "hello world"},
		{name: "multiple newlines are collapsed to a single space", rule: rules[0], value: "hello\n\n\nworld", wantRes: "hello world"},
		{name: "collapse whitespaces and newlines", rule: rules[0], value: "hello\n \n world", wantRes: "hello world"},
		{name: "remove carriage returns", rule: rules[0], value: "hello\r\nworld", wantRes: "hello world"},

		{name: "preserve newlines", rule: rules[5], value: "hello\nworld", wantRes: "hello\nworld"},
		{name: "newlines at the ends are always removed", rule: rules[5], value: "   \nhello\nworl\nd\n", wantRes: "hello\nworl\nd"},
		{name: "preserve multiple newlines", rule: rules[5], value: "hello\n\n\nworld", wantRes: "hello\n\n\nworld"},
		{name: "remove carriage returns but preserve newlines", rule: rules[5], value: "hello\r\nworld", wantRes: "hello\nworld"},
		{name: "whitespaces and newlines, preserving newlines", rule: rules[5], value: "hello \n world-hello   \n\n \n world", wantRes: "hello \nworld-hello \n\n\nworld"},

		{name: "replace and collapse whitespaces with underscore", rule: rules[6], value: "hi   !", wantRes: "hi_!"},
		{name: "replace and collapse newlines with underscore", rule: rules[6], value: "hi  \n \n\n !", wantRes: "hi_!"},
		{name: "replace Unicode spaces with underscore", rule: rules[6], value: "hi !", wantRes: "hi_!"},
		{name: "replace and collapse Unicode spaces with underscore", rule: rules[6], value: "h  \n\r i !", wantRes: "h_i_!"},
		{name: "replace and collapse whitespaces with underscore, trim from ends", rule: rules[6], value: "  hi   ! \n", wantRes: "hi_!"},

		{name: "replace and collapse whitespaces with underscore, keep newlines 1", rule: rules[7], value: "hi   !", wantRes: "hi_!"},
		{name: "replace and collapse whitespaces with underscore, keep newlines 2", rule: rules[7], value: "hi  \n \n\n !", wantRes: "hi_\n\n\n!"},
		{name: "replace Unicode spaces with underscore, keep newlines", rule: rules[7], value: "hi !\n !", wantRes: "hi_!\n!"},
		{name: "replace and collapse Unicode spaces with underscore, keep newlines", rule: rules[7], value: "h  \n\r i !", wantRes: "h_\ni_!"},

		{name: "asciionly allows ASCII characters", rule: rules[8], value: "hello   !\nworld", wantRes: "hello ! world"},
		{name: "asciionly removes all non-ASCII characters 1", rule: rules[8], value: "日本語😊", wantRes: ""},
		// 1️⃣ is made of the character "1" which is preserved (+ \uFE0F \u20E3)
		{name: "asciionly removes all non-ASCII characters 2", rule: rules[8], value: "1️⃣", wantRes: "1"},
		// Same character, in both NFC and NFD forms; it is removed after normalization
		{name: "asciionly operates after normalization to NFC", rule: rules[8], value: "a\u0308 \u00e4", wantRes: ""},

		{name: "normalize to form NFD 1", rule: rules[9], value: "e\u0300", wantRes: "e\u0300"},
		{name: "normalize to form NFD 2", rule: rules[9], value: "\u00e8", wantRes: "e\u0300"},
		{name: "normalize to form NFKC", rule: rules[10], value: "①", wantRes: "1"},
		{name: "normalize to form NFD with asciionly", rule: rules[11], value: "e\u0300", wantRes: "e"},

		{name: "invalid rule: min<1", rule: rulesInvalid[0], value: "foo", wantErr: true},
		{name: "invalid rule: max<1", rule: rulesInvalid[1], value: "foo", wantErr: true},
		{name: "invalid rule: min>max", rule: rulesInvalid[2], value: "foo", wantErr: true},
		{name: "invalid rule: invalid unorm value", rule: rulesInvalid[3], value: "foo", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := stringValidator(tt.rule)
			gotRes, err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringValidator().validator error = %v, wantErr %v (value = %s)", err, tt.wantErr, gotRes)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("stringValidator().validator = %v (%v), want %v (%v)", gotRes, []byte(gotRes), tt.wantRes, []byte(tt.wantRes))
			}
		})
	}
}
