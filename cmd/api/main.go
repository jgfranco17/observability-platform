package main

import (
	"os"

	"github.com/jgfranco17/observability-platform/internal/logging"
	"github.com/jgfranco17/observability-platform/internal/service"
	"github.com/jgfranco17/observability-platform/internal/service/config"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logging.New(os.Stderr, logrus.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	server := service.New(cfg, logger)
	if err := server.Start(); err != nil {
		logger.Fatalf("Failed to start service: %v", err)
	}
}
