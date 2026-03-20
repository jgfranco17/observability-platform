package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jgfranco17/observability-platform/internal/logging"
	"github.com/sirupsen/logrus"
)

type Client struct {
	baseURL    url.URL
	doer       httpDoer
	identifier string
}

// NewClient creates a new observability client with the given base URL.
func NewClient(identifier string, baseURL string) (*Client, error) {
	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	return &Client{
		baseURL:    *parsedURL,
		identifier: identifier,
		doer: &http.Client{
			Timeout: 5 * time.Second,
			Jar:     &jar{},
		},
	}, nil
}

// StartSpan creates a new span and returns an updated context with trace information.
// If the context already contains a trace, the new span will be a child span.
func (c *Client) StartSpan(ctx context.Context, name string, kind SpanKind) (context.Context, *Span) {
	var trace *Trace

	// Check if we have a parent trace in context
	if parentTraceID, parentSpanID, ok := TraceFromContext(ctx); ok {
		trace = newChildTrace(parentTraceID, parentSpanID, name, kind)
	} else {
		trace = newTrace(name, kind)
	}

	span := &Span{
		trace:     trace,
		startTime: trace.StartTime,
		ended:     false,
	}

	// Add trace context to the returned context
	ctx = ContextWithTrace(ctx, trace.TraceID, trace.SpanID)

	return ctx, span
}

// SendTrace sends a completed trace to the backend.
func (c *Client) SendTrace(ctx context.Context, trace *Trace) error {
	resp, err := c.makeJSONRequest(ctx, "/api/v1/traces", []Trace{*trace})
	if err != nil {
		return fmt.Errorf("failed to send trace: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// SendTraces sends multiple traces to the backend in a single request.
func (c *Client) SendTraces(ctx context.Context, traces []Trace) error {
	if len(traces) == 0 {
		return nil
	}
	resp, err := c.makeJSONRequest(ctx, "/api/v1/traces", traces)
	if err != nil {
		return fmt.Errorf("failed to send traces: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// RecordCounter records a counter metric and sends it to the backend.
func (c *Client) RecordCounter(ctx context.Context, name string, value float64, labels MetricLabels) error {
	metric := NewCounter(name, value, labels)
	return c.SendMetric(ctx, metric)
}

// RecordGauge records a gauge metric and sends it to the backend.
func (c *Client) RecordGauge(ctx context.Context, name string, value float64, labels MetricLabels) error {
	metric := NewGauge(name, value, labels)
	return c.SendMetric(ctx, metric)
}

// RecordHistogram records a histogram metric and sends it to the backend.
func (c *Client) RecordHistogram(ctx context.Context, name string, value float64, labels MetricLabels) error {
	metric := NewHistogram(name, value, labels)
	return c.SendMetric(ctx, metric)
}

// SendMetric sends a single metric to the backend.
func (c *Client) SendMetric(ctx context.Context, metric Metric) error {
	resp, err := c.makeJSONRequest(ctx, "/api/v1/metrics", []Metric{metric})
	if err != nil {
		return fmt.Errorf("failed to send metric: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// SendMetrics sends multiple metrics to the backend in a single request.
func (c *Client) SendMetrics(ctx context.Context, metrics []Metric) error {
	if len(metrics) == 0 {
		return nil
	}
	resp, err := c.makeJSONRequest(ctx, "/api/v1/metrics", metrics)
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// Log sends a structured log entry to the backend.
// If the context contains trace information, it will be automatically included.
func (c *Client) Log(ctx context.Context, level EntryLogLevel, message string, attributes map[string]string) error {
	log := NewEntryWithTrace(ctx, level, message, attributes)
	return c.SendLog(ctx, log)
}

// SendLog sends a single log entry to the backend.
func (c *Client) SendLog(ctx context.Context, log Entry) error {
	resp, err := c.makeJSONRequest(ctx, "/api/v1/logs", []Entry{log})
	if err != nil {
		return fmt.Errorf("failed to send log: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// SendLogs sends multiple log entries to the backend in a single request.
func (c *Client) SendLogs(ctx context.Context, logs []Entry) error {
	if len(logs) == 0 {
		return nil
	}
	resp, err := c.makeJSONRequest(ctx, "/api/v1/logs", logs)
	if err != nil {
		return fmt.Errorf("failed to send logs: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// makeJSONRequest is a helper method to send a JSON-encoded POST request to the
// specified platform endpoint.
func (c *Client) makeJSONRequest(ctx context.Context, endpoint string, data any) (*http.Response, error) {
	logger := logging.FromContext(ctx)

	url := c.baseURL.JoinPath(endpoint).String()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		logger.WithFields(logrus.Fields{
			"endpoint": endpoint,
			"status":   resp.StatusCode,
		}).Trace("Request successful")
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		return resp, fmt.Errorf("access denied, service returned HTTP %d", resp.StatusCode)
	default:
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}
