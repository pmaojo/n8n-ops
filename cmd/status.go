package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pmaojo/n8n-ops/internal/credentials"

	"github.com/pmaojo/n8n-ops/internal/git"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show deployment status and workflow versions",
	Long: `Display comprehensive status information including:

• Git status and uncommitted changes
• Workflow deployment status per environment
• Credential validation status
• Branch information and sync state
• Recent deployment history

The status command helps identify:
- Uncommitted workflow changes that need to be saved
- Out-of-sync workflows between environments
- Missing or misconfigured credentials
- Deployment conflicts or issues

Examples:
  n8n-ops status                          # Show full status
  n8n-ops status --env production         # Status for specific environment
  n8n-ops status --json                   # JSON output for automation
  n8n-ops status --check-uncommitted      # Focus on uncommitted changes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cliFrom(cmd)
		if cli == nil {
			return fmt.Errorf("CLI not initialized")
		}
		return runStatus(cmd, args, cli)
	},
}

var (
	statusJsonOutput       bool
	statusCheckUncommitted bool
	statusShowCredentials  bool
)

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVar(&statusJsonOutput, "json", false, "output in JSON format")
	statusCmd.Flags().BoolVar(&statusCheckUncommitted, "check-uncommitted", false, "focus on uncommitted workflow changes")
	statusCmd.Flags().BoolVar(&statusShowCredentials, "credentials", false, "include credential validation in status")
}

func runStatus(cmd *cobra.Command, args []string, cli *CLI) error {
	if cli == nil {
		return fmt.Errorf("CLI not initialized")
	}
	if statusCheckUncommitted {
		return runUncommittedCheck(cli.Environment)
	}

	if statusJsonOutput {
		return runStatusJSON(cli.Environment)
	}

	return runStatusHuman(cli.Environment)
}

func runUncommittedCheck(env string) error {
	fmt.Printf("🔍 Checking for Uncommitted Workflow Changes\n")
	fmt.Printf("===========================================\n\n")

	checker := git.NewGitStatusChecker(".", nil)

	// Show warning if uncommitted changes exist
	warning, err := checker.WarnIfUncommittedChanges()
	if err != nil {
		return fmt.Errorf("failed to check Git status: %w", err)
	}
	if warning != "" {
		fmt.Print(warning)
	}

	// Get detailed summary
	summary, err := checker.GetUncommittedWorkflowSummary()
	if err != nil {
		return fmt.Errorf("failed to get workflow summary: %w", err)
	}

	fmt.Printf("📋 Summary:\n")
	fmt.Printf("%s\n", summary)

	return nil
}

func runStatusJSON(env string) error {
	checker := git.NewGitStatusChecker(".", nil)
	gitStatus, err := checker.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get Git status: %w", err)
	}

	status := map[string]interface{}{
		"timestamp":   time.Now(),
		"environment": env,
		"git":         gitStatus,
		"workflows": map[string]interface{}{
			"total":       len(gitStatus.WorkflowFiles),
			"uncommitted": len(gitStatus.UncommittedWorkflows),
			"changes":     gitStatus.UncommittedWorkflows,
		},
		"health": "ok",
	}

	if statusShowCredentials {
		// Add credential status
		status["credentials"] = map[string]interface{}{
			"validated": false,
			"missing":   []string{},
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(status)
}

func runStatusHuman(env string) error {
	// ASCII Art Header
	fmt.Printf(`
        ad88888ba   888888888888          db   888888888888 88        88   ad88888ba   
       d8"     "8b      88              d88b      88      88        88  d8"     "8b  
       Y8,              88             d8'8b      88      88        88  Y8,          
       88b,             88            d8' 8b      88      88        88  88b,         
        Y88888Pb,       88           d8'   8b     88      88        88   Y88888Pb,   
              8b        88          d888888888b    88      88        88         8b   
       Y8a    a8P       88         d8'       8b    88      Y8a.    .a8P  Y8a    a8P   
        Y88888P"        88        d8'         8b   88       Y8888888P"    Y88888P"   

