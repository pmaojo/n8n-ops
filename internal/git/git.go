package git

import "strings"

// GetEnvironmentFromBranch maps git branches to n8n environments
func GetEnvironmentFromBranch(branch string) string {
	switch {
	case strings.Contains(branch, "dev") || strings.Contains(branch, "develop"):
		return "development"
	case strings.Contains(branch, "staging") || strings.Contains(branch, "stage"):
		return "staging"
	case strings.Contains(branch, "main") || strings.Contains(branch, "master") || strings.Contains(branch, "prod"):
		return "production"
	default:
		return "development" // default to development for feature branches
	}
}

// GetBranchFromEnvironment maps n8n environments to suggested git branches
func GetBranchFromEnvironment(env string) string {
	switch env {
	case "development":
		return "develop"
	case "staging":
		return "staging"
	case "production":
		return "main"
	default:
		return "develop"
	}
}
