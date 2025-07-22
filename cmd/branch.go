package cmd

import (
        "fmt"

        "github.com/spf13/cobra"
        "github.com/n8n-workflows/n8n-ops/internal/git"
        "github.com/n8n-workflows/n8n-ops/internal/utils"
)

var branchCmd = &cobra.Command{
        Use:   "branch",
        Short: "Manage Git branches and environment mappings",
        Long: `Manage Git branches and their relationship with n8n environments.
This command helps you work with branches, switch between them, and understand
how branches map to different n8n environments (development, staging, production).

Examples:
  n8n-ops branch list                    # List all branches
  n8n-ops branch current                 # Show current branch  
  n8n-ops branch create feature-auth     # Create new branch
  n8n-ops branch switch staging          # Switch to staging branch
  n8n-ops branch env                     # Show environment for current branch`,
}

var branchListCmd = &cobra.Command{
        Use:   "list",
        Short: "List all Git branches",
        RunE:  runBranchList,
}

var branchCurrentCmd = &cobra.Command{
        Use:   "current", 
        Short: "Show current Git branch",
        RunE:  runBranchCurrent,
}

var branchCreateCmd = &cobra.Command{
        Use:   "create [branch-name]",
        Short: "Create and switch to new branch",
        Args:  cobra.ExactArgs(1),
        RunE:  runBranchCreate,
}

var branchSwitchCmd = &cobra.Command{
        Use:   "switch [branch-name]", 
        Short: "Switch to existing branch",
        Args:  cobra.ExactArgs(1),
        RunE:  runBranchSwitch,
}

var branchEnvCmd = &cobra.Command{
        Use:   "env",
        Short: "Show environment mapping for current branch",
        RunE:  runBranchEnv,
}

func init() {
        rootCmd.AddCommand(branchCmd)
        branchCmd.AddCommand(branchListCmd)
        branchCmd.AddCommand(branchCurrentCmd)
        branchCmd.AddCommand(branchCreateCmd)
        branchCmd.AddCommand(branchSwitchCmd)
        branchCmd.AddCommand(branchEnvCmd)
}

func runBranchList(cmd *cobra.Command, args []string) error {
        logger := utils.GetLogger()
        
        if !git.IsGitRepository() {
                return fmt.Errorf("not in a Git repository")
        }
        
        branches, err := git.GetBranchList()
        if err != nil {
                return fmt.Errorf("failed to get branch list: %w", err)
        }
        
        currentBranch, err := git.GetCurrentBranch()
        if err != nil {
                logger.Warn("Failed to get current branch", "error", err)
                currentBranch = ""
        }
        
        fmt.Println("Branches:")
        for _, branch := range branches {
                marker := "  "
                if branch == currentBranch {
                        marker = "* "
                }
                env := git.GetEnvironmentFromBranch(branch)
                fmt.Printf("%s%s (→ %s)\n", marker, branch, env)
        }
        
        return nil
}

func runBranchCurrent(cmd *cobra.Command, args []string) error {
        if !git.IsGitRepository() {
                return fmt.Errorf("not in a Git repository")
        }
        
        branch, err := git.GetCurrentBranch()
        if err != nil {
                return fmt.Errorf("failed to get current branch: %w", err)
        }
        
        env := git.GetEnvironmentFromBranch(branch)
        fmt.Printf("Current branch: %s (→ %s environment)\n", branch, env)
        
        return nil
}

func runBranchCreate(cmd *cobra.Command, args []string) error {
        logger := utils.GetLogger()
        branchName := args[0]
        
        if !git.IsGitRepository() {
                return fmt.Errorf("not in a Git repository")
        }
        
        logger.Info("Creating new branch", "branch", branchName)
        
        if err := git.CheckoutBranch(branchName); err != nil {
                return fmt.Errorf("failed to create branch %s: %w", branchName, err)
        }
        
        env := git.GetEnvironmentFromBranch(branchName)
        fmt.Printf("Created and switched to branch: %s (→ %s environment)\n", branchName, env)
        
        return nil
}

func runBranchSwitch(cmd *cobra.Command, args []string) error {
        logger := utils.GetLogger()
        branchName := args[0]
        
        if !git.IsGitRepository() {
                return fmt.Errorf("not in a Git repository")
        }
        
        logger.Info("Switching to branch", "branch", branchName)
        
        if err := git.CheckoutBranch(branchName); err != nil {
                return fmt.Errorf("failed to switch to branch %s: %w", branchName, err)
        }
        
        env := git.GetEnvironmentFromBranch(branchName)
        fmt.Printf("Switched to branch: %s (→ %s environment)\n", branchName, env)
        
        return nil
}

func runBranchEnv(cmd *cobra.Command, args []string) error {
        if !git.IsGitRepository() {
                return fmt.Errorf("not in a Git repository")
        }
        
        branch, err := git.GetCurrentBranch()
        if err != nil {
                return fmt.Errorf("failed to get current branch: %w", err)
        }
        
        env := git.GetEnvironmentFromBranch(branch)
        
        fmt.Printf("Branch: %s\n", branch)
        fmt.Printf("Environment: %s\n", env)
        fmt.Printf("\nBranch → Environment mapping:\n")
        fmt.Printf("  main/master/prod* → production\n")
        fmt.Printf("  staging/stage*    → staging\n") 
        fmt.Printf("  develop/dev*      → development\n")
        fmt.Printf("  feature/*         → development (default)\n")
        
        return nil
}