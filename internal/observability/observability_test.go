package observability

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestSentryIntegration(t *testing.T) {
	// Test Sentry configuration
	config := SentryConfig{
		DSN:         "https://test@sentry.io/123456",
		Environment: "test",
		Release:     "1.0.0",
		SampleRate:  1.0,
	}

	logger := logrus.New()
	sentry := NewSentryIntegration(config, logger)

	if sentry == nil {
		t.Error("Sentry integration should not be nil")
	}

	if sentry.config.Environment != "test" {
		t.Error("Environment should be set correctly")
	}
}

func TestGrafanaIntegration(t *testing.T) {
	// Test Grafana configuration
	config := GrafanaConfig{
		URL:       "http://localhost:3000",
		APIKey:    "test-api-key",
		OrgID:     1,
		Dashboard: "n8n-ops-test",
	}

	logger := logrus.New()
	grafana := NewGrafanaIntegration(config, logger)

	if grafana == nil {
		t.Error("Grafana integration should not be nil")
	}

	if grafana.config.URL != "http://localhost:3000" {
		t.Error("URL should be set correctly")
	}
}

func TestGrafanaMetrics(t *testing.T) {
	// Test metrics structure
	metrics := GrafanaMetrics{
		WorkflowExecutions: 10,
		FailureRate:       0.1,
		SyncOperations:    5,
		ActiveWorkflows:   3,
		ResponseTime:      150.5,
		Environment:       "test",
		Timestamp:         time.Now().Unix(),
	}

	if metrics.WorkflowExecutions != 10 {
		t.Error("WorkflowExecutions should be 10")
	}

	if metrics.FailureRate != 0.1 {
		t.Error("FailureRate should be 0.1")
	}

	if metrics.Environment != "test" {
		t.Error("Environment should be test")
	}
}

func TestSentryConfiguration(t *testing.T) {
	// Test configuration validation
	validConfigs := []SentryConfig{
		{
			DSN:         "https://key@sentry.io/project",
			Environment: "production",
			Release:     "1.0.0",
			SampleRate:  1.0,
		},
		{
			DSN:         "https://key@sentry.io/project",
			Environment: "development",
			Release:     "dev",
			SampleRate:  0.1,
		},
	}

	for _, config := range validConfigs {
		if config.DSN == "" {
			t.Error("DSN should not be empty")
		}
		if config.Environment == "" {
			t.Error("Environment should not be empty")
		}
		if config.SampleRate < 0 || config.SampleRate > 1 {
			t.Error("SampleRate should be between 0 and 1")
		}
	}
}

func TestGrafanaConfiguration(t *testing.T) {
	// Test configuration validation
	validConfigs := []GrafanaConfig{
		{
			URL:       "https://grafana.example.com",
			APIKey:    "grafana-api-key-123",
			OrgID:     1,
			Dashboard: "n8n-ops-monitoring",
		},
		{
			URL:       "http://localhost:3000",
			APIKey:    "local-api-key",
			OrgID:     2,
			Dashboard: "local-dashboard",
		},
	}

	for _, config := range validConfigs {
		if config.URL == "" {
			t.Error("URL should not be empty")
		}
		if config.APIKey == "" {
			t.Error("APIKey should not be empty")
		}
		if config.OrgID <= 0 {
			t.Error("OrgID should be positive")
		}
	}
}