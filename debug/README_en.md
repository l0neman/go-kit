# debug package

The `debug` package provides utilities for debugging, such as throttled log printing.

## Usage

```go
package main

import (
    "log"
    "github.com/l0neman/go-kit/debug"
)

func main() {
    // Intelligent throttled log printing: for a large number of repeated logs, only print 2^n times
    for i := 0; i < 1024; i++ {
        debug.PrintThrottled(func(count int) {
            log.Println("log 1 print count", count)
        })
    }
}
```