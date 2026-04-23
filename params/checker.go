package params

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/l0neman/go-kit/errorx"
	"github.com/l0neman/go-kit/framework/structure"
)

type Visitor struct {
}

// VisitField visits and checks each field
func (v *Visitor) VisitField(ctx *structure.FieldContext) error {
	// Get check expression from tag
	exp := ctx.Tag("check")
	if exp == "" {
		// Skip check if no check tag is present
		return nil
	}

	// Get field name, prefer json tag
	fieldName := ctx.Tag("json")
	if fieldName == "" {
		fieldName = ctx.FieldName()
	} else {
		if idx := strings.Index(fieldName, ","); idx != -1 {
			fieldName = fieldName[:idx]
		}
	}

	rules, err := parseExp(exp)
	if err != nil {
		return fmt.Errorf("field '%s' %v, path > %s", fieldName, err, ctx.Path())
	}

	for _, r := range rules {
		op, valStr := r.op, r.valStr

		// Perform validation based on field type
		switch ctx.Kind() {
		case reflect.String:
			val := ctx.Value().(string)
			switch op {
			case expressNotEmpty:
				if val == "" {
					return fmt.Errorf("field '%s' cannot be empty, path > %s", fieldName, ctx.Path())
				}
			case expressEmpty:
				if val != "" {
					return fmt.Errorf("field '%s' must be empty, path > %s", fieldName, ctx.Path())
				}
			default:
				return fmt.Errorf("field '%s' is of string type, does not support '%s' expression, path > %s", fieldName, op, ctx.Path())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldVal := reflect.ValueOf(ctx.Value()).Int()
			checkVal, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				return fmt.Errorf("field '%s' check value '%s' has incorrect type, expected integer, path > %s", fieldName, valStr, ctx.Path())
			}
			if !compareInt(fieldVal, checkVal, op) {
				return fmt.Errorf("field '%s' value '%d' does not satisfy condition '%s%s', path > %s", fieldName, fieldVal, op, valStr, ctx.Path())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldVal := reflect.ValueOf(ctx.Value()).Uint()
			checkVal, err := strconv.ParseUint(valStr, 10, 64)
			if err != nil {
				return fmt.Errorf("field '%s' check value '%s' has incorrect type, expected unsigned integer, path > %s", fieldName, valStr, ctx.Path())
			}
			if !compareUint(fieldVal, checkVal, op) {
				return fmt.Errorf("field '%s' value '%d' does not satisfy condition '%s%s', path > %s", fieldName, fieldVal, op, valStr, ctx.Path())
			}
		case reflect.Float32, reflect.Float64:
			fieldVal := reflect.ValueOf(ctx.Value()).Float()
			checkVal, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return fmt.Errorf("field '%s' check value '%s' has incorrect type, expected float, path > %s", fieldName, valStr, ctx.Path())
			}

			if !compareFloat(fieldVal, checkVal, op) {
				return fmt.Errorf("field '%s' value '%f' does not satisfy condition '%s%s', path > %s", fieldName, fieldVal, op, valStr, ctx.Path())
			}
		default:
			// For unsupported types, you can choose to ignore or report an error
			// return fmt.Errorf("field '%s' is of unsupported type '%s', cannot perform check, path > %s", fieldName, ctx.Kind(), ctx.Path())
		}
	}

	return nil
}

// compareInt compares integers
func compareInt(a, b int64, op string) bool {
	switch op {
	case expressionGreaterThan:
		return a > b
	case expressLessThan:
		return a < b
	case expressGreaterThanOrEqualTo:
		return a >= b
	case expressLessThanOrEqualTo:
		return a <= b
	case expressEqualTo:
		return a == b
	case expressNotEqualTo:
		return a != b
	default:
		return false
	}
}

// compareUint compares unsigned integers
func compareUint(a, b uint64, op string) bool {
	switch op {
	case expressionGreaterThan:
		return a > b
	case expressLessThan:
		return a < b
	case expressGreaterThanOrEqualTo:
		return a >= b
	case expressLessThanOrEqualTo:
		return a <= b
	case expressEqualTo:
		return a == b
	case expressNotEqualTo:
		return a != b
	default:
		return false
	}
}

// compareFloat compares floating point numbers
func compareFloat(a, b float64, op string) bool {
	switch op {
	case expressionGreaterThan:
		return a > b
	case expressLessThan:
		return a < b
	case expressGreaterThanOrEqualTo:
		return a >= b
	case expressLessThanOrEqualTo:
		return a <= b
	case expressEqualTo:
		return a == b
	case expressNotEqualTo:
		return a != b
	default:
		return false
	}
}

func (v *Visitor) Recursion(ctx *structure.FieldContext) bool {
	return true
}

var expCache sync.Map

type rule struct {
	op     string
	valStr string
}

func parseExp(exp string) ([]rule, error) {
	if cached, ok := expCache.Load(exp); ok {
		return cached.([]rule), nil
	}

	var rules []rule
	expressions := strings.Split(exp, ",")
	for _, singleExp := range expressions {
		singleExp = strings.TrimSpace(singleExp)
		if singleExp == "" {
			continue
		}

		var op, valStr string
		if strings.HasPrefix(singleExp, expressGreaterThanOrEqualTo) {
			op = expressGreaterThanOrEqualTo
			valStr = strings.TrimSpace(singleExp[len(expressGreaterThanOrEqualTo):])
		} else if strings.HasPrefix(singleExp, expressLessThanOrEqualTo) {
			op = expressLessThanOrEqualTo
			valStr = strings.TrimSpace(singleExp[len(expressLessThanOrEqualTo):])
		} else if strings.HasPrefix(singleExp, expressEqualTo) {
			op = expressEqualTo
			valStr = strings.TrimSpace(singleExp[len(expressEqualTo):])
		} else if strings.HasPrefix(singleExp, expressNotEqualTo) {
			op = expressNotEqualTo
			valStr = strings.TrimSpace(singleExp[len(expressNotEqualTo):])
		} else if strings.HasPrefix(singleExp, expressionGreaterThan) {
			op = expressionGreaterThan
			valStr = strings.TrimSpace(singleExp[len(expressionGreaterThan):])
		} else if strings.HasPrefix(singleExp, expressLessThan) {
			op = expressLessThan
			valStr = strings.TrimSpace(singleExp[len(expressLessThan):])
		} else if singleExp == expressNotEmpty {
			op = expressNotEmpty
		} else if singleExp == expressEmpty {
			op = expressEmpty
		} else {
			return nil, fmt.Errorf("invalid check expression '%s'", singleExp)
		}

		rules = append(rules, rule{op: op, valStr: valStr})
	}

	expCache.Store(exp, rules)
	return rules, nil
}

const (
	expressionGreaterThan       = ">"
	expressLessThan             = "<"
	expressGreaterThanOrEqualTo = ">="
	expressLessThanOrEqualTo    = "<="
	expressEqualTo              = "=="
	expressNotEqualTo           = "!="
	expressNotEmpty             = "not_empty"
	expressEmpty                = "is_empty"
)

// CheckFields checks struct fields
// obj: the struct object to check
// Return value: returns error message if there is an error; returns nil if there is no error
// Add check tag to struct fields with expression as value, for example:
// `check:"not_empty"`
// `check:">0,<=100"`
// `check:"==10"`
func CheckFields(obj any) error {
	parser := structure.NewParser(&Visitor{})
	err := parser.Parse(obj)
	if err != nil {
		return errorx.Wrap(err, "check error")
	}

	return nil
}
