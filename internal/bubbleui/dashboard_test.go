package bubbleui

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestUpdateSelectedIndex(t *testing.T) {
	m := newModel(context.Background(), nil, time.Second)
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
	m := newModel(context.Background(), nil, time.Second)
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

func TestNewDashboardUsesRefreshInterval(t *testing.T) {
	d := NewDashboard(nil, 5*time.Second)
	if d.refresh != 5*time.Second {
		t.Fatalf("expected refresh interval 5s, got %s", d.refresh)
	}
}
