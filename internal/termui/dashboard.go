package termui

import (
	"context"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/pmaojo/n8n-ops/internal/client"
	wf "github.com/pmaojo/n8n-ops/internal/workflow"
)

// RunTermuiDashboard starts the interactive dashboard.
// It blocks until the context is canceled or the user presses 'q' or Ctrl+C.
func RunTermuiDashboard(ctx context.Context, c client.Client) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	tbl := NewWorkflowTable()
	metrics := NewMetricsParagraph()
	events := NewEventList()
	grid := BuildGrid(tbl, metrics, events)
	ui.Render(grid)

	start := time.Now()
	ticker := time.NewTicker(3 * time.Second)
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
			updateData(ctx, c, tbl, metrics, events, start)
			ui.Render(grid)
		}
	}
}

func updateData(ctx context.Context, c client.Client, tbl *widgets.Table, p *widgets.Paragraph, l *widgets.List, start time.Time) {
	var workflows []*wf.Workflow
	if c != nil {
		if w, err := c.GetWorkflows(ctx); err == nil {
			workflows = w
		}
	}
	statuses := make([]WorkflowStatus, 0, len(workflows))
	for _, wf := range workflows {
		status := "inactive"
		if wf.Active {
			status = "active"
		}
		statuses = append(statuses, WorkflowStatus{ID: wf.ID, Name: wf.Name, Status: status})
	}
	UpdateWorkflowTable(tbl, statuses)

	m := Metrics{
		Workflows: len(workflows),
		Uptime:    time.Since(start).Truncate(time.Second).String(),
	}
	UpdateMetrics(p, m)

	UpdateEventList(l, []string{time.Now().Format("15:04:05") + " dashboard refreshed"})
}
