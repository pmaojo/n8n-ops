package dashboard

import "time"

// Workflow holds minimal workflow info for the dashboard.
type Workflow struct {
	Name   string
	Status string
	Branch string
}

// Credential holds credential status info for the dashboard.
type Credential struct {
	Name   string
	Status string
}

// UncommittedWorkflow represents a workflow with local changes.
type UncommittedWorkflow struct {
	WorkflowName string
	Status       string
	Environment  string
	FilePath     string
}

// Data aggregates dashboard information.
type Data struct {
	Environment           string
	Status                string
	Workflows             []Workflow
	Credentials           []Credential
	LastSync              time.Time
	HasUncommittedChanges bool
	UncommittedWorkflows  []UncommittedWorkflow
	GitBranch             string
}

// Service defines the behavior required to obtain dashboard data.
type Service interface {
	GetDashboardData(env string) (Data, error)
}
