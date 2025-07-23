package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Executor abstracts git command execution.
type Executor interface {
	CurrentBranch() (string, error)
	Checkout(branch string) error
	Branches() ([]string, error)
	Add(paths ...string) error
	Commit(message string) error
	Push(branch string) error
	StatusPorcelain() (string, error)
	LastCommit() (string, error)
	IsRepository() bool
	Log(branch, format string) (string, error)
	LsTree(branch, path string) ([]string, error)
	RemoteBranches() ([]string, error)
	UserEmail() (string, error)
}

type execExecutor struct {
	workDir string
}

// NewExecutor returns an Executor that runs git commands via os/exec.
func NewExecutor(workDir string) Executor {
	return &execExecutor{workDir: workDir}
}

func (e *execExecutor) gitCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	if e.workDir != "" {
		cmd.Dir = e.workDir
	}
	return cmd
}

func (e *execExecutor) CurrentBranch() (string, error) {
	out, err := e.gitCmd("rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("get current branch: %w", err)
	}
	branch := strings.TrimSpace(string(out))
	if branch == "" {
		return "main", nil
	}
	return branch, nil
}

func (e *execExecutor) Checkout(branch string) error {
	if err := e.gitCmd("show-ref", "--verify", "--quiet", "refs/heads/"+branch).Run(); err != nil {
		return e.gitCmd("checkout", "-b", branch).Run()
	}
	return e.gitCmd("checkout", branch).Run()
}

func (e *execExecutor) Branches() ([]string, error) {
	out, err := e.gitCmd("branch").Output()
	if err != nil {
		return nil, fmt.Errorf("list branches: %w", err)
	}
	var branches []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "* ") {
			line = strings.TrimPrefix(line, "* ")
		}
		branches = append(branches, line)
	}
	return branches, nil
}

func (e *execExecutor) Add(paths ...string) error {
	args := append([]string{"add"}, paths...)
	if err := e.gitCmd(args...).Run(); err != nil {
		return fmt.Errorf("add files: %w", err)
	}
	return nil
}

func (e *execExecutor) Commit(message string) error {
	if err := e.gitCmd("commit", "-m", message).Run(); err != nil {
		return fmt.Errorf("commit changes: %w", err)
	}
	return nil
}

func (e *execExecutor) Push(branch string) error {
	var stderr bytes.Buffer
	cmd := e.gitCmd("push", "origin", branch)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("push branch %s: %w\n%s", branch, err, stderr.String())
	}
	return nil
}

func (e *execExecutor) StatusPorcelain() (string, error) {
	out, err := e.gitCmd("status", "--porcelain").Output()
	if err != nil {
		return "", fmt.Errorf("git status: %w", err)
	}
	return string(out), nil
}

func (e *execExecutor) LastCommit() (string, error) {
	out, err := e.gitCmd("log", "-1", "--pretty=format:%h %s").Output()
	if err != nil {
		return "", fmt.Errorf("last commit: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (e *execExecutor) Log(branch, format string) (string, error) {
	if format == "" {
		format = "%H|%s|%an|%ct"
	}
	out, err := e.gitCmd("log", "-1", "--format="+format, branch).Output()
	if err != nil {
		return "", fmt.Errorf("git log: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (e *execExecutor) LsTree(branch, path string) ([]string, error) {
	args := []string{"ls-tree", "-r", "--name-only", branch}
	if path != "" {
		args = append(args, path)
	}
	out, err := e.gitCmd(args...).Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-tree: %w", err)
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

func (e *execExecutor) RemoteBranches() ([]string, error) {
	out, err := e.gitCmd("branch", "-r", "--format=%(refname:short)").Output()
	if err != nil {
		return nil, fmt.Errorf("list remote branches: %w", err)
	}
	var branches []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "HEAD") {
			continue
		}
		if strings.HasPrefix(line, "origin/") {
			line = strings.TrimPrefix(line, "origin/")
		}
		branches = append(branches, line)
	}
	return branches, nil
}

func (e *execExecutor) IsRepository() bool {
	return e.gitCmd("rev-parse", "--git-dir").Run() == nil
}

func (e *execExecutor) UserEmail() (string, error) {
	out, err := e.gitCmd("config", "user.email").Output()
	if err != nil {
		return "", fmt.Errorf("get git user email: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
