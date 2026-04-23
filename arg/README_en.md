# arg package

The `arg` package provides struct-based command-line argument parsing functionality for common single-layer command scenarios.

It is essentially a wrapper around the `flag` package.

It makes it easy to parse command-line argument values, as well as set default values and help information.

## Usage

Here is a simple example of parsing command-line arguments based on a struct:

```go
package main

import (
    "fmt"
    "github.com/l0neman/go-kit/arg"
)

func main() {
    // Corresponding command input: ./main -host 0.0.0.0 -port 8888 -enable_tls
    // Data types are the same as the flag package: string, bool, int, int64, float64
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
```