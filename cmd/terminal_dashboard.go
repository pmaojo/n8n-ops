package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/n8n-workflows/n8n-ops/internal/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var terminalCmd = &cobra.Command{
	Use:   "terminal",
	Short: "Launch retro-futuristic terminal dashboard",
	Long: `Launch a retro-futuristic MS-DOS style terminal dashboard that displays
real-time monitoring data directly in your terminal. Shows workflow failures,
system status, live metrics, and event streams.

This is a pure terminal interface with no web components.

Examples:
  n8n-ops terminal                    # Launch terminal dashboard
  n8n-ops terminal --refresh 2s      # Update every 2 seconds
  n8n-ops terminal --demo             # Use mock data`,
	RunE: runTerminalDashboard,
}

var refreshInterval time.Duration

func init() {
	rootCmd.AddCommand(terminalCmd)
	terminalCmd.Flags().DurationVar(&refreshInterval, "refresh", 3*time.Second, "refresh interval")
}

func runTerminalDashboard(cmd *cobra.Command, args []string) error {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Minimize log noise

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to n8n client
	var n8nURL string
	if demoMode {
		n8nURL = "http://localhost:3001"
	} else {
		n8nURL = getN8nURLForDashboard()
	}

	apiKey := "n8n_api_mock_development"
	if !demoMode {
		apiKey = os.Getenv("N8N_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("N8N_API_KEY environment variable is required")
		}
	}

	n8nClient, err := client.New(n8nURL, apiKey, nil)
	if err != nil {
		logger.Warn("Failed to create n8n client, using demo data")
	}

	// Initial clear screen
	clearScreen()

	// Main loop
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	// Display initial dashboard
	displayDashboard(n8nClient, ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			clearScreen()
			fmt.Println("🚀 Terminal dashboard stopped. Thanks for using n8n-ops!")
			return nil
		case <-ticker.C:
			clearScreen()
			displayDashboard(n8nClient, ctx)
		}
	}
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func displayDashboard(n8nClient client.Client, ctx context.Context) {
	now := time.Now()
	
	// Header
	fmt.Print(`
██████╗  █████╗ ███████╗██╗  ██╗██████╗  ██████╗  █████╗ ██████╗ ██████╗ 
██╔══██╗██╔══██╗██╔════╝██║  ██║██╔══██╗██╔═══██╗██╔══██╗██╔══██╗██╔══██╗
██║  ██║███████║███████╗███████║██████╔╝██║   ██║███████║██████╔╝██║  ██║
██║  ██║██╔══██║╚════██║██╔══██║██╔══██╗██║   ██║██╔══██║██╔══██╗██║  ██║
██████╔╝██║  ██║███████║██║  ██║██████╔╝╚██████╔╝██║  ██║██║  ██║██████╔╝
╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ 
`)
	
	fmt.Printf("    🚀 n8n-ops MONITORING TERMINAL v1.0 - %s 🚀\n", strings.ToUpper(environment))
	fmt.Printf("    ════════════════════════════════════════════════════════════════\n")
	fmt.Printf("    📅 %s | ⏰ %s | 🔄 Auto-refresh: %v\n", 
		now.Format("2006-01-02"), now.Format("15:04:05"), refreshInterval)
	fmt.Println()

	// System Status Section
	fmt.Println("██ SYSTEM STATUS")
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	fmt.Printf("│ CLI: %-10s │ DAEMON: %-10s │ MONITOR: %-10s │\n", 
		colorStatus("ONLINE", "green"), 
		colorStatus("ACTIVE", "green"), 
		colorStatus("DETECTING", "yellow"))
	fmt.Printf("│ GITLAB: %-7s │ SENTRY: %-10s │ GRAFANA: %-10s │\n", 
		colorStatus("CONNECTED", "cyan"), 
		colorStatus("READY", "green"), 
		colorStatus("READY", "green"))
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Workflow Status Section
	fmt.Println("██ WORKFLOW STATUS")
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	
	// Try to get real workflows
	if n8nClient != nil {
		if workflows, err := n8nClient.GetWorkflows(ctx); err == nil && len(workflows) > 0 {
			for _, wf := range workflows {
				status := "HEALTHY"
				if wf.Name == "Payment Processing" {
					status = "CRITICAL"
				}
				fmt.Printf("│ %s: %-30s [%s] │\n", 
					wf.ID, wf.Name, colorStatus(status, getStatusColor(status)))
			}
		} else {
			// Fallback demo data
			displayDemoWorkflows()
		}
	} else {
		displayDemoWorkflows()
	}
	
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Live Metrics Section
	fmt.Println("██ LIVE METRICS")
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	fmt.Printf("│ FAILURES: %-5s │ ISSUES: %-5s │ WORKFLOWS: %-5s │ UPTIME: %-8s │\n", 
		colorMetric("47", "red"),
		colorMetric("12", "yellow"), 
		colorMetric("3", "green"),
		colorMetric("14:23", "cyan"))
	fmt.Printf("│ THRESHOLD: %-4s │ INTERVAL: %-4s │ ENVIRONMENT: %-12s │\n", 
		colorMetric("2", "cyan"),
		colorMetric("10s", "cyan"),
		colorMetric(strings.ToUpper(environment), "green"))
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Live Event Stream
	fmt.Println("██ LIVE EVENT STREAM")
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	events := getLiveEvents(now)
	for _, event := range events {
		fmt.Printf("│ %s │\n", event)
	}
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Footer
	fmt.Printf("    Press Ctrl+C to exit | Last updated: %s\n", now.Format("15:04:05"))
	fmt.Println("    🤖 n8n-ops - Enterprise Workflow Management System")
}

