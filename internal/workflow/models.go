package workflow

import (
        "encoding/json"
        "fmt"
        "time"
)

// Workflow represents an n8n workflow
type Workflow struct {
        ID        string                 `json:"id,omitempty"`
        Name      string                 `json:"name" validate:"required"`
        Active    bool                   `json:"active"`
        Nodes     []Node                 `json:"nodes" validate:"required,min=1"`
        Connections map[string]interface{} `json:"connections,omitempty"`
        Settings  map[string]interface{} `json:"settings,omitempty"`
        StaticData map[string]interface{} `json:"staticData,omitempty"`
        VersionId  int                   `json:"versionId,omitempty"`
        Tags       []Tag                 `json:"tags,omitempty"`
        
        // Custom fields for CLI management
        SyncMetadata *SyncMetadata `json:"syncMetadata,omitempty"`
}

// Node represents a workflow node
type Node struct {
        ID          string                 `json:"id,omitempty"`
        Name        string                 `json:"name" validate:"required"`
        Type        string                 `json:"type" validate:"required"`
        TypeVersion int                   `json:"typeVersion,omitempty"`
        Position    []float64             `json:"position" validate:"required,len=2"`
        Parameters  map[string]interface{} `json:"parameters,omitempty"`
        Credentials map[string]interface{} `json:"credentials,omitempty"`
        Disabled    bool                  `json:"disabled,omitempty"`
        Notes       string                `json:"notes,omitempty"`
        Color       string                `json:"color,omitempty"`
}

