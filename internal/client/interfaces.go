package client

import (
	"context"

	"github.com/pmaojo/n8n-ops/internal/workflow"
)

// WorkflowReader defines read-only operations for workflows
type WorkflowReader interface {
	GetWorkflows(ctx context.Context) ([]*workflow.Workflow, error)
	GetWorkflow(ctx context.Context, id string) (*workflow.Workflow, error)
	HealthCheck(ctx context.Context) error
}

// WorkflowWriter defines write operations for workflows
type WorkflowWriter interface {
	CreateWorkflow(ctx context.Context, wf *workflow.Workflow) (*workflow.Workflow, error)
	UpdateWorkflow(ctx context.Context, id string, wf *workflow.Workflow) (*workflow.Workflow, error)
	DeleteWorkflow(ctx context.Context, id string) error
}

// WorkflowExecutor defines execution operations
type WorkflowExecutor interface {
	ExecuteWorkflow(ctx context.Context, id string) (*workflow.ExecutionResult, error)
	GetExecution(ctx context.Context, id string) (*workflow.ExecutionResult, error)
	GetExecutions(ctx context.Context, workflowID string, status string, limit int) ([]*workflow.ExecutionResult, error)
}

// Client combines all workflow operations
type Client interface {
	WorkflowReader
	WorkflowWriter
	WorkflowExecutor
}

// HTTPTransport abstracts the HTTP client for testing
type HTTPTransport interface {
	Do(req *Request) (*Response, error)
}

// Request wraps http.Request for abstraction
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// Response wraps http.Response for abstraction
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}
