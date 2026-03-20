package observability

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLog(t *testing.T) {
	level := LogLevelInfo
	message := "test log message"
	attributes := map[string]string{"user_id": "123", "action": "login"}

	before := time.Now()
	log := NewEntry(level, message, attributes)
	after := time.Now()

	assert.Equal(t, level, log.Level)
	assert.Equal(t, message, log.Message)
	assert.Equal(t, attributes, log.Attributes)
	assert.Empty(t, log.TraceID)
	assert.Empty(t, log.SpanID)
	assert.True(t, log.Timestamp.After(before) || log.Timestamp.Equal(before))
	assert.True(t, log.Timestamp.Before(after) || log.Timestamp.Equal(after))
}

func TestNewLog_NilAttributes(t *testing.T) {
	log := NewEntry(LogLevelDebug, "test message", nil)

	assert.Equal(t, LogLevelDebug, log.Level)
	assert.Equal(t, "test message", log.Message)
	assert.Nil(t, log.Attributes)
	assert.NotZero(t, log.Timestamp)
}

func TestNewLogWithTrace_WithContext(t *testing.T) {
	traceID := "test-trace-id"
	spanID := "test-span-id"
	ctx := ContextWithTrace(context.Background(), traceID, spanID)

	log := NewEntryWithTrace(ctx, LogLevelError, "error occurred", map[string]string{"error": "timeout"})

	assert.Equal(t, LogLevelError, log.Level)
	assert.Equal(t, "error occurred", log.Message)
	assert.Equal(t, traceID, log.TraceID)
	assert.Equal(t, spanID, log.SpanID)
	assert.NotNil(t, log.Attributes)
}

func TestNewLogWithTrace_WithoutContext(t *testing.T) {
	ctx := context.Background()

	log := NewEntryWithTrace(ctx, LogLevelWarn, "warning message", nil)

	assert.Equal(t, LogLevelWarn, log.Level)
	assert.Equal(t, "warning message", log.Message)
	assert.Empty(t, log.TraceID)
	assert.Empty(t, log.SpanID)
}

func TestLogLevels(t *testing.T) {
	testCases := []struct {
		name     string
		level    EntryLogLevel
		expected string
	}{
		{name: "debug level", level: LogLevelDebug, expected: "debug"},
		{name: "info level", level: LogLevelInfo, expected: "info"},
		{name: "warn level", level: LogLevelWarn, expected: "warn"},
		{name: "error level", level: LogLevelError, expected: "error"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			log := NewEntry(tc.level, "test", nil)
			assert.Equal(t, tc.expected, string(log.Level))
		})
	}
}

func TestNewLogWithTrace_PreservesAttributes(t *testing.T) {
	ctx := ContextWithTrace(context.Background(), "trace-123", "span-456")
	attributes := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	log := NewEntryWithTrace(ctx, LogLevelInfo, "test message", attributes)

	assert.Equal(t, "value1", log.Attributes["key1"])
	assert.Equal(t, "value2", log.Attributes["key2"])
	assert.Equal(t, "trace-123", log.TraceID)
	assert.Equal(t, "span-456", log.SpanID)
}

func TestLog_EmptyMessage(t *testing.T) {
	log := NewEntry(LogLevelInfo, "", nil)

	assert.Empty(t, log.Message)
	assert.Equal(t, LogLevelInfo, log.Level)
}
