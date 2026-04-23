package collectionx

// Defines an empty struct as the value type for map, since it takes no memory space
type void struct{}

type Set[T comparable] interface {
	Add(value T)
	Remove(value T)
	Contains(value T) bool
	Size() int
	Clear()
	ToSlice() []T
}

// MapSet implements Set
type MapSet[T comparable] struct {
	Values map[T]void `json:"Values"`
}

// NewSet creates a new Set
func NewSet[T comparable]() Set[T] {
	return &MapSet[T]{
		Values: make(map[T]void),
	}
}

func NewSetWithSlice[T comparable](ts []T) Set[T] {
	values := make(map[T]void)

	for _, t := range ts {
		values[t] = void{}
	}

	return &MapSet[T]{
		Values: values,
	}
}

// Add adds an element to the Set
func (ms *MapSet[T]) Add(value T) {
	ms.Values[value] = void{}
}

// Remove removes an element from the Set
func (ms *MapSet[T]) Remove(value T) {
	delete(ms.Values, value)
}

// Contains checks if an element exists in the Set
func (ms *MapSet[T]) Contains(value T) bool {
	_, ok := ms.Values[value]
	return ok
}

// Size returns the number of elements in the Set
func (ms *MapSet[T]) Size() int {
	return len(ms.Values)
}

// Clear clears the Set
func (ms *MapSet[T]) Clear() {
	ms.Values = make(map[T]void)
}

// ToSlice converts elements in the Set to a slice
func (ms *MapSet[T]) ToSlice() []T {
	slice := make([]T, 0, len(ms.Values))
	for key := range ms.Values {
		slice = append(slice, key)
	}

	return slice
}
