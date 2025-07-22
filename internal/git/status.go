package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GitStatus represents the status of files in Git
type GitStatus struct {
	ModifiedFiles   []string `json:"modifiedFiles"`
	UntrackedFiles  []string `json:"untrackedFiles"`
	StagedFiles     []string `json:"stagedFiles"`
	WorkflowFiles   []string `json:"workflowFiles"`
	HasChanges      bool     `json:"hasChanges"`
	CurrentBranch   string   `json:"currentBranch"`
	LastCommit      string   `json:"lastCommit"`
	UncommittedWorkflows []WorkflowChange `json:"uncommittedWorkflows"`
}

// WorkflowChange represents a changed workflow file
type WorkflowChange struct {
	FilePath     string    `json:"filePath"`
	WorkflowName string    `json:"workflowName"`
	Status       string    `json:"status"` // modified, added, deleted
	Environment  string    `json:"environment"`
	LastModified time.Time `json:"lastModified"`
	Branch       string    `json:"branch"`
}

// GitStatusChecker checks Git status for workflow changes
type GitStatusChecker struct {
	WorkingDir string
}

func NewGitStatusChecker(workingDir string) *GitStatusChecker {
	if workingDir == "" {
		workingDir = "."
	}
	return &GitStatusChecker{
		WorkingDir: workingDir,
	}
}

// GetStatus returns current Git status with workflow-specific information
func (gsc *GitStatusChecker) GetStatus() (*GitStatus, error) {
	status := &GitStatus{
		ModifiedFiles:        make([]string, 0),
		UntrackedFiles:       make([]string, 0),
		StagedFiles:          make([]string, 0),
		WorkflowFiles:        make([]string, 0),
		UncommittedWorkflows: make([]WorkflowChange, 0),
	}

	// Get current branch
	branch, err := gsc.getCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}
	status.CurrentBranch = branch

	// Get last commit
	lastCommit, err := gsc.getLastCommit()
	if err != nil {
		// Not a fatal error, might be initial commit
		status.LastCommit = "No commits"
	} else {
		status.LastCommit = lastCommit
	}

	// Get Git status
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = gsc.WorkingDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run git status: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 3 {
			continue
		}

		statusCode := line[:2]
		filePath := strings.TrimSpace(line[2:])

		// Check if it's a workflow file
		isWorkflowFile := gsc.isWorkflowFile(filePath)
		
		if isWorkflowFile {
			status.WorkflowFiles = append(status.WorkflowFiles, filePath)
			
			workflowChange := WorkflowChange{
				FilePath:     filePath,
				WorkflowName: gsc.extractWorkflowName(filePath),
				Environment:  gsc.extractEnvironment(filePath),
				LastModified: time.Now(), // Would get from file stats in real implementation
				Branch:       branch,
			}

			// Parse Git status codes
			switch {
			case strings.Contains(statusCode, "M"): // Modified
				status.ModifiedFiles = append(status.ModifiedFiles, filePath)
				workflowChange.Status = "modified"
			case strings.Contains(statusCode, "A"): // Added
				status.StagedFiles = append(status.StagedFiles, filePath)
				workflowChange.Status = "added"
			case strings.Contains(statusCode, "D"): // Deleted
				workflowChange.Status = "deleted"
			case strings.Contains(statusCode, "?"): // Untracked
				status.UntrackedFiles = append(status.UntrackedFiles, filePath)
				workflowChange.Status = "untracked"
			}

			status.UncommittedWorkflows = append(status.UncommittedWorkflows, workflowChange)
		} else {
			// Handle non-workflow files
			switch {
			case strings.Contains(statusCode, "M"):
				status.ModifiedFiles = append(status.ModifiedFiles, filePath)
			case strings.Contains(statusCode, "A"):
				status.StagedFiles = append(status.StagedFiles, filePath)
			case strings.Contains(statusCode, "?"):
				status.UntrackedFiles = append(status.UntrackedFiles, filePath)
			}
		}
	}

	status.HasChanges = len(status.ModifiedFiles) > 0 || len(status.UntrackedFiles) > 0 || len(status.StagedFiles) > 0

	return status, nil
}

// isWorkflowFile checks if a file path is a workflow JSON file
func (gsc *GitStatusChecker) isWorkflowFile(filePath string) bool {
	// Check if file is in workflows directory and has .json extension
	return strings.Contains(filePath, "workflows/") && strings.HasSuffix(filePath, ".json")
}

// extractWorkflowName extracts workflow name from file path
func (gsc *GitStatusChecker) extractWorkflowName(filePath string) string {
	fileName := filepath.Base(filePath)
	// Remove .json extension and version suffix
	name := strings.TrimSuffix(fileName, ".json")
	
	// Remove version pattern (e.g., -v1.0.1234)
	if idx := strings.LastIndex(name, "-v"); idx != -1 {
		name = name[:idx]
	}
	
	// Convert hyphens to spaces and capitalize
	name = strings.ReplaceAll(name, "-", " ")
	return strings.Title(name)
}

