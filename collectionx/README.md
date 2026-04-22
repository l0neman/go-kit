# collectionx 包

`collectionx` 包提供了额外的集合数据结构和操作工具。

## Set

`Set` 是一个不包含重复元素的泛型集合，本质上是对 `Map` 的封装。

### 用法

可以重新创建一个 `Set`，也可以直接从 `slice` 或可变参数创建。

```go
package main

import (
    "log"
    "github.com/l0neman/go-kit/collectionx"
)

func main() {
    // 1. 创建一个空 Set
    s1 := collectionx.NewSet[string]()
    s1.Add("apple")
    s1.Add("banana")
    s1.Add("apple") // 重复元素无效

    log.Printf("s1 contains 'apple': %v\n", s1.Contains("apple")) // true
    log.Printf("s1 size: %d\n", s1.Size()) // 2

    s1.Remove("banana")
    log.Printf("s1 contains 'banana' after removal: %v\n", s1.Contains("banana")) // false

    // 2. 从 slice 创建 Set (自动去重)
    slice := []int{1, 2, 3, 2, 1}
    s2 := collectionx.NewSetWithSlice(slice)
    log.Printf("s2 values: %v\n", s2.ToSlice()) // [1, 2, 3] (顺序不定)

    // 3. 使用可变参数创建 Set (自动去重)
    s3 := collectionx.SetOf(10, 20, 30, 20)
    log.Printf("s3 contains 20: %v\n", s3.Contains(20)) // true
    log.Printf("s3 size: %d\n", s3.Size()) // 3
}
```

## 切片操作 (Slice Operations)

### `Filter` 与 `FilterSeq`

根据提供的函数过滤切片元素。

```go
numbers := []int{1, 2, 3, 4, 5, 6}
isEven := func(n int) bool { return n%2 == 0 }
evens := collectionx.Filter(numbers, isEven) // a new slice: [2, 4, 6]

// 搭配迭代器：惰性消费过滤流
// for e := range collectionx.FilterSeq(slices.Values(numbers), isEven) { ... }
```

### `Map` 与 `MapSeq`

将切片中的每个元素应用到转换函数，并返回一个新的切片。

```go
numbers := []int{1, 2, 3}
double := func(n int) int { return n * 2 }
doubled := collectionx.Map(numbers, double) // a new slice: [2, 4, 6]
```

### `Contains` / `ContainsAll`

检查切片是否包含某个元素或所有指定元素。

```go
items := []string{"a", "b", "c"}
hasB := collectionx.Contains(items, "b") // true
hasAll := collectionx.ContainsAll(items, []string{"a", "c"}) // true
hasD := collectionx.ContainsAll(items, []string{"a", "d"}) // false
```

### `Merge` / `RemoveAll`

合并两个切片（去重）或从一个切片中移除另一个切片的所有元素。

```go
sliceA := []int{1, 2, 3}
sliceB := []int{3, 4, 5}

merged := collectionx.Merge(sliceA, sliceB) // [1, 2, 3, 4, 5] (顺序不定)

remaining := collectionx.RemoveAll(sliceA, []int{2, 3}) // [1]
```