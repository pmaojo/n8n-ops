package git

import (
	"fmt"
	"strings"
)

var execImpl Executor = NewExecutor("")

// SetExecutor overrides the executor used for Git operations.
func SetExecutor(e Executor) {
	if e != nil {
		execImpl = e
	}
}

// GetCurrentBranch returns the current Git branch name
func GetCurrentBranch() (string, error) {
	branch, err := execImpl.CurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return branch, nil
}

// CheckoutBranch switches to the specified branch
func CheckoutBranch(branch string) error {
	return execImpl.Checkout(branch)
}

// GetBranchList returns all local branches
func GetBranchList() ([]string, error) {
	branches, err := execImpl.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branch list: %w", err)
	}
	return branches, nil
}

// IsGitRepository checks if current directory is a git repository
func IsGitRepository() bool {
	return execImpl.IsRepository()
}

// CommitChanges commits changes to the current branch
func CommitChanges(message string) error {
	if err := execImpl.Add("."); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}
	if err := execImpl.Commit(message); err != nil {
		return err
	}
	return nil
}

// PushBranch pushes the current branch to origin
func PushBranch(branch string) error {
	if err := execImpl.Push(branch); err != nil {
		return err
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
