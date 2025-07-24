package observability

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

// SentryConfig holds configuration for Sentry integration
type SentryConfig struct {
	DSN         string
	Environment string
	Release     string
	SampleRate  float64
}

// SentryIntegration handles error tracking and performance monitoring
type SentryIntegration struct {
	config SentryConfig
	logger *logrus.Logger
	hub    *sentry.Hub
}

// NewSentryIntegration creates a new Sentry integration
func NewSentryIntegration(config SentryConfig, logger *logrus.Logger) *SentryIntegration {
	return &SentryIntegration{
		config: config,
		logger: logger,
	}
}

// Initialize sets up the Sentry SDK
func (s *SentryIntegration) Initialize() error {
	if s.config.DSN == "" {
		return errors.New("sentry DSN is required")
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              s.config.DSN,
		Environment:      s.config.Environment,
		Release:          s.config.Release,
		TracesSampleRate: s.config.SampleRate,
	})
	if err != nil {
		return fmt.Errorf("initialize sentry: %w", err)
	}

	s.hub = sentry.CurrentHub().Clone()

	s.logger.WithFields(logrus.Fields{
		"environment": s.config.Environment,
		"release":     s.config.Release,
		"sample_rate": s.config.SampleRate,
	}).Info("Sentry integration initialized")

	return nil
}

// CaptureWorkflowFailure reports workflow failures to Sentry
func (s *SentryIntegration) CaptureWorkflowFailure(workflowID, workflowName, errMsg string, ctx context.Context) {
	if s.hub == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	s.hub.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("workflow_id", workflowID)
		scope.SetTag("workflow_name", workflowName)
		scope.SetTag("environment", s.config.Environment)
		s.hub.CaptureException(errors.New(errMsg))
	})
}

// CaptureSync reports sync operations
func (s *SentryIntegration) CaptureSync(direction, environment string, workflowCount int, duration time.Duration) {
	if s.hub == nil {
		return
	}

	s.hub.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("direction", direction)
		scope.SetTag("environment", environment)
		scope.SetTag("workflow_count", fmt.Sprint(workflowCount))
		scope.SetTag("duration_ms", fmt.Sprint(duration.Milliseconds()))
		s.hub.CaptureMessage("sync operation")
	})
}

// CapturePerformance starts a performance transaction
func (s *SentryIntegration) CapturePerformance(operation string) interface{} {
	if s.hub == nil {
		return nil
	}
	span := sentry.StartSpan(context.Background(), operation)
	return span
}

// Close flushes pending events
func (s *SentryIntegration) Close() {
	sentry.Flush(2 * time.Second)
}
