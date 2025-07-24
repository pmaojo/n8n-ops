package bubbleui

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	wf "github.com/pmaojo/n8n-ops/internal/workflow"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestUpdateSelectedIndex(t *testing.T) {
	m := newModel(context.Background(), nil, time.Second, nil)
	m.workflows = []WorkflowStatus{
		{ID: "1", Name: "A", Status: "active"},
		{ID: "2", Name: "B", Status: "inactive"},
	}

	mdl, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m2 := mdl.(model)
	if m2.selectedIndex != 1 {
		t.Fatalf("expected index 1, got %d", m2.selectedIndex)
	}

	mdl, _ = m2.Update(tea.KeyMsg{Type: tea.KeyDown})
	m3 := mdl.(model)
	if m3.selectedIndex != 1 {
		t.Fatalf("index should stay at last element, got %d", m3.selectedIndex)
	}

	mdl, _ = m3.Update(tea.KeyMsg{Type: tea.KeyUp})
	m4 := mdl.(model)
	if m4.selectedIndex != 0 {
		t.Fatalf("expected index 0 after up, got %d", m4.selectedIndex)
	}

	mdl, _ = m4.Update(tea.KeyMsg{Type: tea.KeyUp})
	m5 := mdl.(model)
	if m5.selectedIndex != 0 {
		t.Fatalf("index should not go below 0, got %d", m5.selectedIndex)
	}
}

func TestViewHighlightsSelectedRow(t *testing.T) {
	m := newModel(context.Background(), nil, time.Second, nil)
	m.workflows = []WorkflowStatus{
		{ID: "1", Name: "A", Status: "active"},
		{ID: "2", Name: "B", Status: "inactive"},
	}
	m.selectedIndex = 1

	view := m.View()
	highlight := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("2"))
	row := fmt.Sprintf("%-8s %-30s %-10s", "2", "B", "inactive")
	expected := highlight.Render(row)

	if !strings.Contains(view, expected) {
		t.Errorf("highlighted row not found in view")
	}
}

type mockWorkflowReader struct {
	GetWorkflowsFunc func(context.Context) ([]*wf.Workflow, error)
}

func (m *mockWorkflowReader) GetWorkflows(ctx context.Context) ([]*wf.Workflow, error) {
	if m.GetWorkflowsFunc != nil {
		return m.GetWorkflowsFunc(ctx)
	}
	return nil, nil
}

func (m *mockWorkflowReader) GetWorkflow(ctx context.Context, id string) (*wf.Workflow, error) {
	return nil, nil
}

func (m *mockWorkflowReader) HealthCheck(ctx context.Context) error {
	return nil
}

func TestRefreshDataFailureLogs(t *testing.T) {
	logger := logrus.New()
	hook := test.NewLocal(logger)
	reader := &mockWorkflowReader{
		GetWorkflowsFunc: func(ctx context.Context) ([]*wf.Workflow, error) {
			return nil, fmt.Errorf("boom")
		},
	}
	m := newModel(context.Background(), reader, time.Second, logger)
	m.refreshData()

	if len(m.events) == 0 || !strings.Contains(m.events[0], "failed") {
		t.Fatalf("expected failure event, got %v", m.events)
	}
	if len(hook.AllEntries()) == 0 {
		t.Fatal("expected log entry for failure")
	}
	if hook.LastEntry().Message != "failed to fetch workflows" {
		t.Errorf("unexpected log message: %s", hook.LastEntry().Message)
	}
}
