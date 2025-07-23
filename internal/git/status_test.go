package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
}

func initRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init", "-b", "main")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
}

func writeWorkflow(t *testing.T, dir, relative, content string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(relative))
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func TestGetUncommittedWorkflowSummary_NoChanges(t *testing.T) {
	repo := t.TempDir()
	initRepo(t, repo)

	writeWorkflow(t, repo, "workflows/development/test.json", `{}`)
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-m", "init")

	checker := NewGitStatusChecker(repo, NewExecutor(repo))
	summary, err := checker.GetUncommittedWorkflowSummary()
	require.NoError(t, err)
	assert.Equal(t, "All workflows are committed", summary)
}

func TestGetStatusDetectsWorkflowChanges(t *testing.T) {
	repo := t.TempDir()
	initRepo(t, repo)

	writeWorkflow(t, repo, "workflows/development/original.json", `{"name":"orig"}`)
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-m", "initial")

	writeWorkflow(t, repo, "workflows/development/original.json", `{"name":"changed"}`)
	writeWorkflow(t, repo, "workflows/development/new.json", `{"name":"new"}`)

	checker := NewGitStatusChecker(repo, NewExecutor(repo))
	status, err := checker.GetStatus()
	require.NoError(t, err)

	assert.True(t, status.HasChanges)
	assert.Contains(t, status.ModifiedFiles, "workflows/development/original.json")
	assert.Contains(t, status.UntrackedFiles, "workflows/development/new.json")
	assert.Len(t, status.UncommittedWorkflows, 2)

	var modified, untracked bool
	for _, wf := range status.UncommittedWorkflows {
		switch {
		case wf.FilePath == "workflows/development/original.json" && wf.Status == "modified":
			modified = true
		case wf.FilePath == "workflows/development/new.json" && wf.Status == "untracked":
			untracked = true
		}
	}
	assert.True(t, modified)
	assert.True(t, untracked)

	summary, err := checker.GetUncommittedWorkflowSummary()
	require.NoError(t, err)
	assert.Contains(t, summary, "2 uncommitted workflow changes")
	assert.Contains(t, summary, "📝 Original (modified)")
	assert.Contains(t, summary, "❓ New (untracked)")
}
