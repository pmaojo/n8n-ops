package bubbleui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/pmaojo/n8n-ops/internal/utils"
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

func statusColor(t Theme, status string) lipgloss.Color {
	switch utils.StatusColor(status) {
	case "green":
		return t.SuccessColor
	case "red":
		return t.ErrorColor
	case "yellow":
		return t.WarningColor
	default:
		return t.NeutralColor
	}
}

func statusStyle(t Theme, status string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(statusColor(t, status))
}
