package perf

import (
	"fmt"
	"time"
)

// TimeTracker 耗时测试工具，用于测量代码性能。它是单线程的，非并发安全，建议在单个函数体中创建使用
type TimeTracker struct {
	timePoint time.Time
}

// NewTimeTracker 创建一个新的 TimeTracker 实例，并立即记录当前时间作为起始点。
func NewTimeTracker() *TimeTracker {
	return &TimeTracker{timePoint: time.Now()}
}

// RecordTime 记录当前时间，用于重置计时起点。
// 如果需要多次测量不同代码段，可以在每次测量前调用此方法重置时间点。
func (t *TimeTracker) RecordTime() {
	t.timePoint = time.Now()
}

// TrackTime 计算从上次记录时间到当前的耗时，并调用回调函数返回耗时和格式化字符串。
// 计算完成后会自动重置时间点为当前时间，以便进行下一次测量。
// 回调函数参数：elapsed 为原始的 time.Duration，formatElapsedTime 为格式化后的字符串。
func (t *TimeTracker) TrackTime(callback func(elapsed time.Duration, formattedElapsed string)) {
	elapsed := time.Since(t.timePoint)
	formattedTime := t.formatElapsedTime(elapsed)

	callback(elapsed, formattedTime)

	// 重置时间点，为下一次测量做准备
	t.timePoint = time.Now()
}

func (t *TimeTracker) formatElapsedTime(elapsed time.Duration) string {
	totalMs := elapsed.Milliseconds() // 总毫秒数

	// 计算分钟、秒和毫秒部分
	minutes := totalMs / 60000       // 60000 ms = 1 minute
	seconds := (totalMs / 1000) % 60 // 取秒部分
	milliseconds := totalMs % 1000   // 取毫秒部分

	// 根据时间长度格式化输出
	var formattedTime string
	switch {
	case minutes > 0:
		formattedTime = fmt.Sprintf("%d 分钟 %d 秒 %d 毫秒（原始：%d 毫秒）", minutes, seconds, milliseconds, totalMs)
	case seconds > 0:
		formattedTime = fmt.Sprintf("%d 秒 %d 毫秒（原始：%d 毫秒）", seconds, milliseconds, totalMs)
	default:
		formattedTime = fmt.Sprintf("%d 毫秒", milliseconds)
	}

	return formattedTime
}
