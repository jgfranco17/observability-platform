package handlers

import (
	"net/http"
	"time"

	"github.com/jgfranco17/observability-platform/internal/logging"
)

// GetHealth is a simple health check endpoint that returns a JSON
// response with the status and timestamp.
func GetHealthHandler() http.HandlerFunc {
	type healthResponse struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.FromContext(r.Context())
		logger.Debug("Health check requested")

		respondWithJSON(w, http.StatusOK, healthResponse{
			Status:    "healthy",
			Timestamp: time.Now().Format(time.DateTime),
		})
	}
}
