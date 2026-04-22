# errorx 包

`errorx` 包提供用于错误包装和格式化的便捷工具函数。

## 设计目标

错误多级返回时，如果每次直接返回底层函数的错误，通常缺乏上层业务的上下文信息。
例如，一个数据库操作可能返回 `sql.ErrNoRows`，但对于调用者来说，更希望看到如
“获取用户信息失败 > record not found” 这样的完整信息链。

本包以简单的方式为错误添加上下文，同时保留原始错误的类型信息，
以便上层代码可以使用 `errors.Is` 或 `errors.As` 进行判断。

## 核心函数

### `Wrap(err error, desc string) error`

`Wrap` 函数将一个描述性的字符串 `desc` 添加到现有 `err` 的前面，生成一个新的错误。
等效于 `fmt.Errorf("%s > %w", desc, err)`，将原始 `err` “包装”起来，保留原始错误的类型信息。

#### 用法

```go
package main

import (
    "database/sql"
    "errors"
    "log"
    "fmt"
    "github.com/l0neman/go-kit/errorx"
)

// mock a data access function
func fetchUserFromDB(id int) error {
    // 模拟数据库未找到记录
    return sql.ErrNoRows
}

// service layer function
func getUser(id int) error {
    err := fetchUserFromDB(id)
    if err != nil {
        // 使用 Wrap 添加上下文信息
        return errorx.Wrap(err, fmt.Sprintf("failed to get user with id %d", id))
    }
	
    return nil
}

func main() {
    err := getUser(123)

    if err != nil {
        // 打印完整的错误链
        // 输出: failed to get user with id 123 > sql: no rows in result set
		log.Println(err)

        // 因为 Wrap 保留了原始错误，我们可以使用 errors.Is 进行检查
        if errors.Is(err, sql.ErrNoRows) {
			log.Println("Confirmed: The user was not found in the database.")
        }
    }
}
```

### `Wrapf(err error, format string, a ...any) error`

`Wrapf` 函数与 `Wrap` 类似，但允许使用格式化字符串来生成描述信息。等效于 `fmt.Errorf("%s > %w", fmt.Sprintf(format, a...), err)`。

#### 用法

```go
package main

import (
    "database/sql"
    "errors"
    "fmt"
	"log"
    "github.com/l0neman/go-kit/errorx"
)

func getUserWithAge(id int, age int) error {
    // 模拟数据库错误
    err := sql.ErrNoRows
    if err != nil {
        // 使用 Wrapf 添加上下文信息，支持格式化
        return errorx.Wrapf(err, "failed to get user with id %d and age %d", id, age)
    }
	
    return nil
}

func main() {
    err := getUserWithAge(123, 30)

    if err != nil {
        // 打印完整的错误链
        // 输出: failed to get user with id 123 and age 30 > sql: no rows in result set
		log.Println(err)

        // 仍然可以使用 errors.Is 进行类型检查
        if errors.Is(err, sql.ErrNoRows) {
			log.Println("Confirmed: The user was not found in the database.")
        }
    }
}
```

### `Wraps(err error, desc string) string`

`Wraps` 函数用于生成一个格式化的 **字符串**，不进行错误包装。等效于 `fmt.Sprintf("%s > %v", desc, err)`。

当你仅仅需要一个包含上下文的错误描述（例如用于日志记录），而不需要保留原始错误类型时，可以使用此函数。

#### 用法

```go
package main

import (
    "fmt"
    "log"
    "github.com/l0neman/go-kit/errorx"
)

func main() {
    originalErr := fmt.Errorf("timeout")
    
    // 生成一个用于日志的字符串
    logMessage := errorx.Wraps(originalErr, "network operation failed")
    
    // 输出: network operation failed > timeout
    log.Println(logMessage)
}
```

### `Wrapfs(err error, format string, a ...any) string`

`Wrapfs` 函数与 `Swrap` 类似，但允许使用格式化字符串来生成描述信息。它等效于 `fmt.Sprintf("%s > %v", fmt.Sprintf(format, a...), err)`。

#### 用法

```go
package main

import (
	"fmt"
    "log"
    "github.com/l0neman/go-kit/errorx"
)

func main() {
    originalErr := fmt.Errorf("timeout")
    
    // 生成一个用于日志的格式化字符串
    logMessage := errorx.Wrapfs(originalErr, "network operation failed for host %s", "example.com")
    
    // 输出: network operation failed for host example.com > timeout
    log.Println(logMessage)
}
```