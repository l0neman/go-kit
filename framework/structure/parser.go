package structure

// parser 提供了一个通用的 Go 结构体解析工具。
//
// 核心设计是访问者模式（Visitor Pattern)，它将结构体的遍历逻辑与字段的处理逻辑解耦。
//
// 主要组件:
//  - Parser: 负责深度遍历一个结构体，处理递归、指针、切片等复杂情况。
//  - Visitor: 一个接口，用户通过实现它来注入自定义的字段处理逻辑。
//  - FieldContext: 在遍历过程中，为每个字段提供一个上下文对象，它封装了反射的复杂性，
//    并提供了简单易用的 API 来获取标签、读写值等。

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// 定义了几个用于控制解析流程的特殊错误。
// Visitor 可以通过返回这些错误来与 Parser 通信。
var (
	// ErrStop 指示 Parser 停止整个解析过程。
	ErrStop = fmt.Errorf("停止解析")
	// ErrSkipRecursive 指示 Parser 跳过对当前字段的递归解析。
	// 这对于避免无限循环或处理自定义序列化类型很有用。
	ErrSkipRecursive = fmt.Errorf("跳过递归解析")
)

// FieldContext 持有正在被访问字段的所有上下文信息。
// 它提供了一系列辅助方法，以抽象掉反射的复杂性，让用户可以更安全、更便捷地操作字段。
type FieldContext struct {
	fieldValue  reflect.Value
	structField reflect.StructField
	parentValue reflect.Value
	currentPath string
}

// Path 返回当前字段的点分路径，例如："Config.Server.Port"。
func (c *FieldContext) Path() string {
	return c.currentPath
}

// FieldName 返回结构体字段的原始名称。
func (c *FieldContext) FieldName() string {
	return c.structField.Name
}

// Tag 获取指定键 key 的标签 tag 值。如果标签不存在，则返回空字符串。
func (c *FieldContext) Tag(key string) string {
	return c.structField.Tag.Get(key)
}

// Value 将字段的当前值作为 any 返回。
func (c *FieldContext) Value() any {
	if c.fieldValue.IsValid() {
		return c.fieldValue.Interface()
	}

	return nil
}

// Kind 返回字段的 reflect.Kind 类型。
func (c *FieldContext) Kind() reflect.Kind {
	return c.fieldValue.Kind()
}

// Type 返回字段的 reflect.Type 类型。
func (c *FieldContext) Type() reflect.Type {
	return c.fieldValue.Type()
}

// Addr 返回字段值的指针，作为 any。
// 这对于需要传递字段地址给其他库（如标准库 flag）的场景非常有用。
func (c *FieldContext) Addr() any {
	if c.fieldValue.CanAddr() {
		return c.fieldValue.Addr().Interface()
	}

	return nil
}

// IsNil 检查字段的底层值是否为 nil。
// 对指针、map、切片、接口、通道和函数类型有效。
func (c *FieldContext) IsNil() bool {
	if !c.fieldValue.IsValid() {
		return true
	}

	switch c.fieldValue.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface, reflect.Chan, reflect.Func:
		return c.fieldValue.IsNil()
	default:
		// ignore
	}

	return false
}

// CanSet 报告该字段的值是否可以被修改。
// 通常，只有导出的字段（首字母大写）才可以被设置。
func (c *FieldContext) CanSet() bool {
	return c.fieldValue.CanSet()
}

// Set 尝试设置字段的值。
// 它可以处理从字符串到多种基本类型的自动转换。
// 这是Visitor实现初始化或修改字段值的核心方法。
func (c *FieldContext) Set(value any) error {
	if !c.CanSet() {
		return fmt.Errorf("字段 %s 不可设置（可能未导出）", c.FieldName())
	}

	valToSet := reflect.ValueOf(value)
	// 检查类型是否可以直接赋值
	if valToSet.IsValid() && valToSet.Type().AssignableTo(c.fieldValue.Type()) {
		c.fieldValue.Set(valToSet)
		return nil
	}

	// 如果传入的是字符串，尝试进行类型转换
	if strVal, ok := value.(string); ok {
		return c.setString(strVal)
	}

	return fmt.Errorf("无法将类型 %T 的值赋给字段 %s (期望类型 %s)", value, c.FieldName(), c.fieldValue.Type())
}