// extractEnvironment extracts environment from file path
func (gsc *GitStatusChecker) extractEnvironment(filePath string) string {
	if strings.Contains(filePath, "workflows/development/") {
		return "development"
	} else if strings.Contains(filePath, "workflows/staging/") {
		return "staging"
	} else if strings.Contains(filePath, "workflows/production/") {
		return "production"
	}
	return "unknown"
}

// getCurrentBranch gets the current Git branch
func (gsc *GitStatusChecker) getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = gsc.WorkingDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getLastCommit gets the last commit hash and message
func (gsc *GitStatusChecker) getLastCommit() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%h %s")
	cmd.Dir = gsc.WorkingDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetUncommittedWorkflowSummary returns a summary of uncommitted workflow changes
func (gsc *GitStatusChecker) GetUncommittedWorkflowSummary() (string, error) {
	status, err := gsc.GetStatus()
	if err != nil {
		return "", err
	}

	if len(status.UncommittedWorkflows) == 0 {
		return "All workflows are committed", nil
	}

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("⚠️  %d uncommitted workflow changes:\n\n", len(status.UncommittedWorkflows)))

	for _, workflow := range status.UncommittedWorkflows {
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
		default:
			statusIcon = "📄"
		}

		summary.WriteString(fmt.Sprintf("%s %s (%s)\n", statusIcon, workflow.WorkflowName, workflow.Status))
		summary.WriteString(fmt.Sprintf("   Environment: %s\n", workflow.Environment))
		summary.WriteString(fmt.Sprintf("   File: %s\n\n", workflow.FilePath))
	}

	summary.WriteString("💡 Run 'git add . && git commit -m \"Update workflows\"' to commit changes")

	return summary.String(), nil
}

// WarnIfUncommittedChanges prints a warning if there are uncommitted workflow changes
func (gsc *GitStatusChecker) WarnIfUncommittedChanges() error {
	status, err := gsc.GetStatus()
	if err != nil {
		return err
	}

	if len(status.UncommittedWorkflows) > 0 {
		fmt.Printf("\n🚨 WARNING: Uncommitted Workflow Changes Detected!\n")
		fmt.Printf("=====================================\n\n")

		for _, workflow := range status.UncommittedWorkflows {
			var statusColor string
			switch workflow.Status {
			case "modified":
				statusColor = "\033[1;33m" // Yellow
			case "added":
				statusColor = "\033[1;32m" // Green
			case "deleted":
				statusColor = "\033[1;31m" // Red
			case "untracked":
				statusColor = "\033[1;35m" // Magenta
			default:
				statusColor = "\033[0m" // Reset
			}

			fmt.Printf("%s• %s (%s)\033[0m\n", statusColor, workflow.WorkflowName, workflow.Status)
			fmt.Printf("  Environment: %s | File: %s\n", workflow.Environment, workflow.FilePath)
		}

		fmt.Printf("\n💡 Recommendation:\n")
		fmt.Printf("  git add .\n")
		fmt.Printf("  git commit -m \"Update %d workflow(s)\"\n", len(status.UncommittedWorkflows))
		fmt.Printf("  git push origin %s\n\n", status.CurrentBranch)
	}

	return nil
}

// CheckBeforeSync checks for uncommitted changes before sync operations
func (gsc *GitStatusChecker) CheckBeforeSync() error {
	status, err := gsc.GetStatus()
	if err != nil {
		return err
	}

	if len(status.UncommittedWorkflows) > 0 {
		fmt.Printf("⚠️  Found %d uncommitted workflow changes\n", len(status.UncommittedWorkflows))
		fmt.Printf("Sync operations may overwrite local changes\n\n")

		for _, workflow := range status.UncommittedWorkflows {
			fmt.Printf("• %s (%s) - %s\n", workflow.WorkflowName, workflow.Status, workflow.Environment)
		}

		fmt.Printf("\nRecommendation: Commit your changes first:\n")
		fmt.Printf("  git add . && git commit -m \"Save workflow changes\"\n\n")
		
		return fmt.Errorf("uncommitted workflow changes detected - commit before syncing")
	}

	return nil
}

// AutoCommitWorkflows automatically commits workflow changes if enabled
func (gsc *GitStatusChecker) AutoCommitWorkflows(message string) error {
	status, err := gsc.GetStatus()
	if err != nil {
		return err
	}

	if len(status.UncommittedWorkflows) == 0 {
		return nil // Nothing to commit
	}

	// Add workflow files
	for _, workflow := range status.UncommittedWorkflows {
		cmd := exec.Command("git", "add", workflow.FilePath)
		cmd.Dir = gsc.WorkingDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add %s: %w", workflow.FilePath, err)
		}
	}

	// Commit changes
	if message == "" {
		message = fmt.Sprintf("Auto-commit %d workflow changes", len(status.UncommittedWorkflows))
	}

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = gsc.WorkingDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	fmt.Printf("✅ Auto-committed %d workflow changes\n", len(status.UncommittedWorkflows))
	return nil
}