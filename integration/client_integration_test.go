package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/testutil"
	"github.com/pmaojo/n8n-ops/internal/workflow"
)

var stopServer func()
var cleanup func()

func TestMain(m *testing.M) {
	var err error
	cleanup, err = testutil.BuildMockServer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	stopServer, err = testutil.StartMockServer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cleanup()
		os.Exit(1)
	}
	code := m.Run()
	if stopServer != nil {
		stopServer()
	}
	if cleanup != nil {
		cleanup()
	}
	os.Exit(code)
}

func TestClientWorkflowCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c, err := client.New("http://localhost:3001", "n8n_api_mock_development", nil)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	ctx := context.Background()
	wf := &workflow.Workflow{
		Name:   "Integration Workflow",
		Active: false,
		Nodes: []workflow.Node{
			{
				Name:     "Start",
				Type:     "n8n-nodes-base.start",
				Position: []float64{240, 300},
			},
		},
		Connections: map[string]interface{}{},
	}

	created, err := c.CreateWorkflow(ctx, wf)
	if err != nil {
		t.Fatalf("create workflow: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected created workflow to have ID")
	}

	fetched, err := c.GetWorkflow(ctx, created.ID)
	if err != nil {
		t.Fatalf("get workflow: %v", err)
	}
	if fetched.Name != wf.Name {
		t.Fatalf("expected name %s, got %s", wf.Name, fetched.Name)
	}

	created.Name = "Updated via Test"
	updated, err := c.UpdateWorkflow(ctx, created.ID, created)
	if err != nil {
		t.Fatalf("update workflow: %v", err)
	}
	if updated.Name != "Updated via Test" {
		t.Fatalf("expected updated name, got %s", updated.Name)
	}

	if err := c.DeleteWorkflow(ctx, created.ID); err != nil {
		t.Fatalf("delete workflow: %v", err)
	}

	if _, err := c.GetWorkflow(ctx, created.ID); err == nil {
		t.Fatalf("expected error fetching deleted workflow")
	}
}
