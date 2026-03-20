package handlers

import (
	"net/http"

	"github.com/jgfranco17/observability-platform/internal/db"
	"github.com/jgfranco17/observability-platform/internal/logging"
)

func NewObservabilityHandler(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		logger.Info("System info requested")
	}
}
