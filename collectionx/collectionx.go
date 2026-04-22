package collectionx

// IsEmpty 切片是否为空
func IsEmpty[T any](source []T) bool {
	return source == nil || len(source) == 0
}

// Append 将数据添加到切片
func Append[T any](slice *[]T, a ...T) {
	*slice = append(*slice, a...)
}

// SliceOf 快速获取切片
func SliceOf[T any](a ...T) *[]T {
	return &a
}

// SetOf 快速创建 Set
func SetOf[T comparable](ts ...T) Set[T] {
	return NewSetWithSlice(ts)
}
