package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pmaojo/n8n-ops/internal/cliutils"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cliFrom(cmd)
		if cli == nil {
			return fmt.Errorf("CLI not initialized")
		}
		return runSync(cmd, args, cli)
	},
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

func runSync(cmd *cobra.Command, args []string, cli *CLI) error {
	ctx := context.Background()

	n8nClient, cm, err := cliutils.SetupClient(cli.Environment, demoMode, cli.Logger)
	if err != nil {
		var missing cliutils.MissingCredentialError
		if errors.As(err, &missing) {
			envSuffix := strings.ToUpper(cli.Environment)
			fmt.Printf("⚠️  n8n credentials not configured for %s environment\n", cli.Environment)
			fmt.Printf("💡 Set environment variables or use --demo flag:\n")
			fmt.Printf("   export N8N_%s_URL=http://localhost:3001\n", envSuffix)
			fmt.Printf("   export N8N_%s_API_KEY=n8n_api_mock_%s\n", envSuffix, cli.Environment)
			return nil
		}
		return err
	}

	checker := git.NewGitStatusChecker(".", nil)
	svc := isync.NewService(n8nClient, cm, checker, cli.Logger, cli.Environment)

	return svc.Sync(ctx, isync.Options{
		OutputDir: outputDir,
		Force:     force,
		DryRun:    syncDryRun,
	})
}
