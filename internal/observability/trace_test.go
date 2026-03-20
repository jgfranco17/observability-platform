package observability

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTrace(t *testing.T) {
	trace := newTrace("test-operation", SpanKindInternal)

	assert.NotEmpty(t, trace.TraceID)
	assert.NotEmpty(t, trace.SpanID)
	assert.Empty(t, trace.ParentID)
	assert.Equal(t, "test-operation", trace.Name)
	assert.Equal(t, SpanKindInternal, trace.Kind)
	assert.Equal(t, SpanStatusOK, trace.Status)
	assert.NotZero(t, trace.StartTime)
	assert.NotNil(t, trace.Attributes)
}

func TestNewChildTrace(t *testing.T) {
	parentTraceID := "parent-trace-id"
	parentSpanID := "parent-span-id"

	trace := newChildTrace(parentTraceID, parentSpanID, "child-operation", SpanKindClient)

	assert.Equal(t, parentTraceID, trace.TraceID)
	assert.NotEmpty(t, trace.SpanID)
	assert.NotEqual(t, parentSpanID, trace.SpanID)
	assert.Equal(t, parentSpanID, trace.ParentID)
	assert.Equal(t, "child-operation", trace.Name)
	assert.Equal(t, SpanKindClient, trace.Kind)
}

func TestSpan_End(t *testing.T) {
	trace := newTrace("test-operation", SpanKindInternal)
	span := &Span{
		trace:     trace,
		startTime: trace.StartTime,
		ended:     false,
	}

	time.Sleep(10 * time.Millisecond)
	span.End()

	assert.True(t, span.ended)
	assert.NotZero(t, trace.EndTime)
	assert.Greater(t, trace.Duration, time.Duration(0))
	assert.GreaterOrEqual(t, trace.Duration, 10*time.Millisecond)
}

func TestSpan_End_Idempotent(t *testing.T) {
	trace := newTrace("test-operation", SpanKindInternal)
	span := &Span{
		trace:     trace,
		startTime: trace.StartTime,
		ended:     false,
	}

	span.End()
	firstEndTime := trace.EndTime
	firstDuration := trace.Duration

	time.Sleep(10 * time.Millisecond)
	span.End()

	assert.Equal(t, firstEndTime, trace.EndTime)
	assert.Equal(t, firstDuration, trace.Duration)
}

func TestSpan_SetStatus(t *testing.T) {
	trace := newTrace("test-operation", SpanKindInternal)
	span := &Span{
		trace:     trace,
		startTime: trace.StartTime,
	}

	assert.Equal(t, SpanStatusOK, trace.Status)

	span.SetStatus(SpanStatusError)
	assert.Equal(t, SpanStatusError, trace.Status)
}

func TestSpan_SetAttribute(t *testing.T) {
	trace := newTrace("test-operation", SpanKindInternal)
	span := &Span{
		trace:     trace,
		startTime: trace.StartTime,
	}

	span.SetAttribute("key1", "value1")
	span.SetAttribute("key2", "value2")

	assert.Equal(t, "value1", trace.Attributes["key1"])
	assert.Equal(t, "value2", trace.Attributes["key2"])
}

func TestSpan_Getters(t *testing.T) {
	trace := newTrace("test-operation", SpanKindInternal)
	span := &Span{
		trace:     trace,
		startTime: trace.StartTime,
	}

	assert.Equal(t, trace.TraceID, span.TraceID())
	assert.Equal(t, trace.SpanID, span.SpanID())
	assert.Equal(t, trace, span.Trace())
}

func TestContextWithTrace(t *testing.T) {
	ctx := context.Background()
	traceID := "test-trace-id"
	spanID := "test-span-id"

	ctx = ContextWithTrace(ctx, traceID, spanID)

	extractedTraceID, extractedSpanID, ok := TraceFromContext(ctx)
	require.True(t, ok)
	assert.Equal(t, traceID, extractedTraceID)
	assert.Equal(t, spanID, extractedSpanID)
}

func TestTraceFromContext_Empty(t *testing.T) {
	ctx := context.Background()

	_, _, ok := TraceFromContext(ctx)
	assert.False(t, ok)
}

func TestTraceFromContext_InvalidType(t *testing.T) {
	type wrongKey struct{}
	ctx := context.WithValue(context.Background(), wrongKey{}, "some value")

	_, _, ok := TraceFromContext(ctx)
	assert.False(t, ok)
}

func TestSpanKind_Values(t *testing.T) {
	testCases := []struct {
		name string
		kind SpanKind
	}{
		{name: "internal", kind: SpanKindInternal},
		{name: "server", kind: SpanKindServer},
		{name: "client", kind: SpanKindClient},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trace := newTrace("test", tc.kind)
			assert.Equal(t, tc.kind, trace.Kind)
		})
	}
}
