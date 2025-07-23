package cmd

import (
        "github.com/spf13/cobra"
)

// deployCmd es un alias conveniente para sync --to-n8n
var deployCmd = &cobra.Command{
        Use:   "deploy",
        Short: "Deploy workflows to n8n (alias for sync --to-n8n)",
        Long: `Deploy is a convenient alias for 'sync --to-n8n' that uploads local workflow changes to n8n.

This command provides the same intelligent deployment capabilities as sync --to-n8n:
• Git status verification before deployment
• Timestamp-based change detection  
• Conflict resolution with user prompts
• Dry-run mode for safety

Examples:
  n8n-ops deploy --env staging        # Deploy to staging (with Git checks)
  n8n-ops deploy --force              # Deploy ignoring Git status
  n8n-ops deploy --dry-run            # Show what would be deployed`,
        RunE: runDeploy,
}

func init() {
        rootCmd.AddCommand(deployCmd)
        
        // Inherit sync command flags for the deploy alias
        deployCmd.Flags().BoolVarP(&force, "force", "f", false, "force sync, overwriting conflicts")
        deployCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory (default: workflows/<environment>)")
        deployCmd.Flags().StringVarP(&branch, "branch", "b", "", "git branch (defaults to current branch)")
        deployCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "show what would be synced without making changes")
}

func runDeploy(cmd *cobra.Command, args []string) error {
        // Set the --to-n8n flag for sync command
        toN8n = true
        
        // Call the existing sync function with --to-n8n enabled
        return runSync(cmd, args)
}