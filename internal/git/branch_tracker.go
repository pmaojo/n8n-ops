package git

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/n8n-ops/n8n-ops/internal/logging"
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
	Status        string            `json:"status"` // active, merged, stale
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

// BranchTracker manages branch-specific workflow tracking
type BranchTracker struct {
	projectPath string
	logger      *logging.ContextLogger
}

// NewBranchTracker creates a new branch tracker
func NewBranchTracker(projectPath string) *BranchTracker {
	return &BranchTracker{
		projectPath: projectPath,
		logger:      logging.WithField("component", "branch-tracker"),
	}
}

// GetCurrentBranch returns the current Git branch
func (bt *BranchTracker) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = bt.projectPath
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	
	branch := strings.TrimSpace(string(output))
	bt.logger.Debug("Current branch detected: %s", branch)
	
	return branch, nil
}

// GetBranchWorkflows identifies all workflows in the current branch
func (bt *BranchTracker) GetBranchWorkflows(branch string) (*BranchWorkflowInfo, error) {
	bt.logger.Info("Analyzing workflows in branch: %s", branch)
	
	// Get branch information
	commitInfo, err := bt.getBranchCommitInfo(branch)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch commit info: %w", err)
	}
	
	// Find workflow files in the branch
	workflowFiles, err := bt.findWorkflowFiles(branch)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow files: %w", err)
	}
	
	// Determine environment based on branch name
	environment := bt.determineEnvironment(branch)
	
	branchInfo := &BranchWorkflowInfo{
		Branch:        branch,
		WorkflowFiles: workflowFiles,
		LastCommit:    commitInfo.Hash,
		CommitMessage: commitInfo.Message,
		Author:        commitInfo.Author,
		Timestamp:     commitInfo.Timestamp,
		Environment:   environment,
		Status:        bt.determineBranchStatus(branch),
		Metadata:      make(map[string]string),
	}
	
	// Add branch metadata
	branchInfo.Metadata["totalWorkflows"] = fmt.Sprintf("%d", len(workflowFiles))
	branchInfo.Metadata["activeWorkflows"] = fmt.Sprintf("%d", bt.countActiveWorkflows(workflowFiles))
	branchInfo.Metadata["branchType"] = bt.getBranchType(branch)
	
	bt.logger.Info("Branch analysis complete: %d workflows found in %s", 
		len(workflowFiles), branch)
	
	return branchInfo, nil
}

// GetAllBranchesWorkflows returns workflow information for all branches
func (bt *BranchTracker) GetAllBranchesWorkflows() (map[string]*BranchWorkflowInfo, error) {
	bt.logger.Info("Analyzing workflows across all branches")
	
	branches, err := bt.getAllBranches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}
	
	branchWorkflows := make(map[string]*BranchWorkflowInfo)
	
	for _, branch := range branches {
		branchInfo, err := bt.GetBranchWorkflows(branch)
		if err != nil {
			bt.logger.Warn("Failed to analyze branch %s: %v", branch, err)
			continue
		}
		branchWorkflows[branch] = branchInfo
	}
	
	bt.logger.Info("Multi-branch analysis complete: %d branches analyzed", len(branchWorkflows))
	
	return branchWorkflows, nil
}

// CompareWorkflowsBetweenBranches compares workflows between two branches
func (bt *BranchTracker) CompareWorkflowsBetweenBranches(branchA, branchB string) (*WorkflowComparison, error) {
	bt.logger.Info("Comparing workflows between branches: %s vs %s", branchA, branchB)
	
	workflowsA, err := bt.GetBranchWorkflows(branchA)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflows for branch %s: %w", branchA, err)
	}
	
	workflowsB, err := bt.GetBranchWorkflows(branchB)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflows for branch %s: %w", branchB, err)
	}
	
	comparison := bt.compareWorkflows(workflowsA, workflowsB)
	
	bt.logger.Info("Branch comparison complete: %d added, %d modified, %d deleted",
		len(comparison.Added), len(comparison.Modified), len(comparison.Deleted))
	
	return comparison, nil
}

// WorkflowComparison represents the differences between workflows in two branches
type WorkflowComparison struct {
	BranchA   string         `json:"branchA"`
	BranchB   string         `json:"branchB"`
	Added     []WorkflowFile `json:"added"`     // Files in B but not in A
	Deleted   []WorkflowFile `json:"deleted"`   // Files in A but not in B
	Modified  []WorkflowFile `json:"modified"`  // Files that changed between A and B
	Unchanged []WorkflowFile `json:"unchanged"` // Files that are identical
}

// Internal helper methods

type CommitInfo struct {
	Hash      string
	Message   string
	Author    string
	Timestamp time.Time
}

func (bt *BranchTracker) getBranchCommitInfo(branch string) (*CommitInfo, error) {
	// Get commit hash
	cmd := exec.Command("git", "rev-parse", branch)
	cmd.Dir = bt.projectPath
	hashOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	hash := strings.TrimSpace(string(hashOutput))
	
	// Get commit message and author
	cmd = exec.Command("git", "log", "-1", "--format=%s|%an|%ct", branch)
	cmd.Dir = bt.projectPath
	infoOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	parts := strings.Split(strings.TrimSpace(string(infoOutput)), "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected git log format")
	}
	
	timestamp := time.Now() // Default fallback
	if len(parts[2]) > 0 {
		if ts, err := time.Parse("1136239445", parts[2]); err == nil {
			timestamp = ts
		}
	}
	
	return &CommitInfo{
		Hash:      hash,
		Message:   parts[0],
		Author:    parts[1],
		Timestamp: timestamp,
	}, nil
}

