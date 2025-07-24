package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/workflow"
	"github.com/sirupsen/logrus"
)

// BackupInfo stores backup metadata
type BackupInfo struct {
	OriginalFile string    `json:"originalFile"`
	BackupFile   string    `json:"backupFile"`
	Timestamp    time.Time `json:"timestamp"`
	Environment  string    `json:"environment"`
	WorkflowID   string    `json:"workflowId"`
	WorkflowName string    `json:"workflowName"`
}

// fileWatcher abstracts fsnotify.Watcher for easier testing.
type fileWatcher interface {
	Add(name string) error
	Close() error
	Events() <-chan fsnotify.Event
	Errors() <-chan error
}

type fsnotifyWatcher struct{ *fsnotify.Watcher }

func (w *fsnotifyWatcher) Events() <-chan fsnotify.Event { return w.Watcher.Events }
func (w *fsnotifyWatcher) Errors() <-chan error          { return w.Watcher.Errors }

var daemonWatcherFactory = func() (fileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &fsnotifyWatcher{w}, nil
}

var daemonClientFactory = func(url string) (client.Client, error) {
	return client.New(url, "n8n_api_mock_development", nil)
}

// runDaemonMode starts the file watcher and synchronizes workflow changes with n8n.
// The environment parameter determines which environment configuration to use.
func runDaemonMode(env string) {
	runDaemonModeCtx(context.Background(), env)
}

// runDaemonModeCtx allows injecting a context for testing.
func runDaemonModeCtx(ctx context.Context, env string) {
	logEntry := logger.WithFields(logrus.Fields{
		"command": "daemon",
		"env":     env,
	})

	logEntry.Info("Starting n8n-ops daemon mode")

	if language == "es" {
		fmt.Printf("🤖 Modo daemon iniciado - %s\n", env)
		fmt.Printf("👁️ Monitoreando archivos JSON en ./workflows/%s/\n", env)
		fmt.Printf("💾 Creando backups automáticos antes de actualizar workflows\n")
	} else {
		fmt.Printf("🤖 Daemon mode started - %s environment\n", env)
		fmt.Printf("👁️ Watching JSON files in ./workflows/%s/\n", env)
		fmt.Printf("💾 Creating automatic backups before updating workflows\n")
	}

	// Create n8n client with proper URL for demo mode
	var n8nURL string
	if demoMode {
		n8nURL = "http://localhost:3001"
	} else {
		// Use environment-specific URL from config
		switch env {
		case "development":
			n8nURL = "http://localhost:5678"
		case "staging":
			n8nURL = "https://n8n-staging.example.com"
		case "production":
			n8nURL = "https://n8n-prod.example.com"
		default:
			n8nURL = "http://localhost:5678"
		}
	}

	n8nClient, err := daemonClientFactory(n8nURL)
	if err != nil {
		logEntry.WithError(err).Fatal("Failed to create n8n client")
		return
	}

	// Test connection
	if err := testN8nConnectionDaemon(n8nClient); err != nil {
		logEntry.WithError(err).Fatal("Failed to connect to n8n API")
		return
	}

	// Setup file watcher
	watcher, err := daemonWatcherFactory()
	if err != nil {
		logEntry.WithError(err).Fatal("Failed to create file watcher")
		return
	}
	defer watcher.Close()

	// Watch directory
	watchDir := fmt.Sprintf("./workflows/%s", env)
	if err := setupDirectoryWatch(watcher, watchDir); err != nil {
		logEntry.WithError(err).Fatal("Failed to setup directory watch")
		return
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("✅ Connected to n8n API. Daemon ready!\n")
	fmt.Printf("⏹️ Press Ctrl+C to stop daemon\n\n")

	// Main daemon loop
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events():
			if !ok {
				return
			}
			handleFileEvent(event, n8nClient, logEntry, env)

		case err, ok := <-watcher.Errors():
			if !ok {
				return
			}
			logEntry.WithError(err).Error("File watcher error")

		case sig := <-sigChan:
			fmt.Printf("\n🛑 Received signal %v, stopping daemon...\n", sig)
			return
		}
	}
}

