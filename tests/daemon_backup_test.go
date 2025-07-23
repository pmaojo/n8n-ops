package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackupStructure(t *testing.T) {
	tmpDir := t.TempDir()

	wf := TestWorkflow{
		ID:        "backup_test_001",
		Name:      "Backup Test Workflow",
		Active:    true,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(tmpDir, fmt.Sprintf("%s_%s_backup.json", wf.ID, timestamp))
	metaFile := filepath.Join(tmpDir, fmt.Sprintf("%s_%s_backup.meta.json", wf.ID, timestamp))

	data, err := json.MarshalIndent(wf, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(backupFile, data, 0644))

	meta := map[string]interface{}{
		"originalFile": fmt.Sprintf("./workflows/development/%s.json", wf.ID),
		"backupFile":   backupFile,
		"timestamp":    time.Now().Format(time.RFC3339),
		"environment":  "development",
		"workflowId":   wf.ID,
		"workflowName": wf.Name,
	}

	metaData, err := json.MarshalIndent(meta, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(metaFile, metaData, 0644))

	assert.FileExists(t, backupFile)
	assert.FileExists(t, metaFile)

	content, err := os.ReadFile(backupFile)
	require.NoError(t, err)

	var restored TestWorkflow
	require.NoError(t, json.Unmarshal(content, &restored))
	assert.Equal(t, wf.ID, restored.ID)
	assert.Equal(t, wf.Name, restored.Name)
}
