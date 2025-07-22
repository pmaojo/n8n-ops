package cmd

import (
        "fmt"
        "os"
        "path/filepath"
        "time"

        "github.com/spf13/cobra"
        "github.com/spf13/viper"
        "github.com/n8n-workflows/cli/internal/ascii"
        "github.com/n8n-workflows/cli/internal/client"
        "github.com/n8n-workflows/cli/internal/config"
        "github.com/n8n-workflows/cli/internal/git"
        "github.com/n8n-workflows/cli/internal/storage"
        "github.com/n8n-workflows/cli/internal/utils"
        "github.com/n8n-workflows/cli/internal/workflow"
)

var syncCmd = &cobra.Command{
        Use:   "sync",
        Short: "Sync workflows from n8n instance to local filesystem",
        Long: `Sync workflows from the specified n8n environment to the local filesystem.
This command downloads all workflows from the target n8n instance and stores them
in the appropriate environment directory with metadata for tracking changes.

Examples:
  n8n-cli sync --env development    # Sync from development environment
  n8n-cli sync --env staging        # Sync from staging environment
  n8n-cli sync --force              # Force sync, overwriting local changes`,
        RunE: runSync,
}

var (
        force      bool
        outputDir  string
        branch     string
)

func init() {
        rootCmd.AddCommand(syncCmd)
        
        syncCmd.Flags().BoolVarP(&force, "force", "f", false, "force sync, overwriting local changes")
        syncCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory (default: workflows/<environment>)")
        syncCmd.Flags().StringVarP(&branch, "branch", "b", "", "git branch (defaults to current branch)")
}

func runSync(cmd *cobra.Command, args []string) error {
        env := viper.GetString("environment")
        
        // Handle branch logic if Git repository
        if git.IsGitRepository() {
                currentBranch, err := git.GetCurrentBranch()
                if err != nil {
                        logger.Warn("Failed to get current Git branch", "error", err)
                } else {
                        logger.Info("Current Git branch", "branch", currentBranch)
                        
                        // If branch flag is provided, switch to it
                        if branch != "" && branch != currentBranch {
                                logger.Info("Switching to branch", "branch", branch)
                                if err := git.CheckoutBranch(branch); err != nil {
                                        return fmt.Errorf("failed to checkout branch %s: %w", branch, err)
                                }
                                currentBranch = branch
                        }
                        
                        // Auto-detect environment from branch if not explicitly set
                        if viper.GetString("environment") == "development" && env == "development" {
                                detectedEnv := git.GetEnvironmentFromBranch(currentBranch)
                                if detectedEnv != env {
                                        env = detectedEnv
                                        logger.Info("Auto-detected environment from branch", "environment", env, "branch", currentBranch)
                                }
                        }
                }
        }
        
        // Show banner
        fmt.Print(ascii.Banner(env))
        fmt.Print(ascii.CommandHelp("sync"))
        
        logger.Info("Starting workflow sync", "environment", env)

        // Get environment configuration
        envConfig, err := config.GetEnvironmentConfig(env)
        if err != nil {
                return fmt.Errorf("failed to get environment config: %w", err)
        }

        // Initialize n8n client
        n8nClient, err := client.NewN8nClient(envConfig.URL, envConfig.APIKey)
        if err != nil {
                return fmt.Errorf("failed to create n8n client: %w", err)
        }

        // Initialize storage
        db, err := storage.NewSQLiteDB()
        if err != nil {
                return fmt.Errorf("failed to initialize database: %w", err)
        }
        defer db.Close()

        // Set output directory
        if outputDir == "" {
                outputDir = filepath.Join("workflows", env)
        }

        // Create output directory if it doesn't exist
        if err := os.MkdirAll(outputDir, 0755); err != nil {
                return fmt.Errorf("failed to create output directory: %w", err)
        }

        // Get workflows from n8n instance
        logger.Info("Fetching workflows from n8n instance")
        workflows, err := n8nClient.GetWorkflows()
        if err != nil {
                return fmt.Errorf("failed to fetch workflows: %w", err)
        }

        logger.Info("Found workflows", "count", len(workflows))

        // Process each workflow
        syncStats := &SyncStats{}
        for _, wf := range workflows {
                if err := syncWorkflow(wf, outputDir, env, db, syncStats); err != nil {
                        logger.Error("Failed to sync workflow", "workflow", wf.Name, "error", err)
                        if !force {
                                return fmt.Errorf("sync failed for workflow %s: %w", wf.Name, err)
                        }
                }
        }

        // Generate sync metadata
        if err := generateSyncMetadata(outputDir, env, syncStats); err != nil {
                logger.Warn("Failed to generate sync metadata", "error", err)
        }

        logger.Info("Sync completed successfully", 
                "created", syncStats.Created,
                "updated", syncStats.Updated,
                "skipped", syncStats.Skipped,
        )

        return nil
}

