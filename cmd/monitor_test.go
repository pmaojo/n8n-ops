package cmd

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/workflow"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type fakeClient struct {
	healthCalled bool
}

func (f *fakeClient) HealthCheck(ctx context.Context) error                      { f.healthCalled = true; return nil }
func (f *fakeClient) GetWorkflows(context.Context) ([]*workflow.Workflow, error) { return nil, nil }
func (f *fakeClient) GetWorkflow(context.Context, string) (*workflow.Workflow, error) {
	return nil, nil
}
func (f *fakeClient) CreateWorkflow(context.Context, *workflow.Workflow) (*workflow.Workflow, error) {
	return nil, nil
}
func (f *fakeClient) UpdateWorkflow(context.Context, string, *workflow.Workflow) (*workflow.Workflow, error) {
	return nil, nil
}
func (f *fakeClient) DeleteWorkflow(context.Context, string) error { return nil }
func (f *fakeClient) ExecuteWorkflow(context.Context, string) (*workflow.ExecutionResult, error) {
	return nil, nil
}
func (f *fakeClient) GetExecution(context.Context, string) (*workflow.ExecutionResult, error) {
	return nil, nil
}
func (f *fakeClient) GetExecutions(context.Context, string, string, int) ([]*workflow.ExecutionResult, error) {
	return nil, nil
}
func (f *fakeClient) GetCredentials(context.Context) ([]*credentials.N8nCredential, error) {
	return nil, nil
}
func (f *fakeClient) GetCredential(context.Context, string) (*credentials.N8nCredential, error) {
	return nil, nil
}
func (f *fakeClient) CreateCredential(context.Context, *credentials.N8nCredential) (*credentials.N8nCredential, error) {
	return nil, nil
}
func (f *fakeClient) UpdateCredential(context.Context, string, *credentials.N8nCredential) (*credentials.N8nCredential, error) {
	return nil, nil
}
func (f *fakeClient) DeleteCredential(context.Context, string) error { return nil }
func (f *fakeClient) GetCredentialSchema(context.Context, string) (map[string]interface{}, error) {
	return nil, nil
}

func TestGetN8nClientDemo(t *testing.T) {
	fc := &fakeClient{}
	monitorClientFactory = func(url, apiKey string) (client.Client, error) { return fc, nil }
	c, err := getN8nClient(context.Background(), "dev", true)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.True(t, fc.healthCalled)
}

type mockDetector struct{ started bool }

func (m *mockDetector) Start(ctx context.Context) error { m.started = true; <-ctx.Done(); return nil }

func TestStartMonitoringCancels(t *testing.T) {
	det := &mockDetector{}
	logEntry := logrus.New().WithField("test", "startMonitoring")
	sig := make(chan os.Signal, 1)

	go func() {
		time.Sleep(50 * time.Millisecond)
		sig <- syscall.SIGINT
	}()

	assert.NoError(t, startMonitoring(det, logEntry, sig))
	assert.True(t, det.started)
}

func TestSetupIssueManagerDemo(t *testing.T) {
	mgr, id, err := setupIssueManager(true, "", "")
	assert.NoError(t, err)
	assert.NotNil(t, mgr)
	assert.Equal(t, "", id)
}
