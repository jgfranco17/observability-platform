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
	router *chi.Mux
	logger *logrus.Logger
	config config.ServiceSettings
}

func New(ctx context.Context, cfg config.ServiceSettings, logger *logrus.Logger, dbFactory db.DatabaseClientFactory) (*Service, error) {
	database, err := dbFactory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}
	s := &Service{
		router: chi.NewRouter(),
		logger: logger,
		config: cfg,
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

func (s *Service) registerRoutes(database db.DatabaseClient) {
	s.router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Get("/health", handlers.GetHealthHandler())
		apiRouter.Route("/v1", func(v1Router chi.Router) {
			v1Router.Get("/observability", handlers.HandlerGetReports(database))
			v1Router.Post("/observability", handlers.HandlerAddReport(database))
		})
	})
}
