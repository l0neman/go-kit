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
			wantErr: "字段 'Name' 不能为空",
		},
		{
			name: "String Is Empty Fail",
			input: &struct {
				Name string `check:"is_empty"`
			}{"not empty"},
			wantErr: "字段 'Name' 必须为空",
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
			wantErr: "不满足条件 '>=18'",
		},
		{
			name: "Int Greater Than Fail",
			input: &struct {
				Age int `check:"<60"`
			}{60},
			wantErr: "不满足条件 '<60'",
		},
		{
			name: "Int Not Equal Fail",
			input: &struct {
				Code int `check:"!=0"`
			}{0},
			wantErr: "不满足条件 '!=0'",
		},
		{
			name: "Int Equal Fail",
			input: &struct {
				Code int `check:"==1"`
			}{2},
			wantErr: "不满足条件 '==1'",
		},
		{
			name: "Multiple Checks Fail",
			input: &struct {
				Age int `check:">18, <30"`
			}{35},
			wantErr: "不满足条件 '<30'",
		},
		{
			name: "Float Less Than Or Equal Fail",
			input: &struct {
				Price float64 `check:"<=99.9"`
			}{100.0},
			wantErr: "不满足条件 '<=99.9'",
		},
		{
			name: "Nested Struct Fail",
			input: &struct {
				Data Nested
			}{Nested{Value: 5}},
			wantErr: "不满足条件 '>10'",
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
			wantErr: "不满足条件 '>10'",
		},
		{
			name: "Invalid Expression",
			input: &struct {
				Name string `check:"invalid_op"`
			}{"test"},
			wantErr: "检查表达式 'invalid_op' 无效",
		},
		{
			name: "Mismatched Type In Expression",
			input: &struct {
				Age int `check:">abc"`
			}{25},
			wantErr: "检查值 'abc' 类型错误",
		},
		{
			name: "Unsupported Expression For Type",
			input: &struct {
				Name string `check:">10"`
			}{"test"},
			wantErr: "字符串类型，不支持 '>' 表达式",
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
