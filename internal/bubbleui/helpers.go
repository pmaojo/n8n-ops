package bubbleui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func workflowTableRows(workflows []WorkflowStatus) [][]string {
	rows := [][]string{{"ID", "Name", "Status"}}
	for _, wf := range workflows {
		rows = append(rows, []string{wf.ID, wf.Name, wf.Status})
	}
	return rows
}

func metricsText(m Metrics) string {
	return fmt.Sprintf("Failures: %d | Issues: %d | Workflows: %d | Uptime: %s", m.Failures, m.Issues, m.Workflows, m.Uptime)
}

func summaryGauge(active, total int) (percent int, label string) {
	if total == 0 {
		return 0, "0/0"
	}
	percent = int(float64(active) / float64(total) * 100)
	label = fmt.Sprintf("%d/%d active", active, total)
	return
}

func statusColor(status string) lipgloss.Color {
	switch strings.ToLower(status) {
	case "active", "healthy", "online", "ready":
		return lipgloss.Color("2")
	case "inactive", "critical", "offline", "error":
		return lipgloss.Color("1")
	case "warning", "detecting":
		return lipgloss.Color("3")
	default:
		return lipgloss.Color("7")
	}
}

func statusStyle(status string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(statusColor(status))
}
