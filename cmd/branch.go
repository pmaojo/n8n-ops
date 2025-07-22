package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/n8n-ops/n8n-ops/internal/git"
	"github.com/n8n-ops/n8n-ops/internal/logging"
)

var (
	branchListAll    bool
	branchCompare    string
	branchOutputJSON bool
	branchShowActive bool
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage workflow branches with intelligent naming conventions",
	Long: `The branch command helps you manage n8n workflows across Git branches.
It can identify active workflows in specific branches, compare workflows between
branches, and provide insights into your workflow development process.

Features:
- Identify workflows in current or specific branch
- Compare workflows between branches
- List all branches with workflow information
- Show only branches with active workflows
- Export branch workflow data as JSON

Examples:
  n8n-ops branch                           # Show workflows in current branch
  n8n-ops branch --list                    # List all branches with workflows  
  n8n-ops branch --compare main            # Compare current branch with main
  n8n-ops branch --list --active           # Show only branches with active workflows
  n8n-ops branch --json                    # Export as JSON`,
	RunE: runBranchCmd,
}

func init() {
	rootCmd.AddCommand(branchCmd)

	branchCmd.Flags().BoolVar(&branchListAll, "list", false, "List all branches with workflow information")
	branchCmd.Flags().StringVar(&branchCompare, "compare", "", "Compare current branch with specified branch")
	branchCmd.Flags().BoolVar(&branchOutputJSON, "json", false, "Output results as JSON")
	branchCmd.Flags().BoolVar(&branchShowActive, "active", false, "Show only branches with active workflows")
}

func runBranchCmd(cmd *cobra.Command, args []string) error {
	logger := logging.WithFields(map[string]string{
		"command": "branch",
		"env":     environment,
	})

	logger.Info("Starting branch workflow analysis")

	// Initialize branch tracker
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	tracker := git.NewBranchTracker(workingDir)

	// Handle different command modes
	switch {
	case branchListAll:
		return handleListAllBranches(tracker, logger)
	case branchCompare != "":
		return handleCompareBranches(tracker, branchCompare, logger)
	default:
		return handleCurrentBranch(tracker, logger)
	}
}

func handleCurrentBranch(tracker *git.BranchTracker, logger *logging.ContextLogger) error {
	currentBranch, err := tracker.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	logger.Info("Analyzing current branch: %s", currentBranch)

	branchInfo, err := tracker.GetBranchWorkflows(currentBranch)
	if err != nil {
		return fmt.Errorf("failed to analyze branch workflows: %w", err)
	}

	if branchOutputJSON {
		return outputJSON(branchInfo)
	}

	return displayBranchInfo(branchInfo)
}

func handleListAllBranches(tracker *git.BranchTracker, logger *logging.ContextLogger) error {
	logger.Info("Analyzing all branches")

	allBranches, err := tracker.GetAllBranchesWorkflows()
	if err != nil {
		return fmt.Errorf("failed to analyze all branches: %w", err)
	}

	// Filter for active workflows if requested
	if branchShowActive {
		filteredBranches := make(map[string]*git.BranchWorkflowInfo)
		for branch, info := range allBranches {
			activeCount := 0
			for _, workflow := range info.WorkflowFiles {
				if workflow.Active {
					activeCount++
				}
			}
			if activeCount > 0 {
				filteredBranches[branch] = info
			}
		}
		allBranches = filteredBranches
	}

	if branchOutputJSON {
		return outputJSON(allBranches)
	}

	return displayAllBranches(allBranches)
}

func handleCompareBranches(tracker *git.BranchTracker, targetBranch string, logger *logging.ContextLogger) error {
	currentBranch, err := tracker.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	logger.Info("Comparing branches: %s vs %s", currentBranch, targetBranch)

	comparison, err := tracker.CompareWorkflowsBetweenBranches(currentBranch, targetBranch)
	if err != nil {
		return fmt.Errorf("failed to compare branches: %w", err)
	}

	if branchOutputJSON {
		return outputJSON(comparison)
	}

	return displayBranchComparison(comparison)
}

func displayBranchInfo(info *git.BranchWorkflowInfo) error {
	if language == "es" {
		fmt.Printf("🌿 INFORMACIÓN DE RAMA: %s\n", info.Branch)
		fmt.Println("=" + strings.Repeat("=", 50))
	} else {
		fmt.Printf("🌿 BRANCH INFORMATION: %s\n", info.Branch)
		fmt.Println("=" + strings.Repeat("=", 50))
	}

	fmt.Printf("📁 Environment: %s\n", info.Environment)
	fmt.Printf("📝 Last Commit: %s\n", info.LastCommit[:8])
	fmt.Printf("💬 Commit Message: %s\n", info.CommitMessage)
	fmt.Printf("👤 Author: %s\n", info.Author)
	fmt.Printf("⏰ Timestamp: %s\n", info.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("📊 Status: %s\n", info.Status)

	fmt.Println("\n📋 WORKFLOW SUMMARY:")
	fmt.Printf("   Total Workflows: %s\n", info.Metadata["totalWorkflows"])
	fmt.Printf("   Active Workflows: %s\n", info.Metadata["activeWorkflows"])
	fmt.Printf("   Branch Type: %s\n", info.Metadata["branchType"])

	if len(info.WorkflowFiles) > 0 {
		fmt.Println("\n🔧 WORKFLOWS IN THIS BRANCH:")
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tPATH\tSTATUS\tID")
		fmt.Fprintln(w, "----\t----\t------\t--")

		for _, workflow := range info.WorkflowFiles {
			status := "inactive"
			if workflow.Active {
				status = "active"
			}
			
			name := workflow.WorkflowName
			if len(name) > 30 {
				name = name[:27] + "..."
			}
			
			path := filepath.Base(workflow.Path)
			id := workflow.WorkflowID
			if len(id) > 8 {
				id = id[:8] + "..."
			}
			
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, path, status, id)
		}
		
		w.Flush()
	} else {
		fmt.Println("\n📭 No workflows found in this branch")
	}

	return nil
}

