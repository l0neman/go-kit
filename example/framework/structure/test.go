package structure

import (
	"fmt"

	"github.com/l0neman/go-kit/framework/structure"
)

type SubConfigA struct {
	Name string `yaml:"name" default:"SubA-Default-Name"`
	TTL  int    `yaml:"ttl" default:"60"`
}

// SubConfigB is another sub struct, used for pointers and slices
type SubConfigB struct {
	Endpoint string   `yaml:"endpoint" default:"https://api.example.com/v2"`
	Enabled  bool     `yaml:"enabled" default:"true"`
	Params   []string `yaml:"params" default:"[\"a\", \"b\"]"` // default tag is invalid for slices, setField not implemented
}

// ComplexConfig is our main test target
type ComplexConfig struct {
	// Simple fields
	AppName string `yaml:"app_name" default:"Super-App"`
	Version string `yaml:"version" default:"1.0.0"`

	// Nested struct (value type)
	PrimaryDB SubConfigA `yaml:"primary_db"`

	// Pointer to struct (may be nil)
	Cache *SubConfigA `yaml:"cache"`

	// Slice containing struct values
	// Note: For testing, we need to pre-allocate some empty elements in the instance
	ValueSlice []SubConfigA `yaml:"value_slice"`

	// Slice containing struct pointers
	// Note: Also needs pre-allocation
	PointerSlice []*SubConfigB `yaml:"pointer_slice"`

	// Map with struct values
	ValueMap map[string]SubConfigA `yaml:"value_map"`

	// Map with struct pointer values
	PointerMap map[string]*SubConfigB `yaml:"pointer_map"`
}

// Key provides a custom key name for ComplexConfig
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
