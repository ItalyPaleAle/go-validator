package validator

import (
	"reflect"
	"testing"
)

// TestValidate primarily focuses on validating that code compiles are types accepted and returned by Validate and ValidateAny are correct
// It doesn't perform full validation of all cases: that is left to the tests for individual validators
func TestValidate(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			name    string
			val     string
			wantRes string
			wantErr bool
		}{
			{name: "string test", val: "hello world", wantRes: "hello world", wantErr: false},
			{name: "empty string", val: "", wantRes: "", wantErr: false},
		}
		t.Run("Validate", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := Validate(tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !reflect.DeepEqual(gotRes, tt.wantRes) {
						t.Errorf("Validate() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
		t.Run("ValidateAny", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := ValidateAny(tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("ValidateAny() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if reflect.TypeOf(gotRes).String() != reflect.TypeOf(tt.wantRes).String() {
						t.Errorf("ValidateAny() = type %T, want type %T", gotRes, tt.wantRes)
						return
					}
					if !reflect.DeepEqual(gotRes, tt.wantRes) {
						t.Errorf("ValidateAny() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
		t.Run("ValidateAny with pointer", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := ValidateAny(&tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("ValidateAny() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					rt := reflect.TypeOf(gotRes)
					if rt.Kind() != reflect.Pointer {
						t.Error("ValidateAny() did not return a pointer")
						return
					}
					gotRes = reflect.ValueOf(gotRes).Elem().Interface()
					gotResV, ok := gotRes.(string)
					if !ok {
						t.Error("ValidateAny() did not return a pointer to string")
						return
					}
					if !reflect.DeepEqual(gotResV, tt.wantRes) {
						t.Errorf("ValidateAny() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
	})

	t.Run("[]string", func(t *testing.T) {
		tests := []struct {
			name    string
			val     []string
			wantRes []string
			wantErr bool
		}{
			{name: "slice of strings", val: []string{"hello world"}, wantRes: []string{"hello world"}, wantErr: false},
			{name: "empty slice", val: []string{}, wantRes: []string{}, wantErr: false},
			{name: "nil", val: nil, wantRes: nil, wantErr: false},
		}
		t.Run("Validate", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := Validate(tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !reflect.DeepEqual(gotRes, tt.wantRes) {
						t.Errorf("Validate() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
		t.Run("ValidateAny", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := ValidateAny(tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("ValidateAny() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if reflect.TypeOf(gotRes).String() != reflect.TypeOf(tt.wantRes).String() {
						t.Errorf("ValidateAny() = type %T, want type %T", gotRes, tt.wantRes)
						return
					}
					if !reflect.DeepEqual(gotRes, tt.wantRes) {
						t.Errorf("ValidateAny() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
		t.Run("ValidateAny with pointer", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := ValidateAny(&tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("ValidateAny() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					rt := reflect.TypeOf(gotRes)
					if rt.Kind() != reflect.Pointer {
						t.Error("ValidateAny() did not return a pointer")
						return
					}
					gotRes = reflect.ValueOf(gotRes).Elem().Interface()
					gotResV, ok := gotRes.([]string)
					if !ok {
						t.Error("ValidateAny() did not return a pointer to []string")
						return
					}
					if !reflect.DeepEqual(gotResV, tt.wantRes) {
						t.Errorf("ValidateAny() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
	})

	t.Run("map[string]string", func(t *testing.T) {
		tests := []struct {
			name    string
			val     map[string]string
			wantRes map[string]string
			wantErr bool
		}{
			{name: "map of strings", val: map[string]string{"foo": "bar"}, wantRes: map[string]string{"foo": "bar"}, wantErr: false},
			{name: "empty map", val: map[string]string{}, wantRes: map[string]string{}, wantErr: false},
			{name: "nil", val: nil, wantRes: nil, wantErr: false},
		}
		t.Run("Validate", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := Validate(tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !reflect.DeepEqual(gotRes, tt.wantRes) {
						t.Errorf("Validate() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
		t.Run("ValidateAny", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := ValidateAny(tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("ValidateAny() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if reflect.TypeOf(gotRes).String() != reflect.TypeOf(tt.wantRes).String() {
						t.Errorf("ValidateAny() = type %T, want type %T", gotRes, tt.wantRes)
						return
					}
					if !reflect.DeepEqual(gotRes, tt.wantRes) {
						t.Errorf("ValidateAny() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
		t.Run("ValidateAny with pointer", func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotRes, err := ValidateAny(&tt.val, "")
					if (err != nil) != tt.wantErr {
						t.Errorf("ValidateAny() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					rt := reflect.TypeOf(gotRes)
					if rt.Kind() != reflect.Pointer {
						t.Error("ValidateAny() did not return a pointer")
						return
					}
					gotRes = reflect.ValueOf(gotRes).Elem().Interface()
					gotResV, ok := gotRes.(map[string]string)
					if !ok {
						t.Error("ValidateAny() did not return a pointer to map[string]string")
						return
					}
					if !reflect.DeepEqual(gotResV, tt.wantRes) {
						t.Errorf("ValidateAny() = %v, want %v", gotRes, tt.wantRes)
					}
				})
			}
		})
	})

	t.Run("ValidateAny cases", func(t *testing.T) {
		var zeroStr string
		tests := []struct {
			name    string
			val     any
			wantRes any
			wantErr bool
		}{
			{name: "nil", val: nil, wantRes: nil, wantErr: false},
			{name: "pointer to zero string", val: &zeroStr, wantRes: &zeroStr, wantErr: false},
			{name: "unsupported type 1", val: struct{}{}, wantRes: nil, wantErr: true},
			{name: "unsupported type 2", val: t, wantRes: nil, wantErr: true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotRes, err := ValidateAny(tt.val, "")
				if (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(gotRes, tt.wantRes) {
					t.Errorf("Validate() = %v, want %v", gotRes, tt.wantRes)
				}
			})
		}
	})
}
