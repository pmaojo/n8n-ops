package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/workflow"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockWatcher struct {
	events chan fsnotify.Event
	errs   chan error
}

func (m *mockWatcher) Add(name string) error         { return nil }
func (m *mockWatcher) Close() error                  { return nil }
func (m *mockWatcher) Events() <-chan fsnotify.Event { return m.events }
func (m *mockWatcher) Errors() <-chan error          { return m.errs }

type mockClient struct {
	healthCalled bool
	updateCalled bool
}

func (m *mockClient) HealthCheck(ctx context.Context) error                      { m.healthCalled = true; return nil }
func (m *mockClient) GetWorkflows(context.Context) ([]*workflow.Workflow, error) { return nil, nil }
func (m *mockClient) GetWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	return &workflow.Workflow{ID: id, Name: "t"}, nil
}
func (m *mockClient) CreateWorkflow(ctx context.Context, wf *workflow.Workflow) (*workflow.Workflow, error) {
	m.updateCalled = true
	return wf, nil
}
func (m *mockClient) UpdateWorkflow(ctx context.Context, id string, wf *workflow.Workflow) (*workflow.Workflow, error) {
	m.updateCalled = true
	return wf, nil
}
func (m *mockClient) DeleteWorkflow(context.Context, string) error { return nil }
func (m *mockClient) ExecuteWorkflow(context.Context, string) (*workflow.ExecutionResult, error) {
	return nil, nil
}
func (m *mockClient) GetExecution(context.Context, string) (*workflow.ExecutionResult, error) {
	return nil, nil
}
func (m *mockClient) GetExecutions(context.Context, string, string, int) ([]*workflow.ExecutionResult, error) {
	return nil, nil
}
func (m *mockClient) GetCredentials(context.Context) ([]*credentials.N8nCredential, error) {
	return nil, nil
}
func (m *mockClient) GetCredential(context.Context, string) (*credentials.N8nCredential, error) {
	return nil, nil
}
func (m *mockClient) CreateCredential(context.Context, *credentials.N8nCredential) (*credentials.N8nCredential, error) {
	return nil, nil
}
func (m *mockClient) UpdateCredential(context.Context, string, *credentials.N8nCredential) (*credentials.N8nCredential, error) {
	return nil, nil
}
func (m *mockClient) DeleteCredential(context.Context, string) error { return nil }
func (m *mockClient) GetCredentialSchema(context.Context, string) (map[string]interface{}, error) {
	return nil, nil
}

func TestRunDaemonModeCtxProcessesEvent(t *testing.T) {
	logger = logrus.New()
	demoMode = true
	tmp := t.TempDir()
	old, _ := filepath.Abs(".")
	_ = os.Chdir(tmp)
	defer os.Chdir(old)

	env := "dev"
	os.MkdirAll(filepath.Join("workflows", env), 0755)
	file := filepath.Join("workflows", env, "wf.json")
	os.WriteFile(file, []byte(`{"id":"1","name":"t"}`), 0644)

	mw := &mockWatcher{events: make(chan fsnotify.Event, 1), errs: make(chan error)}
	mc := &mockClient{}

	daemonWatcherFactory = func() (fileWatcher, error) { return mw, nil }
	daemonClientFactory = func(string) (client.Client, error) { return mc, nil }

	ctx, cancel := context.WithCancel(context.Background())
	go runDaemonModeCtx(ctx, env)
	mw.events <- fsnotify.Event{Name: file, Op: fsnotify.Write}
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	assert.True(t, mc.healthCalled)
	assert.True(t, mc.updateCalled)
}
