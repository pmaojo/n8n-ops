package termui

import (
	"context"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/pmaojo/n8n-ops/internal/client"
	wf "github.com/pmaojo/n8n-ops/internal/workflow"
)

// Dashboard coordinates widgets and updates for the terminal UI.
type Dashboard struct {
	client  client.WorkflowReader
	refresh time.Duration

	table   *widgets.Table
	metrics *widgets.Paragraph
	events  *widgets.List
	gauge   *widgets.Gauge
	start   time.Time
}

// NewDashboard returns an initialized Dashboard.
func NewDashboard(c client.WorkflowReader, refresh time.Duration) *Dashboard {
	return &Dashboard{
		client:  c,
		refresh: refresh,
		table:   NewWorkflowTable(),
		metrics: NewMetricsParagraph(),
		events:  NewEventList(),
		gauge:   NewSummaryGauge(),
		start:   time.Now(),
	}
}

// Run launches the dashboard event loop.
func (d *Dashboard) Run(ctx context.Context) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	grid := BuildGrid(d.table, d.metrics, d.events, d.gauge)
	ui.Render(grid)

	ticker := time.NewTicker(d.refresh)
	defer ticker.Stop()

	uiEvents := ui.PollEvents()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent && (e.ID == "q" || e.ID == "<C-c>") {
				return nil
			}
		case <-ticker.C:
			d.update(ctx)
			ui.Render(grid)
		}
	}
}

func (d *Dashboard) update(ctx context.Context) {
	workflows := d.fetchWorkflows(ctx)

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
	UpdateWorkflowTable(d.table, statuses)

	UpdateSummaryGauge(d.gauge, active, len(workflows))

	m := Metrics{
		Workflows: len(workflows),
		Uptime:    time.Since(d.start).Truncate(time.Second).String(),
	}
	UpdateMetrics(d.metrics, m)

	UpdateEventList(d.events, []string{time.Now().Format("15:04:05") + " dashboard refreshed"})
}

func (d *Dashboard) fetchWorkflows(ctx context.Context) []*wf.Workflow {
	if d.client == nil {
		return nil
	}
	workflows, err := d.client.GetWorkflows(ctx)
	if err != nil {
		return nil
	}
	return workflows
}

// RunTermuiDashboard starts the interactive dashboard with default settings.
func RunTermuiDashboard(ctx context.Context, c client.Client) error {
	db := NewDashboard(c, 3*time.Second)
	return db.Run(ctx)
}
