package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestWorkflow struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Active      bool                   `json:"active"`
	Nodes       []interface{}          `json:"nodes"`
	Connections map[string]interface{} `json:"connections"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	UpdatedAt   string                 `json:"updatedAt,omitempty"`
}

func waitForEvent(t *testing.T, watcher *fsnotify.Watcher, mask fsnotify.Op, timeout time.Duration) fsnotify.Event {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&mask != 0 {
				return ev
			}
		case err := <-watcher.Errors:
			t.Fatalf("watcher error: %v", err)
		case <-timer.C:
			t.Fatalf("timeout waiting for file event")
		}
	}
}

func TestDaemonFileWatching(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := fsnotify.NewWatcher()
	require.NoError(t, err)
	defer watcher.Close()

	require.NoError(t, watcher.Add(tmpDir))

	filePath := filepath.Join(tmpDir, "test-workflow.json")
	wf := TestWorkflow{ID: "watch_001", Name: "Watch Test"}
	data, err := json.Marshal(wf)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filePath, data, 0644))

	ev := waitForEvent(t, watcher, fsnotify.Create|fsnotify.Write, time.Second)
	assert.Equal(t, filePath, ev.Name)

	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	var readWF TestWorkflow
	require.NoError(t, json.Unmarshal(content, &readWF))
	assert.Equal(t, wf.Name, readWF.Name)
}

func TestDaemonMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := fsnotify.NewWatcher()
	require.NoError(t, err)
	defer watcher.Close()

	require.NoError(t, watcher.Add(tmpDir))

	count := 3
	for i := 0; i < count; i++ {
		wf := TestWorkflow{ID: fmt.Sprintf("multi_%d", i)}
		data, err := json.Marshal(wf)
		require.NoError(t, err)
		file := filepath.Join(tmpDir, fmt.Sprintf("multi-%d.json", i))
		require.NoError(t, os.WriteFile(file, data, 0644))
		waitForEvent(t, watcher, fsnotify.Create|fsnotify.Write, time.Second)
	}

	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)

	jsonFiles := 0
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			jsonFiles++
		}
	}

	assert.Equal(t, count, jsonFiles)
}
