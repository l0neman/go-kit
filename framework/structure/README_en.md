# structure module

The core `structure` module provides a generic, advanced, and easy-to-use Go struct parsing utility.

Based on the **Visitor Pattern**, it cleverly decouples complex deep struct traversal logic (such as handling nested structures, pointer dereferencing, slices, and Maps) from specific field processing logic.

Whether you're analyzing project structures, automatically generating code, or implementing configuration loaders or parameter validators, `structure` strips away the tedious reflection logic, allowing you to focus solely on your business logic.

## Features

* **Shields Reflection Complexity**: Built-in comprehensive heavy type checking, recursive traversal branches, pointer dereferencing, and other reflection handling mechanisms.
* **Extremely Simple Interface Definition**: You only need to implement the `Visitor` interface containing `VisitField` and `Recursion` methods to customize everything.
* **Rich Context Capabilities**: `FieldContext` provides extremely convenient APIs (getting tag values, reading original values via `Value()`, modifying values via `Set()`, checking for nil via `IsNil()`, and cross-type conversion mechanisms (smart string conversion to various basic types, etc.)).
* **Built-in Call Support**: Provides independent reflection call method helper utility `CallMethod`.

## Usage Scenarios

This module is extremely simple to use, requiring only implementing `Visitor`, and offers high flexibility.
Below are 3 application scenarios already used in go-kit.

### 1. Checking Parameters in Structs

You can implement a parameter validator (Validator) based on it. For example, reading the `validate` tag on structs to perform null or boundary checks:

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
    Email    string // Optional
}

type ValidatorVisitor struct {
    Err error
}

func (v *ValidatorVisitor) VisitField(ctx *structure.FieldContext) error {
    rule := ctx.Tag("validate")
    if rule == "required" && (ctx.IsNil() || ctx.Value() == "") {
        v.Err = fmt.Errorf("parameter %s is required", ctx.FieldName())
        return structure.ErrStop // Stop further traversal
    }
    
    if rule == "min=18" {
        if age, ok := ctx.Value().(int); ok && age < 18 {
            v.Err = fmt.Errorf("parameter %s does not meet minimum age requirement", ctx.FieldName())
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

### 2. Command-line Parsing Convenience Utility

When you have complex structs that need to be received as command-line argument configurations, you can dynamically collect parameter lists by parsing tags, or pass the struct's address to a specific parsing framework or standard library for flag binding.

```go
package main

import (
    "fmt"
    "github.com/l0neman/go-kit/framework/structure"
)

type ServerCmd struct {
    Host string `flag:"host" desc:"Server listening address"`
    Port int    `flag:"port" desc:"Server listening port"`
}

type FlagBuilderVisitor struct{}

func (f *FlagBuilderVisitor) VisitField(ctx *structure.FieldContext) error {
    flagName := ctx.Tag("flag")
    if flagName == "" {
        return nil
    }
    
    desc := ctx.Tag("desc")
    
    // In actual applications, you can get the pointer via ctx.Addr() and pass it to flag.StringVar(&ptr...) etc.
    fmt.Printf("Register command-line argument: --%s  (%s) - binding target address: %p\n", flagName, desc, ctx.Addr())
    return nil
}

func (f *FlagBuilderVisitor) Recursion(_ *structure.FieldContext) bool {
    return true
}
```

### 3. Complex Processing, Recursive Assignment, and Calling Struct-specified Methods

Combined with `structure`'s support for `Set(value any)` field assignment and flexible multi-type string conversion, it is very suitable for implementing complex configuration initialization and method triggering flows: for example, assigning default values based on the `default` tag, and automatically calling hooks when specific methods like `Init` are present.

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

// Conventional initialization interface or methods checked via reflection
func (d *DBConfig) Init() {
    fmt.Printf("[DBConfig] Init() called. Type=%s, DSN=%s\n", d.Type, d.DSN)
}

type AppConfig struct {
    Name string   `default:"go-app"`
    DB   DBConfig // Nested struct
}

type InitVisitor struct{}

func (i *InitVisitor) VisitField(ctx *structure.FieldContext) error {
    // Complex struct value injection: handling default value assignment
    if ctx.CanSet() && (ctx.IsNil() || isEmptyValue(ctx.Value())) {
        defaultVal := ctx.Tag("default")
        if defaultVal != "" {
            // Set will automatically handle corresponding data types (e.g., passing a string can automatically convert to corresponding int, bool, etc.)
            if err := ctx.Set(defaultVal); err != nil {
                return fmt.Errorf("failed to set default value %s: %w", ctx.FieldName(), err)
            }
        }
    }

    // Calling complex flow lifecycle functions: when struct contains Init method, call it here
    if ctx.Kind() == reflect.Struct {
        // ctx.Addr() gets the struct pointer to call the pointer method
        if ptr := ctx.Addr(); ptr != nil {
            // CallMethod doesn't care what the target is, dynamically query if there's a corresponding method to call
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