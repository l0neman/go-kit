# stringx 包

`stringx` 包提供了一些用于检查空字符串的便捷辅助函数，作为对 Go 标准 `strings` 包的补充。

## `IsEmpty(str string) bool`

`IsEmpty` 函数检查一个字符串的长度是否为 0。这是一个 `len(str) == 0` 的简单别名。

**注意：** 此函数不认为仅包含空白字符（如空格、制表符）的字符串为空。

## `HasEmpty(strings ...string) bool`

`HasEmpty` 函数接收一个或多个字符串作为参数，如果其中 **任何一个** 字符串为空（根据 `IsEmpty` 的定义），则返回 `true`。

对于一次性校验多个输入字段非常有用。

## 用法示例

```go
package main

import (
    "log"
    "github.com/l0neman/go-kit/stringx"
)

func main() {
    // --- IsEmpty 示例 ---
    log.Printf("'' is empty: %v\n", stringx.IsEmpty("")) // true
    log.Printf("' ' is empty: %v\n", stringx.IsEmpty(" "))  // false
    log.Printf("'hello' is empty: %v\n", stringx.IsEmpty("hello")) // false

    log.Println("---")

    // --- HasEmpty 示例 ---
    // 检查一组有效的输入
    validInputs := []string{"user", "pass123", "token"}
    if stringx.HasEmpty(validInputs...) {
        log.Println("错误：输入中包含空字段")
    } else {
        log.Println("所有输入均有效") // Will print this
    }

    // 检查一组包含空字符串的输入
    invalidInputs := []string{"user", "", "token"} // 第二个是空字符串
    if stringx.HasEmpty(invalidInputs...) {
        log.Println("错误：输入中包含空字段") // Will print this
    } else {
        log.Println("所有输入均有效")
    }
}
```