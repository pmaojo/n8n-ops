package utils

import (
	"fmt"
	"os"
	"strings"
)

// CheckN8nCredentials ensures that the base URL and API key for n8n are present
// either as generic environment variables (N8N_URL, N8N_API_KEY) or as
// environment-specific ones (N8N_<ENV>_URL, N8N_<ENV>_API_KEY).
func CheckN8nCredentials(environment string) error {
	envSuffix := strings.ToUpper(environment)

	url := os.Getenv(fmt.Sprintf("N8N_%s_URL", envSuffix))
	if url == "" {
		url = os.Getenv("N8N_URL")
	}

	apiKey := os.Getenv(fmt.Sprintf("N8N_%s_API_KEY", envSuffix))
	if apiKey == "" {
		apiKey = os.Getenv("N8N_API_KEY")
	}

	missing := []string{}
	if url == "" {
		missing = append(missing, "n8n_url")
	}
	if apiKey == "" {
		missing = append(missing, "n8n_api_key")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required credentials for %s environment: %s", environment, strings.Join(missing, ", "))
	}

	return nil
}
