package utils

import (
	"fmt"
	"strings"
)

// CheckN8nCredentials ensures that the base URL and API key for n8n are present
// either as generic environment variables (N8N_URL, N8N_API_KEY) or as
// environment-specific ones using the preferred `N8N_URL_<ENV>` and
// `N8N_API_KEY_<ENV>` format. Legacy variables `N8N_<ENV>_URL` and
// `N8N_<ENV>_API_KEY` are also supported for backwards compatibility.
func CheckN8nCredentials(provider EnvProvider, environment string) error {
	envSuffix := strings.ToUpper(environment)

	url := provider.Getenv(fmt.Sprintf("N8N_URL_%s", envSuffix))
	if url == "" {
		url = provider.Getenv(fmt.Sprintf("N8N_%s_URL", envSuffix))
	}
	if url == "" {
		url = provider.Getenv("N8N_URL")
	}

	apiKey := provider.Getenv(fmt.Sprintf("N8N_API_KEY_%s", envSuffix))
	if apiKey == "" {
		apiKey = provider.Getenv(fmt.Sprintf("N8N_%s_API_KEY", envSuffix))
	}
	if apiKey == "" {
		apiKey = provider.Getenv("N8N_API_KEY")
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

// EnvVarNames contains the ordered environment variable names for API keys and URLs.
type EnvVarNames struct {
	APIKey []string
	URL    []string
}

// BuildEnvVarNames returns the preferred and legacy environment variable names
// for the provided environment. The first element in each slice is the
// canonical name, followed by fallback names.
func BuildEnvVarNames(environment string) EnvVarNames {
	envUpper := strings.ToUpper(environment)
	aliasMap := map[string]string{
		"development": "DEV",
		"production":  "PROD",
	}
	alias := aliasMap[environment]
	if alias == "" {
		alias = envUpper
	}

	names := EnvVarNames{
		APIKey: []string{
			fmt.Sprintf("N8N_API_KEY_%s", envUpper),
			fmt.Sprintf("N8N_%s_API_KEY", envUpper),
		},
		URL: []string{
			fmt.Sprintf("N8N_URL_%s", envUpper),
			fmt.Sprintf("N8N_%s_URL", envUpper),
		},
	}

	if alias != envUpper {
		names.APIKey = append(names.APIKey, fmt.Sprintf("N8N_API_KEY_%s", alias))
		names.URL = append(names.URL, fmt.Sprintf("N8N_URL_%s", alias))
	}

	return names
}
