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
	for _, wf := range workflows {
		rows = append(rows, []string{wf.ID, wf.Name, wf.Status})
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

// BuildGrid arranges widgets using termui's grid system.
func BuildGrid(tbl *widgets.Table, p *widgets.Paragraph, l *widgets.List) *ui.Grid {
	w, h := ui.TerminalDimensions()
	return BuildGridWithSize(tbl, p, l, w, h)
}

// BuildGridWithSize arranges widgets with explicit dimensions.
func BuildGridWithSize(tbl *widgets.Table, p *widgets.Paragraph, l *widgets.List, width, height int) *ui.Grid {
	grid := ui.NewGrid()
	grid.SetRect(0, 0, width, height)
	grid.Set(
		ui.NewRow(0.5, tbl),
		ui.NewRow(0.25, p),
		ui.NewRow(0.25, l),
	)
	return grid
}

// itoa is a small helper to avoid strconv import in this file.
func itoa(v int) string {
	return fmt.Sprintf("%d", v)
}
