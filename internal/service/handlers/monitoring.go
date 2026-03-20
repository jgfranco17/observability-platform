package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jgfranco17/observability-platform/internal/db"
	"github.com/jgfranco17/observability-platform/internal/logging"
	"github.com/jgfranco17/observability-platform/internal/observability"
	"github.com/sirupsen/logrus"
)

func HandlerAddReport(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)

		var payload []observability.Report
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			errMessage := fmt.Sprintf("invalid request body: %v", err)
			http.Error(w, errMessage, http.StatusBadRequest)
			return
		}
		logger.WithFields(logrus.Fields{
			"count":  len(payload),
			"source": r.RemoteAddr,
		}).Info("Received observability report")

		var errs []error
		for _, report := range payload {
			if err := dbClient.StoreReport(ctx, report); err != nil {
				logger.WithFields(logrus.Fields{
					"id":    report.ID,
					"error": err.Error(),
				}).Warn("Failed to store report")
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			joinedErrs := errors.Join(errs...)
			logger.WithError(joinedErrs).Errorf("Failed to store %d reports", len(errs))
			http.Error(w, "failed to store reports", http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusCreated, jsonResponse{"status": "success"})
	}
}

func HandlerGetReports(dbClient db.DatabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)

		reports, err := dbClient.GetAllReports(ctx)
		if err != nil {
			logger.WithError(err).Error("Failed to retrieve reports")
			http.Error(w, "failed to retrieve reports", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, reports)
	}
}
