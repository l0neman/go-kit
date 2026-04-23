package collectionx

import (
	"slices"
)

// Equals compares if the slices contain the same elements
func Equals[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	countMap := make(map[T]int)

	for _, item := range a {
		countMap[item]++
	}

	for _, item := range b {
		countMap[item]--
		if countMap[item] < 0 {
			return false
		}
	}

	for _, count := range countMap {
		if count != 0 {
			return false
		}
	}

	return true
}

// Contains checks if the slice contains the target element
func Contains[S ~[]T, T comparable](slice S, a T) bool {
	return slices.Index(slice, a) != -1
}

// ContainsAll checks if all elements are contained
func ContainsAll[S ~[]T, T comparable](source, target S) bool {
	elementMap := make(map[T]struct{})
	for _, item := range source {
		elementMap[item] = struct{}{}
	}

	for _, item := range target {
		if _, exists := elementMap[item]; !exists {
			return false
		}
	}

	return true
}

// Merge merges slices
func Merge[S ~[]T, T comparable](source S, target S) S {
	sourceSet := NewSetWithSlice(source)

	result := make(S, 0, len(source)+len(target))
	result = append(result, source...)
	for _, item := range target {
		if !sourceSet.Contains(item) {
			result = append(result, item)
		}
	}

	return result
}

// RemoveAll removes elements equal to target from slice
func RemoveAll[S ~[]T, T comparable](source S, target S) S {
	valueMap := NewSetWithSlice(target)

	result := make(S, 0, len(source))
	for _, item := range source {
		if !valueMap.Contains(item) {
			result = append(result, item)
		}
	}

	return result
}

// Map maps the collection to another type of collection
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}

	return result
}

// Filter filters the collection, keeping only elements that satisfy the condition
func Filter[T any](slice []T, f func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}

	return result
}