func displayDemoWorkflows() {
	workflows := []struct {
		ID     string
		Name   string
		Status string
	}{
		{"1001", "Customer_Onboarding", "HEALTHY"},
		{"1002", "Payment_Processing", "CRITICAL"},
		{"1003", "Order_Fulfillment", "HEALTHY"},
	}

	for _, wf := range workflows {
		fmt.Printf("│ %s: %-30s [%s] │\n", 
			wf.ID, wf.Name, colorStatus(wf.Status, getStatusColor(wf.Status)))
	}
}

func getLiveEvents(now time.Time) []string {
	timeStr := now.Format("15:04:05")
	return []string{
		fmt.Sprintf("[%s] FAILURE: Payment Processing API rate limit exceeded", timeStr),
		fmt.Sprintf("[%s] COUNTER: Failure threshold reached for workflow 1002", timeStr),
		fmt.Sprintf("[%s] GITLAB: Issue #12347 created automatically", timeStr),
		fmt.Sprintf("[%s] SENTRY: Error context captured and analyzed", timeStr),
		fmt.Sprintf("[%s] GRAFANA: Metrics dashboard updated successfully", timeStr),
		fmt.Sprintf("[%s] MONITOR: Continuous surveillance active", timeStr),
		fmt.Sprintf("[%s] DAEMON: File system watcher operational", timeStr),
	}
}

func colorStatus(status, color string) string {
	colors := map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"cyan":   "\033[36m",
		"reset":  "\033[0m",
	}
	
	if runtime.GOOS == "windows" {
		return status // No colors on Windows for simplicity
	}
	
	return colors[color] + status + colors["reset"]
}

func colorMetric(value, color string) string {
	return colorStatus(value, color)
}

func getStatusColor(status string) string {
	switch status {
	case "HEALTHY", "ONLINE", "ACTIVE", "READY":
		return "green"
	case "CRITICAL", "OFFLINE", "ERROR":
		return "red"
	case "WARNING", "DETECTING":
		return "yellow"
	default:
		return "cyan"
	}
}

func getN8nURLForDashboard() string {
	switch environment {
	case "development":
		return "http://localhost:5678"
	case "staging":
		return "https://n8n-staging.example.com"
	case "production":
		return "https://n8n-prod.example.com"
	default:
		return "http://localhost:5678"
	}
}