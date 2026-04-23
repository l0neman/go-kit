# executor package

The `executor` package provides utilities for executing background tasks and concurrent tasks.

## LoopTask

`LoopTask` conveniently creates a task that executes periodically.

### Usage

```go
package main

import (
    "log"
    "time"
    "github.com/l0neman/go-kit/executor"
)

func main() {
    task := func() {
        log.Println("Task is running...")
    }

    looper := executor.NewLoopTask(2 * time.Second)
    // Start() is a blocking call, so we put it in a Goroutine to run in background
    // You can also use the built-in simplified version: executor.GoLoopTask(task, 2 * time.Second)
    go looper.Start(task)

    // Let the main program run for a while
    time.Sleep(10 * time.Second)

    looper.Stop()
    log.Println("Looper stopped.")
}
```