package cmd

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubExecutor struct {
	git.MockExecutor
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestRunBranchCmd_ListAllJSON(t *testing.T) {
	branchListAll = true
	branchCompare = ""
	branchOutputJSON = true
	branchShowActive = false

	exec := &stubExecutor{}
	exec.RemoteBranchesFunc = func() ([]string, error) { return []string{"main"}, nil }
	exec.LogFunc = func(branch, format string) (string, error) { return "abc|msg|user|1", nil }
	exec.LsTreeFunc = func(branch, path string) ([]string, error) { return []string{"workflows/dev/test.json"}, nil }

	svc := git.NewService(exec)
	output := captureStdout(t, func() { require.NoError(t, runBranchCmd(nil, nil, svc)) })
	idx := strings.Index(output, "\n{")
	if idx == -1 {
		t.Fatalf("no json output: %s", output)
	}
	var out map[string]*BranchWorkflowInfo
	require.NoError(t, json.Unmarshal([]byte(output[idx+1:]), &out))
	assert.Len(t, out, 1)
	assert.NotNil(t, out["main"])
}

func TestHandleCurrentBranch_JSON(t *testing.T) {
	branchOutputJSON = true
	exec := &stubExecutor{}
	exec.CurrentBranchFunc = func() (string, error) { return "feat", nil }
	exec.LogFunc = func(branch, format string) (string, error) { return "abc|msg|user|1", nil }
	exec.LsTreeFunc = func(branch, path string) ([]string, error) { return []string{"workflows/dev/test.json"}, nil }
	svc := git.NewService(exec)
	output := captureStdout(t, func() { require.NoError(t, handleCurrentBranch(svc)) })
	idx := strings.Index(output, "\n{")
	require.NotEqual(t, -1, idx)
	var out BranchWorkflowInfo
	require.NoError(t, json.Unmarshal([]byte(output[idx+1:]), &out))
	assert.Equal(t, "feat", out.Branch)
}

func TestHandleCompareBranches_JSON(t *testing.T) {
	branchOutputJSON = true
	exec := &stubExecutor{}
	exec.CurrentBranchFunc = func() (string, error) { return "feat", nil }
	exec.LogFunc = func(branch, format string) (string, error) {
		if branch == "feat" {
			return "abc|msg|user|1", nil
		}
		return "def|msg2|user2|1", nil
	}
	exec.LsTreeFunc = func(branch, path string) ([]string, error) {
		if branch == "feat" {
			return []string{"workflows/dev/a.json"}, nil
		}
		return []string{"workflows/dev/b.json"}, nil
	}
	svc := git.NewService(exec)
	output := captureStdout(t, func() { require.NoError(t, handleCompareBranches(svc, "main")) })
	idx := strings.Index(output, "\n{")
	require.NotEqual(t, -1, idx)
	var out WorkflowComparison
	require.NoError(t, json.Unmarshal([]byte(output[idx+1:]), &out))
	assert.Equal(t, "feat", out.BranchA)
	assert.Equal(t, "main", out.BranchB)
}
