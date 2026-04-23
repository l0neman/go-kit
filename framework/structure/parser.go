package structure

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// parser provides a generic Go struct parsing tool.
//
// The core design is the Visitor Pattern, which decouples the struct traversal logic from field handling logic.
//
// Main components:
//  - Parser: Responsible for deep traversal of a struct, handling recursion, pointers, slices, and other complex cases.
//  - Visitor: An interface, through which users inject custom field handling logic.
//  - FieldContext: During traversal, provides a context object for each field, encapsulating reflection complexity,
//    and provides easy-to-use API for getting tags, reading/writing values, etc.

// Defines several special errors for controlling the parsing flow.
// Visitor can communicate with the Parser by returning these errors.
var (
	// ErrStop instructs the Parser to stop the entire parsing process.
	ErrStop = fmt.Errorf("stop parsing")
	// ErrSkipRecursive instructs the Parser to skip recursive parsing of the current field.
	// This is useful for avoiding infinite loops or handling custom serialization types.
	ErrSkipRecursive = fmt.Errorf("skip recursive parsing")
)

// FieldContext holds all context information for the field being visited.
// It provides a set of helper methods to abstract away reflection complexity, allowing users to manipulate fields more safely and easily.
type FieldContext struct {
	fieldValue  reflect.Value
	structField reflect.StructField
	parentValue reflect.Value
	currentPath string
}

// Path returns the dot-separated path of the current field, e.g. "Config.Server.Port".
func (c *FieldContext) Path() string {
	return c.currentPath
}

// FieldName returns the original name of the struct field.
func (c *FieldContext) FieldName() string {
	return c.structField.Name
}

// Tag gets the tag value for the specified key. If the tag does not exist, it returns an empty string.
func (c *FieldContext) Tag(key string) string {
	return c.structField.Tag.Get(key)
}

// Value returns the current value of the field as any.
func (c *FieldContext) Value() any {
	if c.fieldValue.IsValid() {
		return c.fieldValue.Interface()
	}

	return nil
}

// Kind returns the reflect.Kind type of the field.
func (c *FieldContext) Kind() reflect.Kind {
	return c.fieldValue.Kind()
}

// Type returns the reflect.Type type of the field.
func (c *FieldContext) Type() reflect.Type {
	return c.fieldValue.Type()
}

// Addr returns a pointer to the field value as any.
// This is very useful for scenarios where you need to pass the field address to other libraries (e.g., standard library flag).
func (c *FieldContext) Addr() any {
	if c.fieldValue.CanAddr() {
		return c.fieldValue.Addr().Interface()
	}

	return nil
}

// IsNil checks if the underlying value of the field is nil.
// Valid for pointer, map, slice, interface, channel, and function types.
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

// CanSet reports whether the field's value can be modified.
// Typically, only exported fields (with uppercase first letter) can be set.
func (c *FieldContext) CanSet() bool {
	return c.fieldValue.CanSet()
}

// Set attempts to set the field's value.
// It can handle automatic conversion from strings to various primitive types.
// This is the core method for Visitor implementations to initialize or modify field values.
func (c *FieldContext) Set(value any) error {
	if !c.CanSet() {
		return fmt.Errorf("field %s cannot be set (may not be exported)", c.FieldName())
	}

	valToSet := reflect.ValueOf(value)
	// Check if type can be directly assigned
	if valToSet.IsValid() && valToSet.Type().AssignableTo(c.fieldValue.Type()) {
		c.fieldValue.Set(valToSet)
		return nil
	}

	// If the input is a string, try type conversion
	if strVal, ok := value.(string); ok {
		return c.setString(strVal)
	}

	return fmt.Errorf("cannot assign value of type %T to field %s (expected type %s)", value, c.FieldName(), c.fieldValue.Type())
}

// setString is a helper function for the Set method, specifically handling conversion from strings to target types.
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
			return fmt.Errorf("value %s overflows for field %s", value, c.FieldName())
		}

		c.fieldValue.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		if c.fieldValue.OverflowUint(uintVal) {
			return fmt.Errorf("value %s overflows for field %s", value, c.FieldName())
		}

		c.fieldValue.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		if c.fieldValue.OverflowFloat(floatVal) {
			return fmt.Errorf("value %s overflows for field %s", value, c.FieldName())
		}

		c.fieldValue.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		c.fieldValue.SetBool(boolVal)
	default:
		return fmt.Errorf("conversion from string to field %s (type: %s) not supported", c.FieldName(), c.fieldValue.Kind())
	}

	return nil
}

// Visitor defines the interface for struct field processors.
// Users need to implement this interface to define their own business logic.
type Visitor interface {
	// VisitField is called when Parser traverses to each field.
	// The ctx parameter contains all relevant information and operations for that field.
	// Implementors can control the parsing flow by returning specific errors (e.g., ErrStop).
	VisitField(ctx *FieldContext) error
	// Recursion whether to recurse into this member
	Recursion(ctx *FieldContext) bool
}

// Parser is the core of struct parsing. It holds the visitor list and drives the entire traversal process.
type Parser struct {
	visitor Visitor
}

// NewParser creates a new Parser instance.
// It receives a Visitor and a series of ParserOptions for configuration.
func NewParser(visitor Visitor) *Parser {
	p := &Parser{
		visitor: visitor,
	}

	return p
}

// Parse starts the parsing process for the given object (must be a pointer to a struct).
func (p *Parser) Parse(obj any) error {
	if obj == nil {
		return fmt.Errorf("object to be parsed cannot be nil")
	}

	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("object must be a pointer to a struct")
	}

	// Start traversal from the root of the struct
	return p.walk(val.Elem(), val.Elem().Type().Name())
}

// walk is the internal recursive function, responsible for deep traversal of the value structure.
func (p *Parser) walk(val reflect.Value, path string) error {
	// If it's a pointer, first dereference it
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil // nil pointer cannot be dereferenced further, return directly
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

			// Call all visitors in sequence
			err := p.visitor.VisitField(ctx)
			if err != nil {
				if errors.Is(err, ErrStop) || errors.Is(err, ErrSkipRecursive) {
					return err // Pass control flow errors up
				}

				// For other business errors, return directly
				return fmt.Errorf("error on field %s: %w", ctx.Path(), err)
			}

			if p.visitor.Recursion(ctx) {
				// Key point: After Visitor executes, re-fetch the field value,
				// because Visitor may have initialized a nil pointer.
				updatedFieldVal := val.Field(i)
				if err := p.walk(updatedFieldVal, ctx.currentPath); err != nil {
					if errors.Is(err, ErrStop) {
						return err // If deep recursion requests stop, continue passing up
					}

					// Other recursion errors can be handled as needed, here we choose to return directly
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
		// Note: The order of map iteration is undefined
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

// CallMethod is a standalone utility function for invoking methods on objects via reflection.
// It does not belong to the Parser's core traversal logic, but is closely related to struct operations, so it is placed here as a helper.
func CallMethod(obj any, methodName string, args ...any) ([]reflect.Value, error) {
	val := reflect.ValueOf(obj)
	method := val.MethodByName(methodName)

	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found on type %T", methodName, obj)
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	if method.Type().NumIn() != len(in) {
		return nil, fmt.Errorf("method %s parameter count mismatch: expected %d, got %d", methodName, method.Type().NumIn(), len(in))
	}

	return method.Call(in), nil
}
