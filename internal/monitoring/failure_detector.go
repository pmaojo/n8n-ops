package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/n8n-workflows/n8n-ops/internal/client"
	"github.com/n8n-workflows/n8n-ops/internal/issues"
	"github.com/sirupsen/logrus"
)

// FailureDetector monitors workflow executions and creates issues for failures
type FailureDetector struct {
	n8nClient     client.Client
	issueManager  issues.IssueManager
	logger        *logrus.Logger
	checkInterval time.Duration
	retryThreshold int
}

// NewFailureDetector creates a new failure detector
func NewFailureDetector(n8nClient client.Client, issueManager issues.IssueManager, logger *logrus.Logger) *FailureDetector {
	return &FailureDetector{
		n8nClient:     n8nClient,
		issueManager:  issueManager,
		logger:        logger,
		checkInterval: 1 * time.Minute, // Check every minute
		retryThreshold: 3, // Create issue after 3 consecutive failures
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
		
		// Check recent executions for this workflow
		if err := fd.checkWorkflowExecutions(ctx, workflow.ID, failureCounts, since); err != nil {
			fd.logger.WithFields(logrus.Fields{
				"workflowId": workflow.ID,
				"workflowName": workflow.Name,
			}).WithError(err).Warn("Failed to check workflow executions")
		}
	}
	
	return nil
}

// checkWorkflowExecutions checks executions for a specific workflow
func (fd *FailureDetector) checkWorkflowExecutions(ctx context.Context, workflowID string, failureCounts map[string]int, since time.Time) error {
	// This is a simplified version - in real implementation, you'd get executions from n8n API
	// For now, we'll simulate failure detection based on workflow health
	
	// Simulate getting recent executions (replace with actual API call)
	hasRecentFailure := fd.simulateFailureCheck(workflowID)
	
	if hasRecentFailure {
		failureCounts[workflowID]++
		fd.logger.WithField("workflowId", workflowID).Warn("Workflow failure detected")
		
		// Create issue if threshold reached
		if failureCounts[workflowID] >= fd.retryThreshold {
			if err := fd.createFailureIssue(ctx, workflowID); err != nil {
				return fmt.Errorf("failed to create failure issue: %w", err)
			}
			// Reset counter after creating issue
			failureCounts[workflowID] = 0
		}
	} else {
		// Reset failure count on success
		if failureCounts[workflowID] > 0 {
			fd.logger.WithField("workflowId", workflowID).Info("Workflow recovered")
			if err := fd.handleWorkflowRecovery(ctx, workflowID); err != nil {
				fd.logger.WithError(err).Warn("Failed to handle workflow recovery")
			}
		}
		failureCounts[workflowID] = 0
	}
	
	return nil
}

// simulateFailureCheck simulates checking for workflow failures
// In real implementation, this would query n8n executions API
func (fd *FailureDetector) simulateFailureCheck(workflowID string) bool {
	// Simulate some workflows having occasional failures
	// This would be replaced with actual execution status checking
	return time.Now().Unix()%10 == 0 && workflowID == "test-workflow-1"
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
				"issueId": issue.ID,
				"workflowId": workflowID,
			}).WithError(err).Warn("Failed to update issue with recovery")
			continue
		}
		
		fd.logger.WithFields(logrus.Fields{
			"workflowId": workflowID,
			"issueId": issue.ID,
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