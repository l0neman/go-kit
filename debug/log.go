package debug

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
)

var throttledLog *ThrottledLog

const stackOffset = 4 // stack offset

// for ThrottledLog
func init() {
	throttledLog = &ThrottledLog{callCount: map[string]int{}}
}

// ThrottledLog intelligently handles frequently printed logs
type ThrottledLog struct {
	mu        sync.Mutex
	callCount map[string]int
}

// Print prints the log with throttling optimization
func (sl *ThrottledLog) Print(logFun func(count int)) {
	// Build unique key for caller information
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
	// Get call count
	count := sl.callCount[key]
	count++
	sl.callCount[key] = count
	// 2 4 8 16 ... Use powers of 2 for printing
	return count, count&(count-1) == 0
}

// PrintThrottled prints logs intelligently
func PrintThrottled(logFun func(count int)) {
	throttledLog.Print(logFun)
}
