package debug

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
)

var throttledLog *ThrottledLog

const stackOffset = 4 // 栈偏移

// for ThrottledLog
func init() {
	throttledLog = &ThrottledLog{callCount: map[string]int{}}
}

// ThrottledLog 智能处理频繁打印的日志
type ThrottledLog struct {
	mu        sync.Mutex
	callCount map[string]int
}

// Print 打印日志，根据降频优化
func (sl *ThrottledLog) Print(logFun func(count int)) {
	// 构建调用者信息的唯一 key
	callerKey := makeKeyWithStackInfo(stackOffset)
	sl.mu.Lock()
	defer sl.mu.Unlock()

	count, canLog := sl.countLogAndCanLog(callerKey)
	if canLog {
		logFun(count)
	}
}

func makeKeyWithStackInfo(stackOffset int) string {
	var pcs [1]uintptr
	runtime.Callers(stackOffset, pcs[:])
	frames := runtime.CallersFrames([]uintptr{pcs[0]})
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d-%s", filepath.Base(frame.File), frame.Line, frame.Function)
}

func (sl *ThrottledLog) countLogAndCanLog(key string) (int, bool) {
	// 获取调用次数
	count := sl.callCount[key]
	count++
	sl.callCount[key] = count
	// 2 4 8 16 ... 采用 2 的 n 次方进行打印
	return count, count&(count-1) == 0
}

// PrintThrottled 智能打印日志
func PrintThrottled(logFun func(count int)) {
	throttledLog.Print(logFun)
}
