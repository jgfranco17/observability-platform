package observability

import (
	"context"
	"time"
)

// EntryLogLevel represents the severity level of a log entry.
type EntryLogLevel string

const (
	LogLevelDebug EntryLogLevel = "debug"
	LogLevelInfo  EntryLogLevel = "info"
	LogLevelWarn  EntryLogLevel = "warn"
	LogLevelError EntryLogLevel = "error"
)

// Entry represents a structured log entry with optional trace correlation.
type Entry struct {
	Level      EntryLogLevel     `json:"level"`
	Message    string            `json:"message"`
	TraceID    string            `json:"trace_id,omitempty"`
	SpanID     string            `json:"span_id,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Timestamp  time.Time         `json:"timestamp"`
}

// NewEntry creates a new log entry with the given level and message.
func NewEntry(level EntryLogLevel, message string, attributes map[string]string) Entry {
	return Entry{
		Level:      level,
		Message:    message,
		Attributes: attributes,
		Timestamp:  time.Now(),
	}
}

// NewEntryWithTrace creates a new log entry with trace correlation.
func NewEntryWithTrace(ctx context.Context, level EntryLogLevel, message string, attributes map[string]string) Entry {
	log := NewEntry(level, message, attributes)

	// Extract trace context if available
	if traceID, spanID, ok := TraceFromContext(ctx); ok {
		log.TraceID = traceID
		log.SpanID = spanID
	}

	return log
}
