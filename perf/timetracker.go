package perf

import (
	"fmt"
	"time"
)

// TimeTracker is a performance testing tool for measuring code execution time.
// It is single-threaded and not concurrency-safe. It is recommended to create and use it within a single function body.
type TimeTracker struct {
	timePoint time.Time
}

// NewTimeTracker creates a new TimeTracker instance and immediately records the current time as the starting point.
func NewTimeTracker() *TimeTracker {
	return &TimeTracker{timePoint: time.Now()}
}

// RecordTime records the current time to reset the timer.
// If you need to measure different code segments multiple times, you can call this method before each measurement to reset the time point.
func (t *TimeTracker) RecordTime() {
	t.timePoint = time.Now()
}

// TrackTime calculates the elapsed time from the last recorded time to now, and calls the callback function to return the elapsed time and formatted string.
// After calculation, it will automatically reset the time point to the current time for the next measurement.
// Callback function parameters: elapsed is the raw time.Duration, formatElapsedTime is the formatted string.
func (t *TimeTracker) TrackTime(callback func(elapsed time.Duration, formattedElapsed string)) {
	elapsed := time.Since(t.timePoint)
	formattedTime := t.formatElapsedTime(elapsed)

	callback(elapsed, formattedTime)

	// Reset time point for next measurement
	t.timePoint = time.Now()
}

func (t *TimeTracker) formatElapsedTime(elapsed time.Duration) string {
	totalMs := elapsed.Milliseconds() // Total milliseconds

	// Calculate minutes, seconds and milliseconds parts
	minutes := totalMs / 60000       // 60000 ms = 1 minute
	seconds := (totalMs / 1000) % 60 // Extract seconds
	milliseconds := totalMs % 1000   // Extract milliseconds

	// Format output based on time length
	var formattedTime string
	switch {
	case minutes > 0:
		formattedTime = fmt.Sprintf("%d minutes %d seconds %d ms (raw: %d ms)", minutes, seconds, milliseconds, totalMs)
	case seconds > 0:
		formattedTime = fmt.Sprintf("%d seconds %d ms (raw: %d ms)", seconds, milliseconds, totalMs)
	default:
		formattedTime = fmt.Sprintf("%d ms", milliseconds)
	}

	return formattedTime
}
