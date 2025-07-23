package issues

import (
        "bytes"
        "context"
        "encoding/json"
        "fmt"
        "net/http"
        "time"
)

// IssueManager handles automatic issue creation for workflow failures
type IssueManager interface {
        CreateWorkflowFailureIssue(ctx context.Context, failure *WorkflowFailure) (*Issue, error)
        UpdateIssueWithRecovery(ctx context.Context, issueID int, recovery *RecoveryInfo) error
        GetOpenFailureIssues(ctx context.Context, workflowID string) ([]*Issue, error)
}

// WorkflowFailure represents a workflow execution failure
type WorkflowFailure struct {
        WorkflowID    string                 `json:"workflowId"`
        WorkflowName  string                 `json:"workflowName"`
        ExecutionID   string                 `json:"executionId"`
        Environment   string                 `json:"environment"`
        FailedAt      time.Time              `json:"failedAt"`
        ErrorMessage  string                 `json:"errorMessage"`
        ErrorDetails  map[string]interface{} `json:"errorDetails,omitempty"`
        NodeName      string                 `json:"nodeName,omitempty"`
        RetryCount    int                    `json:"retryCount"`
        PipelineID    string                 `json:"pipelineId,omitempty"`
        PipelineURL   string                 `json:"pipelineUrl,omitempty"`
        CommitSHA     string                 `json:"commitSha,omitempty"`
        Branch        string                 `json:"branch,omitempty"`
}

// RecoveryInfo represents workflow recovery information
type RecoveryInfo struct {
        RecoveredAt   time.Time `json:"recoveredAt"`
        RecoveryType  string    `json:"recoveryType"` // "auto", "manual", "rollback"
        RecoveredBy   string    `json:"recoveredBy,omitempty"`
        NewCommitSHA  string    `json:"newCommitSha,omitempty"`
        Notes         string    `json:"notes,omitempty"`
}

