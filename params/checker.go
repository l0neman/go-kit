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

// VisitField 访问并检查每一个字段
func (v *Visitor) VisitField(ctx *structure.FieldContext) error {
	// 从 tag 中获取检查表达式
	exp := ctx.Tag("check")
	if exp == "" {
		// 如果没有 check 标签，则跳过检查
		return nil
	}

	// 获取字段名，优先使用 json 标签
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
		return fmt.Errorf("字段 '%s' %v，路径 > %s", fieldName, err, ctx.Path())
	}

	for _, r := range rules {
		op, valStr := r.op, r.valStr

		// 根据字段类型执行校验
		switch ctx.Kind() {
		case reflect.String:
			val := ctx.Value().(string)
			switch op {
			case expressNotEmpty:
				if val == "" {
					return fmt.Errorf("字段 '%s' 不能为空，路径 > %s", fieldName, ctx.Path())
				}
			case expressEmpty:
				if val != "" {
					return fmt.Errorf("字段 '%s' 必须为空，路径 > %s", fieldName, ctx.Path())
				}
			default:
				return fmt.Errorf("字段 '%s' 是字符串类型，不支持 '%s' 表达式，路径 > %s", fieldName, op, ctx.Path())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldVal := reflect.ValueOf(ctx.Value()).Int()
			checkVal, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				return fmt.Errorf("字段 '%s' 的检查值 '%s' 类型错误，期望为整数，路径 > %s", fieldName, valStr, ctx.Path())
			}
			if !compareInt(fieldVal, checkVal, op) {
				return fmt.Errorf("字段 '%s' 的值 '%d' 不满足条件 '%s%s'，路径 > %s", fieldName, fieldVal, op, valStr, ctx.Path())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldVal := reflect.ValueOf(ctx.Value()).Uint()
			checkVal, err := strconv.ParseUint(valStr, 10, 64)
			if err != nil {
				return fmt.Errorf("字段 '%s' 的检查值 '%s' 类型错误，期望为无符号整数，路径 > %s", fieldName, valStr, ctx.Path())
			}
			if !compareUint(fieldVal, checkVal, op) {
				return fmt.Errorf("字段 '%s' 的值 '%d' 不满足条件 '%s%s'，路径 > %s", fieldName, fieldVal, op, valStr, ctx.Path())
			}
		case reflect.Float32, reflect.Float64:
			fieldVal := reflect.ValueOf(ctx.Value()).Float()
			checkVal, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return fmt.Errorf("字段 '%s' 的检查值 '%s' 类型错误，期望为浮点数，路径 > %s", fieldName, valStr, ctx.Path())
			}

			if !compareFloat(fieldVal, checkVal, op) {
				return fmt.Errorf("字段 '%s' 的值 '%f' 不满足条件 '%s%s'，路径 > %s", fieldName, fieldVal, op, valStr, ctx.Path())
			}
		default:
			// 对于不支持的类型，可以选择忽略或报错
			// return fmt.Errorf("字段 '%s' 是不支持的类型 '%s'，无法进行检查，路径 > %s", fieldName, ctx.Kind(), ctx.Path())
		}
	}

	return nil
}

// compareInt 比较整数
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

// compareUint 比较无符号整数
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

// compareFloat 比较浮点数
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
			return nil, fmt.Errorf("的检查表达式 '%s' 无效", singleExp)
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

// CheckFields 检查结构体字段
// obj: 要检查的结构体对象
// 返回值: 如果有错误，返回错误信息；如果没有错误，返回 nil
// 结构体字段增加 check 标签，值为表达式，例如:
// `check:"not_empty"`
// `check:">0,<=100"`
// `check:"==10"`
func CheckFields(obj any) error {
	parser := structure.NewParser(&Visitor{})
	err := parser.Parse(obj)
	if err != nil {
		return errorx.Wrap(err, "检查出错")
	}

	return nil
}
