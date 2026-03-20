package db

import (
	"context"
	"sync"

	"github.com/jgfranco17/observability-platform/internal/observability"
)

type internalClient struct {
	storage map[string]observability.Report // In-memory storage for demonstration
	mu      sync.RWMutex                    // Mutex to handle concurrent access
}

var clientSingleton DatabaseClient = &internalClient{
	storage: make(map[string]observability.Report),
}

// StoreReport stores a report in the in-memory storage.
func (ic *internalClient) StoreReport(ctx context.Context, report observability.Report) error {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	if _, exists := ic.storage[report.ID]; exists {
		return ErrConflict
	}
	ic.storage[report.ID] = report
	return nil
}

func (ic *internalClient) GetAllReports(ctx context.Context) ([]observability.Report, error) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	reports := make([]observability.Report, 0, len(ic.storage))
	for _, report := range ic.storage {
		reports = append(reports, report)
	}
	return reports, nil
}
