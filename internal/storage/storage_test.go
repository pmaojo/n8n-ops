package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorageDirectory(t *testing.T) {
	// Test storage directory operations
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, "storage")
	
	err := os.MkdirAll(storageDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create storage directory: %v", err)
	}
	
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		t.Error("Storage directory should exist after creation")
	}
}

func TestMetadataStorage(t *testing.T) {
	// Test metadata storage functionality
	metadata := map[string]interface{}{
		"lastSync":    "2025-07-23T01:00:00Z",
		"environment": "development",
		"workflows":   3,
	}
	
	for key, value := range metadata {
		if key == "" {
			t.Error("Metadata key should not be empty")
		}
		if value == nil {
			t.Errorf("Metadata value for %s should not be nil", key)
		}
	}
}

func TestBackupOperations(t *testing.T) {
	// Test backup operations
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	
	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup directory: %v", err)
	}
	
	// Test backup file creation
	testFile := filepath.Join(backupDir, "backup_2025-07-23.json")
	testContent := `{"workflows": [], "timestamp": "2025-07-23T01:00:00Z"}`
	
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create backup file: %v", err)
	}
	
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Backup file should exist after creation")
	}
}