// Tag represents a workflow tag
type Tag struct {
        ID        string    `json:"id,omitempty"`
        Name      string    `json:"name"`
        CreatedAt time.Time `json:"createdAt,omitempty"`
        UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// SyncMetadata contains metadata about workflow synchronization
type SyncMetadata struct {
        SyncDate    time.Time `json:"syncDate"`
        Environment string    `json:"environment"`
        GitCommit   string    `json:"gitCommit,omitempty"`
        SyncedBy    string    `json:"syncedBy,omitempty"`
        PipelineID  string    `json:"pipelineId,omitempty"`
        PipelineURL string    `json:"pipelineUrl,omitempty"`
}

// Connection represents a connection between nodes
type Connection struct {
        Node  string `json:"node"`
        Type  string `json:"type"`
        Index int    `json:"index"`
}

// NodeConnection represents connections from a specific node output
type NodeConnection struct {
        Main [][]Connection `json:"main,omitempty"`
}

// WorkflowExecution represents a workflow execution
type WorkflowExecution struct {
        ID           string                 `json:"id"`
        Finished     bool                   `json:"finished"`
        Mode         string                 `json:"mode"`
        RetryOf      string                 `json:"retryOf,omitempty"`
        StartedAt    time.Time              `json:"startedAt"`
        StoppedAt    time.Time              `json:"stoppedAt,omitempty"`
        WorkflowData map[string]interface{} `json:"workflowData"`
        Data         map[string]interface{} `json:"data,omitempty"`
}

// WorkflowTemplate represents a workflow template for creation
type WorkflowTemplate struct {
        Name        string                 `json:"name"`
        Description string                 `json:"description"`
        Category    string                 `json:"category"`
        Tags        []string               `json:"tags"`
        Nodes       []Node                 `json:"nodes"`
        Connections map[string]interface{} `json:"connections"`
        Settings    map[string]interface{} `json:"settings,omitempty"`
}

// ValidationResult represents the result of workflow validation
type ValidationResult struct {
        Valid    bool     `json:"valid"`
        Errors   []string `json:"errors,omitempty"`
        Warnings []string `json:"warnings,omitempty"`
}

// DeploymentStatus represents the status of a deployment
type DeploymentStatus struct {
        ID            string    `json:"id"`
        Environment   string    `json:"environment"`
        Status        string    `json:"status"`
        StartTime     time.Time `json:"startTime"`
        EndTime       time.Time `json:"endTime,omitempty"`
        WorkflowCount int       `json:"workflowCount"`
        ErrorMessage  string    `json:"errorMessage,omitempty"`
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
        Environment      string                `json:"environment"`
        TotalWorkflows   int                   `json:"totalWorkflows"`
        SyncedWorkflows  int                   `json:"syncedWorkflows"`
        CreatedFiles     []string              `json:"createdFiles"`
        UpdatedFiles     []string              `json:"updatedFiles"`
        SkippedFiles     []string              `json:"skippedFiles"`
        Errors           []string              `json:"errors,omitempty"`
        SyncMetadata     map[string]interface{} `json:"syncMetadata"`
}

// DeployResult represents the result of a deployment operation
type DeployResult struct {
        Environment       string   `json:"environment"`
        DeploymentID      string   `json:"deploymentId"`
        DeployedWorkflows int      `json:"deployedWorkflows"`
        CreatedWorkflows  int      `json:"createdWorkflows"`
        UpdatedWorkflows  int      `json:"updatedWorkflows"`
        ActivatedWorkflows int      `json:"activatedWorkflows"`
        FailedWorkflows   []string `json:"failedWorkflows,omitempty"`
        Errors            []string `json:"errors,omitempty"`
}

// RollbackResult represents the result of a rollback operation
type RollbackResult struct {
        Environment        string   `json:"environment"`
        RollbackID         string   `json:"rollbackId"`
        TargetDeploymentID string   `json:"targetDeploymentId"`
        RolledBackWorkflows int     `json:"rolledBackWorkflows"`
        FailedWorkflows    []string `json:"failedWorkflows,omitempty"`
        Errors             []string `json:"errors,omitempty"`
}

// WorkflowDiff represents the differences between two workflows
type WorkflowDiff struct {
        WorkflowName string                 `json:"workflowName"`
        HasChanges   bool                   `json:"hasChanges"`
        Changes      map[string]interface{} `json:"changes"`
}

// EnvironmentStatus represents the status of an environment
type EnvironmentStatus struct {
        Name              string    `json:"name"`
        URL               string    `json:"url"`
        Connected         bool      `json:"connected"`
        WorkflowCount     int       `json:"workflowCount"`
        ActiveWorkflows   int       `json:"activeWorkflows"`
        LastSync          time.Time `json:"lastSync,omitempty"`
        LastDeployment    time.Time `json:"lastDeployment,omitempty"`
        ErrorMessage      string    `json:"errorMessage,omitempty"`
}

// ProjectStatus represents the overall project status
type ProjectStatus struct {
        Environments    []EnvironmentStatus `json:"environments"`
        TotalWorkflows  int                `json:"totalWorkflows"`
        LocalFiles      int                `json:"localFiles"`
        LastActivity    time.Time          `json:"lastActivity"`
        GitStatus       GitStatus          `json:"gitStatus"`
}

// GitStatus represents Git repository status
type GitStatus struct {
        Branch          string   `json:"branch"`
        HasChanges      bool     `json:"hasChanges"`
        UntrackedFiles  []string `json:"untrackedFiles"`
        ModifiedFiles   []string `json:"modifiedFiles"`
        LastCommit      string   `json:"lastCommit"`
        LastCommitDate  time.Time `json:"lastCommitDate"`
}

// Custom JSON marshaling for time fields
func (w *Workflow) MarshalJSON() ([]byte, error) {
        type Alias Workflow
        return json.Marshal(&struct {
                *Alias
        }{
                Alias: (*Alias)(w),
        })
}

// UnmarshalJSON implements custom JSON unmarshaling
func (w *Workflow) UnmarshalJSON(data []byte) error {
        type Alias Workflow
        aux := &struct {
                *Alias
        }{
                Alias: (*Alias)(w),
        }
        return json.Unmarshal(data, &aux)
}

// Clone creates a deep copy of the workflow
func (w *Workflow) Clone() *Workflow {
        data, _ := json.Marshal(w)
        var clone Workflow
        json.Unmarshal(data, &clone)
        return &clone
}

// GetNodeByName returns a node by its name
func (w *Workflow) GetNodeByName(name string) *Node {
        for i, node := range w.Nodes {
                if node.Name == name {
                        return &w.Nodes[i]
                }
        }
        return nil
}

// GetNodesByType returns all nodes of a specific type
func (w *Workflow) GetNodesByType(nodeType string) []Node {
        var nodes []Node
        for _, node := range w.Nodes {
                if node.Type == nodeType {
                        nodes = append(nodes, node)
                }
        }
        return nodes
}

// IsActive returns whether the workflow is active
func (w *Workflow) IsActive() bool {
        return w.Active
}

// HasTag checks if the workflow has a specific tag
func (w *Workflow) HasTag(tagName string) bool {
        for _, tag := range w.Tags {
                if tag.Name == tagName {
                        return true
                }
        }
        return false
}

// AddTag adds a tag to the workflow
func (w *Workflow) AddTag(tag Tag) {
        if !w.HasTag(tag.Name) {
                w.Tags = append(w.Tags, tag)
        }
}

// RemoveTag removes a tag from the workflow
func (w *Workflow) RemoveTag(tagName string) {
        for i, tag := range w.Tags {
                if tag.Name == tagName {
                        w.Tags = append(w.Tags[:i], w.Tags[i+1:]...)
                        break
                }
        }
}

// GetStartNodes returns all start nodes (nodes that can begin execution)
func (w *Workflow) GetStartNodes() []Node {
        startNodeTypes := map[string]bool{
                "n8n-nodes-base.start":    true,
                "n8n-nodes-base.webhook":  true,
                "n8n-nodes-base.cron":     true,
                "n8n-nodes-base.interval": true,
                "n8n-nodes-base.manualTrigger": true,
        }
        
        var startNodes []Node
        for _, node := range w.Nodes {
                if startNodeTypes[node.Type] {
                        startNodes = append(startNodes, node)
                }
        }
        return startNodes
}

// Validate performs basic validation on the workflow
func (w *Workflow) Validate() error {
        if w.Name == "" {
                return fmt.Errorf("workflow name is required")
        }
        
        if len(w.Nodes) == 0 {
                return fmt.Errorf("workflow must contain at least one node")
        }
        
        // Check for duplicate node names
        nodeNames := make(map[string]bool)
        for _, node := range w.Nodes {
                if nodeNames[node.Name] {
                        return fmt.Errorf("duplicate node name: %s", node.Name)
                }
                nodeNames[node.Name] = true
                
                if node.Name == "" {
                        return fmt.Errorf("node name is required")
                }
                
                if node.Type == "" {
                        return fmt.Errorf("node type is required for node: %s", node.Name)
                }
        }
        
        return nil
}