// Issue represents a GitLab issue
type Issue struct {
        ID          int       `json:"id"`
        IID         int       `json:"iid"`
        Title       string    `json:"title"`
        Description string    `json:"description"`
        State       string    `json:"state"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
        Labels      []string  `json:"labels"`
        WebURL      string    `json:"web_url"`
}

// GitLabIssueManager implements IssueManager for GitLab
type GitLabIssueManager struct {
        baseURL    string
        projectID  string
        token      string
        httpClient *http.Client
}

// NewGitLabIssueManager creates a new GitLab issue manager
func NewGitLabIssueManager(baseURL, projectID, token string) *GitLabIssueManager {
        return &GitLabIssueManager{
                baseURL:    baseURL,
                projectID:  projectID,
                token:      token,
                httpClient: &http.Client{Timeout: 30 * time.Second},
        }
}

// CreateWorkflowFailureIssue creates a GitLab issue for workflow failure
func (gim *GitLabIssueManager) CreateWorkflowFailureIssue(ctx context.Context, failure *WorkflowFailure) (*Issue, error) {
        title := fmt.Sprintf("🚨 Workflow Failure: %s (%s)", failure.WorkflowName, failure.Environment)
        
        description := gim.buildIssueDescription(failure)
        labels := gim.buildIssueLabels(failure)
        
        issueData := map[string]interface{}{
                "title":       title,
                "description": description,
                "labels":      labels,
        }
        
        jsonData, err := json.Marshal(issueData)
        if err != nil {
                return nil, fmt.Errorf("failed to marshal issue data: %w", err)
        }
        
        url := fmt.Sprintf("%s/api/v4/projects/%s/issues", gim.baseURL, gim.projectID)
        req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
                return nil, fmt.Errorf("failed to create request: %w", err)
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("PRIVATE-TOKEN", gim.token)
        
        resp, err := gim.httpClient.Do(req)
        if err != nil {
                return nil, fmt.Errorf("failed to create issue: %w", err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusCreated {
                return nil, fmt.Errorf("GitLab API returned status: %d", resp.StatusCode)
        }
        
        var issue Issue
        if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
                return nil, fmt.Errorf("failed to decode response: %w", err)
        }
        
        return &issue, nil
}

// UpdateIssueWithRecovery updates an issue when workflow recovers
func (gim *GitLabIssueManager) UpdateIssueWithRecovery(ctx context.Context, issueID int, recovery *RecoveryInfo) error {
        recoveryComment := gim.buildRecoveryComment(recovery)
        
        // Add recovery comment
        commentData := map[string]interface{}{
                "body": recoveryComment,
        }
        
        jsonData, err := json.Marshal(commentData)
        if err != nil {
                return fmt.Errorf("failed to marshal comment data: %w", err)
        }
        
        url := fmt.Sprintf("%s/api/v4/projects/%s/issues/%d/notes", gim.baseURL, gim.projectID, issueID)
        req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
                return fmt.Errorf("failed to create comment request: %w", err)
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("PRIVATE-TOKEN", gim.token)
        
        resp, err := gim.httpClient.Do(req)
        if err != nil {
                return fmt.Errorf("failed to add comment: %w", err)
        }
        defer resp.Body.Close()
        
        // Close the issue if auto-recovered
        if recovery.RecoveryType == "auto" {
                return gim.closeIssue(ctx, issueID)
        }
        
        return nil
}

// GetOpenFailureIssues gets open issues for a specific workflow
func (gim *GitLabIssueManager) GetOpenFailureIssues(ctx context.Context, workflowID string) ([]*Issue, error) {
        url := fmt.Sprintf("%s/api/v4/projects/%s/issues?state=opened&labels=workflow-failure,workflow:%s", 
                gim.baseURL, gim.projectID, workflowID)
        
        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
                return nil, fmt.Errorf("failed to create request: %w", err)
        }
        
        req.Header.Set("PRIVATE-TOKEN", gim.token)
        
        resp, err := gim.httpClient.Do(req)
        if err != nil {
                return nil, fmt.Errorf("failed to get issues: %w", err)
        }
        defer resp.Body.Close()
        
        var issues []*Issue
        if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
                return nil, fmt.Errorf("failed to decode response: %w", err)
        }
        
        return issues, nil
}

// closeIssue closes a GitLab issue
func (gim *GitLabIssueManager) closeIssue(ctx context.Context, issueID int) error {
        updateData := map[string]interface{}{
                "state_event": "close",
        }
        
        jsonData, err := json.Marshal(updateData)
        if err != nil {
                return fmt.Errorf("failed to marshal update data: %w", err)
        }
        
        url := fmt.Sprintf("%s/api/v4/projects/%s/issues/%d", gim.baseURL, gim.projectID, issueID)
        req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
        if err != nil {
                return fmt.Errorf("failed to create close request: %w", err)
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("PRIVATE-TOKEN", gim.token)
        
        resp, err := gim.httpClient.Do(req)
        if err != nil {
                return fmt.Errorf("failed to close issue: %w", err)
        }
        defer resp.Body.Close()
        
        return nil
}

// buildIssueDescription creates a detailed issue description
func (gim *GitLabIssueManager) buildIssueDescription(failure *WorkflowFailure) string {
        description := fmt.Sprintf(`## 🚨 Workflow Execution Failure

**Workflow Details:**
- **Name:** %s
- **ID:** %s
- **Environment:** %s
- **Execution ID:** %s

**Failure Information:**
- **Failed At:** %s
- **Error:** %s
- **Retry Count:** %d`, 
                failure.WorkflowName,
                failure.WorkflowID,
                failure.Environment,
                failure.ExecutionID,
                failure.FailedAt.Format(time.RFC3339),
                failure.ErrorMessage,
                failure.RetryCount)

        if failure.NodeName != "" {
                description += fmt.Sprintf("\n- **Failed Node:** %s", failure.NodeName)
        }

        if failure.PipelineURL != "" {
                description += fmt.Sprintf("\n\n**Pipeline Information:**\n- **Pipeline:** [%s](%s)", 
                        failure.PipelineID, failure.PipelineURL)
        }

        if failure.CommitSHA != "" {
                description += fmt.Sprintf("\n- **Commit:** %s", failure.CommitSHA)
        }

        if failure.Branch != "" {
                description += fmt.Sprintf("\n- **Branch:** %s", failure.Branch)
        }

        description += `

## 🔧 Troubleshooting Steps

1. Check the workflow execution logs in n8n
2. Verify input data and node configurations
3. Test workflow in development environment
4. Check for recent changes in the repository

## 📊 Next Actions

- [ ] Investigate root cause
- [ ] Apply fix or rollback
- [ ] Test in staging environment
- [ ] Deploy to production

---
*This issue was automatically created by n8n-ops when a workflow failure was detected.*`

        return description
}

// buildIssueLabels creates appropriate labels for the failure issue
func (gim *GitLabIssueManager) buildIssueLabels(failure *WorkflowFailure) []string {
        labels := []string{
                "workflow-failure",
                "automated",
                fmt.Sprintf("env:%s", failure.Environment),
                fmt.Sprintf("workflow:%s", failure.WorkflowID),
        }

        // Add severity label based on retry count
        if failure.RetryCount > 5 {
                labels = append(labels, "severity:high")
        } else if failure.RetryCount > 2 {
                labels = append(labels, "severity:medium")
        } else {
                labels = append(labels, "severity:low")
        }

        // Add node-specific label if available
        if failure.NodeName != "" {
                labels = append(labels, fmt.Sprintf("node:%s", failure.NodeName))
        }

        return labels
}

// buildRecoveryComment creates a comment for workflow recovery
func (gim *GitLabIssueManager) buildRecoveryComment(recovery *RecoveryInfo) string {
        comment := fmt.Sprintf(`## ✅ Workflow Recovery Detected

**Recovery Information:**
- **Recovered At:** %s
- **Recovery Type:** %s`,
                recovery.RecoveredAt.Format(time.RFC3339),
                recovery.RecoveryType)

        if recovery.RecoveredBy != "" {
                comment += fmt.Sprintf("\n- **Recovered By:** %s", recovery.RecoveredBy)
        }

        if recovery.NewCommitSHA != "" {
                comment += fmt.Sprintf("\n- **New Commit:** %s", recovery.NewCommitSHA)
        }

        if recovery.Notes != "" {
                comment += fmt.Sprintf("\n\n**Notes:** %s", recovery.Notes)
        }

        comment += "\n\n---\n*This update was automatically generated by n8n-ops.*"

        return comment
}