package cmd

import (
        "fmt"
        "time"

        "github.com/spf13/cobra"
        "github.com/spf13/viper"
        "github.com/n8n-workflows/cli/internal/client"
        "github.com/n8n-workflows/cli/internal/config"
        "github.com/n8n-workflows/cli/internal/storage"
        "github.com/n8n-workflows/cli/internal/utils"
        "github.com/n8n-workflows/cli/internal/workflow"
)

var rollbackCmd = &cobra.Command{
        Use:   "rollback [deployment-id]",
        Short: "Rollback to a previous deployment",
        Long: `Rollback to a previous deployment state.
If no deployment ID is specified, rolls back to the most recent successful deployment.

Examples:
  n8n-cli rollback --env production                    # Rollback to last successful deployment
  n8n-cli rollback deploy_production_1234567890        # Rollback to specific deployment
  n8n-cli rollback --list                             # List available deployments`,
        RunE: runRollback,
}

var (
        listDeployments bool
        confirmRollback bool
)

func init() {
        rootCmd.AddCommand(rollbackCmd)
        
        rollbackCmd.Flags().BoolVar(&listDeployments, "list", false, "list available deployments")
        rollbackCmd.Flags().BoolVarP(&confirmRollback, "yes", "y", false, "confirm rollback without prompt")
}

func runRollback(cmd *cobra.Command, args []string) error {
        env := viper.GetString("environment")
        
        // Initialize storage
        db, err := storage.NewSQLiteDB()
        if err != nil {
                return fmt.Errorf("failed to initialize database: %w", err)
        }
        defer db.Close()

        if listDeployments {
                return listDeploymentHistory(db, env)
        }

        logger.Info("Starting rollback process", "environment", env)

        // Get deployment to rollback to
        var targetDeployment *storage.DeploymentRecord
        if len(args) > 0 {
                // Rollback to specific deployment
                deployment, err := storage.GetDeploymentRecord(db, args[0])
                if err != nil {
                        return fmt.Errorf("failed to get deployment record: %w", err)
                }
                targetDeployment = deployment
        } else {
                // Rollback to last successful deployment
                deployments, err := storage.GetDeploymentHistory(db, env, 2) // Get last 2 deployments
                if err != nil {
                        return fmt.Errorf("failed to get deployment history: %w", err)
                }
                
                if len(deployments) < 2 {
                        return fmt.Errorf("not enough deployment history for rollback")
                }
                
                // Use the second most recent successful deployment
                for i := 1; i < len(deployments); i++ {
                        if deployments[i].Status == "success" {
                                targetDeployment = deployments[i]
                                break
                        }
                }
                
                if targetDeployment == nil {
                        return fmt.Errorf("no successful deployment found for rollback")
                }
        }

        // Confirm rollback
        if !confirmRollback {
                fmt.Printf("Are you sure you want to rollback environment '%s' to deployment '%s' from %s? [y/N]: ", 
                        env, targetDeployment.ID, targetDeployment.StartTime.Format("2006-01-02 15:04:05"))
                
                var response string
                fmt.Scanln(&response)
                if response != "y" && response != "Y" && response != "yes" {
                        fmt.Println("Rollback cancelled")
                        return nil
                }
        }

        // Perform rollback
        return performRollback(targetDeployment, env, db)
}

func listDeploymentHistory(db *storage.SQLiteDB, env string) error {
        deployments, err := storage.GetDeploymentHistory(db, env, 10)
        if err != nil {
                return fmt.Errorf("failed to get deployment history: %w", err)
        }

        if len(deployments) == 0 {
                fmt.Printf("No deployment history found for environment '%s'\n", env)
                return nil
        }

        fmt.Printf("Deployment history for environment '%s':\n\n", env)
        fmt.Printf("%-30s %-10s %-20s %-15s %-10s\n", "Deployment ID", "Status", "Date", "Deployed By", "Workflows")
        fmt.Printf("%-30s %-10s %-20s %-15s %-10s\n", 
                "------------------------------", 
                "----------", 
                "--------------------", 
                "---------------", 
                "----------")

        for _, deployment := range deployments {
                deployedBy := deployment.DeployedBy
                if deployedBy == "" {
                        deployedBy = "unknown"
                }
                
                fmt.Printf("%-30s %-10s %-20s %-15s %-10d\n",
                        deployment.ID,
                        deployment.Status,
                        deployment.StartTime.Format("2006-01-02 15:04:05"),
                        deployedBy,
                        deployment.WorkflowCount,
                )
        }

        return nil
}

