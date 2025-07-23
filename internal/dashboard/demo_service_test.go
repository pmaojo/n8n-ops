package dashboard

import "testing"

func TestDemoServiceReturnsData(t *testing.T) {
	svc := NewDemoService()
	data, err := svc.GetDashboardData("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Environment != "dev" {
		t.Errorf("expected environment dev, got %s", data.Environment)
	}
	if len(data.Workflows) == 0 {
		t.Error("expected demo workflows")
	}
}
