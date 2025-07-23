package cmd

import (
	"testing"

	"github.com/pmaojo/n8n-ops/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCurrentBranchUsesExecutor(t *testing.T) {
	called := false
	exec := &git.MockExecutor{
		CurrentBranchFunc: func() (string, error) {
			called = true
			return "main", nil
		},
	}

	branch, err := getCurrentBranch(exec)
	require.NoError(t, err)
	assert.Equal(t, "main", branch)
	assert.True(t, called)
}

func TestGetBranchWorkflows(t *testing.T) {
	exec := &git.MockExecutor{
		LogFunc: func(branch, format string) (string, error) {
			return "abc123|initial commit|alice|0", nil
		},
		LsTreeFunc: func(branch, path string) ([]string, error) {
			return []string{"workflows/dev/wf.json"}, nil
		},
	}

	info, err := getBranchWorkflows(exec, "main")
	require.NoError(t, err)
	require.Len(t, info.WorkflowFiles, 1)
	assert.Equal(t, "abc123", info.LastCommit)
	assert.Equal(t, "main", info.Branch)
}