func setupDirectoryWatch(watcher fileWatcher, watchDir string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(watchDir, 0755); err != nil {
		return fmt.Errorf("failed to create watch directory: %w", err)
	}

	// Add directory to watcher
	if err := watcher.Add(watchDir); err != nil {
		return fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	return nil
}

// handleFileEvent processes file system events produced by the watcher.
// The env parameter specifies the target environment for workflow operations.
func handleFileEvent(event fsnotify.Event, n8nClient client.Client, logEntry *logrus.Entry, env string) {
	// Only process JSON files
	if !strings.HasSuffix(event.Name, ".json") {
		return
	}

	// Only process write events (ignore create/remove for now)
	if event.Op&fsnotify.Write == 0 {
		return
	}

	logEntry.WithField("file", event.Name).Info("JSON file modified")
	fmt.Printf("📝 File changed: %s\n", filepath.Base(event.Name))

	// Process the file change
	if err := processJSONFileChange(event.Name, n8nClient, logEntry, env); err != nil {
		logEntry.WithError(err).Error("Failed to process JSON file change")
		fmt.Printf("❌ Error processing %s: %v\n", filepath.Base(event.Name), err)
	}
}

// processJSONFileChange reads the modified JSON workflow and updates it in n8n.
// The env parameter is used when creating backups of the current workflow.
func processJSONFileChange(filePath string, n8nClient client.Client, logEntry *logrus.Entry, env string) error {
	// Read and parse JSON file
	workflowData, err := readJSONWorkflow(filePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON workflow: %w", err)
	}

	// Validate workflow
	workflowID, ok := workflowData["id"].(string)
	if !ok || workflowID == "" {
		return fmt.Errorf("workflow ID is required")
	}

	// Create backup of existing workflow
	if err := createWorkflowBackup(workflowID, n8nClient, env); err != nil {
		logEntry.WithError(err).Warn("Failed to create workflow backup")
	}

	// Update workflow in n8n
	if err := updateWorkflowInN8n(workflowData, n8nClient, logEntry); err != nil {
		return fmt.Errorf("failed to update workflow in n8n: %w", err)
	}

	workflowName := workflowData["name"].(string)
	fmt.Printf("✅ Workflow '%s' updated in n8n\n", workflowName)
	return nil
}

func readJSONWorkflow(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var workflow map[string]interface{}
	if err := json.Unmarshal(data, &workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// createWorkflowBackup saves the current workflow state from n8n to a file.
// The env parameter determines the backup directory and metadata information.
func createWorkflowBackup(workflowID string, n8nClient client.Client, env string) error {
	// Get current workflow from n8n
	ctx := context.Background()
	currentWorkflow, err := n8nClient.GetWorkflow(ctx, workflowID)
	if err != nil {
		// If workflow doesn't exist, no backup needed
		return nil
	}

	// Create backups directory
	backupDir := fmt.Sprintf("./backups/%s", env)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupFile := fmt.Sprintf("%s/%s_%s_backup.json", backupDir, workflowID, timestamp)

	// Save current workflow as backup
	backupData, err := json.MarshalIndent(currentWorkflow, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(backupFile, backupData, 0644); err != nil {
		return err
	}

	// Save backup metadata
	backupInfo := BackupInfo{
		OriginalFile: fmt.Sprintf("./workflows/%s/%s.json", env, workflowID),
		BackupFile:   backupFile,
		Timestamp:    time.Now(),
		Environment:  env,
		WorkflowID:   workflowID,
		WorkflowName: currentWorkflow.Name,
	}

	metadataFile := fmt.Sprintf("%s/%s_%s_backup.meta.json", backupDir, workflowID, timestamp)
	metadataData, err := json.MarshalIndent(backupInfo, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(metadataFile, metadataData, 0644); err != nil {
		return err
	}

	fmt.Printf("💾 Backup created: %s\n", filepath.Base(backupFile))
	return nil
}

func updateWorkflowInN8n(jsonWorkflow map[string]interface{}, n8nClient client.Client, logEntry *logrus.Entry) error {
	workflowID := jsonWorkflow["id"].(string)

	// Convert map to workflow struct
	jsonData, err := json.Marshal(jsonWorkflow)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	var wf workflow.Workflow
	if err := json.Unmarshal(jsonData, &wf); err != nil {
		return fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	// Set updated timestamp in original JSON if not provided
	if _, exists := jsonWorkflow["updatedAt"]; !exists || jsonWorkflow["updatedAt"] == "" {
		jsonWorkflow["updatedAt"] = time.Now().Format(time.RFC3339)
	}

	// Check if workflow exists
	ctx := context.Background()
	_, err = n8nClient.GetWorkflow(ctx, workflowID)
	if err != nil {
		// Workflow doesn't exist, create new one
		logEntry.WithField("workflowId", workflowID).Info("Creating new workflow")
		_, err = n8nClient.CreateWorkflow(ctx, &wf)
		return err
	}

	// Workflow exists, update it
	logEntry.WithField("workflowId", workflowID).Info("Updating existing workflow")
	_, err = n8nClient.UpdateWorkflow(ctx, workflowID, &wf)
	return err
}

func testN8nConnectionDaemon(n8nClient client.Client) error {
	ctx := context.Background()
	err := n8nClient.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	fmt.Printf("🔗 Connected to n8n API\n")
	return nil
}
