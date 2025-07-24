package git

import "fmt"

// Service provides Git operations backed by an Executor.
type Service struct {
	executor Executor
}

// NewService constructs a Service using the given Executor. A default
// executor is used when nil is provided.
func NewService(exec Executor) *Service {
	if exec == nil {
		exec = NewExecutor("")
	}
	return &Service{executor: exec}
}

// GetCurrentBranch returns the current Git branch name.
func (s *Service) GetCurrentBranch() (string, error) {
	branch, err := s.executor.CurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return branch, nil
}

// CheckoutBranch switches to the specified branch.
func (s *Service) CheckoutBranch(branch string) error {
	return s.executor.Checkout(branch)
}

// GetBranchList returns all local branches.
func (s *Service) GetBranchList() ([]string, error) {
	branches, err := s.executor.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branch list: %w", err)
	}
	return branches, nil
}

// IsGitRepository checks if current directory is a git repository.
func (s *Service) IsGitRepository() bool {
	return s.executor.IsRepository()
}

// CommitChanges commits changes to the current branch.
func (s *Service) CommitChanges(message string) error {
	if err := s.executor.Add("."); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}
	if err := s.executor.Commit(message); err != nil {
		return err
	}
	return nil
}

// PushBranch pushes the current branch to origin.
func (s *Service) PushBranch(branch string) error {
	if err := s.executor.Push(branch); err != nil {
		return err
	}
	return nil
}

// Log proxies to the underlying executor's Log method.
func (s *Service) Log(branch, format string) (string, error) {
	return s.executor.Log(branch, format)
}

// LsTree proxies to the underlying executor's LsTree method.
func (s *Service) LsTree(branch, path string) ([]string, error) {
	return s.executor.LsTree(branch, path)
}

// RemoteBranches proxies to the underlying executor's RemoteBranches method.
func (s *Service) RemoteBranches() ([]string, error) {
	return s.executor.RemoteBranches()
}

// UserEmail proxies to the underlying executor's UserEmail method.
func (s *Service) UserEmail() (string, error) {
	return s.executor.UserEmail()
}
