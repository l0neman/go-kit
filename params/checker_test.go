package params

import (
	"strings"
	"testing"
)

func TestCheckFields(t *testing.T) {
	type Nested struct {
		Value int `check:">10"`
	}

	testCases := []struct {
		name    string
		input   interface{}
		wantErr string // Substring of the expected error message. Empty if no error is expected.
	}{
		{
			name: "Valid Case",
			input: &struct {
				Name  string  `check:"not_empty"`
				Age   int     `check:">=18,<100"`
				Score float64 `check:">0.0"`
			}{"John", 25, 99.5},
			wantErr: "",
		},
		{
			name: "String Not Empty Fail",
			input: &struct {
				Name string `check:"not_empty"`
			}{""},
			wantErr: "field 'Name' cannot be empty",
		},
		{
			name: "String Is Empty Fail",
			input: &struct {
				Name string `check:"is_empty"`
			}{"not empty"},
			wantErr: "field 'Name' must be empty",
		},
		{
			name: "String Is Empty Success",
			input: &struct {
				Name string `check:"is_empty"`
			}{""},
			wantErr: "",
		},
		{
			name: "Int Less Than Fail",
			input: &struct {
				Age int `check:">=18"`
			}{17},
			wantErr: "does not satisfy condition '>=18'",
		},
		{
			name: "Int Greater Than Fail",
			input: &struct {
				Age int `check:"<60"`
			}{60},
			wantErr: "does not satisfy condition '<60'",
		},
		{
			name: "Int Not Equal Fail",
			input: &struct {
				Code int `check:"!=0"`
			}{0},
			wantErr: "does not satisfy condition '!=0'",
		},
		{
			name: "Int Equal Fail",
			input: &struct {
				Code int `check:"==1"`
			}{2},
			wantErr: "does not satisfy condition '==1'",
		},
		{
			name: "Multiple Checks Fail",
			input: &struct {
				Age int `check:">18, <30"`
			}{35},
			wantErr: "does not satisfy condition '<30'",
		},
		{
			name: "Float Less Than Or Equal Fail",
			input: &struct {
				Price float64 `check:"<=99.9"`
			}{100.0},
			wantErr: "does not satisfy condition '<=99.9'",
		},
		{
			name: "Nested Struct Fail",
			input: &struct {
				Data Nested
			}{Nested{Value: 5}},
			wantErr: "does not satisfy condition '>10'",
		},
		{
			name: "Nested Struct Success",
			input: &struct {
				Data Nested
			}{Nested{Value: 20}},
			wantErr: "",
		},
		{
			name: "Pointer To Nested Struct Fail",
			input: &struct {
				Data *Nested
			}{&Nested{Value: 5}},
			wantErr: "does not satisfy condition '>10'",
		},
		{
			name: "Invalid Expression",
			input: &struct {
				Name string `check:"invalid_op"`
			}{"test"},
			wantErr: "invalid check expression 'invalid_op'",
		},
		{
			name: "Mismatched Type In Expression",
			input: &struct {
				Age int `check:">abc"`
			}{25},
			wantErr: "check value 'abc' has incorrect type",
		},
		{
			name: "Unsupported Expression For Type",
			input: &struct {
				Name string `check:">10"`
			}{"test"},
			wantErr: "is of string type, does not support '>' expression",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckFields(tc.input)

			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("CheckFields() error = %v, wantErr %v", err, tc.wantErr)
				}
			} else {
				if err == nil {
					t.Errorf("CheckFields() error = nil, wantErr %v", tc.wantErr)
				} else if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("CheckFields() error = %v, wantErr substring %v", err.Error(), tc.wantErr)
				}
			}
		})
	}
}
