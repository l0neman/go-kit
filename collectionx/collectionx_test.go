package collectionx

import (
	"reflect"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	type args[T any] struct {
		source []T
	}

	tests := []struct {
		name string
		args interface{}
		want bool
	}{
		{
			name: "nil slice is empty",
			args: args[int]{source: nil},
			want: true,
		},
		{
			name: "empty slice is empty",
			args: args[int]{source: []int{}},
			want: true,
		},
		{
			name: "non-empty slice is not empty",
			args: args[int]{source: []int{1, 2, 3}},
			want: false,
		},
		{
			name: "single element slice is not empty",
			args: args[string]{source: []string{"a"}},
			want: false,
		},
		{
			name: "empty string slice",
			args: args[string]{source: []string{}},
			want: true,
		},
		{
			name: "nil string slice",
			args: args[string]{source: nil},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[int]:
				if got := IsEmpty(args.source); got != tt.want {
					t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
				}
			case args[string]:
				if got := IsEmpty(args.source); got != tt.want {
					t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAppend(t *testing.T) {
	type args[T any] struct {
		initial []T
		toAdd   []T
	}

	tests := []struct {
		name     string
		args     interface{}
		expected interface{}
	}{
		{
			name:     "append to empty slice",
			args:     args[int]{initial: []int{}, toAdd: []int{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "append to non-empty slice",
			args:     args[int]{initial: []int{1}, toAdd: []int{2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "append empty slice",
			args:     args[int]{initial: []int{1, 2}, toAdd: []int{}},
			expected: []int{1, 2},
		},
		{
			name:     "append single element",
			args:     args[string]{initial: []string{"a"}, toAdd: []string{"b", "c"}},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "append to nil slice",
			args:     args[int]{initial: nil, toAdd: []int{1, 2}},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[int]:
				result := args.initial
				Append(&result, args.toAdd...)
				if !Equals(result, tt.expected.([]int)) {
					t.Errorf("Append() result = %v, want %v", result, tt.expected.([]int))
				}
			case args[string]:
				result := args.initial
				Append(&result, args.toAdd...)
				if !Equals(result, tt.expected.([]string)) {
					t.Errorf("Append() result = %v, want %v", result, tt.expected.([]string))
				}
			}
		})
	}
}

func TestSliceOf(t *testing.T) {
	type args[T any] struct {
		values []T
	}

	tests := []struct {
		name     string
		args     interface{}
		expected interface{}
	}{
		{
			name:     "create slice with single value",
			args:     args[int]{values: []int{1}},
			expected: []int{1},
		},
		{
			name:     "create slice with multiple values",
			args:     args[int]{values: []int{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "create empty slice",
			args:     args[string]{values: []string{}},
			expected: []string{},
		},
		{
			name:     "create slice with string values",
			args:     args[string]{values: []string{"a", "b", "c"}},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "create slice with duplicate values",
			args:     args[int]{values: []int{1, 1, 2}},
			expected: []int{1, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[int]:
				result := SliceOf(args.values...)
				if !reflect.DeepEqual(*result, tt.expected.([]int)) {
					t.Errorf("SliceOf() result = %v, want %v", *result, tt.expected.([]int))
				}
			case args[string]:
				result := SliceOf(args.values...)
				if !reflect.DeepEqual(*result, tt.expected.([]string)) {
					t.Errorf("SliceOf() result = %v, want %v", *result, tt.expected.([]string))
				}
			}
		})
	}
}

func TestSetOfFunction(t *testing.T) {
	type args[T comparable] struct {
		values []T
	}

	tests := []struct {
		name     string
		args     interface{}
		expected []interface{}
	}{
		{
			name:     "create set with single value",
			args:     args[int]{values: []int{1}},
			expected: []interface{}{1},
		},
		{
			name:     "create set with multiple values",
			args:     args[int]{values: []int{1, 2, 3}},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "create set with duplicate values (should be removed)",
			args:     args[int]{values: []int{1, 1, 2}},
			expected: []interface{}{1, 2},
		},
		{
			name:     "create set with string values",
			args:     args[string]{values: []string{"a", "b", "c"}},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name:     "create empty set",
			args:     args[int]{values: []int{}},
			expected: []interface{}{},
		},
		{
			name:     "create set with duplicate string values",
			args:     args[string]{values: []string{"a", "a", "b"}},
			expected: []interface{}{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[int]:
				set := SetOf(args.values...)
				resultSlice := set.ToSlice()

				// Convert expected to int slice for comparison
				expectedInt := make([]int, len(tt.expected))
				for i, v := range tt.expected {
					expectedInt[i] = v.(int)
				}

				if !ContainsAll(resultSlice, expectedInt) || !ContainsAll(expectedInt, resultSlice) {
					t.Errorf("SetOf() result = %v, want %v", resultSlice, expectedInt)
				}

				// Check that size is correct (duplicates should be removed)
				expectedSize := len(expectedInt)
				if set.Size() != expectedSize {
					t.Errorf("SetOf() size = %d, want %d", set.Size(), expectedSize)
				}
			case args[string]:
				set := SetOf(args.values...)
				resultSlice := set.ToSlice()

				// Convert expected to string slice for comparison
				expectedString := make([]string, len(tt.expected))
				for i, v := range tt.expected {
					expectedString[i] = v.(string)
				}

				if !ContainsAll(resultSlice, expectedString) || !ContainsAll(expectedString, resultSlice) {
					t.Errorf("SetOf() result = %v, want %v", resultSlice, expectedString)
				}

				// Check that size is correct (duplicates should be removed)
				expectedSize := len(expectedString)
				if set.Size() != expectedSize {
					t.Errorf("SetOf() size = %d, want %d", set.Size(), expectedSize)
				}
			}
		})
	}
}
