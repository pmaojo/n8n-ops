package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/git"
	isync "github.com/pmaojo/n8n-ops/internal/sync"
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
	ctx := context.Background()

	cm := credentials.NewCredentialManager(environment)

	var (
		apiURL    string
		apiKey    string
		n8nClient client.Client
		err       error
	)

	if demoMode {
		n8nClient = client.NewDemoN8nClient()
	} else {
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
			return nil
		}
		n8nClient, err = client.New(apiURL, apiKey, nil)
		if err != nil {
			return fmt.Errorf("failed to create n8n client: %w", err)
		}
	}

	checker := git.NewGitStatusChecker(".", nil)
	svc := isync.NewService(n8nClient, cm, checker, logger, environment)

	return svc.Sync(ctx, isync.Options{
		OutputDir: outputDir,
		Force:     force,
		DryRun:    syncDryRun,
	})
}