func performRollback(targetDeployment *storage.DeploymentRecord, env string, db *storage.SQLiteDB) error {
        logger.Info("Performing rollback", 
                "target_deployment", targetDeployment.ID,
                "environment", env,
        )

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

        // Get workflows that were deployed in the target deployment
        workflows, err := storage.GetWorkflowsByDeployment(db, targetDeployment.ID)
        if err != nil {
                return fmt.Errorf("failed to get workflows for deployment: %w", err)
        }

        if len(workflows) == 0 {
                return fmt.Errorf("no workflows found for target deployment")
        }

        // Create rollback deployment record
        rollbackID := fmt.Sprintf("rollback_%s_%d", env, time.Now().Unix())
        rollbackDeployment := &storage.DeploymentRecord{
                ID:          rollbackID,
                Environment: env,
                Status:      "in_progress",
                StartTime:   time.Now(),
                GitCommit:   targetDeployment.GitCommit,
                DeployedBy:  fmt.Sprintf("rollback_to_%s", targetDeployment.ID),
        }

        if err := storage.CreateDeploymentRecord(db, rollbackDeployment); err != nil {
                return fmt.Errorf("failed to create rollback deployment record: %w", err)
        }

        // Rollback each workflow
        rollbackCount := 0
        for _, workflow := range workflows {
                logger.Info("Rolling back workflow", "name", workflow.Name, "id", workflow.ID)
                
                // Load workflow from file if it exists
                if workflow.FilePath != "" {
                        wf, err := loadWorkflowFromRecord(workflow)
                        if err != nil {
                                logger.Warn("Failed to load workflow from file, skipping", "workflow", workflow.Name, "error", err)
                                continue
                        }

                        // Deploy the workflow
                        if _, err := n8nClient.UpdateWorkflow(workflow.ID, wf); err != nil {
                                logger.Error("Failed to rollback workflow", "workflow", workflow.Name, "error", err)
                                // Continue with other workflows instead of failing completely
                                continue
                        }

                        // Activate/deactivate based on original state
                        if wf.Active {
                                if err := n8nClient.ActivateWorkflow(workflow.ID); err != nil {
                                        logger.Warn("Failed to activate workflow during rollback", "workflow", workflow.Name, "error", err)
                                }
                        } else {
                                if err := n8nClient.DeactivateWorkflow(workflow.ID); err != nil {
                                        logger.Warn("Failed to deactivate workflow during rollback", "workflow", workflow.Name, "error", err)
                                }
                        }

                        rollbackCount++
                        logger.Info("Successfully rolled back workflow", "name", workflow.Name)
                }
        }

        // Update rollback deployment record
        rollbackDeployment.Status = "success"
        rollbackDeployment.EndTime = time.Now()
        rollbackDeployment.WorkflowCount = rollbackCount
        if err := storage.UpdateDeploymentRecord(db, rollbackDeployment); err != nil {
                logger.Warn("Failed to update rollback deployment record", "error", err)
        }

        logger.Info("Rollback completed successfully",
                "rollback_id", rollbackID,
                "workflows_rolled_back", rollbackCount,
                "target_deployment", targetDeployment.ID,
        )

        fmt.Printf("✅ Rollback completed successfully\n")
        fmt.Printf("Rolled back %d workflows to deployment %s\n", rollbackCount, targetDeployment.ID)
        fmt.Printf("Rollback deployment ID: %s\n", rollbackID)

        return nil
}

func loadWorkflowFromRecord(record *storage.WorkflowRecord) (*workflow.Workflow, error) {
        return utils.LoadWorkflowFromFile(record.FilePath)
}
