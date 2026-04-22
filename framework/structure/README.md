# structure 模块

核心模块 `structure` 提供了一个通用、高级且易于使用的 Go 结构体解析工具。

基于 **访问者模式 (Visitor Pattern)**，巧妙地将复杂的结构体深度遍历
（如处理嵌套、指针解引用、切片和 Map 等复杂结构）逻辑与具体的字段处理逻辑解耦。

无论是分析项目结构、自动生成代码，还是实现配置加载器、参数校验器，
`structure` 都能将枯燥的反射逻辑剥离，让你只专注于你的业务逻辑。

## 特性

* **屏蔽反射复杂性**：内置了完整的繁重类型判断、递归遍历分支、取值指针解引用等反射处理机制。
* **极其简单的接口定义**：你只需要实现含有 `VisitField` 和 `Recursion` 方法的 `Visitor` 接口即可定制一切。
* **丰富的上下文能力**：`FieldContext` 提供了极为简便的 API（获取 Tag 设值、读取原值 `Value()`、修改赋值 `Set()`、判断空值 `IsNil()` 以及跨类型转换机制（string 智能转换为各种基础类型等））。
* **内置调用支持**：提供独立的反射调用方法辅助工具 `CallMethod`。

## 应用场景

这个模块使用起来极简，仅需实现 `Visitor`，具有很高的灵活性。
以下列举 go-kit 中已使用到的 3 个应用场景。

### 1. 检查结构体中的参数

可以基于它实现一个参数校验器（Validator）。例如读取结构体的 `validate` tag 进行对应要求的空值或者边界检查：

```go
package main

import (
    "errors"
    "fmt"
    "github.com/l0neman/go-kit/framework/structure"
)

type User struct {
    Username string `validate:"required"`
    Age      int    `validate:"min=18"`
    Email    string // 选填
}

type ValidatorVisitor struct {
    Err error
}

func (v *ValidatorVisitor) VisitField(ctx *structure.FieldContext) error {
    rule := ctx.Tag("validate")
    if rule == "required" && (ctx.IsNil() || ctx.Value() == "") {
        v.Err = fmt.Errorf("参数 %s 是必填项", ctx.FieldName())
        return structure.ErrStop // 停止后续遍历
    }
    
    if rule == "min=18" {
        if age, ok := ctx.Value().(int); ok && age < 18 {
            v.Err = fmt.Errorf("参数 %s 不满足最小年龄要求", ctx.FieldName())
            return structure.ErrStop
        }
    }

    return nil
}

func (v *ValidatorVisitor) Recursion(_ *structure.FieldContext) bool {
    return true
}

func Validate(obj any) error {
    visitor := &ValidatorVisitor{}
    parser := structure.NewParser(visitor)
    if err := parser.Parse(obj); err != nil {
        return err
    }

    return visitor.Err
}
```

### 2. 命令行解析便利工具

有复杂的结构需要被当作命令行参数配置接收时，可以动态通过解析 tag 收集参数列表，
或将结构体的地址传给特定的解析框架或标准库做 flag 绑定。

```go
package main

import (
    "fmt"
    "github.com/l0neman/go-kit/framework/structure"
)

type ServerCmd struct {
    Host string `flag:"host" desc:"服务器监听地址"`
    Port int    `flag:"port" desc:"服务器监听端口"`
}

type FlagBuilderVisitor struct{}

func (f *FlagBuilderVisitor) VisitField(ctx *structure.FieldContext) error {
    flagName := ctx.Tag("flag")
    if flagName == "" {
        return nil
    }
    
    desc := ctx.Tag("desc")
    
    // 在实际应用中，可以通过 ctx.Addr() 获取指针并传递给 flag.StringVar(&ptr...) 等方法
    fmt.Printf("注册命令行参数: --%s  (%s) - 绑定目标地址: %p\n", flagName, desc, ctx.Addr())
    return nil
}

func (f *FlagBuilderVisitor) Recursion(_ *structure.FieldContext) bool {
    return true
}
```

### 3. 复杂处理，递归赋值和调用结构体规定方法

结合 `structure` 对 `Set(value any)` 字段赋值和多类型字符串转换的支持包容性，十分适合实现复杂的配置初始化与方法触发流程：例如根据 `default` 标签赋予默认值，并根据是否含有 `Init` 等特定方法进行钩子的自动调用。

```go
package main

import (
    "fmt"
    "github.com/l0neman/go-kit/framework/structure"
    "reflect"
)

type DBConfig struct {
    Type string `default:"mysql"`
    DSN  string `default:"root:123456@tcp(127.0.0.1:3306)/test"`
}

// 约定的初始化接口或根据反射检查存在的方法
func (d *DBConfig) Init() {
    fmt.Printf("[DBConfig] Init() called. Type=%s, DSN=%s\n", d.Type, d.DSN)
}

type AppConfig struct {
    Name string   `default:"go-app"`
    DB   DBConfig // 嵌套结构
}

type InitVisitor struct{}

func (i *InitVisitor) VisitField(ctx *structure.FieldContext) error {
    // 复杂的结构体值注入：处理默认值赋予
    if ctx.CanSet() && (ctx.IsNil() || isEmptyValue(ctx.Value())) {
        defaultVal := ctx.Tag("default")
        if defaultVal != "" {
            // Set 会自动处理对应数据类型（比如传入字符串能自动转化到对应的 int、bool 等）
            if err := ctx.Set(defaultVal); err != nil {
                return fmt.Errorf("设置默认值失败 %s: %w", ctx.FieldName(), err)
            }
        }
    }

    // 复杂流程生命周期函数的调用：当结构体包含了 Init 方法时在此处调用
    if ctx.Kind() == reflect.Struct {
        // ctx.Addr() 获取到结构体的指针从而调用指针方法
        if ptr := ctx.Addr(); ptr != nil {
            // CallMethod 不关心目标是什么对象，动态查询是否有对应方法进行调用
            _, _ = structure.CallMethod(ptr, "Init")
        }
    }

    return nil
}

func (i *InitVisitor) Recursion(_ *structure.FieldContext) bool {
    return true
}

func isEmptyValue(v any) bool {
    return v == "" || v == 0 || v == false || v == nil
}
```

