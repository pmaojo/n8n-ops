package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	r.Close()
	return string(out)
}

func TestRenderWorkflowStatus(t *testing.T) {
	workflows := []WorkflowInfo{
		{ID: "1", Name: "Test", Status: "HEALTHY"},
		{ID: "2", Name: "Payment", Status: "CRITICAL"},
	}

	output := captureOutput(func() {
		renderWorkflowStatus(workflows)
	})
	if !strings.Contains(output, "Payment") || !strings.Contains(output, "CRITICAL") {
		t.Errorf("unexpected workflow output: %s", output)
	}
}

func TestRenderMetrics(t *testing.T) {
	m := TerminalMetrics{
		Failures:    5,
		Issues:      1,
		Workflows:   2,
		Uptime:      "10",
		Threshold:   1,
		Interval:    "5s",
		Environment: "DEV",
	}

	output := captureOutput(func() {
		renderMetrics(m)
	})
	if !strings.Contains(output, "DEV") || !strings.Contains(output, "5s") {
		t.Errorf("unexpected metrics output: %s", output)
	}
}
