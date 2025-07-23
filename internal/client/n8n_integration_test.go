//go:build integration
// +build integration

package client

import (
	"context"
	"testing"
	"time"

	"github.com/pmaojo/n8n-ops/internal/workflow"
)

// Integration tests - run with: go test -tags=integration

func TestIntegrationHealthCheck(t *testing.T) {
	client, err := New("http://localhost:3001", "n8n_api_mock_development", nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
}

func TestIntegrationGetWorkflows(t *testing.T) {
	client, err := New("http://localhost:3001", "n8n_api_mock_development", nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	workflows, err := client.GetWorkflows(ctx)
	if err != nil {
		t.Errorf("GetWorkflows failed: %v", err)
		return
	}

	t.Logf("Retrieved %d workflows", len(workflows))
}

func TestIntegrationCreateAndDeleteWorkflow(t *testing.T) {
	client, err := New("http://localhost:3001", "n8n_api_mock_development", nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create test workflow
	testWorkflow := &workflow.Workflow{
		Name:   "Test Integration Workflow",
		Active: false,
		Nodes: []workflow.Node{
			{
				Name:       "Start",
				Type:       "n8n-nodes-base.start",
				Position:   []float64{240, 300},
				Parameters: make(map[string]interface{}),
			},
		},
		Settings: make(map[string]interface{}),
	}

	created, err := client.CreateWorkflow(ctx, testWorkflow)
	if err != nil {
		t.Errorf("CreateWorkflow failed: %v", err)
		return
	}

	t.Logf("Created workflow with ID: %s", created.ID)

	// Clean up - delete the workflow
	err = client.DeleteWorkflow(ctx, created.ID)
	if err != nil {
		t.Errorf("DeleteWorkflow failed: %v", err)
	}
}

// Benchmark tests against mock server
func BenchmarkHealthCheckMock(b *testing.B) {
	client, err := New("http://localhost:3001", "n8n_api_mock_development", nil)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.HealthCheck(ctx)
	}
}

func BenchmarkGetWorkflowsMock(b *testing.B) {
	client, err := New("http://localhost:3001", "n8n_api_mock_development", nil)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.GetWorkflows(ctx)
	}
}
