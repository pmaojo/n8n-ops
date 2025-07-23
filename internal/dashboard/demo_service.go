package dashboard

import "time"

// DemoService provides hard-coded dashboard data for demo purposes.
type DemoService struct{}

// NewDemoService returns a new DemoService.
func NewDemoService() *DemoService {
	return &DemoService{}
}

// GetDashboardData returns static data used for the demo dashboard.
func (d *DemoService) GetDashboardData(env string) (Data, error) {
	// In the future, replace with real implementations.
	data := Data{
		Environment: env,
		Status:      "connected",
		LastSync:    time.Now(),
		GitBranch:   "main",
		Workflows: []Workflow{
			{Name: "Customer Onboarding Process", Status: "active", Branch: "main"},
			{Name: "Email Notification System", Status: "active", Branch: "main"},
			{Name: "Payment Processing", Status: "active", Branch: "main"},
		},
		Credentials: []Credential{
			{Name: "N8N API", Status: "missing"},
			{Name: "SMTP", Status: "missing"},
			{Name: "PostgreSQL", Status: "missing"},
			{Name: "Stripe", Status: "missing"},
		},
		HasUncommittedChanges: false,
		UncommittedWorkflows:  nil,
	}

	return data, nil
}
