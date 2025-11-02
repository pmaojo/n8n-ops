package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/cliutils"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for workflows that need synchronization",
	Long: `Check if there are workflows in n8n that have been modified but not yet
synchronized to the local filesystem. This helps detect changes made in the
n8n web interface that haven't been committed to Git.

Examples:
  n8n-ops check --env development    # Check development environment
  n8n-ops check --env production --json  # JSON output for scripting
  n8n-ops check --env staging --quiet    # Silent mode (exit code only)`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := cliFrom(cmd)
		exitCode, err := runCheck(cmd, args, cli)
		if err != nil && !quiet {
			fmt.Printf("❌ %v\n", err)
		}
		os.Exit(exitCode)
	},
}

var (
	quiet         bool
	jsonOutput    bool
	failIfChanges bool
	alertOnly     bool
)

type CheckResult struct {
	Environment    string           `json:"environment"`
	LastSync       time.Time        `json:"lastSync"`
	TotalWorkflows int              `json:"totalWorkflows"`
	Synchronized   int              `json:"synchronized"`
	Modified       int              `json:"modified"`
	Workflows      WorkflowStatuses `json:"workflows"`
}

type WorkflowStatuses struct {
	Synchronized []WorkflowStatus `json:"synchronized"`
	Modified     []WorkflowStatus `json:"modified"`
}

type WorkflowStatus struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"` // "sync", "modified"
	LocalVersion  int       `json:"localVersion,omitempty"`
	RemoteVersion int       `json:"remoteVersion,omitempty"`
	LastModified  time.Time `json:"lastModified,omitempty"`
	TimeAgo       string    `json:"timeAgo,omitempty"`
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVar(&quiet, "quiet", false, "quiet mode - only return exit code")
	checkCmd.Flags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	checkCmd.Flags().BoolVar(&failIfChanges, "fail-if-changes", false, "exit with code 1 if changes detected")
	checkCmd.Flags().BoolVar(&alertOnly, "alert-only", false, "only show alerts, don't suggest actions")
}

func runCheck(cmd *cobra.Command, args []string, cli *CLI) (int, error) {
	if cli == nil {
		return 1, fmt.Errorf("CLI not initialized")
	}
	cli.Logger.WithField("env", cli.Environment).Info("Checking workflow sync status")

	n8nClient, _, err := cliutils.SetupClient(cli.Environment, demoMode, cli.Logger)
	if err != nil {
		return 1, fmt.Errorf("error configuring n8n client: %w", err)
	}

	result, err := checkWorkflowSync(cli.Environment, n8nClient)
	if err != nil {
		return 1, fmt.Errorf("error checking sync status: %w", err)
	}

	if jsonOutput {
		printCheckResultJSON(result)
	} else if !quiet {
		printCheckResultTable(result)
	}

	if result.Modified > 0 {
		return 1, nil
	}

	return 0, nil
}

func checkWorkflowSync(env string, workflowClient client.WorkflowReader) (*CheckResult, error) {
	// Read local workflows
	workflowDir := fmt.Sprintf("./workflows/%s", env)
	localWorkflows, err := getLocalWorkflows(workflowDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read local workflows: %w", err)
	}

	if demoMode {
		return checkWorkflowSyncDemo(env, localWorkflows)
	}

	if workflowClient == nil {
		return nil, fmt.Errorf("workflow client is not configured")
	}

	return checkWorkflowSyncReal(env, localWorkflows, workflowClient)
}

