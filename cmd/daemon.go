package cmd

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "os"
        "os/signal"
        "path/filepath"
        "strings"
        "syscall"
        "time"

        "github.com/fsnotify/fsnotify"
        "github.com/n8n-workflows/n8n-ops/internal/client"
        "github.com/sirupsen/logrus"
        "gopkg.in/yaml.v2"
)

// WorkflowYAML represents a workflow in YAML format
type WorkflowYAML struct {
        ID          string                 `yaml:"id" json:"id"`
        Name        string                 `yaml:"name" json:"name"`
        Active      bool                   `yaml:"active" json:"active"`
        Nodes       []interface{}          `yaml:"nodes" json:"nodes"`
        Connections map[string]interface{} `yaml:"connections" json:"connections"`
        Settings    map[string]interface{} `yaml:"settings,omitempty" json:"settings,omitempty"`
        Tags        []string               `yaml:"tags,omitempty" json:"tags,omitempty"`
        UpdatedAt   string                 `yaml:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// BackupInfo stores backup metadata
type BackupInfo struct {
        OriginalFile string    `json:"originalFile"`
        BackupFile   string    `json:"backupFile"`
        Timestamp    time.Time `json:"timestamp"`
        Environment  string    `json:"environment"`
        WorkflowID   string    `json:"workflowId"`
        WorkflowName string    `json:"workflowName"`
}

func runDaemonMode() {
        logger := logrus.WithFields(logrus.Fields{
                "command": "daemon",
                "env":     environment,
        })

        logger.Info("Starting n8n-ops daemon mode")

        if language == "es" {
                fmt.Printf("🤖 Modo daemon iniciado - %s\n", environment)
                fmt.Printf("👁️ Monitoreando archivos YAML en ./workflows/%s/\n", environment)
                fmt.Printf("💾 Creando backups automáticos antes de actualizar workflows\n")
        } else {
                fmt.Printf("🤖 Daemon mode started - %s environment\n", environment)
                fmt.Printf("👁️ Watching YAML files in ./workflows/%s/\n", environment)
                fmt.Printf("💾 Creating automatic backups before updating workflows\n")
        }

        // Create n8n client
        n8nClient, err := client.NewN8nClient(environment, "mock-api-key")
        if err != nil {
                logger.WithError(err).Fatal("Failed to create n8n client")
                return
        }

        // Test connection
        if err := testN8nConnectionDaemon(n8nClient); err != nil {
                logger.WithError(err).Fatal("Failed to connect to n8n API")
                return
        }

        // Setup file watcher
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
                logger.WithError(err).Fatal("Failed to create file watcher")
                return
        }
        defer watcher.Close()

        // Watch directory
        watchDir := fmt.Sprintf("./workflows/%s", environment)
        if err := setupDirectoryWatch(watcher, watchDir); err != nil {
                logger.WithError(err).Fatal("Failed to setup directory watch")
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
                case event, ok := <-watcher.Events:
                        if !ok {
                                return
                        }
                        handleFileEvent(event, n8nClient, logger)

                case err, ok := <-watcher.Errors:
                        if !ok {
                                return
                        }
                        logger.WithError(err).Error("File watcher error")

                case sig := <-sigChan:
                        fmt.Printf("\n🛑 Received signal %v, stopping daemon...\n", sig)
                        return
                }
        }
}

func setupDirectoryWatch(watcher *fsnotify.Watcher, watchDir string) error {
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

func handleFileEvent(event fsnotify.Event, n8nClient *client.N8nClient, logger *logrus.Entry) {
        // Only process YAML files
        if !strings.HasSuffix(event.Name, ".yaml") && !strings.HasSuffix(event.Name, ".yml") {
                return
        }

        // Only process write events (ignore create/remove for now)
        if event.Op&fsnotify.Write == 0 {
                return
        }

        logger.WithField("file", event.Name).Info("YAML file modified")
        fmt.Printf("📝 File changed: %s\n", filepath.Base(event.Name))

        // Process the file change
        if err := processYAMLFileChange(event.Name, n8nClient, logger); err != nil {
                logger.WithError(err).Error("Failed to process YAML file change")
                fmt.Printf("❌ Error processing %s: %v\n", filepath.Base(event.Name), err)
        }
}

func processYAMLFileChange(filePath string, n8nClient *client.N8nClient, logger *logrus.Entry) error {
        // Read and parse YAML file
        workflow, err := readYAMLWorkflow(filePath)
        if err != nil {
                return fmt.Errorf("failed to read YAML workflow: %w", err)
        }

        // Validate workflow
        if workflow.ID == "" {
                return fmt.Errorf("workflow ID is required")
        }

        // Create backup of existing workflow
        if err := createWorkflowBackup(workflow.ID, n8nClient); err != nil {
                logger.WithError(err).Warn("Failed to create workflow backup")
        }

        // Convert YAML to JSON for n8n API
        jsonWorkflow, err := convertYAMLToJSON(workflow)
        if err != nil {
                return fmt.Errorf("failed to convert YAML to JSON: %w", err)
        }

        // Update workflow in n8n
        if err := updateWorkflowInN8n(jsonWorkflow, n8nClient, logger); err != nil {
                return fmt.Errorf("failed to update workflow in n8n: %w", err)
        }

        fmt.Printf("✅ Workflow '%s' updated in n8n\n", workflow.Name)
        return nil
}

func readYAMLWorkflow(filePath string) (*WorkflowYAML, error) {
        data, err := ioutil.ReadFile(filePath)
        if err != nil {
                return nil, err
        }

        var workflow WorkflowYAML
        if err := yaml.Unmarshal(data, &workflow); err != nil {
                return nil, err
        }

        return &workflow, nil
}

func createWorkflowBackup(workflowID string, n8nClient *client.N8nClient) error {
        // Get current workflow from n8n
        currentWorkflow, err := n8nClient.GetWorkflow(workflowID)
        if err != nil {
                // If workflow doesn't exist, no backup needed
                return nil
        }

        // Create backups directory
        backupDir := fmt.Sprintf("./backups/%s", environment)
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

        if err := ioutil.WriteFile(backupFile, backupData, 0644); err != nil {
                return err
        }

        // Save backup metadata
        backupInfo := BackupInfo{
                OriginalFile: fmt.Sprintf("./workflows/%s/%s.yaml", environment, workflowID),
                BackupFile:   backupFile,
                Timestamp:    time.Now(),
                Environment:  environment,
                WorkflowID:   workflowID,
                WorkflowName: currentWorkflow.Name,
        }

        metadataFile := fmt.Sprintf("%s/%s_%s_backup.meta.json", backupDir, workflowID, timestamp)
        metadataData, err := json.MarshalIndent(backupInfo, "", "  ")
        if err != nil {
                return err
        }

        if err := ioutil.WriteFile(metadataFile, metadataData, 0644); err != nil {
                return err
        }

        fmt.Printf("💾 Backup created: %s\n", filepath.Base(backupFile))
        return nil
}

func convertYAMLToJSON(workflow *WorkflowYAML) (map[string]interface{}, error) {
        // Convert struct to map for JSON serialization
        jsonData, err := json.Marshal(workflow)
        if err != nil {
                return nil, err
        }

        var jsonWorkflow map[string]interface{}
        if err := json.Unmarshal(jsonData, &jsonWorkflow); err != nil {
                return nil, err
        }

        // Set updatedAt timestamp if not provided
        if workflow.UpdatedAt == "" {
                jsonWorkflow["updatedAt"] = time.Now().Format(time.RFC3339)
        }

        return jsonWorkflow, nil
}

func updateWorkflowInN8n(jsonWorkflow map[string]interface{}, n8nClient *client.N8nClient, logger *logrus.Entry) error {
        workflowID := jsonWorkflow["id"].(string)

        // Check if workflow exists
        existingWorkflow, err := n8nClient.GetWorkflow(workflowID)
        if err != nil {
                // Workflow doesn't exist, create new one
                logger.WithField("workflowId", workflowID).Info("Creating new workflow")
                return n8nClient.CreateWorkflow(jsonWorkflow)
        }

        // Workflow exists, update it
        logger.WithField("workflowId", workflowID).Info("Updating existing workflow")
        return n8nClient.UpdateWorkflow(workflowID, jsonWorkflow)
}

func testN8nConnectionDaemon(n8nClient *client.N8nClient) error {
        health, err := n8nClient.HealthCheck()
        if err != nil {
                return fmt.Errorf("health check failed: %w", err)
        }

        if !health["healthy"].(bool) {
                return fmt.Errorf("n8n API is not healthy: status code %v", health["status_code"])
        }

        fmt.Printf("🔗 Connected to n8n API: %s\n", n8nClient.BaseURL())
        return nil
}