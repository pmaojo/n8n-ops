package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pmaojo/n8n-ops/internal/observability"
	"github.com/sirupsen/logrus"
)

// Config configures the REST server.
type Config struct {
	Enabled bool
	Port    int
}

// Service exposes status and metrics endpoints.
type Service struct {
	cfg     Config
	logger  *logrus.Logger
	grafana *observability.GrafanaIntegration
	srv     *http.Server
}

// NewService creates a new Service.
func NewService(cfg Config, logger *logrus.Logger, g *observability.GrafanaIntegration) *Service {
	if logger == nil {
		logger = logrus.New()
	}
	return &Service{cfg: cfg, logger: logger, grafana: g}
}

// Start begins serving HTTP requests until the context is canceled.
func (s *Service) Start(ctx context.Context) error {
	if !s.cfg.Enabled {
		s.logger.Info("REST server disabled")
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", s.statusHandler)
	mux.HandleFunc("/api/metrics", s.metricsHandler)

	s.srv = &http.Server{
		Addr:    ":" + strconv.Itoa(s.cfg.Port),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.srv.Shutdown(shutdownCtx)
	}()

	s.logger.Infof("REST server listening on %d", s.cfg.Port)
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Service) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *Service) statusHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	}
	s.writeJSON(w, resp)
}

func (s *Service) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if s.grafana == nil {
		http.NotFound(w, r)
		return
	}
	metrics := s.grafana.LatestMetrics()
	s.writeJSON(w, metrics)
}
