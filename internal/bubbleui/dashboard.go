package bubbleui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pmaojo/n8n-ops/internal/client"
	wf "github.com/pmaojo/n8n-ops/internal/workflow"
)

// Dashboard implements a Bubble Tea terminal dashboard.
type Dashboard struct {
	model *model
}

// Ensure Dashboard satisfies the ui.DashboardUI interface.
var _ interface{ Run(context.Context) error } = (*Dashboard)(nil)

// NewDashboard initializes a Bubble Tea dashboard.
func NewDashboard(c client.WorkflowReader, refresh time.Duration) *Dashboard {
	m := newModel(c, refresh)
	return &Dashboard{model: &m}
}

// Run starts the Bubble Tea program.
func (d *Dashboard) Run(ctx context.Context) error {
	p := tea.NewProgram(d.model, tea.WithContext(ctx))
	_, err := p.Run()
	return err
}

type tickMsg time.Time

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type model struct {
	ctx     context.Context
	client  client.WorkflowReader
	refresh time.Duration
	start   time.Time

	width  int
	height int

	workflows []WorkflowStatus
	metrics   Metrics
	events    []string
	summary   Summary
}

func newModel(c client.WorkflowReader, refresh time.Duration) model {
	return model{client: c, refresh: refresh, start: time.Now()}
}

func (m model) Init() tea.Cmd {
	return tickCmd(m.refresh)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		m.refreshData()
		return m, tickCmd(m.refresh)
	}
	return m, nil
}

func (m *model) refreshData() {
	workflows := m.fetchWorkflows()
	statuses := make([]WorkflowStatus, 0, len(workflows))
	active := 0
	for _, wf := range workflows {
		status := "inactive"
		if wf.Active {
			status = "active"
			active++
		}
		statuses = append(statuses, WorkflowStatus{ID: wf.ID, Name: wf.Name, Status: status})
	}
	m.workflows = statuses
	m.metrics = Metrics{
		Workflows: len(workflows),
		Uptime:    time.Since(m.start).Truncate(time.Second).String(),
	}
	m.summary = Summary{Active: active, Inactive: len(workflows) - active}
	m.events = []string{time.Now().Format("15:04:05") + " dashboard refreshed"}
}

func (m *model) fetchWorkflows() []*wf.Workflow {
	if m.client == nil {
		return nil
	}
	workflows, err := m.client.GetWorkflows(context.Background())
	if err != nil {
		return nil
	}
	return workflows
}

func (m model) View() string {
	rows := workflowTableRows(m.workflows)
	percent, label := summaryGauge(m.summary.Active, m.summary.Active+m.summary.Inactive)
	barWidth := m.width - len(label) - 5
	if barWidth < 0 {
		barWidth = 0
	}
	filled := int(float64(barWidth) * float64(percent) / 100)
	gauge := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	bar := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(strings.Repeat("█", filled)) + strings.Repeat(" ", barWidth-filled)
	gaugeStr := gauge.Render("[" + bar + "] " + label)

	var b strings.Builder
	for _, r := range rows {
		b.WriteString(fmt.Sprintf("%-8s %-30s %-10s\n", r[0], r[1], r[2]))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("Workflows"),
		b.String(),
		lipgloss.NewStyle().Bold(true).Render("Metrics"),
		metricsText(m.metrics),
		lipgloss.NewStyle().Bold(true).Render("Active Workflows"),
		gaugeStr,
		lipgloss.NewStyle().Bold(true).Render("Events"),
		strings.Join(m.events, "\n"),
	)
}
