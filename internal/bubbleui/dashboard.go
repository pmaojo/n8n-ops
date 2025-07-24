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
	"github.com/sirupsen/logrus"
)

// Dashboard implements a Bubble Tea terminal dashboard.
type Dashboard struct {
	client  client.WorkflowReader
	refresh time.Duration
	logger  logrus.FieldLogger
	model   *model
}

// Ensure Dashboard satisfies the ui.DashboardUI interface.
var _ interface{ Run(context.Context) error } = (*Dashboard)(nil)

// NewDashboard initializes a Bubble Tea dashboard.
func NewDashboard(c client.WorkflowReader, refresh time.Duration, logger logrus.FieldLogger) *Dashboard {
	return &Dashboard{client: c, refresh: refresh, logger: logger}
}

// Run starts the Bubble Tea program.
func (d *Dashboard) Run(ctx context.Context) error {
	m := newModel(ctx, d.client, d.refresh, d.logger)
	d.model = &m
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
	logger  logrus.FieldLogger
	start   time.Time

	width  int
	height int

	workflows     []WorkflowStatus
	metrics       Metrics
	events        []string
	summary       Summary
	selectedIndex int
	filterText    string
	inputMode     bool
	viewingDetails bool
	workflowDetail *wf.Workflow
}

func newModel(ctx context.Context, c client.WorkflowReader, refresh time.Duration, logger logrus.FieldLogger) model {
	return model{
		ctx:           ctx,
		client:        c,
		refresh:       refresh,
		logger:        logger,
		start:         time.Now(),
		selectedIndex: 0,
		filterText:    "",
		inputMode:     false,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd(m.refresh)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			switch msg.Type {
			case tea.KeyEnter:
				m.inputMode = false
			case tea.KeyEsc:
				m.inputMode = false
				m.filterText = ""
			case tea.KeyBackspace:
				if len(m.filterText) > 0 {
					m.filterText = m.filterText[:len(m.filterText)-1]
				}
			default:
				if msg.Type == tea.KeyRunes {
					m.filterText += string(msg.Runes)
				}
			}
			filtered := filterWorkflows(m.workflows, m.filterText)
			if m.selectedIndex >= len(filtered) && len(filtered) > 0 {
				m.selectedIndex = len(filtered) - 1
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "/":
			m.inputMode = true
			m.filterText = ""
		case "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down":
			if m.selectedIndex < len(filterWorkflows(m.workflows, m.filterText))-1 {
				m.selectedIndex++
			}
		case "enter":
			if m.viewingDetails {
				m.viewingDetails = false
				m.workflowDetail = nil
			} else {
				m.loadWorkflowDetail()
			}
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
	workflows, err := m.fetchWorkflows()
	if err != nil {
		if m.logger != nil {
			m.logger.WithError(err).Error("failed to fetch workflows")
		}
		m.workflows = nil
		m.metrics = Metrics{Uptime: time.Since(m.start).Truncate(time.Second).String()}
		m.summary = Summary{}
		m.events = []string{time.Now().Format("15:04:05") + " failed to fetch workflows"}
		return
	}
	statuses := make([]WorkflowStatus, 0, len(workflows))
	active := 0
	for _, wf := range workflows {
		status := "inactive"
		if wf.Active {
			status = "active"
			active++
		}
		statuses = append(statuses, WorkflowStatus{ID: wf.ID, Name: wf.Name, Status: status, UpdatedAt: wf.UpdatedAt})
	}
	m.workflows = statuses
	m.metrics = Metrics{
		Workflows: len(workflows),
		Uptime:    time.Since(m.start).Truncate(time.Second).String(),
	}
	m.summary = Summary{Active: active, Inactive: len(workflows) - active}
	m.events = []string{time.Now().Format("15:04:05") + " dashboard refreshed"}
}

func (m *model) fetchWorkflows() ([]*wf.Workflow, error) {
	if m.client == nil {
		return nil, nil
	}
	if m.ctx == nil {
		m.ctx = context.Background()
	}
	workflows, err := m.client.GetWorkflows(m.ctx)
	if err != nil {
		return nil, err
	}
	return workflows, nil
}

func (m *model) loadWorkflowDetail() {
	if m.client == nil || m.selectedIndex >= len(m.workflows) {
		return
	}
	if m.ctx == nil {
		m.ctx = context.Background()
	}
	wfID := m.workflows[m.selectedIndex].ID
	wf, err := m.client.GetWorkflow(m.ctx, wfID)
	if err != nil {
		return
	}
	m.workflowDetail = wf
	m.viewingDetails = true
}

func (m model) View() string {
	filtered := filterWorkflows(m.workflows, m.filterText)
	rows := workflowTableRows(filtered)
	if m.viewingDetails && m.workflowDetail != nil {
		return renderDetailView(m.workflowDetail)
	}
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
	highlight := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("2"))
	for i, r := range rows {
		status := statusStyle(r[2]).Render(r[2])
		row := fmt.Sprintf("%-8s %-30s %-10s", r[0], r[1], status)
		if i == m.selectedIndex+1 { // +1 because rows include header
			row = highlight.Render(row)
		}
		b.WriteString(row + "\n")
	}
	sections := []string{
		lipgloss.NewStyle().Bold(true).Render("Workflows"),
		b.String(),
	}

	if len(filtered) > 0 && m.selectedIndex < len(filtered) {
		wf := filtered[m.selectedIndex]
		details := fmt.Sprintf("Name: %s\nStatus: %s\nLast Updated: %s", wf.Name, wf.Status, wf.UpdatedAt.Format(time.RFC3339))
		sections = append(sections,
			lipgloss.NewStyle().Bold(true).Render("Details"),
			details,
		)
	}

	if m.inputMode || m.filterText != "" {
		sections = append(sections,
			lipgloss.NewStyle().Bold(true).Render("Filter"),
			"Filter: "+m.filterText,
		)
	}

	sections = append(sections,
		lipgloss.NewStyle().Bold(true).Render("Metrics"),
		metricsText(m.metrics),
		lipgloss.NewStyle().Bold(true).Render("Active Workflows"),
		gaugeStr,
		lipgloss.NewStyle().Bold(true).Render("Events"),
		strings.Join(m.events, "\n"),
	)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func renderDetailView(wf *wf.Workflow) string {
	highlight := lipgloss.NewStyle().Bold(true)
	var sections []string
	sections = append(sections, highlight.Render("Workflow Details"))
	sections = append(sections,
		fmt.Sprintf("ID: %s", wf.ID),
		fmt.Sprintf("Name: %s", wf.Name),
		fmt.Sprintf("Active: %t", wf.Active),
		fmt.Sprintf("Last Updated: %s", wf.UpdatedAt.Format(time.RFC3339)),
	)

	if len(wf.Nodes) > 0 {
		sections = append(sections, highlight.Render("Nodes"))
		for _, n := range wf.Nodes {
			sections = append(sections, fmt.Sprintf("- %s (%s)", n.Name, n.Type))
		}
	}

	if len(wf.Tags) > 0 {
		var tags []string
		for _, t := range wf.Tags {
			tags = append(tags, t.Name)
		}
		sections = append(sections, highlight.Render("Tags"), strings.Join(tags, ", "))
	}

	sections = append(sections, "", "Press Enter to return")
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
