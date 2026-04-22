# debug 包

`debug` 包提供了用于调试的工具，例如稀疏日志打印。

## 用法

```go
package main

import (
    "log"
    "github.com/l0neman/go-kit/debug"
)

func main() {
    // 智能节流打印日志，对于大量重复的日志，只打印 2 次方次数的日志
    for i := 0; i < 1024; i++ {
        debug.PrintThrottled(func(count int) {
            log.Println("log 1 打印次数", count)
        })
    }
}
```
