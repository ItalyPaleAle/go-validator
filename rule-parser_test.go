package validator

import (
	"reflect"
	"testing"
)

func Test_parseParams(t *testing.T) {
	type args struct {
		rule string
	}
	tests := []struct {
		name    string
		args    args
		wantRes map[string]string
		wantErr bool
	}{
		{
			name:    "empty rule",
			args:    args{rule: ""},
			wantRes: map[string]string{},
		},
		{
			args:    args{rule: "foo=bar"},
			wantRes: map[string]string{"foo": "bar"},
		},
		{
			args:    args{rule: "foo"},
			wantRes: map[string]string{"foo": ""},
		},
		{
			args:    args{rule: "foo=bar,2=1"},
			wantRes: map[string]string{"foo": "bar", "2": "1"},
		},
		{
			args:    args{rule: "foo=bar,2=1,me"},
			wantRes: map[string]string{"foo": "bar", "2": "1", "me": ""},
		},
		{
			args:    args{rule: "foo=bar,me,2=1"},
			wantRes: map[string]string{"foo": "bar", "me": "", "2": "1"},
		},
		{
			args:    args{rule: "foo=(bar=1,me),2=1"},
			wantRes: map[string]string{"foo": "bar=1,me", "2": "1"},
		},
		{
			args:    args{rule: "(foo),bar"},
			wantRes: map[string]string{"foo": "", "bar": ""},
		},
		{
			args:    args{rule: "foo,bar=((foo))"},
			wantRes: map[string]string{"foo": "", "bar": "(foo)"},
		},
		{
			args:    args{rule: "(foo),bar=((foo))"},
			wantRes: map[string]string{"foo": "", "bar": "(foo)"},
		},
		{
			args:    args{rule: "ciao,mondo,bar=((foo))"},
			wantRes: map[string]string{"ciao": "", "mondo": "", "bar": "(foo)"},
		},
		{
			args:    args{rule: "(ciao,mondo),bar=((foo))"},
			wantRes: map[string]string{"ciao,mondo": "", "bar": "(foo)"},
		},
		{
			args:    args{rule: "((ciao,mondo)),bar=((foo))"},
			wantRes: map[string]string{"(ciao,mondo)": "", "bar": "(foo)"},
		},
		{
			args:    args{rule: "=foo"},
			wantErr: true,
		},
		{
			args:    args{rule: "foo=bar,=1"},
			wantErr: true,
		},
		{
			args:    args{rule: "foo)"},
			wantErr: true,
		},
		{
			args:    args{rule: "foo)bar"},
			wantErr: true,
		},
		{
			args:    args{rule: "((foo)"},
			wantErr: true,
		},
		{
			args:    args{rule: "()"},
			wantErr: true,
		},
		{
			args:    args{rule: "foo<>"},
			wantRes: map[string]string{"foo<>": ""},
		},
		{
			args:    args{rule: "(aa=bb,cc=dd)"},
			wantRes: map[string]string{"aa=bb,cc=dd": ""},
		},
		{
			args:    args{rule: "🐱"},
			wantRes: map[string]string{"🐱": ""},
		},
		{
			args:    args{rule: "🐶,😱"},
			wantRes: map[string]string{"🐶": "", "😱": ""},
		},
		{
			args:    args{rule: "😃=😍"},
			wantRes: map[string]string{"😃": "😍"},
		},
	}
	for _, tt := range tests {
		name := tt.name
		if name == "" {
			name = tt.args.rule
		}
		t.Run(name, func(t *testing.T) {
			gotRes, err := parseParams(tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRuleTree() error = %v, wantErr %v (value = %v)", err, tt.wantErr, gotRes)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("parseRuleTree() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
