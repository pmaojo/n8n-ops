package bubbleui

import "strings"

// filterWorkflows returns workflows whose name contains the filter text.
func filterWorkflows(workflows []WorkflowStatus, filter string) []WorkflowStatus {
	if filter == "" {
		return workflows
	}
	lower := strings.ToLower(filter)
	result := make([]WorkflowStatus, 0, len(workflows))
	for _, wf := range workflows {
		if strings.Contains(strings.ToLower(wf.Name), lower) {
			result = append(result, wf)
		}
	}
	return result
}
