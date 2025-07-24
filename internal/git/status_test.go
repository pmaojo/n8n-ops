package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testExecutor embeds MockExecutor and records git commands for assertions.
// It avoids filesystem side effects enabling isolated tests.
type testExecutor struct {
	*MockExecutor
	StatusOutput string
	Added        []string
	Commits      []string
	Last         string
}

func (te *testExecutor) StatusPorcelain() (string, error) { return te.StatusOutput, nil }
func (te *testExecutor) Add(paths ...string) error {
	te.Added = append(te.Added, paths...)
	return nil
}
func (te *testExecutor) Commit(msg string) error {
	te.Commits = append(te.Commits, msg)
	return nil
}
func (te *testExecutor) LastCommit() (string, error) { return te.Last, nil }

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

func TestGetStatusWithMockExecutor(t *testing.T) {
	exec := &testExecutor{
		MockExecutor: &MockExecutor{CurrentBranchFunc: func() (string, error) { return "feature/test", nil }},
		StatusOutput: " M workflows/development/mod.json\n?? workflows/staging/new.json\nA  workflows/production/add.json\n?? README.md\n",
		Last:         "abc123 init",
	}

	checker := NewGitStatusChecker("", exec)

	status, err := checker.GetStatus()
	require.NoError(t, err)
	assert.Equal(t, "feature/test", status.CurrentBranch)
	assert.Equal(t, "abc123 init", status.LastCommit)
	assert.True(t, status.HasChanges)
	assert.Contains(t, status.ModifiedFiles, "workflows/development/mod.json")
	assert.Contains(t, status.UntrackedFiles, "README.md")
	assert.Contains(t, status.StagedFiles, "workflows/production/add.json")
	require.Len(t, status.UncommittedWorkflows, 3)

	var foundMod, foundNew, foundAdd bool
	for _, wf := range status.UncommittedWorkflows {
		switch wf.FilePath {
		case "workflows/development/mod.json":
			foundMod = wf.Status == "modified" && wf.Environment == "development"
		case "workflows/staging/new.json":
			foundNew = wf.Status == "untracked" && wf.Environment == "staging"
		case "workflows/production/add.json":
			foundAdd = wf.Status == "added" && wf.Environment == "production"
		}
	}
	assert.True(t, foundMod)
	assert.True(t, foundNew)
	assert.True(t, foundAdd)
}

func TestGetUncommittedWorkflowSummaryWithMockExecutor(t *testing.T) {
	exec := &testExecutor{
		MockExecutor: &MockExecutor{CurrentBranchFunc: func() (string, error) { return "feat", nil }},
		StatusOutput: " M workflows/development/a.json\n?? workflows/development/b.json\n",
	}

	checker := NewGitStatusChecker("", exec)
	summary, err := checker.GetUncommittedWorkflowSummary()
	require.NoError(t, err)
	assert.Contains(t, summary, "2 uncommitted workflow changes")
	assert.Contains(t, summary, "📝 A (modified)")
	assert.Contains(t, summary, "❓ B (untracked)")
}

func TestCheckBeforeSyncWarns(t *testing.T) {
	exec := &testExecutor{
		MockExecutor: &MockExecutor{CurrentBranchFunc: func() (string, error) { return "dev", nil }},
		StatusOutput: " M workflows/development/mod.json\n",
	}
	checker := NewGitStatusChecker("", exec)
	out, err := checker.CheckBeforeSync()
	assert.Error(t, err)
	assert.Contains(t, out, "uncommitted workflow changes")
}

func TestAutoCommitWorkflows(t *testing.T) {
	exec := &testExecutor{
		MockExecutor: &MockExecutor{CurrentBranchFunc: func() (string, error) { return "dev", nil }},
		StatusOutput: " M workflows/development/mod.json\n?? workflows/development/new.json\n",
	}
	checker := NewGitStatusChecker("", exec)
	msg, err := checker.AutoCommitWorkflows(context.Background(), "commit msg")
	require.NoError(t, err)
	assert.Contains(t, msg, "Auto-committed")
	assert.ElementsMatch(t, []string{"workflows/development/mod.json", "workflows/development/new.json"}, exec.Added)
	require.Len(t, exec.Commits, 1)
	assert.Equal(t, "commit msg", exec.Commits[0])
}
