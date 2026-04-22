package perf

import (
	"strings"
	"testing"
	"time"
)

func TestTimeTracker_NewTimeTracker(t *testing.T) {
	tracker := NewTimeTracker()
	if tracker == nil {
		t.Error("NewTimeTracker 返回了 nil")
		return
	}

	// 检查是否成功初始化了 timePoint
	if tracker.timePoint.IsZero() {
		t.Error("TimeTracker 的 timePoint 未正确初始化")
	}
}

func TestTimeTracker_RecordTime(t *testing.T) {
	tracker := NewTimeTracker()

	// 记录初始时间
	tracker.RecordTime()

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	// 再次记录时间
	tracker.RecordTime()

	// 验证时间点是否已更新
	if tracker.timePoint.IsZero() {
		t.Error("RecordTime 未正确更新时间点")
	}
}

func TestTimeTracker_TrackTime(t *testing.T) {
	tracker := NewTimeTracker()

	// 等待一小段时间以确保有可测量的时间差
	time.Sleep(50 * time.Millisecond)

	// 使用回调函数捕获结果
	var capturedElapsed time.Duration
	var capturedFormatted string

	tracker.TrackTime(func(elapsed time.Duration, formattedElapsed string) {
		capturedElapsed = elapsed
		capturedFormatted = formattedElapsed
	})

	// 验证是否捕获到了非零的耗时
	if capturedElapsed <= 0 {
		t.Errorf("期望捕获到正的耗时，但得到: %v", capturedElapsed)
	}

	// 验证格式化后的字符串不为空
	if capturedFormatted == "" {
		t.Error("格式化后的时间为空")
	}

	// 验证格式化后的字符串包含原始毫秒数
	if !strings.HasSuffix(capturedFormatted, "毫秒") {
		t.Error("格式化后的字符串不包含原始毫秒数信息")
	}
}

func TestTimeTracker_formatElapsedTime(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "毫秒级别",
			duration: 150 * time.Millisecond,
			expected: "150 毫秒",
		},
		{
			name:     "秒级别",
			duration: 1250 * time.Millisecond, // 1秒250毫秒
			expected: "1 秒 250 毫秒（原始：1250 毫秒）",
		},
		{
			name:     "分钟级别",
			duration: 61500 * time.Millisecond, // 1分1秒500毫秒
			expected: "1 分钟 1 秒 500 毫秒（原始：61500 毫秒）",
		},
		{
			name:     "超过1分钟",
			duration: 125 * time.Second, // 2分5秒
			expected: "2 分钟 5 秒 0 毫秒（原始：125000 毫秒）",
		},
	}

	tracker := NewTimeTracker()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tracker.formatElapsedTime(tt.duration)
			if result != tt.expected {
				t.Errorf("formatElapsedTime() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}
