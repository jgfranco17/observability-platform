package db

import (
	"context"
	"sync"

	"github.com/jgfranco17/observability-platform/internal/observability"
)

// DatabaseClient interface for database operations
type DatabaseClient interface {
	StoreReport(ctx context.Context, report observability.Report) error
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
