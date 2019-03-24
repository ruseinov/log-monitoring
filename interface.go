package log_monitoring

import (
	"context"
	"time"
)

// LogEntry represents data that is being parsed from each log entry
type LogEntry struct {
	Method  string
	Path    string
	Status  string
	Section string
	Proto   string
	Bytes   int
}

// Parser represents an interface for parsing log entries line by line
type Parser interface {
	Parse(line string) (*LogEntry, error)
}

// Counter represents a key-value pair of string -> integer
type Counter struct {
	Key   string
	Value int64
}

// Stats represents interesting traffic statistics
type Stats struct {
	TopSections []Counter
	TopStatuses []Counter
	TopMethods  []Counter
	TotalBytes  int64
}

// Monitor processes LogEntries and sends out stats for a given interval
type Monitor interface {
	// Run processes long entries and sends out stats for a given rollup interval
	Run(ctx context.Context, rollupInterval time.Duration) (chan<- *LogEntry, <-chan *Stats)
}

// Alerter monitors RPS and alerts whenever the threshold is crossed or back to normal
type Alerter interface {
	// Incr increments the request count by 1
	Incr()
	// Run watches the value that's being incremented and sends out alerts every 2 minutes as specified by the threshold
	Run(ctx context.Context, rpsThreshold int64) <-chan *Alert
}

// Printer simply gets alerts and stats and prints them in a readable manner
type Printer interface {
	// Run watches alert and stat channels for changes to print
	Run(ctx context.Context, alertChan <-chan *Alert, statsChan <-chan *Stats)
}
