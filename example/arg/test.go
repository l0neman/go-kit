package arg

import (
	"fmt"

	"github.com/l0neman/go-kit/arg"
)

func Test() {
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
