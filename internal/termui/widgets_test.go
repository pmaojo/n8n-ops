package termui

import "testing"

func TestUpdateWorkflowTable(t *testing.T) {
	tbl := NewWorkflowTable()
	data := []WorkflowStatus{{ID: "1", Name: "Test", Status: "active"}}
	UpdateWorkflowTable(tbl, data)
	if len(tbl.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(tbl.Rows))
	}
	if tbl.Rows[1][1] != "Test" {
		t.Errorf("unexpected row value: %v", tbl.Rows[1])
	}
}

func TestUpdateMetrics(t *testing.T) {
	p := NewMetricsParagraph()
	m := Metrics{Failures: 1, Issues: 2, Workflows: 3, Uptime: "1s"}
	UpdateMetrics(p, m)
	if p.Text == "" {
		t.Fatal("metrics text should not be empty")
	}
}

func TestUpdateEventList(t *testing.T) {
	l := NewEventList()
	UpdateEventList(l, []string{"one", "two"})
	if len(l.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(l.Rows))
	}
}

func TestBuildGrid(t *testing.T) {
	grid := BuildGridWithSize(NewWorkflowTable(), NewMetricsParagraph(), NewEventList(), NewSummaryGauge(), 80, 24)
	if grid == nil {
		t.Fatal("grid should not be nil")
	}
}

func TestUpdateSummaryGauge(t *testing.T) {
	g := NewSummaryGauge()
	UpdateSummaryGauge(g, 1, 2)
	if g.Percent != 50 {
		t.Fatalf("expected 50 percent, got %d", g.Percent)
	}
}
