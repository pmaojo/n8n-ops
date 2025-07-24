package bubbleui

import (
	"bytes"
	"context"
	"testing"
	"time"

	wf "github.com/pmaojo/n8n-ops/internal/workflow"
)

type stubRetroClient struct{}

func (stubRetroClient) GetWorkflows(ctx context.Context) ([]*wf.Workflow, error) {
	return []*wf.Workflow{{ID: "1", Name: "A", Active: true, UpdatedAt: time.Now()}}, nil
}
func (stubRetroClient) GetWorkflow(ctx context.Context, id string) (*wf.Workflow, error) {
	return nil, nil
}
func (stubRetroClient) HealthCheck(ctx context.Context) error { return nil }

func TestNewRetroDashboardUsesRefreshInterval(t *testing.T) {
	d := NewRetroDashboard(nil, 5*time.Second, nil)
	if d.refresh != 5*time.Second {
		t.Fatalf("expected 5s refresh, got %s", d.refresh)
	}
}

func TestRetroRenderOutputsWorkflow(t *testing.T) {
	buf := &bytes.Buffer{}
	d := NewRetroDashboard(stubRetroClient{}, time.Second, nil)
	d.out = buf
	d.render(context.Background())
	if !bytes.Contains(buf.Bytes(), []byte("A")) {
		t.Fatalf("expected workflow name in output: %s", buf.String())
	}
}
