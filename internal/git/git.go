package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetCurrentBranch returns the current Git branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	
	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "main", nil // fallback to main if no branch detected
	}
	
	return branch, nil
}

// CheckoutBranch switches to the specified branch
func CheckoutBranch(branch string) error {
	// Check if branch exists
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	if err := cmd.Run(); err != nil {
		// Branch doesn't exist, create it
		cmd = exec.Command("git", "checkout", "-b", branch)
	} else {
		// Branch exists, switch to it
		cmd = exec.Command("git", "checkout", branch)
	}
	
	return cmd.Run()
}

// GetBranchList returns all local branches
func GetBranchList() ([]string, error) {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get branch list: %w", err)
	}
	
	var branches []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Remove the * indicator for current branch
		if strings.HasPrefix(line, "* ") {
			line = strings.TrimPrefix(line, "* ")
		}
		branches = append(branches, line)
	}
	
	return branches, nil
}

// IsGitRepository checks if current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// CommitChanges commits changes to the current branch
func CommitChanges(message string) error {
	// Add all changes
	addCmd := exec.Command("git", "add", ".")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}
	
	// Commit changes
	commitCmd := exec.Command("git", "commit", "-m", message)
	return commitCmd.Run()
}

// PushBranch pushes the current branch to origin
func PushBranch(branch string) error {
	cmd := exec.Command("git", "push", "origin", branch)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push branch %s: %w\nOutput: %s", branch, err, stderr.String())
	}
	
	return nil
}

// GetEnvironmentFromBranch maps git branches to n8n environments
func GetEnvironmentFromBranch(branch string) string {
	switch {
	case strings.Contains(branch, "dev") || strings.Contains(branch, "develop"):
		return "development"
	case strings.Contains(branch, "staging") || strings.Contains(branch, "stage"):
		return "staging"
	case strings.Contains(branch, "main") || strings.Contains(branch, "master") || strings.Contains(branch, "prod"):
		return "production"
	default:
		return "development" // default to development for feature branches
	}
}

// GetBranchFromEnvironment maps n8n environments to suggested git branches
func GetBranchFromEnvironment(env string) string {
	switch env {
	case "development":
		return "develop"
	case "staging":
		return "staging"
	case "production":
		return "main"
	default:
		return "develop"
	}
}