type SyncStats struct {
        Created int
        Updated int
        Skipped int
}

func syncWorkflow(wf *workflow.Workflow, outputDir, env string, db *storage.SQLiteDB, stats *SyncStats) error {
        // Generate filename
        filename := fmt.Sprintf("%s_%s.json", 
                utils.SanitizeFilename(wf.Name), 
                wf.ID,
        )
        filepath := filepath.Join(outputDir, filename)

        // Check if file exists and has changes
        existing, err := storage.GetWorkflowRecord(db, wf.ID, env)
        if err != nil && err != storage.ErrRecordNotFound {
                return fmt.Errorf("failed to check existing workflow: %w", err)
        }

        // Check if workflow has changed
        currentHash := utils.HashWorkflow(wf)
        if existing != nil && existing.Hash == currentHash && !force {
                logger.Debug("Workflow unchanged, skipping", "workflow", wf.Name)
                stats.Skipped++
                return nil
        }

        // Add sync metadata
        wf.SyncMetadata = &workflow.SyncMetadata{
                SyncDate:    time.Now(),
                Environment: env,
                GitCommit:   os.Getenv("CI_COMMIT_SHA"),
                SyncedBy:    os.Getenv("GITLAB_USER_EMAIL"),
        }

        // Write workflow to file
        if err := utils.WriteWorkflowToFile(wf, filepath); err != nil {
                return fmt.Errorf("failed to write workflow to file: %w", err)
        }

        // Update database record
        record := &storage.WorkflowRecord{
                ID:           wf.ID,
                Name:         wf.Name,
                Environment:  env,
                FilePath:     filepath,
                Hash:         currentHash,
                LastSync:     time.Now(),
                Version:      wf.VersionId,
        }

        if existing != nil {
                if err := storage.UpdateWorkflowRecord(db, record); err != nil {
                        return fmt.Errorf("failed to update workflow record: %w", err)
                }
                stats.Updated++
                logger.Info("Updated workflow", "workflow", wf.Name, "file", filename)
        } else {
                if err := storage.CreateWorkflowRecord(db, record); err != nil {
                        return fmt.Errorf("failed to create workflow record: %w", err)
                }
                stats.Created++
                logger.Info("Created workflow", "workflow", wf.Name, "file", filename)
        }

        return nil
}

func generateSyncMetadata(outputDir, env string, stats *SyncStats) error {
        metadata := map[string]interface{}{
                "lastSync":        time.Now().Format(time.RFC3339),
                "environment":     env,
                "totalWorkflows":  stats.Created + stats.Updated + stats.Skipped,
                "created":         stats.Created,
                "updated":         stats.Updated,
                "skipped":         stats.Skipped,
                "gitCommit":       os.Getenv("CI_COMMIT_SHA"),
                "pipelineId":      os.Getenv("CI_PIPELINE_ID"),
                "pipelineUrl":     os.Getenv("CI_PIPELINE_URL"),
                "syncedBy":        os.Getenv("GITLAB_USER_EMAIL"),
        }

        metadataPath := filepath.Join(outputDir, "_sync_metadata.json")
        return utils.WriteJSONFile(metadata, metadataPath)
}
