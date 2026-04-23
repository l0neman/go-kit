package perf

import (
	"strings"
	"testing"
	"time"
)

func TestTimeTracker_NewTimeTracker(t *testing.T) {
	tracker := NewTimeTracker()
	if tracker == nil {
		t.Error("NewTimeTracker returned nil")
		return
	}

	// Check if timePoint was successfully initialized
	if tracker.timePoint.IsZero() {
		t.Error("TimeTracker's timePoint was not correctly initialized")
	}
}

func TestTimeTracker_RecordTime(t *testing.T) {
	tracker := NewTimeTracker()

	// Record initial time
	tracker.RecordTime()

	// Wait a short period of time
	time.Sleep(10 * time.Millisecond)

	// Record time again
	tracker.RecordTime()

	// Verify that the time point has been updated
	if tracker.timePoint.IsZero() {
		t.Error("RecordTime did not correctly update the time point")
	}
}

func TestTimeTracker_TrackTime(t *testing.T) {
	tracker := NewTimeTracker()

	// Wait a short period of time to ensure measurable time difference
	time.Sleep(50 * time.Millisecond)

	// Use callback function to capture result
	var capturedElapsed time.Duration
	var capturedFormatted string

	tracker.TrackTime(func(elapsed time.Duration, formattedElapsed string) {
		capturedElapsed = elapsed
		capturedFormatted = formattedElapsed
	})

	// Verify that a non-zero elapsed time was captured
	if capturedElapsed <= 0 {
		t.Errorf("Expected positive elapsed time, got: %v", capturedElapsed)
	}

	// Verify the formatted string is not empty
	if capturedFormatted == "" {
		t.Error("Formatted time is empty")
	}

	// Verify the formatted string contains original milliseconds
	if !strings.HasSuffix(capturedFormatted, "ms") {
		t.Error("Formatted string does not contain original milliseconds information")
	}
}

func TestTimeTracker_formatElapsedTime(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "millisecond level",
			duration: 150 * time.Millisecond,
			expected: "150 ms",
		},
		{
			name:     "second level",
			duration: 1250 * time.Millisecond, // 1 second 250 milliseconds
			expected: "1 second 250 ms (raw: 1250 ms)",
		},
		{
			name:     "minute level",
			duration: 61500 * time.Millisecond, // 1 minute 1 second 500 milliseconds
			expected: "1 minute 1 second 500 ms (raw: 61500 ms)",
		},
		{
			name:     "over 1 minute",
			duration: 125 * time.Second, // 2 minutes 5 seconds
			expected: "2 minutes 5 seconds 0 ms (raw: 125000 ms)",
		},
	}

	tracker := NewTimeTracker()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tracker.formatElapsedTime(tt.duration)
			if result != tt.expected {
				t.Errorf("formatElapsedTime() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
