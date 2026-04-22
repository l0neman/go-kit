# perf 包

`perf` 包提供了用于性能监控和分析的工具。

## TimeTracker

`TimeTracker` 可以方便地测量和记录代码块的执行时间，并输出人类可读的格式。

### 核心用法

`TimeTracker` 的工作流程如下：

1.  通过 `perf.NewTimeTracker()` 创建一个实例，此时它会自动记录一个初始时间点。
2.  在代码执行到需要测量终点时，调用 `tracker.TrackTime()` 方法，并传入一个回调函数。
3.  `TrackTime` 会计算从初始时间点到当前时刻的耗时，并将原始的 `time.Duration` 值和格式化后的字符串传给你的回调函数。

你也可以在任何时候调用 `tracker.RecordTime()` 来重置初始时间点。

### 示例

```go
package main

import (
    "log"
    "time"
    "github.com/l0neman/go-kit/perf"
)

func main() {
    // 1. 创建一个新的 TimeTracker，开始计时
    tracker := perf.NewTimeTracker()
    log.Println("开始执行一些耗时操作...")

    // 模拟一个耗时操作
    time.Sleep(150 * time.Millisecond)

    // 2. 结束计时，并通过回调函数获取结果
    tracker.TrackTime(func(elapsed time.Duration, formattedElapsed string) {
        log.Printf("操作完成!\n")
        log.Printf("格式化后的耗时: %s\n", formattedElapsed)
        log.Printf("原始耗时 (纳秒): %d\n", elapsed.Nanoseconds())
    })

    log.Println("\n--- 测量分段耗时 ---")
    // 3. 重置计时器来测量一个新的代码块
    tracker.RecordTime()
    time.Sleep(65 * time.Millisecond)
    tracker.TrackTime(func(_ time.Duration, formatted string) {
        log.Printf("第一部分耗时: %s\n", formatted)
    })

    tracker.RecordTime()
    time.Sleep(85 * time.Millisecond)
    tracker.TrackTime(func(_ time.Duration, formatted string) {
        log.Printf("第二部分耗时: %s\n", formatted)
    })
}
```

### 输出格式

`TimeTracker` 会根据耗时长度智能地选择输出格式：

- **毫秒级别**: `150 毫秒`
- **秒级别**: `1 秒 250 毫秒（原始：1250 毫秒）`
- **分钟级别**: `1 分钟 1 秒 500 毫秒（原始：61500 毫秒）`
