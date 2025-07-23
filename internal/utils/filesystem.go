package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/n8n-workflows/n8n-ops/internal/workflow"
)

// SanitizeFilename sanitizes a string to be safe for use as a filename
func SanitizeFilename(name string) string {
	// Replace invalid characters with underscores
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := reg.ReplaceAllString(name, "_")
	
	// Replace spaces with underscores
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	
	// Remove consecutive underscores
	reg2 := regexp.MustCompile(`_+`)
	sanitized = reg2.ReplaceAllString(sanitized, "_")
	
	// Trim underscores from start and end
	sanitized = strings.Trim(sanitized, "_")
	
	// Ensure it's not empty
	if sanitized == "" {
		sanitized = "workflow"
	}
	
	// Limit length to avoid filesystem issues
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}
	
	return sanitized
}

// WriteWorkflowToFile writes a workflow to a JSON file
func WriteWorkflowToFile(wf *workflow.Workflow, filepath string) error {
	// Create directory if it doesn't exist
	dir := filepath[:strings.LastIndex(filepath, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Marshal workflow to JSON with proper formatting
	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// LoadWorkflowFromFile loads a workflow from a JSON file
func LoadWorkflowFromFile(filepath string) (*workflow.Workflow, error) {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("workflow file does not exist: %s", filepath)
	}
	
	// Read file content
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}
	
	// Unmarshal JSON
	var wf workflow.Workflow
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}
	
	return &wf, nil
}

// WriteJSONFile writes any data structure to a JSON file
func WriteJSONFile(data interface{}, filepath string) error {
	// Create directory if it doesn't exist
	dir := filepath[:strings.LastIndex(filepath, "/")]
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	
	// Marshal to JSON with proper formatting
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// LoadJSONFile loads JSON data from a file into the provided interface
func LoadJSONFile(filepath string, v interface{}) error {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filepath)
	}
	
	// Read file content
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	// Unmarshal JSON
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	
	return nil
}

// HashWorkflow generates a hash of the workflow for change detection
func HashWorkflow(wf *workflow.Workflow) string {
	// Create a copy without sync metadata for consistent hashing
	wfCopy := *wf
	wfCopy.SyncMetadata = nil
	
	// Marshal to JSON
	data, err := json.Marshal(wfCopy)
	if err != nil {
		// Fallback to simple hash if marshaling fails
		return fmt.Sprintf("error_%d", time.Now().Unix())
	}
	
	// Generate SHA256 hash
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

// BackupFile creates a backup of a file with timestamp
func BackupFile(filepath string) error {
	// Check if original file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("original file does not exist: %s", filepath)
	}
	
	// Generate backup filename
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup_%s", filepath, timestamp)
	
	// Read original file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}
	
	// Write backup file
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}
	
	return nil
}

// EnsureDirectory ensures a directory exists, creating it if necessary
func EnsureDirectory(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// FileExists checks if a file exists
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// DirectoryExists checks if a directory exists
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.IsDir()
}

// GetFileModTime gets the modification time of a file
func GetFileModTime(filepath string) (time.Time, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// GetFileSize gets the size of a file in bytes
func GetFileSize(filepath string) (int64, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ListWorkflowFiles lists all workflow JSON files in a directory
func ListWorkflowFiles(dir string, recursive bool) ([]string, error) {
	var files []string
	
	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if !info.IsDir() && isWorkflowFile(path) {
				files = append(files, path)
			}
			
			return nil
		})
		return files, err
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		
		for _, entry := range entries {
			if !entry.IsDir() {
				fullPath := filepath.Join(dir, entry.Name())
				if isWorkflowFile(fullPath) {
					files = append(files, fullPath)
				}
			}
		}
		
		return files, nil
	}
}

// isWorkflowFile checks if a file is likely a workflow JSON file
func isWorkflowFile(path string) bool {
	// Must be a JSON file
	if !strings.HasSuffix(strings.ToLower(path), ".json") {
		return false
	}
	
	// Exclude metadata files and other non-workflow files
	basename := filepath.Base(path)
	if strings.HasPrefix(basename, "_") || strings.HasPrefix(basename, ".") {
		return false
	}
	
	excludedPrefixes := []string{
		"deployment-report",
		"sync-metadata",
		"backup",
		"temp",
	}
	
	for _, prefix := range excludedPrefixes {
		if strings.HasPrefix(strings.ToLower(basename), prefix) {
			return false
		}
	}
	
	return true
}

// CleanupOldBackups removes backup files older than the specified duration
func CleanupOldBackups(dir string, maxAge time.Duration) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Check if it's a backup file
		if strings.Contains(info.Name(), ".backup_") {
			// Check age
			if time.Since(info.ModTime()) > maxAge {
				if err := os.Remove(path); err != nil {
					return fmt.Errorf("failed to remove old backup %s: %w", path, err)
				}
			}
		}
		
		return nil
	})
}

// CreateTemporaryFile creates a temporary file with the given prefix
func CreateTemporaryFile(prefix string, content []byte) (string, error) {
	tempFile, err := os.CreateTemp("", prefix+"*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()
	
	if _, err := tempFile.Write(content); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	
	return tempFile.Name(), nil
}

// CopyFile copies a file from source to destination
func CopyFile(src, dst string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	
	// Create destination directory if needed
	dstDir := filepath.Dir(dst)
	if err := EnsureDirectory(dstDir); err != nil {
		return err
	}
	
	// Write destination file
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}
	
	return nil
}

// MoveFile moves a file from source to destination
func MoveFile(src, dst string) error {
	// Try rename first (fast if on same filesystem)
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	
	// Fallback to copy and delete
	if err := CopyFile(src, dst); err != nil {
		return err
	}
	
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to remove source file after copy: %w", err)
	}
	
	return nil
}

// GetRelativePath returns the relative path from base to target
func GetRelativePath(base, target string) (string, error) {
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}
	
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	
	return filepath.Rel(absBase, absTarget)
}

// ValidateFilePath validates that a file path is safe and within expected boundaries
func ValidateFilePath(path string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	
	// Check for path traversal attempts
	if strings.Contains(absPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}
	
	// Additional security checks can be added here
	return nil
}
