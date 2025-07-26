package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// GrafanaConfig holds configuration for Grafana integration
type GrafanaConfig struct {
	URL       string
	APIKey    string
	OrgID     int
	Dashboard string
}

// GrafanaMetrics represents metrics data structure
type GrafanaMetrics struct {
	WorkflowExecutions int     `json:"workflow_executions"`
	FailureRate        float64 `json:"failure_rate"`
	SyncOperations     int     `json:"sync_operations"`
	ActiveWorkflows    int     `json:"active_workflows"`
	ResponseTime       float64 `json:"response_time_ms"`
	Environment        string  `json:"environment"`
	Timestamp          int64   `json:"timestamp"`
}

// GrafanaIntegration handles metrics and dashboard integration
type GrafanaIntegration struct {
	config      GrafanaConfig
	client      *http.Client
	logger      *logrus.Logger
	metricsChan chan GrafanaMetrics
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewGrafanaIntegration creates a new Grafana integration
func NewGrafanaIntegration(config GrafanaConfig, logger *logrus.Logger) *GrafanaIntegration {
	return &GrafanaIntegration{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:      logger,
		metricsChan: make(chan GrafanaMetrics, 100),
	}
}

// Initialize sets up Grafana connection and starts metrics collection
func (g *GrafanaIntegration) Initialize(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	g.ctx, g.cancel = context.WithCancel(ctx)

	if err := g.testConnection(); err != nil {
		g.logger.WithError(err).Error("Failed to connect to Grafana")
		return err
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		g.metricsCollector()
	}()

	g.logger.Info("Grafana integration initialized")
	return nil
}

// testConnection verifies Grafana API connectivity
func (g *GrafanaIntegration) testConnection() error {
	req, err := http.NewRequest("GET", g.config.URL+"/api/health", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+g.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("grafana health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// SendWorkflowMetrics sends workflow execution metrics
func (g *GrafanaIntegration) SendWorkflowMetrics(workflowID, environment string, success bool, duration time.Duration) {
	failure := 0.0
	if !success {
		failure = 1.0
	}

	metrics := GrafanaMetrics{
		WorkflowExecutions: 1,
		FailureRate:        failure,
		ResponseTime:       float64(duration.Milliseconds()),
		Environment:        environment,
		Timestamp:          time.Now().Unix(),
	}

	select {
	case g.metricsChan <- metrics:
	default:
		g.logger.Warn("Metrics channel full, dropping metrics")
	}
}

// SendSyncMetrics sends sync operation metrics
func (g *GrafanaIntegration) SendSyncMetrics(direction, environment string, workflowCount int, duration time.Duration) {
	metrics := GrafanaMetrics{
		SyncOperations:  1,
		ActiveWorkflows: workflowCount,
		ResponseTime:    float64(duration.Milliseconds()),
		Environment:     environment,
		Timestamp:       time.Now().Unix(),
	}

	select {
	case g.metricsChan <- metrics:
	default:
		g.logger.Warn("Metrics channel full, dropping metrics")
	}
}

// metricsCollector processes and sends metrics to Grafana
func (g *GrafanaIntegration) metricsCollector() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	var batchMetrics []GrafanaMetrics

	for {
		select {
		case <-g.ctx.Done():
			if len(batchMetrics) > 0 {
				g.sendMetricsBatch(batchMetrics)
			}
			return
		case metric := <-g.metricsChan:
			batchMetrics = append(batchMetrics, metric)

		case <-ticker.C:
			if len(batchMetrics) > 0 {
				g.sendMetricsBatch(batchMetrics)
				batchMetrics = nil
			}
		}
	}
}

// sendMetricsBatch sends a batch of metrics to Grafana
func (g *GrafanaIntegration) sendMetricsBatch(metrics []GrafanaMetrics) {
	payload := map[string]interface{}{
		"dashboard": g.config.Dashboard,
		"metrics":   metrics,
		"timestamp": time.Now().Unix(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		g.logger.WithError(err).Error("Failed to marshal metrics")
		return
	}

	req, err := http.NewRequest("POST", g.config.URL+"/api/annotations", bytes.NewBuffer(jsonData))
	if err != nil {
		g.logger.WithError(err).Error("Failed to create request")
		return
	}

	req.Header.Set("Authorization", "Bearer "+g.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		g.logger.WithError(err).Error("Failed to send metrics to Grafana")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		g.logger.WithField("count", len(metrics)).Debug("Metrics sent to Grafana")
	} else {
		g.logger.WithField("status", resp.StatusCode).Error("Failed to send metrics")
	}
}

// CreateDashboard creates a default n8n-ops dashboard
func (g *GrafanaIntegration) CreateDashboard(ctx context.Context) error {
	dashboard := map[string]interface{}{
		"dashboard": map[string]interface{}{
			"title": "n8n-ops Monitoring",
			"tags":  []string{"n8n", "automation", "workflows"},
			"panels": []map[string]interface{}{
				{
					"title": "Workflow Executions",
					"type":  "graph",
					"targets": []map[string]interface{}{
						{"expr": "rate(workflow_executions_total[5m])"},
					},
				},
				{
					"title": "Failure Rate",
					"type":  "singlestat",
					"targets": []map[string]interface{}{
						{"expr": "rate(workflow_failures_total[5m]) / rate(workflow_executions_total[5m])"},
					},
				},
				{
					"title": "Sync Operations",
					"type":  "graph",
					"targets": []map[string]interface{}{
						{"expr": "rate(sync_operations_total[5m])"},
					},
				},
			},
		},
		"overwrite": true,
	}

	jsonData, err := json.Marshal(dashboard)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", g.config.URL+"/api/dashboards/db", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+g.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		g.logger.Info("Grafana dashboard created successfully")
		return nil
	}

	return fmt.Errorf("failed to create dashboard, status: %d", resp.StatusCode)
}

// Close stops the metrics collector
func (g *GrafanaIntegration) Close() {
	if g.cancel != nil {
		g.cancel()
	}
	g.wg.Wait()
}
