package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/pmaojo/n8n-ops/internal/git"
)

var (
	branchListAll    bool
	branchCompare    string
	branchOutputJSON bool
	branchShowActive bool
)

// BranchWorkflowInfo represents workflow information tied to a specific branch
type BranchWorkflowInfo struct {
	Branch        string            `json:"branch"`
	WorkflowFiles []WorkflowFile    `json:"workflowFiles"`
	LastCommit    string            `json:"lastCommit"`
	CommitMessage string            `json:"commitMessage"`
	Author        string            `json:"author"`
	Timestamp     time.Time         `json:"timestamp"`
	Environment   string            `json:"environment"`
	Status        string            `json:"status"`
	Metadata      map[string]string `json:"metadata"`
}

// WorkflowFile represents a workflow file in a branch
type WorkflowFile struct {
	Path         string    `json:"path"`
	WorkflowID   string    `json:"workflowId"`
	WorkflowName string    `json:"workflowName"`
	Active       bool      `json:"active"`
	LastModified time.Time `json:"lastModified"`
	FileHash     string    `json:"fileHash"`
}

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
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		return runBranchCmd(cmd, args, git.NewService(git.NewExecutor(cwd)))
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)

	branchCmd.Flags().BoolVar(&branchListAll, "list", false, "List all branches with workflow information")
	branchCmd.Flags().StringVar(&branchCompare, "compare", "", "Compare current branch with specified branch")
	branchCmd.Flags().BoolVar(&branchOutputJSON, "json", false, "Output results as JSON")
	branchCmd.Flags().BoolVar(&branchShowActive, "active", false, "Show only branches with active workflows")
}

func runBranchCmd(cmd *cobra.Command, args []string, svc *git.Service) error {
	if language == "es" {
		fmt.Println("🌿 Iniciando análisis de ramas...")
	} else {
		fmt.Println("🌿 Starting branch workflow analysis...")
	}

	// Handle different command modes
	switch {
	case branchListAll:
		return handleListAllBranches(svc)
	case branchCompare != "":
		return handleCompareBranches(svc, branchCompare)
	default:
		return handleCurrentBranch(svc)
	}
}

func handleCurrentBranch(svc *git.Service) error {
	currentBranch, err := getCurrentBranch(svc)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	fmt.Printf("Analyzing current branch: %s\n", currentBranch)

	branchInfo, err := getBranchWorkflows(svc, currentBranch)
	if err != nil {
		return fmt.Errorf("failed to analyze branch workflows: %w", err)
	}

	if branchOutputJSON {
		return outputJSON(branchInfo)
	}

	return displayBranchInfo(branchInfo)
}

func handleListAllBranches(svc *git.Service) error {
	fmt.Println("Analyzing all branches...")

	allBranches, err := getAllBranchesWorkflows(svc)
	if err != nil {
		return fmt.Errorf("failed to analyze all branches: %w", err)
	}

	// Filter for active workflows if requested
	if branchShowActive {
		filteredBranches := make(map[string]*BranchWorkflowInfo)
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

func handleCompareBranches(svc *git.Service, targetBranch string) error {
	currentBranch, err := getCurrentBranch(svc)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	fmt.Printf("Comparing branches: %s vs %s\n", currentBranch, targetBranch)

	branchInfoA, err := getBranchWorkflows(svc, currentBranch)
	if err != nil {
		return fmt.Errorf("failed to get workflows for branch %s: %w", currentBranch, err)
	}

	branchInfoB, err := getBranchWorkflows(svc, targetBranch)
	if err != nil {
		return fmt.Errorf("failed to get workflows for branch %s: %w", targetBranch, err)
	}

	comparison := compareBranchWorkflows(branchInfoA, branchInfoB)

	if branchOutputJSON {
		return outputJSON(comparison)
	}

	return displayBranchComparison(comparison)
}

// Helper functions for branch operations
func getCurrentBranch(svc *git.Service) (string, error) {
	return svc.GetCurrentBranch()
}

func getBranchWorkflows(svc *git.Service, branch string) (*BranchWorkflowInfo, error) {
	output, err := svc.Log(branch, "%H|%s|%an|%ct")
	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(output), "|")
	if len(parts) < 4 {
		return nil, fmt.Errorf("unexpected git log format")
	}

	workflowFiles, err := findWorkflowFiles(svc, branch)
	if err != nil {
		return nil, err
	}

	// Determine environment
	environment := determineEnvironment(branch)

	branchInfo := &BranchWorkflowInfo{
		Branch:        branch,
		WorkflowFiles: workflowFiles,
		LastCommit:    parts[0],
		CommitMessage: parts[1],
		Author:        parts[2],
		Timestamp:     time.Now(), // Simplified
		Environment:   environment,
		Status:        "active",
		Metadata:      make(map[string]string),
	}

	// Add metadata
	branchInfo.Metadata["totalWorkflows"] = fmt.Sprintf("%d", len(workflowFiles))
	branchInfo.Metadata["activeWorkflows"] = fmt.Sprintf("%d", countActiveWorkflows(workflowFiles))
	branchInfo.Metadata["branchType"] = getBranchType(branch)

	return branchInfo, nil
}

func getAllBranchesWorkflows(svc *git.Service) (map[string]*BranchWorkflowInfo, error) {
	branches, err := getAllBranches(svc)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*BranchWorkflowInfo)
	for _, branch := range branches {
		branchInfo, err := getBranchWorkflows(svc, branch)
		if err != nil {
			fmt.Printf("Warning: Failed to analyze branch %s: %v\n", branch, err)
			continue
		}
		result[branch] = branchInfo
	}

	return result, nil
}

