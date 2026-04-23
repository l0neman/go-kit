package collectionx

// IsEmpty checks if the slice is empty
func IsEmpty[T any](source []T) bool {
	return source == nil || len(source) == 0
}

// Append adds data to the slice
func Append[T any](slice *[]T, a ...T) {
	*slice = append(*slice, a...)
}

// SliceOf quickly obtains a slice
func SliceOf[T any](a ...T) *[]T {
	return &a
}

// SetOf quickly creates a Set
func SetOf[T comparable](ts ...T) Set[T] {
	return NewSetWithSlice(ts)
}
