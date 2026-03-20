package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jgfranco17/observability-platform/internal/db"
	"github.com/jgfranco17/observability-platform/internal/logging"
	"github.com/jgfranco17/observability-platform/internal/observability"
	"github.com/sirupsen/logrus"
)

// HandlerAddTraces handles POST requests to store traces.
func HandlerAddTraces(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)

		var payload []observability.Trace
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			errMessage := fmt.Sprintf("invalid request body: %v", err)
			http.Error(w, errMessage, http.StatusBadRequest)
			return
		}

		logger.WithFields(logrus.Fields{
			"count":  len(payload),
			"source": r.RemoteAddr,
		}).Info("Received traces")

		if err := dbClient.StoreTraces(ctx, payload); err != nil {
			logger.WithError(err).Error("Failed to store traces")
			http.Error(w, "failed to store traces", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, jsonResponse{"status": "success", "count": fmt.Sprintf("%d", len(payload))})
	}
}

// HandlerGetTraces handles GET requests to retrieve all traces.
func HandlerGetTraces(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		traceID := r.URL.Query().Get("trace_id")

		var traces []observability.Trace
		var err error
		if traceID != "" {
			traces, err = dbClient.GetTraces(ctx, traceID)
		} else {
			traces, err = dbClient.GetAllTraces(ctx)
		}
		if err != nil {
			logger.WithError(err).Error("Failed to retrieve traces")
			http.Error(w, "failed to retrieve traces", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, traces)
	}
}

// HandlerAddMetrics handles POST requests to store metrics.
func HandlerAddMetrics(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)

		var payload []observability.Metric
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			errMessage := fmt.Sprintf("invalid request body: %v", err)
			http.Error(w, errMessage, http.StatusBadRequest)
			return
		}

		logger.WithFields(logrus.Fields{
			"count":  len(payload),
			"source": r.RemoteAddr,
		}).Info("Received metrics")

		if err := dbClient.StoreMetrics(ctx, payload); err != nil {
			logger.WithError(err).Error("Failed to store metrics")
			http.Error(w, "failed to store metrics", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, jsonResponse{"status": "success", "count": fmt.Sprintf("%d", len(payload))})
	}
}

// HandlerGetMetrics handles GET requests to retrieve metrics.
func HandlerGetMetrics(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		name := r.URL.Query().Get("name")

		var metrics []observability.Metric
		var err error

		if name != "" {
			metrics, err = dbClient.GetMetrics(ctx, name)
		} else {
			metrics, err = dbClient.GetAllMetrics(ctx)
		}

		if err != nil {
			logger.WithError(err).Error("Failed to retrieve metrics")
			http.Error(w, "failed to retrieve metrics", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, metrics)
	}
}

// HandlerAddLogs handles POST requests to store logs.
func HandlerAddLogs(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)

		var payload []observability.Entry
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			errMessage := fmt.Sprintf("invalid request body: %v", err)
			http.Error(w, errMessage, http.StatusBadRequest)
			return
		}

		logger.WithFields(logrus.Fields{
			"count":  len(payload),
			"source": r.RemoteAddr,
		}).Info("Received logs")

		if err := dbClient.StoreLogs(ctx, payload); err != nil {
			logger.WithError(err).Error("Failed to store logs")
			http.Error(w, "failed to store logs", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusCreated, jsonResponse{"status": "success", "count": fmt.Sprintf("%d", len(payload))})
	}
}

// HandlerGetLogs handles GET requests to retrieve logs.
func HandlerGetLogs(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)

		traceID := r.URL.Query().Get("trace_id")

		var logs []observability.Entry
		var err error

		if traceID != "" {
			logs, err = dbClient.GetLogs(ctx, traceID)
		} else {
			logs, err = dbClient.GetAllLogs(ctx)
		}

		if err != nil {
			logger.WithError(err).Error("Failed to retrieve logs")
			http.Error(w, "failed to retrieve logs", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, logs)
	}
}