func displayAllBranches(branches map[string]*git.BranchWorkflowInfo) error {
	if language == "es" {
		fmt.Println("🌳 ANÁLISIS DE TODAS LAS RAMAS")
		fmt.Println("=" + strings.Repeat("=", 50))
	} else {
		fmt.Println("🌳 ALL BRANCHES ANALYSIS")
		fmt.Println("=" + strings.Repeat("=", 50))
	}

	if len(branches) == 0 {
		fmt.Println("📭 No branches with workflows found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "BRANCH\tENV\tWORKFLOWS\tACTIVE\tLAST COMMIT\tAUTHOR")
	fmt.Fprintln(w, "------\t---\t---------\t------\t-----------\t------")

	for _, info := range branches {
		totalWorkflows := len(info.WorkflowFiles)
		activeWorkflows := 0
		for _, wf := range info.WorkflowFiles {
			if wf.Active {
				activeWorkflows++
			}
		}

		branch := info.Branch
		if len(branch) > 20 {
			branch = branch[:17] + "..."
		}

		author := info.Author
		if len(author) > 15 {
			author = author[:12] + "..."
		}

		commitHash := info.LastCommit
		if len(commitHash) > 8 {
			commitHash = commitHash[:8]
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\t%s\n",
			branch, info.Environment, totalWorkflows, activeWorkflows, commitHash, author)
	}

	w.Flush()

	// Summary
	totalBranches := len(branches)
	totalWorkflows := 0
	totalActive := 0

	for _, info := range branches {
		totalWorkflows += len(info.WorkflowFiles)
		for _, wf := range info.WorkflowFiles {
			if wf.Active {
				totalActive++
			}
		}
	}

	fmt.Printf("\n📊 SUMMARY: %d branches, %d workflows, %d active\n", 
		totalBranches, totalWorkflows, totalActive)

	return nil
}

func displayBranchComparison(comparison *git.WorkflowComparison) error {
	if language == "es" {
		fmt.Printf("🔍 COMPARACIÓN DE RAMAS: %s vs %s\n", comparison.BranchA, comparison.BranchB)
		fmt.Println("=" + strings.Repeat("=", 60))
	} else {
		fmt.Printf("🔍 BRANCH COMPARISON: %s vs %s\n", comparison.BranchA, comparison.BranchB)
		fmt.Println("=" + strings.Repeat("=", 60))
	}

	fmt.Printf("📊 SUMMARY:\n")
	fmt.Printf("   ✅ Added: %d workflows\n", len(comparison.Added))
	fmt.Printf("   📝 Modified: %d workflows\n", len(comparison.Modified))
	fmt.Printf("   ❌ Deleted: %d workflows\n", len(comparison.Deleted))
	fmt.Printf("   ⚪ Unchanged: %d workflows\n", len(comparison.Unchanged))

	if len(comparison.Added) > 0 {
		fmt.Println("\n✅ ADDED WORKFLOWS:")
		for _, wf := range comparison.Added {
			status := "inactive"
			if wf.Active {
				status = "active"
			}
			fmt.Printf("   + %s (%s) [%s]\n", wf.WorkflowName, filepath.Base(wf.Path), status)
		}
	}

	if len(comparison.Modified) > 0 {
		fmt.Println("\n📝 MODIFIED WORKFLOWS:")
		for _, wf := range comparison.Modified {
			status := "inactive"
			if wf.Active {
				status = "active"
			}
			fmt.Printf("   ~ %s (%s) [%s]\n", wf.WorkflowName, filepath.Base(wf.Path), status)
		}
	}

	if len(comparison.Deleted) > 0 {
		fmt.Println("\n❌ DELETED WORKFLOWS:")
		for _, wf := range comparison.Deleted {
			fmt.Printf("   - %s (%s)\n", wf.WorkflowName, filepath.Base(wf.Path))
		}
	}

	// Recommendations
	fmt.Println("\n💡 RECOMMENDATIONS:")
	if len(comparison.Added) > 0 || len(comparison.Modified) > 0 {
		fmt.Println("   • Test workflows in development environment before merging")
		fmt.Println("   • Update workflow documentation if needed")
	}
	if len(comparison.Deleted) > 0 {
		fmt.Println("   • Verify deleted workflows are intentional")
		fmt.Println("   • Consider creating backups before merging")
	}
	if len(comparison.Added) == 0 && len(comparison.Modified) == 0 && len(comparison.Deleted) == 0 {
		fmt.Println("   • No workflow changes detected - branches are in sync")
	}

	return nil
}

func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}