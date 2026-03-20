package observability

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SpanStatus represents the status of a span execution.
type SpanStatus string

const (
	SpanStatusOK    SpanStatus = "ok"
	SpanStatusError SpanStatus = "error"
)

// SpanKind represents the type of span.
type SpanKind string

const (
	SpanKindInternal SpanKind = "internal"
	SpanKindServer   SpanKind = "server"
	SpanKindClient   SpanKind = "client"
)

type TraceAttributes map[string]string

// Trace represents a distributed trace with timing information.
type Trace struct {
	TraceID    string          `json:"trace_id"`
	SpanID     string          `json:"span_id"`
	ParentID   string          `json:"parent_span_id,omitempty"`
	Name       string          `json:"name"`
	Kind       SpanKind        `json:"kind"`
	StartTime  time.Time       `json:"start_time"`
	EndTime    time.Time       `json:"end_time,omitempty"`
	Duration   time.Duration   `json:"duration_ns"`
	Attributes TraceAttributes `json:"attributes,omitempty"`
	Status     SpanStatus      `json:"status"`
}

// Span represents an active span that can be ended.
type Span struct {
	trace     *Trace
	startTime time.Time
	ended     bool
}

// End finalizes the span by setting the end time and calculating duration.
func (s *Span) End() {
	if s.ended {
		return
	}
	s.ended = true
	s.trace.EndTime = time.Now()
	s.trace.Duration = s.trace.EndTime.Sub(s.startTime)
}

// SetStatus sets the span status (ok or error).
func (s *Span) SetStatus(status SpanStatus) {
	s.trace.Status = status
}

// SetAttribute adds a key-value attribute to the span.
func (s *Span) SetAttribute(key, value string) {
	if s.trace.Attributes == nil {
		s.trace.Attributes = make(map[string]string)
	}
	s.trace.Attributes[key] = value
}

// TraceID returns the trace ID of the span.
func (s *Span) TraceID() string {
	return s.trace.TraceID
}

// SpanID returns the span ID.
func (s *Span) SpanID() string {
	return s.trace.SpanID
}

// Trace returns the underlying trace data.
func (s *Span) Trace() *Trace {
	return s.trace
}

// traceContextKey is the key for storing trace context in context.Context.
type traceContextKey struct{}

// traceContext holds trace and span IDs for context propagation.
type traceContext struct {
	traceID string
	spanID  string
}

// ContextWithTrace returns a new context with the trace context embedded.
func ContextWithTrace(ctx context.Context, traceID, spanID string) context.Context {
	return context.WithValue(ctx, traceContextKey{}, &traceContext{
		traceID: traceID,
		spanID:  spanID,
	})
}

// TraceFromContext extracts the trace ID and span ID from context, if present.
func TraceFromContext(ctx context.Context) (traceID, spanID string, ok bool) {
	tc, ok := ctx.Value(traceContextKey{}).(*traceContext)
	if !ok {
		return "", "", false
	}
	return tc.traceID, tc.spanID, true
}

// newTrace creates a new trace with a unique trace and span ID.
func newTrace(name string, kind SpanKind) *Trace {
	return &Trace{
		TraceID:    uuid.New().String(),
		SpanID:     uuid.New().String(),
		Name:       name,
		Kind:       kind,
		StartTime:  time.Now(),
		Attributes: make(TraceAttributes),
		Status:     SpanStatusOK,
	}
}

// newChildTrace creates a new trace as a child of an existing trace.
func newChildTrace(parentTraceID, parentSpanID, name string, kind SpanKind) *Trace {
	return &Trace{
		TraceID:    parentTraceID,
		SpanID:     uuid.New().String(),
		ParentID:   parentSpanID,
		Name:       name,
		Kind:       kind,
		StartTime:  time.Now(),
		Attributes: make(map[string]string),
		Status:     SpanStatusOK,
	}
}
