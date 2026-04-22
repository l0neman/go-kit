# arg 包

`arg` 包提供了基于结构体的命令行参数解析功能，满足常见的单层命令场景。

本质上是对 flag 包地包装。

可以很方便地解析命令行参数值，以及设置参数默认值和帮助信息。

## 用法

以下是基于结构体解析命令行参数的简单示例：

```go
package main

import (
    "fmt"
    "github.com/l0neman/go-kit/arg"
)

func main() {
    // 对应令行输入为: ./main -host 0.0.0.0 -port 8888 -enable_tls
    // 数据类型与 flag 包一致，string、bool、int、in64、float64
    type Config struct {
        Host      string `name:"host" default:"127.0.0.1" help:"Server host"`
        Port      int    `name:"port" default:"8080" help:"Server port"`
        EnableTLS bool   `name:"enable_tls" default:"false" help:"Enable TLS"`
    }

    ptr := &Config{}
    err := arg.Parse(ptr)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", ptr)
}
```
