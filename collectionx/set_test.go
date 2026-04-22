package collectionx

import (
	"testing"
)

func TestNewSet(t *testing.T) {
	s := NewSet[int]()

	// 验证初始状态
	if s.Size() != 0 {
		t.Errorf("Initial size should be 0, got %d", s.Size())
	}

	if s.Contains(1) {
		t.Error("New set should not contain any elements")
	}
}

func TestSet_AddAndContains(t *testing.T) {
	s := NewSet[int]()

	// 测试添加元素
	s.Add(1)
	s.Add(2)
	s.Add(3)

	// 测试包含元素
	if !s.Contains(1) {
		t.Error("Set should contain 1")
	}
	if !s.Contains(2) {
		t.Error("Set should contain 2")
	}
	if !s.Contains(3) {
		t.Error("Set should contain 3")
	}

	// 测试不包含的元素
	if s.Contains(4) {
		t.Error("Set should not contain 4")
	}
}

func TestSet_Remove(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	// 删除元素
	s.Remove(2)

	// 验证删除后的状态
	if !s.Contains(1) {
		t.Error("Set should still contain 1")
	}
	if s.Contains(2) {
		t.Error("Set should not contain 2 after removal")
	}
	if !s.Contains(3) {
		t.Error("Set should still contain 3")
	}
}

func TestSet_Size(t *testing.T) {
	s := NewSet[int]()

	// 初始大小应为 0
	if s.Size() != 0 {
		t.Errorf("Initial size should be 0, got %d", s.Size())
	}

	// 添加元素后
	s.Add(1)
	s.Add(2)
	if s.Size() != 2 {
		t.Errorf("Size should be 2, got %d", s.Size())
	}

	// 删除元素后
	s.Remove(1)
	if s.Size() != 1 {
		t.Errorf("Size should be 1 after removal, got %d", s.Size())
	}
}

func TestSet_Clear(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)

	// 清空集合
	s.Clear()

	// 验证清空后的状态
	if s.Size() != 0 {
		t.Errorf("Size should be 0 after clear, got %d", s.Size())
	}
	if s.Contains(1) || s.Contains(2) {
		t.Error("Set should be empty after clear")
	}
}

func TestSet_ToSlice(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	slice := s.ToSlice()

	// 验证切片长度
	if len(slice) != 3 {
		t.Errorf("Slice length should be 3, got %d", len(slice))
	}

	// 验证切片包含所有元素
	found := make(map[int]bool)
	for _, v := range slice {
		found[v] = true
	}

	if !found[1] || !found[2] || !found[3] {
		t.Error("Slice should contain all elements from the set")
	}
}

func TestNewSetWithSlice(t *testing.T) {
	// 使用切片创建集合
	s := NewSetWithSlice([]int{1, 2, 3, 3}) // 注意有重复元素

	// 验证集合大小（应该自动去重）
	if s.Size() != 3 {
		t.Errorf("Set size should be 3 (duplicates removed), got %d", s.Size())
	}

	// 验证包含的元素
	if !s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
		t.Error("Set should contain elements 1, 2, and 3")
	}
}

func TestSetOf(t *testing.T) {
	// 使用可变参数创建集合
	s := SetOf(1, 2, 3, 3) // 注意有重复元素

	// 验证集合大小（应该自动去重）
	if s.Size() != 3 {
		t.Errorf("Set size should be 3 (duplicates removed), got %d", s.Size())
	}

	// 验证包含的元素
	if !s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
		t.Error("Set should contain elements 1, 2, and 3")
	}
}

func TestSet_StringType(t *testing.T) {
	s := NewSet[string]()
	s.Add("apple")
	s.Add("banana")

	if !s.Contains("apple") {
		t.Error("Set should contain 'apple'")
	}
	if !s.Contains("banana") {
		t.Error("Set should contain 'banana'")
	}
	if s.Contains("orange") {
		t.Error("Set should not contain 'orange'")
	}
}

func TestSet_FloatType(t *testing.T) {
	s := NewSet[float64]()
	s.Add(3.14)
	s.Add(2.71)

	if !s.Contains(3.14) {
		t.Error("Set should contain 3.14")
	}
	if !s.Contains(2.71) {
		t.Error("Set should contain 2.71")
	}
	if s.Contains(1.618) {
		t.Error("Set should not contain 1.618")
	}
}
