package main

import (
	"context"
	"os"

	"github.com/jgfranco17/observability-platform/internal/config"
	"github.com/jgfranco17/observability-platform/internal/db"
	"github.com/jgfranco17/observability-platform/internal/logging"
	"github.com/jgfranco17/observability-platform/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logging.New(os.Stderr, logrus.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatalf("Failed to load config")
	}

	server, err := service.New(ctx, cfg, logger, db.NewClient)
	if err != nil {
		logger.WithError(err).Fatalf("Failed to create service")
	}
	if err := server.Start(); err != nil {
		logger.WithError(err).Fatalf("Failed to start service")
	}
}
