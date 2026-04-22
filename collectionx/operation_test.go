package collectionx

import (
	"reflect"
	"testing"
)

func TestEquals(t *testing.T) {
	type args[T comparable] struct {
		a []T
		b []T
	}

	tests := []struct {
		name string
		args interface{}
		want bool
	}{
		{
			name: "equal slices with same order",
			args: args[int]{a: []int{1, 2, 3}, b: []int{1, 2, 3}},
			want: true,
		},
		{
			name: "equal slices with different order",
			args: args[int]{a: []int{1, 2, 3}, b: []int{3, 2, 1}},
			want: true,
		},
		{
			name: "different lengths",
			args: args[int]{a: []int{1, 2, 3}, b: []int{1, 2}},
			want: false,
		},
		{
			name: "same length but different elements",
			args: args[int]{a: []int{1, 2, 3}, b: []int{1, 2, 4}},
			want: false,
		},
		{
			name: "both empty slices",
			args: args[int]{a: []int{}, b: []int{}},
			want: true,
		},
		{
			name: "slices with duplicate elements",
			args: args[int]{a: []int{1, 1, 2, 3}, b: []int{3, 2, 1, 1}},
			want: true,
		},
		{
			name: "different counts of same elements",
			args: args[int]{a: []int{1, 1, 2}, b: []int{1, 2, 2}},
			want: false,
		},
		{
			name: "string slices - equal",
			args: args[string]{a: []string{"a", "b", "c"}, b: []string{"c", "a", "b"}},
			want: true,
		},
		{
			name: "string slices - different",
			args: args[string]{a: []string{"a", "b", "c"}, b: []string{"a", "b", "d"}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[int]:
				if got := Equals(args.a, args.b); got != tt.want {
					t.Errorf("Equals() = %v, want %v", got, tt.want)
				}
			case args[string]:
				if got := Equals(args.a, args.b); got != tt.want {
					t.Errorf("Equals() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args[S ~[]T, T comparable] struct {
		slice []T
		item  T
	}

	tests := []struct {
		name string
		args interface{}
		want bool
	}{
		{
			name: "contains int element",
			args: args[[]int, int]{slice: []int{1, 2, 3, 4}, item: 3},
			want: true,
		},
		{
			name: "does not contain int element",
			args: args[[]int, int]{slice: []int{1, 2, 3, 4}, item: 5},
			want: false,
		},
		{
			name: "empty slice does not contain",
			args: args[[]int, int]{slice: []int{}, item: 1},
			want: false,
		},
		{
			name: "contains string element",
			args: args[[]string, string]{slice: []string{"a", "b", "c"}, item: "b"},
			want: true,
		},
		{
			name: "does not contain string element",
			args: args[[]string, string]{slice: []string{"a", "b", "c"}, item: "d"},
			want: false,
		},
		{
			name: "contains duplicate element",
			args: args[[]int, int]{slice: []int{1, 2, 2, 3}, item: 2},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[[]int, int]:
				if got := Contains(args.slice, args.item); got != tt.want {
					t.Errorf("Contains() = %v, want %v", got, tt.want)
				}
			case args[[]string, string]:
				if got := Contains(args.slice, args.item); got != tt.want {
					t.Errorf("Contains() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestContainsAll(t *testing.T) {
	type args[S ~[]T, T comparable] struct {
		source []T
		target []T
	}

	tests := []struct {
		name string
		args interface{}
		want bool
	}{
		{
			name: "source contains all target elements",
			args: args[[]int, int]{source: []int{1, 2, 3, 4, 5}, target: []int{2, 4}},
			want: true,
		},
		{
			name: "source does not contain all target elements",
			args: args[[]int, int]{source: []int{1, 2, 3}, target: []int{2, 4}},
			want: false,
		},
		{
			name: "target is empty",
			args: args[[]int, int]{source: []int{1, 2, 3}, target: []int{}},
			want: true,
		},
		{
			name: "source is empty but target is not",
			args: args[[]int, int]{source: []int{}, target: []int{1}},
			want: false,
		},
		{
			name: "both empty",
			args: args[[]int, int]{source: []int{}, target: []int{}},
			want: true,
		},
		{
			name: "source contains all elements including duplicates",
			args: args[[]int, int]{source: []int{1, 2, 3, 4}, target: []int{1, 1, 2}},
			want: true,
		},
		{
			name: "string slices - contains all",
			args: args[[]string, string]{source: []string{"a", "b", "c", "d"}, target: []string{"b", "d"}},
			want: true,
		},
		{
			name: "string slices - does not contain all",
			args: args[[]string, string]{source: []string{"a", "b", "c"}, target: []string{"b", "e"}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[[]int, int]:
				if got := ContainsAll(args.source, args.target); got != tt.want {
					t.Errorf("ContainsAll() = %v, want %v", got, tt.want)
				}
			case args[[]string, string]:
				if got := ContainsAll(args.source, args.target); got != tt.want {
					t.Errorf("ContainsAll() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMerge(t *testing.T) {
	type args[S ~[]T, T comparable] struct {
		source []T
		target []T
	}

	tests := []struct {
		name string
		args interface{}
		want []int
	}{
		{
			name: "merge with no duplicates",
			args: args[[]int, int]{source: []int{1, 2, 3}, target: []int{4, 5, 6}},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "merge with duplicates",
			args: args[[]int, int]{source: []int{1, 2, 3}, target: []int{3, 4, 5}},
			want: []int{1, 2, 3, 4, 5}, // 3 appears only once in result
		},
		{
			name: "merge with empty source",
			args: args[[]int, int]{source: []int{}, target: []int{1, 2, 3}},
			want: []int{1, 2, 3},
		},
		{
			name: "merge with empty target",
			args: args[[]int, int]{source: []int{1, 2, 3}, target: []int{}},
			want: []int{1, 2, 3},
		},
		{
			name: "merge with both empty",
			args: args[[]int, int]{source: []int{}, target: []int{}},
			want: []int{},
		},
		{
			name: "merge with overlapping elements",
			args: args[[]int, int]{source: []int{1, 2}, target: []int{2, 3, 1, 4}},
			want: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[[]int, int]:
				got := Merge(args.source, args.target)
				if !Equals(got, tt.want) {
					t.Errorf("Merge() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRemoveAll(t *testing.T) {
	type args[S ~[]T, T comparable] struct {
		source []T
		target []T
	}

	tests := []struct {
		name string
		args interface{}
		want []int
	}{
		{
			name: "remove some elements",
			args: args[[]int, int]{source: []int{1, 2, 3, 4, 5}, target: []int{2, 4}},
			want: []int{1, 3, 5},
		},
		{
			name: "remove all elements",
			args: args[[]int, int]{source: []int{1, 2, 3}, target: []int{1, 2, 3}},
			want: []int{},
		},
		{
			name: "remove none - no matches",
			args: args[[]int, int]{source: []int{1, 2, 3}, target: []int{4, 5}},
			want: []int{1, 2, 3},
		},
		{
			name: "remove from empty slice",
			args: args[[]int, int]{source: []int{}, target: []int{1, 2}},
			want: []int{},
		},
		{
			name: "remove empty target from source",
			args: args[[]int, int]{source: []int{1, 2, 3}, target: []int{}},
			want: []int{1, 2, 3},
		},
		{
			name: "remove with duplicates in target",
			args: args[[]int, int]{source: []int{1, 2, 3, 4, 5}, target: []int{2, 2, 4}},
			want: []int{1, 3, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch args := tt.args.(type) {
			case args[[]int, int]:
				got := RemoveAll(args.source, args.target)
				if !Equals(got, tt.want) {
					t.Errorf("RemoveAll() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMap(t *testing.T) {
	// Test function to double an integer
	double := func(x int) int { return x * 2 }

	// Test function to convert int to string
	intToString := func(x int) string { return string(rune('0' + x)) }

	tests := []struct {
		name     string
		input    []int
		function interface{}
		want     interface{}
	}{
		{
			name:     "map integers to doubled values",
			input:    []int{1, 2, 3, 4},
			function: double,
			want:     []int{2, 4, 6, 8},
		},
		{
			name:     "map integers to strings",
			input:    []int{1, 2, 3},
			function: intToString,
			want:     []string{"1", "2", "3"},
		},
		{
			name:     "map empty slice",
			input:    []int{},
			function: double,
			want:     []int{},
		},
		{
			name:     "map with identity function",
			input:    []int{5, 10, 15},
			function: func(x int) int { return x },
			want:     []int{5, 10, 15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch fn := tt.function.(type) {
			case func(int) int:
				got := Map(tt.input, fn)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Map() = %v, want %v", got, tt.want)
				}
			case func(int) string:
				got := Map(tt.input, fn)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Map() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestFilter(t *testing.T) {
	// Test function to check if number is even
	isEven := func(x int) bool { return x%2 == 0 }

	// Test function to check if number is greater than 3
	greaterThan3 := func(x int) bool { return x > 3 }

	tests := []struct {
		name     string
		input    []int
		function func(int) bool
		want     []int
	}{
		{
			name:     "filter even numbers",
			input:    []int{1, 2, 3, 4, 5, 6},
			function: isEven,
			want:     []int{2, 4, 6},
		},
		{
			name:     "filter numbers greater than 3",
			input:    []int{1, 2, 3, 4, 5},
			function: greaterThan3,
			want:     []int{4, 5},
		},
		{
			name:     "filter with no matches",
			input:    []int{1, 2, 3},
			function: func(x int) bool { return x > 10 },
			want:     []int{},
		},
		{
			name:     "filter with all matches",
			input:    []int{2, 4, 6},
			function: isEven,
			want:     []int{2, 4, 6},
		},
		{
			name:     "filter empty slice",
			input:    []int{},
			function: isEven,
			want:     []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.input, tt.function)
			if !Equals(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
