package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/issues"
	"github.com/sirupsen/logrus"
)

// FailureDetector monitors workflow executions and creates issues for failures
type FailureDetector struct {
	n8nClient      client.Client
	issueManager   issues.IssueManager
	logger         *logrus.Logger
	checkInterval  time.Duration
	retryThreshold int
}

// NewFailureDetector creates a new failure detector
func NewFailureDetector(n8nClient client.Client, issueManager issues.IssueManager, logger *logrus.Logger) *FailureDetector {
	return &FailureDetector{
		n8nClient:      n8nClient,
		issueManager:   issueManager,
		logger:         logger,
		checkInterval:  1 * time.Minute, // Check every minute
		retryThreshold: 3,               // Create issue after 3 consecutive failures
	}
}

// Start begins monitoring workflow executions
func (fd *FailureDetector) Start(ctx context.Context) error {
	fd.logger.Info("Starting workflow failure detection")

	ticker := time.NewTicker(fd.checkInterval)
	defer ticker.Stop()

	// Track failure counts per workflow
	failureCounts := make(map[string]int)
	lastCheckTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			fd.logger.Info("Stopping failure detection")
			return ctx.Err()

		case <-ticker.C:
			if err := fd.checkForFailures(ctx, failureCounts, lastCheckTime); err != nil {
				fd.logger.WithError(err).Error("Failed to check for workflow failures")
			}
			lastCheckTime = time.Now()
		}
	}
}

// checkForFailures checks for recent workflow failures
func (fd *FailureDetector) checkForFailures(ctx context.Context, failureCounts map[string]int, since time.Time) error {
	// Get all workflows
	workflows, err := fd.n8nClient.GetWorkflows(ctx)
	if err != nil {
		return fmt.Errorf("failed to get workflows: %w", err)
	}

	for _, workflow := range workflows {
		// Skip inactive workflows
		if !workflow.Active {
			continue
		}

		// Check recent executions for this workflow using n8n API
		hasFailure, err := fd.checkRecentExecutions(ctx, workflow.ID)
		if err != nil {
			fd.logger.WithFields(logrus.Fields{
				"workflowId":   workflow.ID,
				"workflowName": workflow.Name,
			}).WithError(err).Warn("Failed to check workflow executions")
			continue
		}

		if hasFailure {
			failureCounts[workflow.ID]++
			fd.logger.WithField("workflowId", workflow.ID).Warn("Workflow failure detected")

			// Create issue if threshold reached
			if failureCounts[workflow.ID] >= fd.retryThreshold {
				if err := fd.createFailureIssue(ctx, workflow.ID); err != nil {
					fd.logger.WithError(err).Error("Failed to create failure issue")
				}
				// Reset counter after creating issue
				failureCounts[workflow.ID] = 0
			}
		} else {
			// Reset failure count on success
			if failureCounts[workflow.ID] > 0 {
				fd.logger.WithField("workflowId", workflow.ID).Info("Workflow recovered")
				if err := fd.handleWorkflowRecovery(ctx, workflow.ID); err != nil {
					fd.logger.WithError(err).Warn("Failed to handle workflow recovery")
				}
			}
			failureCounts[workflow.ID] = 0
		}
	}

	return nil
}

// checkRecentExecutions checks recent executions for failures using n8n API
func (fd *FailureDetector) checkRecentExecutions(ctx context.Context, workflowID string) (bool, error) {
	// Get recent executions for this workflow (last 10)
	executions, err := fd.n8nClient.GetExecutions(ctx, workflowID, "", 10)
	if err != nil {
		return false, fmt.Errorf("failed to get executions: %w", err)
	}

	if len(executions) == 0 {
		return false, nil // No executions, assume healthy
	}

	// Check last execution for failure
	lastExecution := executions[0]
	if lastExecution.Status == "error" {
		fd.logger.WithFields(logrus.Fields{
			"workflowId":  workflowID,
			"executionId": lastExecution.ID,
			"error":       lastExecution.Error,
		}).Warn("Workflow execution failure detected")
		return true, nil
	}

	return false, nil
}

// createFailureIssue creates a GitLab issue for workflow failure
func (fd *FailureDetector) createFailureIssue(ctx context.Context, workflowID string) error {
	// Get workflow details
	workflow, err := fd.n8nClient.GetWorkflow(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get workflow details: %w", err)
	}

	// Create failure info
	failure := &issues.WorkflowFailure{
		WorkflowID:   workflowID,
		WorkflowName: workflow.Name,
		ExecutionID:  fmt.Sprintf("exec_%d", time.Now().Unix()),
		Environment:  "development", // This would come from config
		FailedAt:     time.Now(),
		ErrorMessage: "Multiple consecutive execution failures detected",
		RetryCount:   fd.retryThreshold,
		Branch:       "main", // This would come from git info
	}

	// Create issue
	issue, err := fd.issueManager.CreateWorkflowFailureIssue(ctx, failure)
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	fd.logger.WithFields(logrus.Fields{
		"workflowId": workflowID,
		"issueId":    issue.ID,
		"issueURL":   issue.WebURL,
	}).Info("Created failure issue for workflow")

	fmt.Printf("🚨 Issue created for workflow failure: %s\n", issue.WebURL)

	return nil
}

// handleWorkflowRecovery handles workflow recovery
func (fd *FailureDetector) handleWorkflowRecovery(ctx context.Context, workflowID string) error {
	// Check if there are open issues for this workflow
	openIssues, err := fd.issueManager.GetOpenFailureIssues(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get open issues: %w", err)
	}

	// Update issues with recovery info
	recovery := &issues.RecoveryInfo{
		RecoveredAt:  time.Now(),
		RecoveryType: "auto",
		Notes:        "Workflow executions are now successful",
	}

	for _, issue := range openIssues {
		if err := fd.issueManager.UpdateIssueWithRecovery(ctx, issue.ID, recovery); err != nil {
			fd.logger.WithFields(logrus.Fields{
				"issueId":    issue.ID,
				"workflowId": workflowID,
			}).WithError(err).Warn("Failed to update issue with recovery")
			continue
		}

		fd.logger.WithFields(logrus.Fields{
			"workflowId": workflowID,
			"issueId":    issue.ID,
		}).Info("Updated issue with workflow recovery")
	}

	if len(openIssues) > 0 {
		fmt.Printf("✅ Workflow %s recovered - updated %d issues\n", workflowID, len(openIssues))
	}

	return nil
}

// SetCheckInterval sets the failure check interval
func (fd *FailureDetector) SetCheckInterval(interval time.Duration) {
	fd.checkInterval = interval
}

// SetRetryThreshold sets the failure threshold before creating issues
func (fd *FailureDetector) SetRetryThreshold(threshold int) {
	fd.retryThreshold = threshold
}
