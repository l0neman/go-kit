# params 包

`params` 包提供了基于结构体标签（Struct Tag）的声明式参数校验器。
允许将校验规则直接定义在结构体字段上，极大地简化了数据校验逻辑。

## 核心功能 `CheckFields`

核心函数是 `CheckFields(v interface{}) error`。
接收一个结构体（或指向结构体的指针）作为参数，并根据字段上的 `check` 标签递归地校验所有字段。

- 如果所有字段都满足规则，函数返回 `nil`。
- 如果有任何字段不满足规则，函数会返回一个描述第一个失败原因的 `error`。

## 用法

下面是一个校验 API 请求体的例子：

```go
package main

import (
    "fmt"
    "github.com/l0neman/go-kit/params"
)

// Profile 是一个嵌套的结构体
type Profile struct {
    Score float64 `check:">0,<100"`
}

// CreateUserRequest 定义了用户的请求体
// 使用 `check` 标签来定义每个字段的校验规则
type CreateUserRequest struct {
    Username string  `check:"not_empty"`
    Age      int     `check:">=18,<120"`
    Email    string  `check:"not_empty"`
    Profile  *Profile // 支持嵌套结构体指针
}

func main() {
    // --- 案例 1: 一个有效的请求 ---
    validRequest := &CreateUserRequest{
        Username: "johndoe",
        Age:      30,
        Email:    "johndoe@example.com",
        Profile:  &Profile{Score: 99.5},
    }

    if err := params.CheckFields(validRequest); err != nil {
        fmt.Printf("Valid request failed validation: %v\n", err)
    } else {
        fmt.Println("Valid request passed validation!")
    }

    // --- 案例 2: 一个无效的请求 (年龄不符合) ---
    invalidRequest := &CreateUserRequest{
        Username: "jane",
        Age:      17, // 年龄小于 18
        Email:    "jane@example.com",
        Profile:  &Profile{Score: 80},
    }

    if err := params.CheckFields(invalidRequest); err != nil {
        // 错误输出会非常清晰，例如: "字段 'Age' 的值 '17' 不满足条件 '>=18'"
        fmt.Printf("Invalid request failed validation: %v\n", err)
    }
}
```

## 支持的规则

多个规则可以用逗号 `,` 分隔。

| 规则 | 说明 | 类型 |
| :--- | :--- | :--- |
| `not_empty` | 字段不能为空字符串。 | `string` |
| `is_empty` | 字段必须为空字符串。 | `string` |
| `>value` | 必须大于 `value`。 | 数字类型 (`int`, `float64`, etc.) |
| `<value` | 必须小于 `value`。 | 数字类型 |
| `>=value` | 必须大于或等于 `value`。 | 数字类型 |
| `<=value` | 必须小于或等于 `value`。 | 数字类型 |
| `==value` | 必须等于 `value`。 | 数字类型 |
| `!=value` | 必须不等于 `value`。 | 数字类型 |
