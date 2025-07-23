package utils

import (
        "bufio"
        "fmt"
        "os"
        "path/filepath"
        "strings"
)

// GetN8nCredentials retrieves n8n URL and API key for the specified environment
// using a cascading approach to check multiple naming conventions
func GetN8nCredentials(environment string) (apiURL, apiKey string) {
        envSuffix := strings.ToUpper(environment)
        
        // Try environment-specific variables first (N8N_DEVELOPMENT_URL, N8N_DEVELOPMENT_API_KEY)
        apiURL = os.Getenv(fmt.Sprintf("N8N_%s_URL", envSuffix))
        apiKey = os.Getenv(fmt.Sprintf("N8N_%s_API_KEY", envSuffix))
        
        // Fallback to short forms for common environments
        if apiURL == "" || apiKey == "" {
                switch environment {
                case "development":
                        if apiURL == "" {
                                apiURL = os.Getenv("N8N_DEV_URL")
                        }
                        if apiKey == "" {
                                apiKey = os.Getenv("N8N_DEV_API_KEY")
                        }
                case "staging":
                        if apiURL == "" {
                                apiURL = os.Getenv("N8N_STAGING_URL")
                        }
                        if apiKey == "" {
                                apiKey = os.Getenv("N8N_STAGING_API_KEY")
                        }
                case "production":
                        if apiURL == "" {
                                apiURL = os.Getenv("N8N_PROD_URL")
                        }
                        if apiKey == "" {
                                apiKey = os.Getenv("N8N_PROD_API_KEY")
                        }
                }
        }
        
        // Final fallback to generic variables (for backward compatibility)
        if apiURL == "" {
                apiURL = os.Getenv("N8N_URL")
        }
        if apiKey == "" {
                apiKey = os.Getenv("N8N_API_KEY")
        }
        
        return apiURL, apiKey
}

// GetGitLabCredentials retrieves GitLab token and project ID
func GetGitLabCredentials(environment string) (token, projectID string) {
        envSuffix := strings.ToUpper(environment)
        
        // Try environment-specific variables first
        token = os.Getenv(fmt.Sprintf("GITLAB_TOKEN_%s", envSuffix))
        projectID = os.Getenv(fmt.Sprintf("GITLAB_PROJECT_ID_%s", envSuffix))
        
        // Fallback to generic variables
        if token == "" {
                token = os.Getenv("GITLAB_TOKEN")
        }
        if projectID == "" {
                projectID = os.Getenv("GITLAB_PROJECT_ID")
                if projectID == "" {
                        projectID = os.Getenv("GITLAB_PROJECT")
                }
        }
        
        return token, projectID
}

// CheckN8nCredentials returns error message if credentials are missing
func CheckN8nCredentials(environment string) error {
        apiURL, apiKey := GetN8nCredentials(environment)
        
        if apiURL == "" || apiKey == "" {
                envSuffix := strings.ToUpper(environment)
                return fmt.Errorf(`n8n credentials not configured for %s environment

Set environment variables:
  export N8N_%s_URL=http://localhost:3001
  export N8N_%s_API_KEY=your_api_key_here

Or use short forms:
  export N8N_DEV_URL=http://localhost:3001  (for development)
  export N8N_DEV_API_KEY=your_api_key_here

Or use --demo flag for testing`, environment, envSuffix, envSuffix)
        }
        
        return nil
}

// LoadEnvFile loads environment variables from a .env file
func LoadEnvFile(filename string) error {
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
                if os.Getenv(key) == "" {
                        os.Setenv(key, value)
                }
        }
        
        return scanner.Err()
}