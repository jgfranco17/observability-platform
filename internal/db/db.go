package db

import (
	"context"
	"sync"

	"github.com/jgfranco17/observability-platform/internal/observability"
)

// DatabaseClient interface for database operations
type DatabaseClient interface {
	// Trace storage
	StoreTraces(ctx context.Context, traces []observability.Trace) error
	GetTraces(ctx context.Context, traceID string) ([]observability.Trace, error)
	GetAllTraces(ctx context.Context) ([]observability.Trace, error)

	// Metric storage
	StoreMetrics(ctx context.Context, metrics []observability.Metric) error
	GetMetrics(ctx context.Context, name string) ([]observability.Metric, error)
	GetAllMetrics(ctx context.Context) ([]observability.Metric, error)

	// Log storage
	StoreLogs(ctx context.Context, logs []observability.Entry) error
	GetLogs(ctx context.Context, traceID string) ([]observability.Entry, error)
	GetAllLogs(ctx context.Context) ([]observability.Entry, error)
}

type DatabaseClientFactory func(ctx context.Context) (DatabaseClient, error)

var mu sync.RWMutex

// NewClient creates a new database client. Locks provides thread safety for
// singleton access.
func NewClient(_ context.Context) (DatabaseClient, error) {
	mu.Lock()
	defer mu.Unlock()

	return clientSingleton, nil
}
