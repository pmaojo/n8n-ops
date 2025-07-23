package termui

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// NewWorkflowTable constructs a table for workflow status information.
func NewWorkflowTable() *widgets.Table {
	tbl := widgets.NewTable()
	tbl.Title = "Workflows"
	tbl.RowSeparator = false
	tbl.Rows = [][]string{{"ID", "Name", "Status"}}
	tbl.TextStyle = ui.NewStyle(ui.ColorWhite)
	return tbl
}

// UpdateWorkflowTable sets workflow data on the given table.
func UpdateWorkflowTable(tbl *widgets.Table, workflows []WorkflowStatus) {
	rows := [][]string{{"ID", "Name", "Status"}}
	tbl.RowStyles = map[int]ui.Style{}
	for i, wf := range workflows {
		rows = append(rows, []string{wf.ID, wf.Name, wf.Status})
		style := ui.NewStyle(ui.ColorRed)
		if wf.Status == "active" {
			style = ui.NewStyle(ui.ColorGreen)
		}
		tbl.RowStyles[i+1] = style
	}
	tbl.Rows = rows
}

// NewMetricsParagraph creates a paragraph widget for metrics display.
func NewMetricsParagraph() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Title = "Metrics"
	p.TextStyle = ui.NewStyle(ui.ColorWhite)
	return p
}

// UpdateMetrics updates the metrics paragraph with metric values.
func UpdateMetrics(p *widgets.Paragraph, m Metrics) {
	p.Text =
		"Failures: " + itoa(m.Failures) +
			" | Issues: " + itoa(m.Issues) +
			" | Workflows: " + itoa(m.Workflows) +
			" | Uptime: " + m.Uptime
}

// NewEventList returns a list widget for log events.
func NewEventList() *widgets.List {
	l := widgets.NewList()
	l.Title = "Events"
	l.WrapText = false
	return l
}

// UpdateEventList sets the event rows of the given list.
func UpdateEventList(l *widgets.List, events []string) {
	l.Rows = events
}

// NewSummaryGauge creates a gauge indicating active workflow ratio.
func NewSummaryGauge() *widgets.Gauge {
	g := widgets.NewGauge()
	g.Title = "Active Workflows"
	g.BarColor = ui.ColorGreen
	g.LabelStyle = ui.NewStyle(ui.ColorBlack)
	return g
}

// UpdateSummaryGauge sets gauge values using active and total workflow counts.
func UpdateSummaryGauge(g *widgets.Gauge, active, total int) {
	if total == 0 {
		g.Percent = 0
		g.Label = "0/0"
		return
	}
	g.Percent = int(float64(active) / float64(total) * 100)
	g.Label = fmt.Sprintf("%d/%d active", active, total)
}

// BuildGrid arranges widgets using termui's grid system.
func BuildGrid(tbl *widgets.Table, p *widgets.Paragraph, l *widgets.List, g *widgets.Gauge) *ui.Grid {
	w, h := ui.TerminalDimensions()
	return BuildGridWithSize(tbl, p, l, g, w, h)
}

// BuildGridWithSize arranges widgets with explicit dimensions.
func BuildGridWithSize(tbl *widgets.Table, p *widgets.Paragraph, l *widgets.List, g *widgets.Gauge, width, height int) *ui.Grid {
	grid := ui.NewGrid()
	grid.SetRect(0, 0, width, height)
	grid.Set(
		ui.NewRow(0.4, tbl),
		ui.NewRow(0.2, p),
		ui.NewRow(0.2, g),
		ui.NewRow(0.2, l),
	)
	return grid
}

// itoa is a small helper to avoid strconv import in this file.
func itoa(v int) string {
	return fmt.Sprintf("%d", v)
}
