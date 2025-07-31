package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/kdimtriCP/gh-inspector/internal/github"
	"github.com/kdimtriCP/gh-inspector/internal/metrics"
	_ "github.com/kdimtriCP/gh-inspector/swagger"
)

type Server struct {
	router          chi.Router
	analyzer        github.Analyzer
	httpServer      *http.Server
	config          *Config
	metricsRecorder metrics.Recorder
}

type Config struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Port:         8080,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func New(analyzer github.Analyzer, config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	metricsRecorder := NewMetricsRecorder()

	s := &Server{
		router:          chi.NewRouter(),
		analyzer:        analyzer,
		config:          config,
		metricsRecorder: metricsRecorder,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(loggingMiddleware)
	s.router.Use(corsMiddleware)

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Post("/score", s.handleScore)
	})

	s.router.Get("/health", s.handleHealth)
	s.router.Handle("/metrics", promhttp.Handler())

	// Swagger documentation
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	s.router.Get("/", s.handleRoot)
}

func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	fmt.Printf("Server starting on port %d...\n", s.config.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *Server) MetricsRecorder() metrics.Recorder {
	return s.metricsRecorder
}