// setString 是 Set 方法的辅助函数，专门处理从字符串到目标类型的转换。
func (c *FieldContext) setString(value string) error {
	switch c.fieldValue.Kind() {
	case reflect.String:
		c.fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		if c.fieldValue.OverflowInt(intVal) {
			return fmt.Errorf("值 %s 对于字段 %s 溢出", value, c.FieldName())
		}

		c.fieldValue.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		if c.fieldValue.OverflowUint(uintVal) {
			return fmt.Errorf("值 %s 对于字段 %s 溢出", value, c.FieldName())
		}

		c.fieldValue.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		if c.fieldValue.OverflowFloat(floatVal) {
			return fmt.Errorf("值 %s 对于字段 %s 溢出", value, c.FieldName())
		}

		c.fieldValue.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		c.fieldValue.SetBool(boolVal)
	default:
		return fmt.Errorf("不支持从字符串到字段 %s（类型: %s）的转换", c.FieldName(), c.fieldValue.Kind())
	}

	return nil
}

// Visitor 定义了结构体字段处理器的接口。
// 用户需要实现这个接口来定义自己的业务逻辑。
type Visitor interface {
	// VisitField 在 Parser 遍历到每个字段时被调用。
	// ctx 参数包含了该字段的所有相关信息和操作方法。
	// 实现者可以通过返回特定的错误（如 ErrStop）来控制解析流程。
	VisitField(ctx *FieldContext) error
	// Recursion 是否递归此成员
	Recursion(ctx *FieldContext) bool
}

// Parser 是结构体解析的核心。它持有访问者列表并负责驱动整个遍历过程。
type Parser struct {
	visitor Visitor
}

// NewParser 创建一个新的 Parser 实例。
// 它接收一个 Visitor 切片和一系列的 ParserOption 作为配置。
func NewParser(visitor Visitor) *Parser {
	p := &Parser{
		visitor: visitor,
	}

	return p
}

// Parse 对给定的对象（必须是指向结构体的指针）开始解析过程。
func (p *Parser) Parse(obj any) error {
	if obj == nil {
		return fmt.Errorf("待解析的对象不能为nil")
	}

	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("对象必须是指向结构体的指针")
	}

	// 从结构体的根节点开始遍历
	return p.walk(val.Elem(), val.Elem().Type().Name())
}

// walk 是内部的递归函数，负责深度遍历值的结构。
func (p *Parser) walk(val reflect.Value, path string) error {
	// 如果是指针，先解引用
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil // nil 指针无法继续深入，直接返回
		}

		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			structField := typ.Field(i)

			if !fieldVal.CanSet() {
				continue
			}

			ctx := &FieldContext{
				fieldValue:  fieldVal,
				structField: structField,
				parentValue: val,
				currentPath: fmt.Sprintf("%s.%s", path, structField.Name),
			}

			// 依次调用所有访问者
			err := p.visitor.VisitField(ctx)
			if err != nil {
				if errors.Is(err, ErrStop) || errors.Is(err, ErrSkipRecursive) {
					return err // 将控制流错误向上传递
				}

				// 对于其他业务错误，直接返回
				return fmt.Errorf("在字段 %s 上发生错误: %w", ctx.Path(), err)
			}

			if p.visitor.Recursion(ctx) {
				// 关键点：在 Visitor 执行后，重新获取字段的值，
				// 因为 Visitor 可能已经初始化了一个 nil 指针。
				updatedFieldVal := val.Field(i)
				if err := p.walk(updatedFieldVal, ctx.currentPath); err != nil {
					if errors.Is(err, ErrStop) {
						return err // 如果深层递归要求停止，则继续向上传递
					}

					// 其他递归错误可以根据需要处理，这里选择直接返回
					return err
				}
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if err := p.walk(val.Index(i), fmt.Sprintf("%s[%d]", path, i)); err != nil {
				if errors.Is(err, ErrStop) {
					return err
				}

				return err
			}
		}

	case reflect.Map:
		// 注意: 遍历 map 的顺序是不确定的
		for _, key := range val.MapKeys() {
			if err := p.walk(val.MapIndex(key), fmt.Sprintf("%s[%v]", path, key.Interface())); err != nil {
				if errors.Is(err, ErrStop) {
					return err
				}

				return err
			}
		}
	default:
		// ignore
	}

	return nil
}

// CallMethod 是一个独立的工具函数，用于通过反射调用对象的方法。
// 它不属于 Parser 的核心遍历逻辑，但与结构体操作密切相关，因此放在这里作为辅助。
func CallMethod(obj any, methodName string, args ...any) ([]reflect.Value, error) {
	val := reflect.ValueOf(obj)
	method := val.MethodByName(methodName)

	if !method.IsValid() {
		return nil, fmt.Errorf("在类型 %T 上未找到方法 %s", obj, methodName)
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	if method.Type().NumIn() != len(in) {
		return nil, fmt.Errorf("方法 %s 参数数量不匹配：需要 %d, 提供了 %d", methodName, method.Type().NumIn(), len(in))
	}

	return method.Call(in), nil
}
