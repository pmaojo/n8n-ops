package bubbleui

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/utils"
	"github.com/sirupsen/logrus"
)

// RetroDashboard renders a simple text dashboard refreshed periodically.
type RetroDashboard struct {
	client  client.WorkflowReader
	refresh time.Duration
	logger  logrus.FieldLogger
	start   time.Time
	out     io.Writer
}

// Ensure RetroDashboard satisfies the ui.DashboardUI interface.
var _ interface{ Run(context.Context) error } = (*RetroDashboard)(nil)

// NewRetroDashboard initializes a RetroDashboard writing to stdout.
func NewRetroDashboard(c client.WorkflowReader, refresh time.Duration, logger logrus.FieldLogger) *RetroDashboard {
	return &RetroDashboard{
		client:  c,
		refresh: refresh,
		logger:  logger,
		start:   time.Now(),
		out:     os.Stdout,
	}
}

// Run starts the periodic rendering loop.
func (d *RetroDashboard) Run(ctx context.Context) error {
	ticker := time.NewTicker(d.refresh)
	defer ticker.Stop()

	if err := utils.ClearTerminalScreen(); err == nil {
		d.render(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := utils.ClearTerminalScreen(); err != nil {
				if d.logger != nil {
					d.logger.WithError(err).Warn("clear screen")
				}
			}
			d.render(ctx)
		}
	}
}

func (d *RetroDashboard) render(ctx context.Context) {
	workflows, summary := d.loadWorkflows(ctx)
	metrics := Metrics{Workflows: len(workflows), Uptime: time.Since(d.start).Truncate(time.Second).String()}
	events := []string{time.Now().Format("15:04:05") + " dashboard refreshed"}

	fmt.Fprintln(d.out, lipgloss.NewStyle().Bold(true).Render("n8n-ops Retro Dashboard"))
	fmt.Fprintln(d.out)
	d.renderWorkflows(workflows)

	percent, label := summaryGauge(summary.Active, summary.Active+summary.Inactive)
	barWidth := 30
	filled := int(float64(barWidth) * float64(percent) / 100)
	gauge := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	bar := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(strings.Repeat("█", filled)) + strings.Repeat(" ", barWidth-filled)
	fmt.Fprintln(d.out, gauge.Render("["+bar+"] "+label))
	fmt.Fprintln(d.out)
	fmt.Fprintln(d.out, metricsText(metrics))
	fmt.Fprintln(d.out)
	for _, e := range events {
		fmt.Fprintln(d.out, e)
	}
}

func (d *RetroDashboard) renderWorkflows(wf []WorkflowStatus) {
	rows := workflowTableRows(wf)
	for i, r := range rows {
		status := r[2]
		if i != 0 {
			status = statusStyle(r[2]).Render(r[2])
		}
		fmt.Fprintf(d.out, "%-8s %-30s %-10s\n", r[0], r[1], status)
	}
	fmt.Fprintln(d.out)
}

func (d *RetroDashboard) loadWorkflows(ctx context.Context) ([]WorkflowStatus, Summary) {
	if d.client == nil {
		return nil, Summary{}
	}
	wfs, err := d.client.GetWorkflows(ctx)
	if err != nil {
		if d.logger != nil {
			d.logger.WithError(err).Error("failed to fetch workflows")
		}
		return nil, Summary{}
	}
	statuses := make([]WorkflowStatus, 0, len(wfs))
	active := 0
	for _, wf := range wfs {
		status := "inactive"
		if wf.Active {
			status = "active"
			active++
		}
		statuses = append(statuses, WorkflowStatus{ID: wf.ID, Name: wf.Name, Status: status, UpdatedAt: wf.UpdatedAt})
	}
	return statuses, Summary{Active: active, Inactive: len(wfs) - active}
}
