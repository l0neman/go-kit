package arg

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/l0neman/go-kit/errorx"
	"github.com/l0neman/go-kit/framework/structure"
)

type StructVisitor struct {
	FlagSet *flag.FlagSet
}

func (v *StructVisitor) VisitField(ctx *structure.FieldContext) error {
	var (
		nameTag    = ctx.Tag("name")
		defaultTag = ctx.Tag("default")
		helpTag    = ctx.Tag("help")
	)

	flagName := nameTag
	if flagName == "" {
		flagName = toSnakeCase(ctx.FieldName())
	}

	addr := ctx.Addr()
	if addr == nil {
		return fmt.Errorf("不支持的成员 > %s", ctx.FieldName())
	}

	switch ctx.Kind() {
	case reflect.String:
		v.FlagSet.StringVar(addr.(*string), flagName, defaultTag, helpTag)
	case reflect.Bool:
		defaultVal, _ := strconv.ParseBool(defaultTag)
		v.FlagSet.BoolVar(addr.(*bool), flagName, defaultVal, helpTag)
	case reflect.Int:
		defaultVal, _ := strconv.Atoi(defaultTag)
		v.FlagSet.IntVar(addr.(*int), flagName, defaultVal, helpTag)
	case reflect.Int64:
		defaultVal, _ := strconv.ParseInt(defaultTag, 10, 64)
		v.FlagSet.Int64Var(addr.(*int64), flagName, defaultVal, helpTag)
	case reflect.Float64:
		defaultVal, _ := strconv.ParseFloat(defaultTag, 64)
		v.FlagSet.Float64Var(addr.(*float64), flagName, defaultVal, helpTag)
	default:
		// ignore
	}

	return nil
}

func (v *StructVisitor) Recursion(*structure.FieldContext) bool {
	return false
}

func Parse(argsStructPtr any) error {
	v := reflect.ValueOf(argsStructPtr)
	if v.Kind() != reflect.Ptr {
		return errors.New("必须传入结构体指针")
	}

	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("必须传入结构体指针")
	}

	structType := elem.Type()
	// 使用结构体名称的小写形式作为 FlagSet 的名称
	flagSetName := strings.ToLower(structType.Name())

	flagSet := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
	visitor := &StructVisitor{
		FlagSet: flagSet,
	}

	structParser := structure.NewParser(visitor)
	err := structParser.Parse(argsStructPtr)
	if err != nil {
		return errorx.Wrap(err, "解析结构体出错")
	}

	return flagSet.Parse(os.Args[1:])
}

// toSnakeCase 将驼峰式命名（CamelCase）字符串转换为蛇形命名（snake_case）。
// 例如: "LocalConf" -> "local_conf", "Master" -> "master"
func toSnakeCase(str string) string {
	var builder strings.Builder
	for i, r := range str {
		// 如果是大写字母且不是第一个字符，则在其前面加下划线
		if i > 0 && unicode.IsUpper(r) {
			builder.WriteRune('_')
		}

		builder.WriteRune(unicode.ToLower(r))
	}

	return builder.String()
}
