# collectionx package

The `collectionx` package provides additional collection data structures and operation utilities.

## Set

`Set` is a generic collection that does not contain duplicate elements, essentially a wrapper around `Map`.

### Usage

You can create a new `Set`, or create one directly from a slice or variadic parameters.

```go
package main

import (
    "log"
    "github.com/l0neman/go-kit/collectionx"
)

func main() {
    // 1. Create an empty Set
    s1 := collectionx.NewSet[string]()
    s1.Add("apple")
    s1.Add("banana")
    s1.Add("apple") // Duplicate element has no effect

    log.Printf("s1 contains 'apple': %v\n", s1.Contains("apple")) // true
    log.Printf("s1 size: %d\n", s1.Size()) // 2

    s1.Remove("banana")
    log.Printf("s1 contains 'banana' after removal: %v\n", s1.Contains("banana")) // false

    // 2. Create Set from slice (automatically deduplicates)
    slice := []int{1, 2, 3, 2, 1}
    s2 := collectionx.NewSetWithSlice(slice)
    log.Printf("s2 values: %v\n", s2.ToSlice()) // [1, 2, 3] (order not guaranteed)

    // 3. Create Set using variadic parameters (automatically deduplicates)
    s3 := collectionx.SetOf(10, 20, 30, 20)
    log.Printf("s3 contains 20: %v\n", s3.Contains(20)) // true
    log.Printf("s3 size: %d\n", s3.Size()) // 3
}
```

## Slice Operations

### `Filter` and `FilterSeq`

Filters slice elements based on the provided function.

```go
numbers := []int{1, 2, 3, 4, 5, 6}
isEven := func(n int) bool { return n%2 == 0 }
evens := collectionx.Filter(numbers, isEven) // a new slice: [2, 4, 6]

// Use with iterator: lazy consuming filtered stream
// for e := range collectionx.FilterSeq(slices.Values(numbers), isEven) { ... }
```

### `Map` and `MapSeq`

Applies a transformation function to each element in the slice and returns a new slice.

```go
numbers := []int{1, 2, 3}
double := func(n int) int { return n * 2 }
doubled := collectionx.Map(numbers, double) // a new slice: [2, 4, 6]
```

### `Contains` / `ContainsAll`

Checks if a slice contains a specific element or all specified elements.

```go
items := []string{"a", "b", "c"}
hasB := collectionx.Contains(items, "b") // true
hasAll := collectionx.ContainsAll(items, []string{"a", "c"}) // true
hasD := collectionx.ContainsAll(items, []string{"a", "d"}) // false
```

### `Merge` / `RemoveAll`

Merges two slices (deduplicates) or removes all elements of one slice from another.

```go
sliceA := []int{1, 2, 3}
sliceB := []int{3, 4, 5}

merged := collectionx.Merge(sliceA, sliceB) // [1, 2, 3, 4, 5] (order not guaranteed)

remaining := collectionx.RemoveAll(sliceA, []int{2, 3}) // [1]
```