# perf package

The `perf` package provides utilities for performance monitoring and analysis.

## TimeTracker

`TimeTracker` can conveniently measure and record the execution time of code blocks, and output in human-readable format.

### Basic Usage

The workflow of `TimeTracker` is as follows:

1.  Create an instance via `perf.NewTimeTracker()`, which automatically records a starting time point.
2.  When the code execution reaches the measurement endpoint, call the `tracker.TrackTime()` method and pass in a callback function.
3.  `TrackTime` calculates the elapsed time from the initial time point to the current moment, and passes the raw `time.Duration` value and the formatted string to your callback function.

You can also call `tracker.RecordTime()` at any time to reset the initial time point.

### Example

```go
package main

import (
    "log"
    "time"
    "github.com/l0neman/go-kit/perf"
)

func main() {
    // 1. Create a new TimeTracker and start timing
    tracker := perf.NewTimeTracker()
    log.Println("Starting some time-consuming operations...")

    // Simulate a time-consuming operation
    time.Sleep(150 * time.Millisecond)

    // 2. End timing and get results via callback function
    tracker.TrackTime(func(elapsed time.Duration, formattedElapsed string) {
        log.Printf("Operation completed!\n")
        log.Printf("Formatted elapsed time: %s\n", formattedElapsed)
        log.Printf("Raw elapsed time (nanoseconds): %d\n", elapsed.Nanoseconds())
    })

    log.Println("\n--- Measuring segmented execution time ---")
    // 3. Reset the timer to measure a new code block
    tracker.RecordTime()
    time.Sleep(65 * time.Millisecond)
    tracker.TrackTime(func(_ time.Duration, formatted string) {
        log.Printf("First part duration: %s\n", formatted)
    })

    tracker.RecordTime()
    time.Sleep(85 * time.Millisecond)
    tracker.TrackTime(func(_ time.Duration, formatted string) {
        log.Printf("Second part duration: %s\n", formatted)
    })
}
```

### Output Format

`TimeTracker` intelligently selects the output format based on the elapsed duration:

- **Milliseconds level**: `150 ms`
- **Seconds level**: `1 second 250 ms (original: 1250 ms)`
- **Minutes level**: `1 minute 1 second 500 ms (original: 61500 ms)`