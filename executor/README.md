# executor 包

`executor` 包提供了用于执行后台任务和并发任务的工具。

## LoopTask

`LoopTask` 可以方便地创建一个周期性执行的任务。

### 用法

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
    // Start() 为阻塞调用，因此我们开启 Goroutine 以将其放入后台
    // 也可以直接用内置简化的：executor.GoLoopTask(task, 2 * time.Second) 代替
    go looper.Start(task)

    // 让主程序运行一段时间
    time.Sleep(10 * time.Second)

    looper.Stop()
	log.Println("Looper stopped.")
}
```
