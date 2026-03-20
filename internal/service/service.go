package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jgfranco17/observability-platform/internal/config"
	"github.com/jgfranco17/observability-platform/internal/db"
	"github.com/jgfranco17/observability-platform/internal/logging"
	"github.com/jgfranco17/observability-platform/internal/service/handlers"
	"github.com/sirupsen/logrus"
)

type Service struct {
	router  *chi.Mux
	logger  *logrus.Logger
	config  config.ServiceSettings
	apiKeys map[string]bool // Fast lookup for API keys
}

func New(ctx context.Context, cfg config.ServiceSettings, logger *logrus.Logger, dbFactory db.DatabaseClientFactory) (*Service, error) {
	database, err := dbFactory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}

	// Build API key lookup map
	apiKeys := make(map[string]bool)
	for _, key := range cfg.Auth.APIKeys {
		if key != "" {
			apiKeys[key] = true
		}
	}

	s := &Service{
		router:  chi.NewRouter(),
		logger:  logger,
		config:  cfg,
		apiKeys: apiKeys,
	}
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Recoverer)
	s.router.Use(s.loggerMiddleware)
	s.router.Use(s.errorLoggerMiddleware)
	s.registerRoutes(database)
	return s, nil
}

func (s *Service) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.logger.WithFields(logrus.Fields{
		"host": s.config.Host,
		"port": s.config.Port,
	}).Info("Server starting")
	return http.ListenAndServe(addr, s.router)
}

func (s *Service) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.AddToContext(r.Context(), s.logger)
		s.logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Received request")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) errorLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		if status := ww.Status(); status >= http.StatusBadRequest {
			s.logger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"status": status,
			}).Error("Request encountered an error")
		}
	})
}

// authMiddleware validates API keys when authentication is enabled.
// Expects API key in X-API-Key header.
func (s *Service) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth if disabled
		if !s.config.Auth.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			s.logger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"remote": r.RemoteAddr,
			}).Warn("Request missing API key")
			http.Error(w, "unauthorized: missing API key", http.StatusUnauthorized)
			return
		}

		if !s.apiKeys[apiKey] {
			s.logger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"remote": r.RemoteAddr,
			}).Warn("Request with invalid API key")
			http.Error(w, "unauthorized: invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Service) registerRoutes(database db.DatabaseClient) {
	s.router.Route("/api", func(apiRouter chi.Router) {
		// Health endpoint - no auth required
		apiRouter.Get("/health", handlers.GetHealthHandler())

		// v1 endpoints - with auth
		apiRouter.Route("/v1", func(v1Router chi.Router) {
			// Apply auth middleware to all v1 routes
			v1Router.Use(s.authMiddleware)

			// Trace endpoints
			v1Router.Post("/traces", handlers.HandlerAddTraces(database))
			v1Router.Get("/traces", handlers.HandlerGetTraces(database))

			// Metric endpoints
			v1Router.Post("/metrics", handlers.HandlerAddMetrics(database))
			v1Router.Get("/metrics", handlers.HandlerGetMetrics(database))

			// Log endpoints
			v1Router.Post("/logs", handlers.HandlerAddLogs(database))
			v1Router.Get("/logs", handlers.HandlerGetLogs(database))
		})
	})
}
