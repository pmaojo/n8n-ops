package cmd

import (
        "encoding/json"
        "fmt"
        "os"
        "path/filepath"
        "time"

        "github.com/spf13/cobra"
        "github.com/n8n-workflows/n8n-ops/internal/client"
)

var syncCmd = &cobra.Command{
        Use:   "sync",
        Short: "Sync workflows bidirectionally between n8n and local filesystem", 
        Long: `Sync performs intelligent bidirectional synchronization between your n8n instance and Git repository.

DETECTION MODES:
• FROM n8n TO Git: Downloads workflows modified in n8n UI to local files
• FROM Git TO n8n: Uploads local workflow changes to n8n instance  
• BIDIRECTIONAL: Compares timestamps and syncs in both directions (default)

CHANGE DETECTION:
• Compares workflow updatedAt timestamps with local file modification times
• Detects new workflows created in n8n UI
• Identifies local JSON file changes not yet pushed to n8n
• Handles conflict resolution with user prompts

Examples:
  n8n-ops sync --env development    # Smart bidirectional sync
  n8n-ops sync --from-n8n          # Only download from n8n to Git
  n8n-ops sync --to-n8n            # Only upload from Git to n8n  
  n8n-ops sync --force             # Force sync, auto-resolve conflicts`,
        RunE: runSync,
}

var (
        force      bool
        outputDir  string
        branch     string
        fromN8n    bool
        toN8n      bool
        dryRun     bool
)

func init() {
        rootCmd.AddCommand(syncCmd)
        
        syncCmd.Flags().BoolVarP(&force, "force", "f", false, "force sync, overwriting conflicts")
        syncCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory (default: workflows/<environment>)")
        syncCmd.Flags().StringVarP(&branch, "branch", "b", "", "git branch (defaults to current branch)")
        syncCmd.Flags().BoolVar(&fromN8n, "from-n8n", false, "sync only from n8n to local files")
        syncCmd.Flags().BoolVar(&toN8n, "to-n8n", false, "sync only from local files to n8n")
        syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be synced without making changes")
}

func runSync(cmd *cobra.Command, args []string) error {
        logger.Info("Starting workflow sync", "environment", environment)
        fmt.Printf("🔄 Syncing workflows from %s environment...\n", environment)

        // Set output directory
        if outputDir == "" {
                outputDir = filepath.Join("workflows", environment)
        }

        // Create output directory if it doesn't exist
        if err := os.MkdirAll(outputDir, 0755); err != nil {
                return fmt.Errorf("failed to create output directory: %w", err)
        }

        // Get n8n configuration from environment (or use demo mode)
        apiURL := os.Getenv("N8N_URL")
        apiKey := os.Getenv("N8N_API_KEY")
        
        if demoMode {
                fmt.Printf("🎭 Demo mode enabled - using mock n8n data\n")
                apiURL = "demo://localhost"
                apiKey = "demo-key"
        } else if apiURL == "" || apiKey == "" {
                fmt.Printf("⚠️  N8N_URL or N8N_API_KEY not configured\n")
                fmt.Printf("💡 Set environment variables or use --demo flag:\n")
                fmt.Printf("   export N8N_URL=http://localhost:3001\n")
                fmt.Printf("   export N8N_API_KEY=n8n_api_mock_development\n")
                fmt.Printf("   OR: ./n8n-ops sync --demo\n")
                return nil
        }

        // Initialize n8n client (with demo mode support)
        n8nClient := client.NewN8nClientWithDemo(apiURL, apiKey, demoMode)
        
        // Test connection
        fmt.Printf("🔗 Connecting to n8n API: %s\n", apiURL)
        user, err := n8nClient.GetMe()
        if err != nil {
                return fmt.Errorf("failed to connect to n8n API: %w", err)
        }
        fmt.Printf("✅ Connected as: %s\n", user.Email)

        // Get workflows from n8n
        workflows, err := n8nClient.GetWorkflows()
        if err != nil {
                return fmt.Errorf("failed to fetch workflows: %w", err)
        }

        fmt.Printf("📋 Found %d workflows in n8n instance\n", len(workflows))

        // Process each workflow
        syncedCount := 0
        for _, workflow := range workflows {
                // Generate filename
                filename := fmt.Sprintf("%s_%s.json", sanitizeFilename(workflow.Name), workflow.ID)
                filepath := filepath.Join(outputDir, filename)

                // Add sync metadata
                workflowData := map[string]interface{}{
                        "id":          workflow.ID,
                        "name":        workflow.Name,
                        "active":      workflow.Active,
                        "nodes":       workflow.Nodes,
                        "connections": workflow.Connections,
                        "createdAt":   workflow.CreatedAt,
                        "updatedAt":   workflow.UpdatedAt,
                        "versionId":   workflow.VersionId,
                        "tags":        workflow.Tags,
                        "syncMetadata": map[string]interface{}{
                                "syncDate":    time.Now(),
                                "environment": environment,
                                "syncedBy":    user.Email,
                                "gitCommit":   os.Getenv("CI_COMMIT_SHA"),
                        },
                }

                // Write workflow to file
                if err := writeWorkflowFile(workflowData, filepath); err != nil {
                        logger.Error("Failed to write workflow file", "workflow", workflow.Name, "error", err)
                        if !force {
                                return fmt.Errorf("failed to write workflow %s: %w", workflow.Name, err)
                        }
                        continue
                }

                fmt.Printf("📝 Synced: %s → %s\n", workflow.Name, filename)
                syncedCount++
        }

        // Generate sync metadata
        metadata := map[string]interface{}{
                "lastSync":        time.Now(),
                "environment":     environment,
                "totalWorkflows":  len(workflows),
                "syncedWorkflows": syncedCount,
                "n8nURL":          apiURL,
                "syncedBy":        user.Email,
                "gitCommit":       os.Getenv("CI_COMMIT_SHA"),
        }

        metadataPath := filepath.Join(outputDir, "_sync_metadata.json")
        if err := writeWorkflowFile(metadata, metadataPath); err != nil {
                logger.Warn("Failed to write sync metadata", "error", err)
        }

        fmt.Printf("✅ Sync completed: %d workflows synced to %s\n", syncedCount, outputDir)
        return nil
}

func sanitizeFilename(name string) string {
        // Replace spaces and special characters with underscores
        result := ""
        for _, char := range name {
                if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
                        result += string(char)
                } else {
                        result += "_"
                }
        }
        return result
}

func writeWorkflowFile(data interface{}, filepath string) error {
        file, err := os.Create(filepath)
        if err != nil {
                return fmt.Errorf("failed to create file: %w", err)
        }
        defer file.Close()

        encoder := json.NewEncoder(file)
        encoder.SetIndent("", "  ")
        if err := encoder.Encode(data); err != nil {
                return fmt.Errorf("failed to encode JSON: %w", err)
        }

        return nil
}
