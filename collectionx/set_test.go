package collectionx

import (
	"testing"
)

func TestNewSet(t *testing.T) {
	s := NewSet[int]()

	// Verify initial state
	if s.Size() != 0 {
		t.Errorf("Initial size should be 0, got %d", s.Size())
	}

	if s.Contains(1) {
		t.Error("New set should not contain any elements")
	}
}

func TestSet_AddAndContains(t *testing.T) {
	s := NewSet[int]()

	// Test adding elements
	s.Add(1)
	s.Add(2)
	s.Add(3)

	// Test containing elements
	if !s.Contains(1) {
		t.Error("Set should contain 1")
	}
	if !s.Contains(2) {
		t.Error("Set should contain 2")
	}
	if !s.Contains(3) {
		t.Error("Set should contain 3")
	}

	// Test elements not contained
	if s.Contains(4) {
		t.Error("Set should not contain 4")
	}
}

func TestSet_Remove(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	// Remove element
	s.Remove(2)

	// Verify state after removal
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

	// Initial size should be 0
	if s.Size() != 0 {
		t.Errorf("Initial size should be 0, got %d", s.Size())
	}

	// After adding elements
	s.Add(1)
	s.Add(2)
	if s.Size() != 2 {
		t.Errorf("Size should be 2, got %d", s.Size())
	}

	// After removing elements
	s.Remove(1)
	if s.Size() != 1 {
		t.Errorf("Size should be 1 after removal, got %d", s.Size())
	}
}

func TestSet_Clear(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)

	// Clear set
	s.Clear()

	// Verify state after clearing
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

	// Verify slice length
	if len(slice) != 3 {
		t.Errorf("Slice length should be 3, got %d", len(slice))
	}

	// Verify slice contains all elements
	found := make(map[int]bool)
	for _, v := range slice {
		found[v] = true
	}

	if !found[1] || !found[2] || !found[3] {
		t.Error("Slice should contain all elements from the set")
	}
}

func TestNewSetWithSlice(t *testing.T) {
	// Create set from slice
	s := NewSetWithSlice([]int{1, 2, 3, 3}) // Note: has duplicate elements

	// Verify set size (should automatically remove duplicates)
	if s.Size() != 3 {
		t.Errorf("Set size should be 3 (duplicates removed), got %d", s.Size())
	}

	// Verify contained elements
	if !s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
		t.Error("Set should contain elements 1, 2, and 3")
	}
}

func TestSetOf(t *testing.T) {
	// Create set from variadic parameters
	s := SetOf(1, 2, 3, 3) // Note: has duplicate elements

	// Verify set size (should automatically remove duplicates)
	if s.Size() != 3 {
		t.Errorf("Set size should be 3 (duplicates removed), got %d", s.Size())
	}

	// Verify contained elements
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
