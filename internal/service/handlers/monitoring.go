package handlers

import (
	"context"
	"net/http"

	"github.com/jgfranco17/observability-platform/internal/logging"
)

type Checker interface {
	Check(ctx context.Context) error
}

func GetSystemInfo(checkers []Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		logger.Info("System info requested")

		var errs []error
		for _, checker := range checkers {
			if err := checker.Check(ctx); err != nil {
				errs = append(errs, err)
				logger.WithError(err).Warn("Checker failed, proceeding to next")
				continue
			}
		}
	}
}
