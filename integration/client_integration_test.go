package integration_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/workflow"
)

func startMockServer(t *testing.T) func() {
	t.Helper()

	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = filepath.Join("..", "mock-n8n-server")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start mock server: %v", err)
	}

	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	if err := waitForServer("http://localhost:3001/health", 5*time.Second); err != nil {
		cmd.Process.Kill()
		<-done
		t.Fatalf("mock server did not start: %v", err)
	}

	return func() {
		cmd.Process.Kill()
		<-done
	}
}

func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for server")
		}
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestClientWorkflowCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	stop := startMockServer(t)
	defer stop()

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
