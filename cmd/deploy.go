package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/n8n-workflows/n8n-ops/internal/client"
	"github.com/n8n-workflows/n8n-ops/internal/config"
	"github.com/n8n-workflows/n8n-ops/internal/storage"
	"github.com/n8n-workflows/n8n-ops/internal/utils"
	"github.com/n8n-workflows/n8n-ops/internal/workflow"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [workflow-file...]",
	Short: "Deploy workflows to n8n instance",
	Long: `Deploy workflow files to the specified n8n environment.
If no files are specified, all workflows in the environment directory will be deployed.
The command validates workflows before deployment and provides rollback capability.

Examples:
  n8n-ops deploy --env staging                    # Deploy all workflows to staging
  n8n-ops deploy workflow1.json --env production  # Deploy specific workflow
  n8n-ops deploy --dry-run --env staging          # Preview deployment without executing`,
	RunE: runDeploy,
}

var (
	dryRun    bool
	skipValidation bool
	workflowDir string
)

func init() {
	rootCmd.AddCommand(deployCmd)
	
	deployCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview deployment without executing")
	deployCmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "skip workflow validation")
	deployCmd.Flags().StringVarP(&workflowDir, "dir", "d", "", "workflow directory (default: workflows/<environment>)")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	env := viper.GetString("environment")
	logger.Info("Starting workflow deployment", "environment", env, "dry-run", dryRun)

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

	// Determine workflow files to deploy
	var workflowFiles []string
	if len(args) > 0 {
		// Deploy specific files
		workflowFiles = args
	} else {
		// Deploy all workflows in environment directory
		if workflowDir == "" {
			workflowDir = filepath.Join("workflows", env)
		}
		
		files, err := getWorkflowFiles(workflowDir)
		if err != nil {
			return fmt.Errorf("failed to get workflow files: %w", err)
		}
		workflowFiles = files
	}

	if len(workflowFiles) == 0 {
		logger.Info("No workflow files found to deploy")
		return nil
	}

	// Validate workflows if not skipped
	if !skipValidation {
		logger.Info("Validating workflows before deployment")
		for _, file := range workflowFiles {
			if err := workflow.ValidateWorkflowFile(file); err != nil {
				return fmt.Errorf("validation failed for %s: %w", file, err)
			}
		}
		logger.Info("All workflows validated successfully")
	}

	if dryRun {
		logger.Info("Dry run - would deploy workflows", "files", workflowFiles)
		return nil
	}

	// Create deployment record
	deploymentID := fmt.Sprintf("deploy_%s_%d", env, time.Now().Unix())
	deployment := &storage.DeploymentRecord{
		ID:          deploymentID,
		Environment: env,
		Status:      "in_progress",
		StartTime:   time.Now(),
		GitCommit:   os.Getenv("CI_COMMIT_SHA"),
		DeployedBy:  os.Getenv("GITLAB_USER_EMAIL"),
	}

	if err := storage.CreateDeploymentRecord(db, deployment); err != nil {
		return fmt.Errorf("failed to create deployment record: %w", err)
	}

	// Deploy workflows
	deployStats := &DeployStats{}
	var deployedWorkflows []string

	for _, file := range workflowFiles {
		workflowID, err := deployWorkflowFile(file, n8nClient, db, env, deployStats)
		if err != nil {
			// Update deployment status to failed
			deployment.Status = "failed"
			deployment.EndTime = time.Now()
			deployment.ErrorMessage = err.Error()
			storage.UpdateDeploymentRecord(db, deployment)
			
			return fmt.Errorf("deployment failed at workflow %s: %w", file, err)
		}
		deployedWorkflows = append(deployedWorkflows, workflowID)
	}

	// Update deployment record to success
	deployment.Status = "success"
	deployment.EndTime = time.Now()
	deployment.WorkflowCount = deployStats.Created + deployStats.Updated
	if err := storage.UpdateDeploymentRecord(db, deployment); err != nil {
		logger.Warn("Failed to update deployment record", "error", err)
	}

	// Generate deployment report
	if err := generateDeploymentReport(deploymentID, env, deployStats, deployedWorkflows); err != nil {
		logger.Warn("Failed to generate deployment report", "error", err)
	}

	logger.Info("Deployment completed successfully",
		"deployment_id", deploymentID,
		"created", deployStats.Created,
		"updated", deployStats.Updated,
		"activated", deployStats.Activated,
	)

	return nil
}

