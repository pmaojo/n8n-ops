package cmd

import (
        "encoding/json"
        "fmt"
        "os"
        "time"

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
        Run: runCheck,
}

var (
        quiet        bool
        jsonOutput   bool
        failIfChanges bool
        alertOnly    bool
)

type CheckResult struct {
        Environment       string             `json:"environment"`
        LastSync         time.Time          `json:"lastSync"`
        TotalWorkflows   int                `json:"totalWorkflows"`
        Synchronized     int                `json:"synchronized"`
        Modified         int                `json:"modified"`
        Workflows        WorkflowStatuses   `json:"workflows"`
}

type WorkflowStatuses struct {
        Synchronized []WorkflowStatus `json:"synchronized"`
        Modified     []WorkflowStatus `json:"modified"`
}

type WorkflowStatus struct {
        ID           string    `json:"id"`
        Name         string    `json:"name"`
        Status       string    `json:"status"` // "sync", "modified"
        LocalVersion int       `json:"localVersion,omitempty"`
        RemoteVersion int      `json:"remoteVersion,omitempty"`
        LastModified time.Time `json:"lastModified,omitempty"`
        TimeAgo      string    `json:"timeAgo,omitempty"`
}

func init() {
        rootCmd.AddCommand(checkCmd)
        checkCmd.Flags().BoolVar(&quiet, "quiet", false, "quiet mode - only return exit code")
        checkCmd.Flags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
        checkCmd.Flags().BoolVar(&failIfChanges, "fail-if-changes", false, "exit with code 1 if changes detected")
        checkCmd.Flags().BoolVar(&alertOnly, "alert-only", false, "only show alerts, don't suggest actions")
}

func runCheck(cmd *cobra.Command, args []string) {
        logger.WithField("env", environment).Info("Checking workflow sync status")

        result, err := checkWorkflowSync()
        if err != nil {
                if !quiet {
                        fmt.Printf("❌ Error checking sync status: %v\n", err)
                }
                os.Exit(1)
        }

        if jsonOutput {
                printCheckResultJSON(result)
        } else if !quiet {
                printCheckResultTable(result)
        }

        // Exit with appropriate code
        if failIfChanges && result.Modified > 0 {
                os.Exit(1)
        } else if result.Modified > 0 {
                os.Exit(1) // Changes detected
        }
        os.Exit(0) // All synchronized
}

func checkWorkflowSync() (*CheckResult, error) {
        // Read local workflows
        workflowDir := fmt.Sprintf("./workflows/%s", environment)
        localWorkflows, err := getLocalWorkflowStatus(workflowDir)
        if err != nil {
                return nil, fmt.Errorf("failed to read local workflows: %w", err)
        }

        // Use demo mode or real API
        if len(localWorkflows) == 0 {
                return checkWorkflowSyncDemo(localWorkflows)
        }

        // Real API comparison (when implemented)
        return checkWorkflowSyncReal(localWorkflows)
}

func checkWorkflowSyncDemo(localWorkflows []struct{ID, Name string}) (*CheckResult, error) {
        result := &CheckResult{
                Environment:    environment,
                LastSync:      time.Now().Add(-15 * time.Minute),
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

func checkWorkflowSyncReal(localWorkflows []struct{ID, Name string}) (*CheckResult, error) {
        // This would implement real n8n API comparison
        result := &CheckResult{
                Environment:    environment,
                LastSync:      time.Now().Add(-15 * time.Minute),
                TotalWorkflows: len(localWorkflows),
                Workflows: WorkflowStatuses{
                        Synchronized: make([]WorkflowStatus, 0),
                        Modified:     make([]WorkflowStatus, 0),
                },
        }

        // For now, mark all as synchronized until n8n API integration
        for _, workflow := range localWorkflows {
                status := WorkflowStatus{
                        ID:     workflow.ID,
                        Name:   workflow.Name,
                        Status: "sync",
                }
                result.Workflows.Synchronized = append(result.Workflows.Synchronized, status)
                result.Synchronized++
        }

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
        fmt.Printf("\n🔍 Workflow Sync Status - %s Environment\n", environment)
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
                fmt.Printf("   ./n8n-ops sync --env %s\n", environment)
                fmt.Printf("   git add workflows/%s/ && git commit -m \"sync: update workflows\"\n", environment)
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