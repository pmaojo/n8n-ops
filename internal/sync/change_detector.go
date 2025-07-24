package sync

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pmaojo/n8n-ops/internal/workflow"
)

// ChangeType represents the classification of a detected change.
// Values include "new", "modified", "deleted", "conflict" and "none".
// It is exported to allow other packages to reason about change results.
type ChangeType string

const (
	ChangeTypeNew      ChangeType = "new"
	ChangeTypeModified ChangeType = "modified"
	ChangeTypeDeleted  ChangeType = "deleted"
	ChangeTypeConflict ChangeType = "conflict"
	ChangeTypeNone     ChangeType = "none"
)

// Change describes the difference between a local workflow file and the
// corresponding workflow retrieved from n8n. The Direction field indicates
// where the most recent change occurred.
type Change struct {
	Type       ChangeType `json:"type"`
	WorkflowID string     `json:"workflowId"`
	Name       string     `json:"name"`
	LocalFile  string     `json:"localFile"`
	LocalTime  time.Time  `json:"localTime"`
	RemoteTime time.Time  `json:"remoteTime"`
	LocalHash  string     `json:"localHash"`
	RemoteHash string     `json:"remoteHash"`
	Direction  string     `json:"direction"` // "to_n8n", "from_n8n", "bidirectional"
}

// ChangeDetector inspects workflow files and their remote counterparts in n8n
// to determine if any differences exist. It contains the local workflows
// directory and the target environment name used for reporting.
type ChangeDetector struct {
	WorkflowsDir string
	Environment  string
}

// NewChangeDetector returns a ChangeDetector configured with the given
// workflows directory and environment name. It performs no IO, allowing the
// caller to manage dependencies explicitly.
func NewChangeDetector(workflowsDir, environment string) *ChangeDetector {
	return &ChangeDetector{
		WorkflowsDir: workflowsDir,
		Environment:  environment,
	}
}

// DetectChanges compares local workflows with remote n8n workflows
func (cd *ChangeDetector) DetectChanges(remoteWorkflows []workflow.Workflow) ([]Change, error) {
	var changes []Change

	// Get local workflows
	localWorkflows, err := cd.getLocalWorkflows()
	if err != nil {
		return nil, fmt.Errorf("failed to get local workflows: %w", err)
	}

	// Create maps for easier comparison
	localMap := make(map[string]LocalWorkflow)
	remoteMap := make(map[string]workflow.Workflow)

	for _, lw := range localWorkflows {
		localMap[lw.ID] = lw
	}

	for _, rw := range remoteWorkflows {
		remoteMap[rw.ID] = rw
	}

	// Check for new or modified remote workflows
	for _, remote := range remoteWorkflows {
		local, existsLocal := localMap[remote.ID]

		if !existsLocal {
			// New workflow in n8n (not in local)
			changes = append(changes, Change{
				Type:       ChangeTypeNew,
				WorkflowID: remote.ID,
				Name:       remote.Name,
				RemoteTime: remote.UpdatedAt,
				Direction:  "from_n8n",
			})
		} else {
			// Compare timestamps and content
			change := cd.compareWorkflows(local, remote)
			if change.Type != ChangeTypeNone {
				changes = append(changes, change)
			}
		}
	}

	// Check for deleted remote workflows (exist locally but not in n8n)
	for _, local := range localWorkflows {
		if _, existsRemote := remoteMap[local.ID]; !existsRemote {
			changes = append(changes, Change{
				Type:       ChangeTypeDeleted,
				WorkflowID: local.ID,
				Name:       local.Name,
				LocalFile:  local.FilePath,
				LocalTime:  local.ModifiedAt,
				Direction:  "from_n8n",
			})
		}
	}

	return changes, nil
}

