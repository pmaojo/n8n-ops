package bubbleui

import "testing"

func TestWorkflowTableRows(t *testing.T) {
	rows := workflowTableRows([]WorkflowStatus{{ID: "1", Name: "Test", Status: "active"}})
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[1][1] != "Test" {
		t.Errorf("unexpected row value: %v", rows[1])
	}
}

func TestMetricsText(t *testing.T) {
	txt := metricsText(Metrics{Failures: 1, Issues: 2, Workflows: 3, Uptime: "1s"})
	if txt == "" {
		t.Fatal("metrics text should not be empty")
	}
}

func TestSummaryGauge(t *testing.T) {
	p, _ := summaryGauge(1, 2)
	if p != 50 {
		t.Fatalf("expected 50 percent, got %d", p)
	}
}
