package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/workflow"
	"github.com/stretchr/testify/require"
)

type fakeWorkflowClient struct {
	workflows []*workflow.Workflow
}

func (f fakeWorkflowClient) GetWorkflows(ctx context.Context) ([]*workflow.Workflow, error) {
	return f.workflows, nil
}

func (f fakeWorkflowClient) GetWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	for _, wf := range f.workflows {
		if wf.ID == id {
			return wf, nil
		}
	}
	return nil, client.ErrNotFound
}

func (f fakeWorkflowClient) HealthCheck(ctx context.Context) error {
	return nil
}

func TestCheckWorkflowSyncRealDetectsModifiedWorkflows(t *testing.T) {
	now := time.Now().Add(-2 * time.Hour)
	local := []WorkflowData{{ID: "1", Name: "Test", VersionID: 1}}
	remote := []*workflow.Workflow{{
		ID:        "1",
		Name:      "Test",
		VersionId: 2,
		UpdatedAt: now,
	}}

	result, err := checkWorkflowSyncReal("development", local, fakeWorkflowClient{workflows: remote})
	require.NoError(t, err)

	require.Equal(t, 1, result.TotalWorkflows)
	require.Equal(t, 1, result.Modified)
	require.Zero(t, result.Synchronized)
	require.Len(t, result.Workflows.Modified, 1)
	status := result.Workflows.Modified[0]
	require.Equal(t, "1", status.ID)
	require.Equal(t, "Test", status.Name)
	require.Equal(t, "modified", status.Status)
	require.Equal(t, 1, status.LocalVersion)
	require.Equal(t, 2, status.RemoteVersion)
	require.WithinDuration(t, now, status.LastModified, time.Second)
}

func TestCheckWorkflowSyncRealDetectsSynchronizedWorkflows(t *testing.T) {
	now := time.Now().Add(-30 * time.Minute)
	local := []WorkflowData{{ID: "2", Name: "Sync", VersionID: 3}}
	remote := []*workflow.Workflow{{
		ID:        "2",
		Name:      "Sync",
		VersionId: 3,
		UpdatedAt: now,
	}}

	result, err := checkWorkflowSyncReal("development", local, fakeWorkflowClient{workflows: remote})
	require.NoError(t, err)

	require.Equal(t, 1, result.TotalWorkflows)
	require.Equal(t, 1, result.Synchronized)
	require.Zero(t, result.Modified)
	require.Len(t, result.Workflows.Synchronized, 1)
	status := result.Workflows.Synchronized[0]
	require.Equal(t, "2", status.ID)
	require.Equal(t, "Sync", status.Name)
	require.Equal(t, "sync", status.Status)
	require.Equal(t, 3, status.LocalVersion)
	require.Equal(t, 3, status.RemoteVersion)
	require.WithinDuration(t, now, status.LastModified, time.Second)
}

func TestCheckWorkflowSyncRealDetectsMissingLocalWorkflows(t *testing.T) {
	now := time.Now().Add(-45 * time.Minute)
	remote := []*workflow.Workflow{{
		ID:        "3",
		Name:      "Remote Only",
		VersionId: 1,
		UpdatedAt: now,
	}}

	result, err := checkWorkflowSyncReal("development", nil, fakeWorkflowClient{workflows: remote})
	require.NoError(t, err)

	require.Equal(t, 1, result.TotalWorkflows)
	require.Equal(t, 1, result.Modified)
	require.Zero(t, result.Synchronized)
	require.Len(t, result.Workflows.Modified, 1)
	status := result.Workflows.Modified[0]
	require.Equal(t, "3", status.ID)
	require.Equal(t, "Remote Only", status.Name)
	require.Equal(t, "modified", status.Status)
	require.Equal(t, 0, status.LocalVersion)
	require.Equal(t, 1, status.RemoteVersion)
	require.WithinDuration(t, now, status.LastModified, time.Second)
}
