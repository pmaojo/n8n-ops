package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetGitLabCredentials retrieves GitLab token and project ID
func GetGitLabCredentials(provider EnvProvider, environment string) (token, projectID string) {
	envSuffix := strings.ToUpper(environment)

	// Try environment-specific variables first
	token = provider.Getenv(fmt.Sprintf("GITLAB_TOKEN_%s", envSuffix))
	projectID = provider.Getenv(fmt.Sprintf("GITLAB_PROJECT_ID_%s", envSuffix))

	// Fallback to generic variables
	if token == "" {
		token = provider.Getenv("GITLAB_TOKEN")
	}
	if projectID == "" {
		projectID = provider.Getenv("GITLAB_PROJECT_ID")
		if projectID == "" {
			projectID = provider.Getenv("GITLAB_PROJECT")
		}
	}

	return token, projectID
}

// LoadEnvFile loads environment variables from a .env file
func LoadEnvFile(provider EnvProvider, filename string) error {
	if !filepath.IsAbs(filename) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		filename = filepath.Join(wd, filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// .env file is optional
			return nil
		}
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		// Only set if not already in environment
		if provider.Getenv(key) == "" {
			provider.Setenv(key, value)
		}
	}

	return scanner.Err()
}
