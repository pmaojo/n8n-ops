package monitoring

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pmaojo/n8n-ops/internal/issues"
	"github.com/pmaojo/n8n-ops/internal/workflow"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestFailureThreshold(t *testing.T) {
	// Test consecutive failure threshold logic
	threshold := 2
	failures := 3

	if failures < threshold {
		t.Error("Should trigger issue creation when threshold exceeded")
	}

	failures = 1
	if failures >= threshold {
		t.Error("Should not trigger with failures below threshold")
	}
}

func TestMonitoringContext(t *testing.T) {
	// Test context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Simulate monitoring loop with context
	select {
	case <-ctx.Done():
		// Context cancellation works as expected
	case <-time.After(200 * time.Millisecond):
		t.Error("Context should have been cancelled")
	}
}

func TestWorkflowFailureData(t *testing.T) {
	// Test workflow failure data structure
	failedExecution := map[string]interface{}{
		"id":         "exec_1001",
		"status":     "error",
		"workflowId": "1001",
	}

	if failedExecution["status"] != "error" {
		t.Error("Failed execution should have error status")
	}

	if failedExecution["id"] == "" {
		t.Error("Failed execution should have an ID")
	}
}

// mockClient implements client.Client with only the methods needed for tests.
type mockClient struct{ wf *workflow.Workflow }

func (m *mockClient) GetWorkflows(ctx context.Context) ([]*workflow.Workflow, error) { return nil, nil }
func (m *mockClient) GetWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	return m.wf, nil
}
func (m *mockClient) HealthCheck(ctx context.Context) error { return nil }
func (m *mockClient) CreateWorkflow(ctx context.Context, wf *workflow.Workflow) (*workflow.Workflow, error) {
	return nil, nil
}
func (m *mockClient) UpdateWorkflow(ctx context.Context, id string, wf *workflow.Workflow) (*workflow.Workflow, error) {
	return nil, nil
}
func (m *mockClient) DeleteWorkflow(ctx context.Context, id string) error { return nil }
func (m *mockClient) ExecuteWorkflow(ctx context.Context, id string) (*workflow.ExecutionResult, error) {
	return nil, nil
}
func (m *mockClient) GetExecution(ctx context.Context, id string) (*workflow.ExecutionResult, error) {
	return nil, nil
}
func (m *mockClient) GetExecutions(ctx context.Context, workflowID, status string, limit int) ([]*workflow.ExecutionResult, error) {
	return nil, nil
}

type mockIssueManager struct{ open map[string][]*issues.Issue }

func (m *mockIssueManager) CreateWorkflowFailureIssue(ctx context.Context, failure *issues.WorkflowFailure) (*issues.Issue, error) {
	return &issues.Issue{ID: 1, WebURL: "https://example.com/issues/1"}, nil
}
func (m *mockIssueManager) UpdateIssueWithRecovery(ctx context.Context, issueID int, recovery *issues.RecoveryInfo) error {
	return nil
}
func (m *mockIssueManager) GetOpenFailureIssues(ctx context.Context, workflowID string) ([]*issues.Issue, error) {
	return m.open[workflowID], nil
}

func TestCreateFailureIssueLogs(t *testing.T) {
	logger := logrus.New()
	hook := test.NewLocal(logger)

	wf := &workflow.Workflow{ID: "1001", Name: "Test"}
	client := &mockClient{wf: wf}
	mgr := &mockIssueManager{}

	fd := NewFailureDetector(client, mgr, logger)

	if err := fd.createFailureIssue(context.Background(), wf.ID); err != nil {
		t.Fatalf("createFailureIssue: %v", err)
	}

	last := hook.LastEntry()
	expected := fmt.Sprintf("🚨 Issue created for workflow failure: %s", "https://example.com/issues/1")
	if last == nil || last.Message != expected {
		t.Fatalf("expected log %q, got %v", expected, last)
	}
}

func TestHandleWorkflowRecoveryLogs(t *testing.T) {
	logger := logrus.New()
	hook := test.NewLocal(logger)

	wf := &workflow.Workflow{ID: "1001", Name: "Test"}
	client := &mockClient{wf: wf}
	mgr := &mockIssueManager{open: map[string][]*issues.Issue{
		"1001": {{ID: 1, WebURL: "https://example.com/issues/1"}},
	}}

	fd := NewFailureDetector(client, mgr, logger)

	if err := fd.handleWorkflowRecovery(context.Background(), wf.ID); err != nil {
		t.Fatalf("handleWorkflowRecovery: %v", err)
	}

	last := hook.LastEntry()
	expected := fmt.Sprintf("✅ Workflow %s recovered - updated %d issues", wf.ID, 1)
	if last == nil || last.Message != expected {
		t.Fatalf("expected log %q, got %v", expected, last)
	}
}