type DeployStats struct {
	Created   int
	Updated   int
	Activated int
}

func getWorkflowFiles(dir string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(path, ".json") && !strings.HasPrefix(info.Name(), "_") {
			files = append(files, path)
		}
		
		return nil
	})
	
	return files, err
}

func deployWorkflowFile(file string, n8nClient *client.N8nClient, db *storage.SQLiteDB, env string, stats *DeployStats) (string, error) {
	// Load workflow from file
	wf, err := utils.LoadWorkflowFromFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to load workflow from file: %w", err)
	}

	logger.Info("Deploying workflow", "name", wf.Name, "file", file)

	// Check if workflow exists
	existingID, err := n8nClient.FindWorkflowByName(wf.Name)
	if err != nil {
		return "", fmt.Errorf("failed to check existing workflow: %w", err)
	}

	var deployedID string
	if existingID != "" {
		// Update existing workflow
		deployedWf, err := n8nClient.UpdateWorkflow(existingID, wf)
		if err != nil {
			return "", fmt.Errorf("failed to update workflow: %w", err)
		}
		deployedID = deployedWf.ID
		stats.Updated++
		logger.Info("Updated workflow", "name", wf.Name, "id", deployedID)
	} else {
		// Create new workflow
		deployedWf, err := n8nClient.CreateWorkflow(wf)
		if err != nil {
			return "", fmt.Errorf("failed to create workflow: %w", err)
		}
		deployedID = deployedWf.ID
		stats.Created++
		logger.Info("Created workflow", "name", wf.Name, "id", deployedID)
	}

	// Activate workflow if it was active
	if wf.Active {
		if err := n8nClient.ActivateWorkflow(deployedID); err != nil {
			logger.Warn("Failed to activate workflow", "name", wf.Name, "error", err)
		} else {
			stats.Activated++
			logger.Info("Activated workflow", "name", wf.Name)
		}
	}

	// Update workflow record in database
	record := &storage.WorkflowRecord{
		ID:           deployedID,
		Name:         wf.Name,
		Environment:  env,
		FilePath:     file,
		Hash:         utils.HashWorkflow(wf),
		LastDeploy:   time.Now(),
		Version:      wf.VersionId,
	}

	// Try to update existing record, create if not found
	if err := storage.UpdateWorkflowRecord(db, record); err == storage.ErrRecordNotFound {
		if err := storage.CreateWorkflowRecord(db, record); err != nil {
			logger.Warn("Failed to create workflow record", "error", err)
		}
	} else if err != nil {
		logger.Warn("Failed to update workflow record", "error", err)
	}

	return deployedID, nil
}

func generateDeploymentReport(deploymentID, env string, stats *DeployStats, workflowIDs []string) error {
	report := map[string]interface{}{
		"deployment_id":   deploymentID,
		"environment":     env,
		"deploy_date":     time.Now().Format(time.RFC3339),
		"git_commit":      os.Getenv("CI_COMMIT_SHA"),
		"pipeline_url":    os.Getenv("CI_PIPELINE_URL"),
		"deployed_by":     os.Getenv("GITLAB_USER_EMAIL"),
		"statistics": map[string]int{
			"created":   stats.Created,
			"updated":   stats.Updated,
			"activated": stats.Activated,
		},
		"workflow_ids": workflowIDs,
	}

	filename := fmt.Sprintf("deployment-report-%s-%s.json", env, deploymentID)
	return utils.WriteJSONFile(report, filename)
}
