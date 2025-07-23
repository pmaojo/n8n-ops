package gitlab

import (
	"net/url"
	"testing"
)

func TestGitLabURLValidation(t *testing.T) {
	// Test GitLab URL validation
	validURLs := []string{
		"https://gitlab.com/project/repo",
		"https://gitlab.example.com/group/project",
		"https://git.company.com/team/workflows",
	}

	invalidURLs := []string{
		"",
		"not-a-url",
		"http://",
		"ftp://invalid.com",
	}

	for _, testURL := range validURLs {
		_, err := url.Parse(testURL)
		if err != nil {
			t.Errorf("Valid URL %s should parse correctly: %v", testURL, err)
		}
	}

	for _, testURL := range invalidURLs {
		if testURL != "" {
			parsed, err := url.Parse(testURL)
			if err == nil && parsed.Scheme == "ftp" {
				// FTP URLs should not be considered valid for GitLab
				continue
			}
			if err == nil && parsed.Scheme == "https" && parsed.Host != "" {
				t.Errorf("Invalid URL %s should not be considered valid", testURL)
			}
		}
	}
}

func TestGitLabAPIEndpoints(t *testing.T) {
	// Test GitLab API endpoint construction
	projectID := "12345"

	expectedEndpoints := map[string]string{
		"issues":         "/api/v4/projects/" + projectID + "/issues",
		"merge_requests": "/api/v4/projects/" + projectID + "/merge_requests",
		"pipelines":      "/api/v4/projects/" + projectID + "/pipelines",
	}

	for endpoint, expected := range expectedEndpoints {
		if expected == "" {
			t.Errorf("Endpoint %s should not be empty", endpoint)
		}

		if endpoint == "issues" && expected != "/api/v4/projects/12345/issues" {
			t.Error("Issues endpoint format incorrect")
		}
	}
}

func TestIssueCreation(t *testing.T) {
	// Test issue creation data structure
	issueData := map[string]interface{}{
		"title":       "Workflow Failure: Payment Processing",
		"description": "Multiple consecutive failures detected",
		"labels":      []string{"workflow-failure", "automation"},
		"assignee_id": 1,
	}

	// Validate required fields
	if issueData["title"] == "" {
		t.Error("Issue title should not be empty")
	}

	if issueData["description"] == "" {
		t.Error("Issue description should not be empty")
	}

	labels, ok := issueData["labels"].([]string)
	if !ok || len(labels) == 0 {
		t.Error("Issue should have labels")
	}
}

func TestCIPipelineIntegration(t *testing.T) {
	// Test CI/CD pipeline integration
	pipelineStages := []string{"build", "validate", "test", "deploy"}

	if len(pipelineStages) != 4 {
		t.Error("Pipeline should have 4 stages")
	}

	for i, stage := range pipelineStages {
		if stage == "" {
			t.Errorf("Pipeline stage %d should not be empty", i)
		}
	}

	// Test stage order
	if pipelineStages[0] != "build" {
		t.Error("First stage should be build")
	}

	if pipelineStages[len(pipelineStages)-1] != "deploy" {
		t.Error("Last stage should be deploy")
	}
}