// WorkflowData represents a workflow from the filesystem
type WorkflowData struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	VersionID int       `json:"versionId"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// getLocalWorkflows reads workflow files from the specified directory
func getLocalWorkflows(workflowDir string) ([]WorkflowData, error) {
	var workflows []WorkflowData

	// Check if directory exists
	if _, err := os.Stat(workflowDir); os.IsNotExist(err) {
		return workflows, nil // Return empty slice if directory doesn't exist
	}

	// Read all JSON files in the workflow directory
	files, err := os.ReadDir(workflowDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(workflowDir, file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue // Skip files that can't be read
			}

			var workflow WorkflowData
			if err := json.Unmarshal(data, &workflow); err != nil {
				continue // Skip files that aren't valid JSON
			}

			workflows = append(workflows, workflow)
		}
	}

	return workflows, nil
}

func checkWorkflowSyncDemo(env string, localWorkflows []WorkflowData) (*CheckResult, error) {
	result := &CheckResult{
		Environment:    env,
		LastSync:       time.Now().Add(-15 * time.Minute),
		TotalWorkflows: 3, // Demo has 3 workflows
		Workflows: WorkflowStatuses{
			Synchronized: make([]WorkflowStatus, 0),
			Modified:     make([]WorkflowStatus, 0),
		},
	}

	// Demo workflows status
	demoWorkflows := []WorkflowStatus{
		{
			ID:            "1001",
			Name:          "Customer Onboarding",
			Status:        "modified",
			LocalVersion:  14,
			RemoteVersion: 15,
			LastModified:  time.Now().Add(-30 * time.Minute),
			TimeAgo:       "30 minutes ago",
		},
		{
			ID:     "1002",
			Name:   "Payment Processing",
			Status: "sync",
		},
		{
			ID:            "1003",
			Name:          "Order Fulfillment",
			Status:        "modified",
			LocalVersion:  7,
			RemoteVersion: 8,
			LastModified:  time.Now().Add(-5 * time.Minute),
			TimeAgo:       "5 minutes ago",
		},
	}

	for _, workflow := range demoWorkflows {
		if workflow.Status == "modified" {
			result.Workflows.Modified = append(result.Workflows.Modified, workflow)
			result.Modified++
		} else {
			result.Workflows.Synchronized = append(result.Workflows.Synchronized, workflow)
			result.Synchronized++
		}
	}

	return result, nil
}

func checkWorkflowSyncReal(env string, localWorkflows []WorkflowData, workflowClient client.WorkflowReader) (*CheckResult, error) {
	result := &CheckResult{
		Environment: env,
		LastSync:    time.Now(),
		Workflows: WorkflowStatuses{
			Synchronized: make([]WorkflowStatus, 0),
			Modified:     make([]WorkflowStatus, 0),
		},
	}

	if workflowClient == nil {
		return result, fmt.Errorf("workflow client is not configured")
	}

	remoteWorkflows, err := workflowClient.GetWorkflows(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fetch remote workflows: %w", err)
	}

	localByID := make(map[string]WorkflowData)
	for _, wf := range localWorkflows {
		if wf.ID == "" {
			continue
		}
		localByID[wf.ID] = wf
	}

	seen := make(map[string]struct{})

	for _, remote := range remoteWorkflows {
		if remote == nil || remote.ID == "" {
			continue
		}
		if _, processed := seen[remote.ID]; processed {
			continue
		}
		seen[remote.ID] = struct{}{}

		status := WorkflowStatus{
			ID:            remote.ID,
			Name:          remote.Name,
			RemoteVersion: remote.VersionId,
		}
		if !remote.UpdatedAt.IsZero() {
			status.LastModified = remote.UpdatedAt
			status.TimeAgo = formatTimeAgo(remote.UpdatedAt)
		}

		if local, ok := localByID[remote.ID]; ok {
			status.LocalVersion = local.VersionID

			if local.VersionID == remote.VersionId {
				status.Status = "sync"
				result.Workflows.Synchronized = append(result.Workflows.Synchronized, status)
				result.Synchronized++
			} else {
				status.Status = "modified"
				result.Workflows.Modified = append(result.Workflows.Modified, status)
				result.Modified++
			}
			continue
		}

		status.Status = "modified"
		result.Workflows.Modified = append(result.Workflows.Modified, status)
		result.Modified++
	}

	for id, local := range localByID {
		if _, ok := seen[id]; ok {
			continue
		}

		status := WorkflowStatus{
			ID:            local.ID,
			Name:          local.Name,
			Status:        "modified",
			LocalVersion:  local.VersionID,
			RemoteVersion: 0,
		}
		if !local.UpdatedAt.IsZero() {
			status.LastModified = local.UpdatedAt
			status.TimeAgo = formatTimeAgo(local.UpdatedAt)
		}

		result.Workflows.Modified = append(result.Workflows.Modified, status)
		result.Modified++
		seen[id] = struct{}{}
	}

	result.TotalWorkflows = len(seen)

	return result, nil
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	}
	return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
}

func printCheckResultTable(result *CheckResult) {
	fmt.Printf("\n🔍 Workflow Sync Status - %s Environment\n", result.Environment)
	fmt.Printf("Last sync: %s\n\n", result.LastSync.Format("2006-01-02 15:04:05"))

	if result.Modified == 0 {
		fmt.Printf("✅ All workflows are synchronized (%d/%d)\n",
			result.Synchronized, result.TotalWorkflows)
		fmt.Printf("No action needed\n")
		return
	}

	fmt.Printf("📊 Sync Status Summary:\n")
	fmt.Printf("   ✅ Synchronized: %d workflows\n", result.Synchronized)
	fmt.Printf("   ⚠️  Modified: %d workflows\n", result.Modified)
	fmt.Printf("   📝 Total: %d workflows\n\n", result.TotalWorkflows)

	if len(result.Workflows.Modified) > 0 {
		fmt.Printf("Workflows Modified in n8n:\n")
		for _, wf := range result.Workflows.Modified {
			fmt.Printf("   ⚠️  %s (v%d → v%d) - %s\n",
				wf.Name, wf.LocalVersion, wf.RemoteVersion, wf.TimeAgo)
		}
		fmt.Println()
	}

	if !alertOnly {
		fmt.Printf("💡 To sync changes:\n")
		fmt.Printf("   ./n8n-ops sync --env %s\n", result.Environment)
		fmt.Printf("   git add workflows/%s/ && git commit -m \"sync: update workflows\"\n", result.Environment)
	}
}

func printCheckResultJSON(result *CheckResult) {
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}
	fmt.Println(string(output))
}
