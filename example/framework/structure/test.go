package structure

import (
	"fmt"

	"github.com/l0neman/go-kit/framework/structure"
)

type SubConfigA struct {
	Name string `yaml:"name" default:"SubA-Default-Name"`
	TTL  int    `yaml:"ttl" default:"60"`
}

// SubConfigB 是另一个子结构体，用于指针和切片
type SubConfigB struct {
	Endpoint string   `yaml:"endpoint" default:"https://api.example.com/v2"`
	Enabled  bool     `yaml:"enabled" default:"true"`
	Params   []string `yaml:"params" default:"[\"a\", \"b\"]"` // default 标签对切片无效，setField 未实现
}

// ComplexConfig 是我们的主要测试目标
type ComplexConfig struct {
	// 简单字段
	AppName string `yaml:"app_name" default:"Super-App"`
	Version string `yaml:"version" default:"1.0.0"`

	// 嵌套结构体（值类型）
	PrimaryDB SubConfigA `yaml:"primary_db"`

	// 指向结构体的指针（可能为 nil）
	Cache *SubConfigA `yaml:"cache"`

	// 包含结构体值的切片
	// 注意：为了测试，我们需要在实例中预先分配一些空的元素
	ValueSlice []SubConfigA `yaml:"value_slice"`

	// 包含结构体指针的切片
	// 注意：同样需要预分配
	PointerSlice []*SubConfigB `yaml:"pointer_slice"`

	// 值为结构体的 Map
	ValueMap map[string]SubConfigA `yaml:"value_map"`

	// 值为结构体指针的 Map
	PointerMap map[string]*SubConfigB `yaml:"pointer_map"`
}

// Key 为 ComplexConfig 提供一个自定义的键名
func (c *ComplexConfig) Key() string {
	return "complex_application_config"
}

type TestVisitor struct {
}

func (t *TestVisitor) VisitField(ctx *structure.FieldContext) error {
	if ctx.IsNil() {
		fmt.Printf("VisitField %s, -> nil\n", ctx.FieldName())
	} else {
		fmt.Printf("VisitField %s, -> %v\n", ctx.FieldName(), ctx.Value())
	}

	return nil
}

func (t *TestVisitor) Recursion(_ *structure.FieldContext) bool {
	return true
}

func Test() {
	conf := &ComplexConfig{}
	parser := structure.NewParser(&TestVisitor{})
	err := parser.Parse(conf)
	if err != nil {
		panic(err)
	}
}
