package integration_test

import (
	"context"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/testutil"
)

func TestClientExecutions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	stop, err := testutil.StartMockServer()
	if err != nil {
		t.Fatalf("start mock server: %v", err)
	}
	defer stop()

	c, err := client.New("http://localhost:3001", "n8n_api_mock_development", nil)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	ctx := context.Background()
	execs, err := c.GetExecutions(ctx, "1001", "", 1)
	if err != nil || len(execs) == 0 {
		t.Fatalf("get executions: %v", err)
	}

	exec, err := c.GetExecution(ctx, execs[0].ID)
	if err != nil {
		t.Fatalf("get execution: %v", err)
	}
	if exec.WorkflowID != "1001" {
		t.Fatalf("unexpected workflow id %s", exec.WorkflowID)
	}
}
