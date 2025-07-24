package watch

import (
	"context"
	"testing"
	"time"

	isync "github.com/pmaojo/n8n-ops/internal/sync"
	"github.com/pmaojo/n8n-ops/internal/workflow"
	"github.com/sirupsen/logrus"
)

type mockClient struct{ workflows []*workflow.Workflow }

func (m *mockClient) GetWorkflows(ctx context.Context) ([]*workflow.Workflow, error) {
	return m.workflows, nil
}
func (m *mockClient) GetWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	return nil, nil
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

type fakeGitChecker struct{ calls int }

func (f *fakeGitChecker) AutoCommitWorkflows(message string) (string, error) {
	f.calls++
	return "", nil
}

type fakeSyncer struct{ calls int }

func (f *fakeSyncer) Sync(ctx context.Context, opts isync.Options) error {
	f.calls++
	return nil
}

func TestCheckChanges_AutoSync(t *testing.T) {
	wf := &workflow.Workflow{ID: "1", Name: "Test"}
	client := &mockClient{workflows: []*workflow.Workflow{wf}}
	gitChecker := &fakeGitChecker{}
	syncer := &fakeSyncer{}
	svc := NewService(client, nil, gitChecker, syncer, logrus.New(), "dev")

	state := map[string]time.Time{"1": time.Now().Add(-2 * time.Minute)}
	err := svc.checkChanges(state, Options{AutoSync: true})
	if err != nil {
		t.Fatalf("checkChanges returned error: %v", err)
	}
	if syncer.calls != 1 {
		t.Errorf("expected sync to be called once, got %d", syncer.calls)
	}
	if gitChecker.calls != 0 {
		t.Errorf("did not expect commit when AutoCommit disabled")
	}
}

func TestCheckChanges_AutoCommit(t *testing.T) {
	wf := &workflow.Workflow{ID: "1", Name: "Test"}
	client := &mockClient{workflows: []*workflow.Workflow{wf}}
	gitChecker := &fakeGitChecker{}
	svc := NewService(client, nil, gitChecker, nil, logrus.New(), "dev")

	state := map[string]time.Time{"1": time.Now().Add(-2 * time.Minute)}
	err := svc.checkChanges(state, Options{AutoCommit: true})
	if err != nil {
		t.Fatalf("checkChanges returned error: %v", err)
	}
	if gitChecker.calls != 1 {
		t.Errorf("expected commit to be called once, got %d", gitChecker.calls)
	}
}
