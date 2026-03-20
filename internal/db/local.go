package db

import (
	"context"
	"sync"

	"github.com/jgfranco17/observability-platform/internal/observability"
)

// internalClient is an in-memory database implementation.
type internalClient struct {
	traces  []observability.Trace  // In-memory trace storage
	metrics []observability.Metric // In-memory metric storage
	logs    []observability.Entry  // In-memory log storage
	mu      sync.RWMutex           // Mutex to handle concurrent access
}

var clientSingleton DatabaseClient = &internalClient{
	traces:  make([]observability.Trace, 0),
	metrics: make([]observability.Metric, 0),
	logs:    make([]observability.Entry, 0),
}

// StoreTraces stores traces in memory.
func (ic *internalClient) StoreTraces(ctx context.Context, traces []observability.Trace) error {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	ic.traces = append(ic.traces, traces...)
	return nil
}

// GetTraces retrieves traces by trace ID.
func (ic *internalClient) GetTraces(ctx context.Context, traceID string) ([]observability.Trace, error) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	var result []observability.Trace
	for _, trace := range ic.traces {
		if trace.TraceID == traceID {
			result = append(result, trace)
		}
	}
	return result, nil
}

// GetAllTraces retrieves all traces.
func (ic *internalClient) GetAllTraces(ctx context.Context) ([]observability.Trace, error) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	return ic.traces, nil
}

// StoreMetrics stores metrics in memory.
func (ic *internalClient) StoreMetrics(ctx context.Context, metrics []observability.Metric) error {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	ic.metrics = append(ic.metrics, metrics...)
	return nil
}

// GetMetrics retrieves metrics by name.
func (ic *internalClient) GetMetrics(ctx context.Context, name string) ([]observability.Metric, error) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	var result []observability.Metric
	for _, metric := range ic.metrics {
		if metric.Name == name {
			result = append(result, metric)
		}
	}
	return result, nil
}

// GetAllMetrics retrieves all metrics.
func (ic *internalClient) GetAllMetrics(ctx context.Context) ([]observability.Metric, error) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	return ic.metrics, nil
}

// StoreLogs stores logs in memory.
func (ic *internalClient) StoreLogs(ctx context.Context, logs []observability.Entry) error {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	ic.logs = append(ic.logs, logs...)
	return nil
}

// GetLogs retrieves logs by trace ID.
func (ic *internalClient) GetLogs(ctx context.Context, traceID string) ([]observability.Entry, error) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	var result []observability.Entry
	for _, log := range ic.logs {
		if log.TraceID == traceID {
			result = append(result, log)
		}
	}
	return result, nil
}

// GetAllLogs retrieves all logs.
func (ic *internalClient) GetAllLogs(ctx context.Context) ([]observability.Entry, error) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	return ic.logs, nil
}