// LocalWorkflow represents a workflow loaded from the filesystem. The Content
// field holds the raw JSON payload used to compute a deterministic hash.
type LocalWorkflow struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	FilePath   string    `json:"filePath"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Hash       string    `json:"hash"`
	Content    []byte    `json:"-"`
}

func (cd *ChangeDetector) getLocalWorkflows() ([]LocalWorkflow, error) {
	var workflows []LocalWorkflow

	if _, err := os.Stat(cd.WorkflowsDir); os.IsNotExist(err) {
		return workflows, nil // No workflows directory yet
	}

	err := filepath.Walk(cd.WorkflowsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			workflow, err := cd.parseLocalWorkflow(path, info)
			if err != nil {
				return fmt.Errorf("failed to parse workflow %s: %w", path, err)
			}
			workflows = append(workflows, workflow)
		}

		return nil
	})

	return workflows, err
}

func (cd *ChangeDetector) parseLocalWorkflow(filePath string, info os.FileInfo) (LocalWorkflow, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return LocalWorkflow{}, err
	}

	var workflow workflow.Workflow
	if err := json.Unmarshal(content, &workflow); err != nil {
		return LocalWorkflow{}, err
	}

	// Calculate content hash
	hash := fmt.Sprintf("%x", md5.Sum(content))

	return LocalWorkflow{
		ID:         workflow.ID,
		Name:       workflow.Name,
		FilePath:   filePath,
		ModifiedAt: info.ModTime(),
		Hash:       hash,
		Content:    content,
	}, nil
}

func (cd *ChangeDetector) compareWorkflows(local LocalWorkflow, remote workflow.Workflow) Change {
	// Calculate remote workflow hash
	remoteContent, err := json.Marshal(remote)
	if err != nil {
		return Change{Type: ChangeTypeNone}
	}
	remoteHash := fmt.Sprintf("%x", md5.Sum(remoteContent))

	// Compare content hashes
	if local.Hash == remoteHash {
		return Change{Type: ChangeTypeNone}
	}

	// Content differs - check timestamps to determine direction
	timeDiff := remote.UpdatedAt.Sub(local.ModifiedAt)

	if timeDiff > 30*time.Second {
		// Remote is significantly newer
		return Change{
			Type:       ChangeTypeModified,
			WorkflowID: remote.ID,
			Name:       remote.Name,
			LocalFile:  local.FilePath,
			LocalTime:  local.ModifiedAt,
			RemoteTime: remote.UpdatedAt,
			LocalHash:  local.Hash,
			RemoteHash: remoteHash,
			Direction:  "from_n8n",
		}
	} else if timeDiff < -30*time.Second {
		// Local is significantly newer
		return Change{
			Type:       ChangeTypeModified,
			WorkflowID: remote.ID,
			Name:       remote.Name,
			LocalFile:  local.FilePath,
			LocalTime:  local.ModifiedAt,
			RemoteTime: remote.UpdatedAt,
			LocalHash:  local.Hash,
			RemoteHash: remoteHash,
			Direction:  "to_n8n",
		}
	} else {
		// Timestamps are close - potential conflict
		return Change{
			Type:       ChangeTypeConflict,
			WorkflowID: remote.ID,
			Name:       remote.Name,
			LocalFile:  local.FilePath,
			LocalTime:  local.ModifiedAt,
			RemoteTime: remote.UpdatedAt,
			LocalHash:  local.Hash,
			RemoteHash: remoteHash,
			Direction:  "bidirectional",
		}
	}
}

// SaveChangeReport saves the detected changes to a JSON file for CI/CD integration
func (cd *ChangeDetector) SaveChangeReport(changes []Change, reportPath string) error {
	report := map[string]interface{}{
		"timestamp":    time.Now(),
		"environment":  cd.Environment,
		"totalChanges": len(changes),
		"changes":      changes,
		"summary": map[string]int{
			"new":       countChangesByType(changes, ChangeTypeNew),
			"modified":  countChangesByType(changes, ChangeTypeModified),
			"deleted":   countChangesByType(changes, ChangeTypeDeleted),
			"conflicts": countChangesByType(changes, ChangeTypeConflict),
		},
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(reportPath, data, 0644)
}

func countChangesByType(changes []Change, changeType ChangeType) int {
	count := 0
	for _, change := range changes {
		if change.Type == changeType {
			count++
		}
	}
	return count
}