func (bt *BranchTracker) findWorkflowFiles(branch string) ([]WorkflowFile, error) {
	// Get list of .json files in workflows directory for the branch
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", branch, "workflows/")
	cmd.Dir = bt.projectPath
	
	output, err := cmd.Output()
	if err != nil {
		bt.logger.Debug("No workflows directory found in branch %s", branch)
		return []WorkflowFile{}, nil
	}
	
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var workflowFiles []WorkflowFile
	
	for _, file := range files {
		if strings.HasSuffix(file, ".json") {
			workflowFile, err := bt.analyzeWorkflowFile(branch, file)
			if err != nil {
				bt.logger.Warn("Failed to analyze workflow file %s: %v", file, err)
				continue
			}
			workflowFiles = append(workflowFiles, *workflowFile)
		}
	}
	
	return workflowFiles, nil
}

func (bt *BranchTracker) analyzeWorkflowFile(branch, filePath string) (*WorkflowFile, error) {
	// Get file content from specific branch
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", branch, filePath))
	cmd.Dir = bt.projectPath
	
	content, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	// Parse workflow JSON to extract metadata
	var workflow map[string]interface{}
	if err := json.Unmarshal(content, &workflow); err != nil {
		return nil, fmt.Errorf("invalid workflow JSON: %w", err)
	}
	
	// Extract workflow information
	workflowID := ""
	workflowName := filepath.Base(filePath)
	active := false
	
	if id, ok := workflow["id"].(string); ok {
		workflowID = id
	}
	if name, ok := workflow["name"].(string); ok {
		workflowName = name
	}
	if activeFlag, ok := workflow["active"].(bool); ok {
		active = activeFlag
	}
	
	// Get file hash for change detection
	cmd = exec.Command("git", "rev-parse", fmt.Sprintf("%s:%s", branch, filePath))
	cmd.Dir = bt.projectPath
	hashOutput, _ := cmd.Output()
	fileHash := strings.TrimSpace(string(hashOutput))
	
	return &WorkflowFile{
		Path:         filePath,
		WorkflowID:   workflowID,
		WorkflowName: workflowName,
		Active:       active,
		LastModified: time.Now(), // Could be extracted from git log if needed
		FileHash:     fileHash,
	}, nil
}

func (bt *BranchTracker) getAllBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r", "--format=%(refname:short)")
	cmd.Dir = bt.projectPath
	
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var branches []string
	
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		// Skip HEAD and origin/ prefix
		if !strings.Contains(branch, "HEAD") && strings.HasPrefix(branch, "origin/") {
			branches = append(branches, strings.TrimPrefix(branch, "origin/"))
		}
	}
	
	return branches, nil
}

func (bt *BranchTracker) determineEnvironment(branch string) string {
	branch = strings.ToLower(branch)
	
	switch {
	case strings.Contains(branch, "production") || strings.Contains(branch, "prod") || branch == "main" || branch == "master":
		return "production"
	case strings.Contains(branch, "staging") || strings.Contains(branch, "stage"):
		return "staging"
	case strings.Contains(branch, "development") || strings.Contains(branch, "develop") || strings.Contains(branch, "dev"):
		return "development"
	case strings.HasPrefix(branch, "feature/") || strings.HasPrefix(branch, "hotfix/"):
		return "development"
	default:
		return "development"
	}
}

func (bt *BranchTracker) determineBranchStatus(branch string) string {
	// This could be enhanced to check if branch is merged, stale, etc.
	// For now, return basic status
	return "active"
}

func (bt *BranchTracker) countActiveWorkflows(workflows []WorkflowFile) int {
	count := 0
	for _, wf := range workflows {
		if wf.Active {
			count++
		}
	}
	return count
}

func (bt *BranchTracker) getBranchType(branch string) string {
	switch {
	case strings.HasPrefix(branch, "feature/"):
		return "feature"
	case strings.HasPrefix(branch, "hotfix/"):
		return "hotfix"
	case strings.HasPrefix(branch, "release/"):
		return "release"
	case strings.HasPrefix(branch, "experiment/"):
		return "experiment"
	case branch == "main" || branch == "master":
		return "main"
	case branch == "develop" || branch == "development":
		return "develop"
	default:
		return "other"
	}
}

func (bt *BranchTracker) compareWorkflows(branchA, branchB *BranchWorkflowInfo) *WorkflowComparison {
	comparison := &WorkflowComparison{
		BranchA: branchA.Branch,
		BranchB: branchB.Branch,
	}
	
	// Create maps for efficient comparison
	filesA := make(map[string]WorkflowFile)
	filesB := make(map[string]WorkflowFile)
	
	for _, file := range branchA.WorkflowFiles {
		filesA[file.Path] = file
	}
	
	for _, file := range branchB.WorkflowFiles {
		filesB[file.Path] = file
	}
	
	// Find added, deleted, modified, and unchanged files
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