package collectionx

import (
	"slices"
)

// Equals 比较是否存在相同的元素
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

// Contains 是否包含目标元素
func Contains[S ~[]T, T comparable](slice S, a T) bool {
	return slices.Index(slice, a) != -1
}

// ContainsAll 包含所有元素
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

// Merge 合并切片
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

// RemoveAll 移除切片中包含的与 target 相等的元素
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

// Map 映射操作，将集合映射为另一种类型的集合
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}

	return result
}

// Filter 过滤操作，仅保留满足条件的元素
func Filter[T any](slice []T, f func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}

	return result
}