func findWorkflowFiles(svc *git.Service, branch string) ([]WorkflowFile, error) {
	files, err := svc.LsTree(branch, "workflows/")
	if err != nil {
		return []WorkflowFile{}, nil // No workflows directory
	}
	var workflowFiles []WorkflowFile

	for _, file := range files {
		if strings.HasSuffix(file, ".json") {
			wf := WorkflowFile{
				Path:         file,
				WorkflowName: filepath.Base(file),
				Active:       false, // Simplified
				LastModified: time.Now(),
				FileHash:     fmt.Sprintf("%x", len(file)), // Simplified hash
			}
			workflowFiles = append(workflowFiles, wf)
		}
	}

	return workflowFiles, nil
}

func getAllBranches(svc *git.Service) ([]string, error) {
	branches, err := svc.RemoteBranches()
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func determineEnvironment(branch string) string {
	branch = strings.ToLower(branch)

	switch {
	case strings.Contains(branch, "production") || branch == "main" || branch == "master":
		return "production"
	case strings.Contains(branch, "staging"):
		return "staging"
	default:
		return "development"
	}
}

func getBranchType(branch string) string {
	switch {
	case strings.HasPrefix(branch, "feature/"):
		return "feature"
	case strings.HasPrefix(branch, "hotfix/"):
		return "hotfix"
	case branch == "main" || branch == "master":
		return "main"
	default:
		return "other"
	}
}

func countActiveWorkflows(workflows []WorkflowFile) int {
	count := 0
	for _, wf := range workflows {
		if wf.Active {
			count++
		}
	}
	return count
}

type WorkflowComparison struct {
	BranchA   string         `json:"branchA"`
	BranchB   string         `json:"branchB"`
	Added     []WorkflowFile `json:"added"`
	Deleted   []WorkflowFile `json:"deleted"`
	Modified  []WorkflowFile `json:"modified"`
	Unchanged []WorkflowFile `json:"unchanged"`
}

func compareBranchWorkflows(branchA, branchB *BranchWorkflowInfo) *WorkflowComparison {
	comparison := &WorkflowComparison{
		BranchA: branchA.Branch,
		BranchB: branchB.Branch,
	}

	filesA := make(map[string]WorkflowFile)
	filesB := make(map[string]WorkflowFile)

	for _, file := range branchA.WorkflowFiles {
		filesA[file.Path] = file
	}

	for _, file := range branchB.WorkflowFiles {
		filesB[file.Path] = file
	}

	// Find differences
	for path, fileB := range filesB {
		if fileA, exists := filesA[path]; exists {
			if fileA.FileHash != fileB.FileHash {
				comparison.Modified = append(comparison.Modified, fileB)
			} else {
				comparison.Unchanged = append(comparison.Unchanged, fileB)
			}
		} else {
			comparison.Added = append(comparison.Added, fileB)
		}
	}

	for path, fileA := range filesA {
		if _, exists := filesB[path]; !exists {
			comparison.Deleted = append(comparison.Deleted, fileA)
		}
	}

	return comparison
}

func displayBranchInfo(info *BranchWorkflowInfo) error {
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

func displayAllBranches(branches map[string]*BranchWorkflowInfo) error {
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

func displayBranchComparison(comparison *WorkflowComparison) error {
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
