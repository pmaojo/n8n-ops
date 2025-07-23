package bubbleui

// WorkflowStatus describes minimal workflow data for display.
type WorkflowStatus struct {
	ID     string
	Name   string
	Status string
}

// Metrics holds dashboard metric data.
type Metrics struct {
	Failures  int
	Issues    int
	Workflows int
	Uptime    string
}

// Summary provides a high level count of active vs inactive workflows.
type Summary struct {
	Active   int
	Inactive int
}
