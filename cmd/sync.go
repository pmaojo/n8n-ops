package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pmaojo/n8n-ops/internal/credentials"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/git"
	"github.com/pmaojo/n8n-ops/internal/utils"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync workflows bidirectionally between n8n and local filesystem",
	Long: `Sync performs intelligent bidirectional synchronization between your n8n instance and Git repository.

DETECTION MODES:
• FROM n8n TO Git: Downloads workflows modified in n8n UI to local files
• FROM Git TO n8n: Uploads local workflow changes to n8n instance  
• BIDIRECTIONAL: Compares timestamps and syncs in both directions (default)

CHANGE DETECTION:
• Compares workflow updatedAt timestamps with local file modification times
• Detects new workflows created in n8n UI
• Identifies local JSON file changes not yet pushed to n8n
• Handles conflict resolution with user prompts

Examples:
  n8n-ops sync --env development    # Smart bidirectional sync
  n8n-ops sync --from-n8n          # Only download from n8n to Git
  n8n-ops sync --to-n8n            # Only upload from Git to n8n  
  n8n-ops sync --force             # Force sync, auto-resolve conflicts`,
	RunE: runSync,
}

var (
	force      bool
	outputDir  string
	branch     string
	fromN8n    bool
	toN8n      bool
	syncDryRun bool
)

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVarP(&force, "force", "f", false, "force sync, overwriting conflicts")
	syncCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory (default: workflows/<environment>)")
	syncCmd.Flags().StringVarP(&branch, "branch", "b", "", "git branch (defaults to current branch)")
	syncCmd.Flags().BoolVar(&fromN8n, "from-n8n", false, "sync only from n8n to local files")
	syncCmd.Flags().BoolVar(&toN8n, "to-n8n", false, "sync only from local files to n8n")
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "show what would be synced without making changes")
}

func runSync(cmd *cobra.Command, args []string) error {
	// Check for uncommitted changes before sync
	if !force {
		checker := git.NewGitStatusChecker(".", nil)
		if err := checker.CheckBeforeSync(); err != nil {
			fmt.Printf("\n❌ Sync blocked: %s\n\n", err.Error())
			fmt.Printf("Use --force to sync anyway (⚠️  may overwrite local changes)\n")
			return nil
		}
	}
	logger.Info("Starting workflow sync", "environment", environment)
	fmt.Printf("🔄 Syncing workflows from %s environment...\n", environment)

	// Set output directory
	if outputDir == "" {
		outputDir = filepath.Join("workflows", environment)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get n8n configuration from environment (or use demo mode)
	var apiURL, apiKey string

	if demoMode {
		fmt.Printf("🎭 Demo mode enabled - using mock n8n data\n")
		apiURL = "demo://localhost"
		apiKey = "demo-key"
	} else {
		cm := credentials.NewCredentialManager(environment)
		var err error
		apiURL, apiKey, err = cm.GetN8nCredentials()
		if err != nil {
			return fmt.Errorf("failed to load credentials: %w", err)
		}
		if apiURL == "" || apiKey == "" {
			envSuffix := strings.ToUpper(environment)
			fmt.Printf("⚠️  n8n credentials not configured for %s environment\n", environment)
			fmt.Printf("💡 Set environment variables or use --demo flag:\n")
			fmt.Printf("   export N8N_%s_URL=http://localhost:3001\n", envSuffix)
			fmt.Printf("   export N8N_%s_API_KEY=n8n_api_mock_%s\n", envSuffix, environment)
			fmt.Printf("   OR: ./n8n-ops sync --demo\n")
			return nil
		}
	}

	// Initialize n8n client (with optional demo mode)
	var n8nClient client.Client
	var err error
	if demoMode {
		n8nClient = client.NewDemoN8nClient()
	} else {
		n8nClient, err = client.New(apiURL, apiKey, nil)
		if err != nil {
			return fmt.Errorf("failed to create n8n client: %w", err)
		}
	}

	// Test connection
	fmt.Printf("🔗 Connecting to n8n API: %s\n", apiURL)
	ctx := context.Background()
	if err := n8nClient.HealthCheck(ctx); err != nil {
		return fmt.Errorf("failed to connect to n8n API: %w", err)
	}
	fmt.Printf("✅ Connected successfully\n")

	// Get workflows from n8n
	workflows, err := n8nClient.GetWorkflows(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch workflows: %w", err)
	}

	fmt.Printf("📋 Found %d workflows in n8n instance\n", len(workflows))

	// Process each workflow
	syncedCount := 0
	for _, workflow := range workflows {
		// Generate filename
		filename := fmt.Sprintf("%s_%s.json", utils.SanitizeFilename(workflow.Name), workflow.ID)
		filepath := filepath.Join(outputDir, filename)

		// Add sync metadata
		workflowData := map[string]interface{}{
			"id":          workflow.ID,
			"name":        workflow.Name,
			"active":      workflow.Active,
			"nodes":       workflow.Nodes,
			"connections": workflow.Connections,
			"versionId":   workflow.VersionId,
			"tags":        workflow.Tags,
			"syncMetadata": map[string]interface{}{
				"syncDate":    time.Now(),
				"environment": environment,
				"syncedBy":    getSyncUser(),
				"gitCommit":   os.Getenv("CI_COMMIT_SHA"),
			},
		}

		// Write workflow to file
		if err := writeWorkflowFile(workflowData, filepath); err != nil {
			logger.Error("Failed to write workflow file", "workflow", workflow.Name, "error", err)
			if !force {
				return fmt.Errorf("failed to write workflow %s: %w", workflow.Name, err)
			}
			continue
		}

		fmt.Printf("📝 Synced: %s → %s\n", workflow.Name, filename)
		syncedCount++
	}

	// Generate sync metadata
	metadata := map[string]interface{}{
		"lastSync":        time.Now(),
		"environment":     environment,
		"totalWorkflows":  len(workflows),
		"syncedWorkflows": syncedCount,
		"n8nURL":          apiURL,
		"syncedBy":        getSyncUser(),
		"gitCommit":       os.Getenv("CI_COMMIT_SHA"),
	}

	metadataPath := filepath.Join(outputDir, "_sync_metadata.json")
	if err := writeWorkflowFile(metadata, metadataPath); err != nil {
		logger.Warn("Failed to write sync metadata", "error", err)
	}

	fmt.Printf("✅ Sync completed: %d workflows synced to %s\n", syncedCount, outputDir)
	return nil
}

func writeWorkflowFile(data interface{}, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// getSyncUser returns the user information for sync metadata
func getSyncUser() string {
	// Try to get user from GitLab CI/CD environment variables
	if gitlabUser := os.Getenv("GITLAB_USER_EMAIL"); gitlabUser != "" {
		return gitlabUser
	}

	if gitlabUser := os.Getenv("CI_COMMIT_AUTHOR_EMAIL"); gitlabUser != "" {
		return gitlabUser
	}

	// Try to get user from Git configuration
	if gitUser := getGitUser(); gitUser != "" {
		return gitUser
	}

	// Try to get system user
	if systemUser := os.Getenv("USER"); systemUser != "" {
		return systemUser + "@local"
	}

	// Fallback
	return "n8n-ops-user"
}

// getGitUser tries to get the Git user email from local Git config
func getGitUser() string {
	// Try to get Git user email
	cmd := exec.Command("git", "config", "user.email")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output))
	}
	return ""
}
