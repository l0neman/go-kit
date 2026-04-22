package collectionx

// 定义一个空结构体作为 map 的值类型，因为它不占用内存空间
type void struct{}

type Set[T comparable] interface {
	Add(value T)
	Remove(value T)
	Contains(value T) bool
	Size() int
	Clear()
	ToSlice() []T
}

// MapSet 实现 Set
type MapSet[T comparable] struct {
	Values map[T]void `json:"Values"`
}

// NewSet 创建一个新的 Set
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

// Add 向 Set 中添加元素
func (ms *MapSet[T]) Add(value T) {
	ms.Values[value] = void{}
}

// Remove 从 Set 中删除元素
func (ms *MapSet[T]) Remove(value T) {
	delete(ms.Values, value)
}

// Contains 检查元素是否存在于 Set 中
func (ms *MapSet[T]) Contains(value T) bool {
	_, ok := ms.Values[value]
	return ok
}

// Size 返回 Set 中的元素数量
func (ms *MapSet[T]) Size() int {
	return len(ms.Values)
}

// Clear 清空 Set
func (ms *MapSet[T]) Clear() {
	ms.Values = make(map[T]void)
}

// ToSlice 将 Set 中的元素转换为切片
func (ms *MapSet[T]) ToSlice() []T {
	slice := make([]T, 0, len(ms.Values))
	for key := range ms.Values {
		slice = append(slice, key)
	}

	return slice
}