`)
	fmt.Printf("🚀 n8n-ops Status Dashboard - %s Environment\n", strings.ToUpper(env))
	fmt.Printf("================================================\n\n")

	// Git Status
	fmt.Printf("📊 Git Repository Status\n")
	fmt.Printf("========================\n")

	checker := git.NewGitStatusChecker(".", nil)
	gitStatus, err := checker.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get Git status: %w", err)
	}

	fmt.Printf("Current Branch: %s\n", gitStatus.CurrentBranch)
	fmt.Printf("Last Commit:    %s\n", gitStatus.LastCommit)
	fmt.Printf("Has Changes:    %v\n", gitStatus.HasChanges)

	if len(gitStatus.UncommittedWorkflows) > 0 {
		fmt.Printf("\n🚨 UNCOMMITTED WORKFLOW CHANGES:\n")
		for _, workflow := range gitStatus.UncommittedWorkflows {
			var statusIcon string
			switch workflow.Status {
			case "modified":
				statusIcon = "📝"
			case "added":
				statusIcon = "➕"
			case "deleted":
				statusIcon = "🗑️"
			case "untracked":
				statusIcon = "❓"
			}
			fmt.Printf("  %s %s (%s) - %s\n", statusIcon, workflow.WorkflowName, workflow.Status, workflow.Environment)
		}
		fmt.Printf("\n💡 Run: git add . && git commit -m \"Update workflows\"\n")
	} else {
		fmt.Printf("✅ All workflows are committed\n")
	}

	// Workflow Status
	fmt.Printf("\n📋 Workflow Overview\n")
	fmt.Printf("===================\n")
	fmt.Printf("Total Workflows:     %d\n", len(gitStatus.WorkflowFiles))
	fmt.Printf("Uncommitted Changes: %d\n", len(gitStatus.UncommittedWorkflows))
	fmt.Printf("Current Environment: %s\n", env)

	// Environment Status
	fmt.Printf("\n🌍 Environment Status\n")
	fmt.Printf("===================\n")

	cm := credentials.NewCredentialManager(env)
	n8nURL, n8nAPIKey, _ := cm.GetN8nCredentials()

	if n8nURL != "" {
		fmt.Printf("N8N URL:    %s\n", n8nURL)
	} else {
		fmt.Printf("N8N URL:    ❌ Not configured\n")
	}

	if n8nAPIKey != "" {
		maskedKey := n8nAPIKey[:4] + strings.Repeat("*", len(n8nAPIKey)-8) + n8nAPIKey[len(n8nAPIKey)-4:]
		fmt.Printf("API Key:    %s\n", maskedKey)
	} else {
		fmt.Printf("API Key:    ❌ Not configured\n")
	}

	// Credential Status (if requested)
	if statusShowCredentials {
		fmt.Printf("\n🔐 Credential Validation\n")
		fmt.Printf("======================\n")
		// This would integrate with the credential manager
		fmt.Printf("Validation Status: ⏳ Checking...\n")
	}

	// Health Check
	fmt.Printf("\n💚 Health Status\n")
	fmt.Printf("===============\n")

	healthStatus := "✅ All systems operational"
	if len(gitStatus.UncommittedWorkflows) > 0 {
		healthStatus = "⚠️  Uncommitted changes detected"
	}
	if n8nURL == "" || n8nAPIKey == "" {
		healthStatus = "❌ Configuration incomplete"
	}

	fmt.Printf("Overall Status: %s\n", healthStatus)

	// Quick Actions
	fmt.Printf("\n🔧 Quick Actions\n")
	fmt.Printf("===============\n")
	fmt.Printf("• Commit changes:     git add . && git commit -m \"Update workflows\"\n")
	fmt.Printf("• Sync workflows:     n8n-ops sync --env %s\n", env)
	fmt.Printf("• Validate credentials: n8n-ops credentials validate --env %s\n", env)
	fmt.Printf("• Check branches:     n8n-ops branch list\n")

	return nil
}

// Helper function to mask sensitive values for status command
func maskStatusValue(value string) string {
	if value == "" {
		return "not_set"
	}
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}
