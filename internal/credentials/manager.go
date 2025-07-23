package credentials

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// CredentialManager handles environment-specific credentials
type CredentialManager struct {
	ConfigPath  string
	Environment string
}

// EnvironmentCredentials holds all credentials for a specific environment
type EnvironmentCredentials struct {
	N8nURL      string            `json:"n8n_url" yaml:"n8n_url"`
	N8nAPIKey   string            `json:"n8n_api_key" yaml:"n8n_api_key"`
	GitLabToken string            `json:"gitlab_token" yaml:"gitlab_token"`
	CustomCreds map[string]string `json:"custom_credentials" yaml:"custom_credentials"`
}

// WorkflowCredentials represents credentials used within workflows
type WorkflowCredentials struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Environment string            `json:"environment"`
	Data        map[string]string `json:"data"`
}

func NewCredentialManager(environment string) *CredentialManager {
	configPath := os.Getenv("N8N_OPS_CONFIG")
	if configPath == "" {
		configPath = filepath.Join(os.Getenv("HOME"), ".n8n-ops.yaml")
	}

	return &CredentialManager{
		ConfigPath:  configPath,
		Environment: environment,
	}
}

// GetEnvironmentCredentials retrieves credentials for the current environment
func (cm *CredentialManager) GetEnvironmentCredentials() (*EnvironmentCredentials, error) {
	viper.SetConfigFile(cm.ConfigPath)
	viper.SetConfigType("yaml")

	// Set environment variable precedence
	viper.AutomaticEnv()
	viper.SetEnvPrefix("N8N_OPS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return cm.getFromEnvironmentVariables(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	envKey := fmt.Sprintf("environments.%s", cm.Environment)

	creds := &EnvironmentCredentials{
		CustomCreds: make(map[string]string),
	}

	// Get credentials with environment variable fallback
	creds.N8nURL = cm.getCredential(envKey+".n8n_url", "N8N_URL")
	creds.N8nAPIKey = cm.getCredential(envKey+".n8n_api_key", "N8N_API_KEY")
	creds.GitLabToken = cm.getCredential(envKey+".gitlab_token", "GITLAB_TOKEN")

	// Get custom credentials
	customCredsKey := envKey + ".custom_credentials"
	if viper.IsSet(customCredsKey) {
		customCreds := viper.GetStringMapString(customCredsKey)
		for key, value := range customCreds {
			creds.CustomCreds[key] = value
		}
	}

	return creds, nil
}

// getCredential tries config file first, then environment variables
func (cm *CredentialManager) getCredential(configKey, envKey string) string {
	// Try config file first
	if viper.IsSet(configKey) {
		return viper.GetString(configKey)
	}

	// Fallback to environment variable
	envValue := os.Getenv(envKey)
	if envValue != "" {
		return envValue
	}

	// Try environment-specific environment variable
	envSpecificKey := fmt.Sprintf("%s_%s", envKey, strings.ToUpper(cm.Environment))
	return os.Getenv(envSpecificKey)
}

// getFromEnvironmentVariables creates credentials from env vars when no config file exists
func (cm *CredentialManager) getFromEnvironmentVariables() *EnvironmentCredentials {
	return &EnvironmentCredentials{
		N8nURL:      cm.getCredential("", "N8N_URL"),
		N8nAPIKey:   cm.getCredential("", "N8N_API_KEY"),
		GitLabToken: cm.getCredential("", "GITLAB_TOKEN"),
		CustomCreds: make(map[string]string),
	}
}

// ValidateCredentials ensures required credentials are present
func (cm *CredentialManager) ValidateCredentials(creds *EnvironmentCredentials) error {
	var missing []string

	if creds.N8nURL == "" {
		missing = append(missing, "n8n_url")
	}
	if creds.N8nAPIKey == "" {
		missing = append(missing, "n8n_api_key")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required credentials for %s environment: %s",
			cm.Environment, strings.Join(missing, ", "))
	}

	return nil
}

// SetCredential stores a credential for the current environment
func (cm *CredentialManager) SetCredential(key, value string) error {
	viper.SetConfigFile(cm.ConfigPath)
	viper.SetConfigType("yaml")

	// Read existing config or create new
	if err := viper.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	envKey := fmt.Sprintf("environments.%s.%s", cm.Environment, key)
	viper.Set(envKey, value)

	return viper.WriteConfig()
}

// ListCredentials shows available credential sources
func (cm *CredentialManager) ListCredentials() (map[string]string, error) {
	creds, err := cm.GetEnvironmentCredentials()
	if err != nil {
		return nil, err
	}

	sources := make(map[string]string)

	// Check sources for each credential
	sources["n8n_url"] = cm.getCredentialSource("n8n_url", creds.N8nURL)
	sources["n8n_api_key"] = cm.getCredentialSource("n8n_api_key", creds.N8nAPIKey)
	sources["gitlab_token"] = cm.getCredentialSource("gitlab_token", creds.GitLabToken)

	for key := range creds.CustomCreds {
		sources[key] = cm.getCredentialSource(key, creds.CustomCreds[key])
	}

	return sources, nil
}

func (cm *CredentialManager) getCredentialSource(key, value string) string {
	if value == "" {
		return "not_set"
	}

	envKey := fmt.Sprintf("environments.%s.%s", cm.Environment, key)
	if viper.IsSet(envKey) {
		return "config_file"
	}

	envVarKey := fmt.Sprintf("N8N_%s", strings.ToUpper(key))
	if os.Getenv(envVarKey) != "" {
		return "environment_variable"
	}

	envSpecificKey := fmt.Sprintf("%s_%s", envVarKey, strings.ToUpper(cm.Environment))
	if os.Getenv(envSpecificKey) != "" {
		return "env_specific_variable"
	}

	return "unknown"
}

// GetWorkflowCredentials retrieves credentials that workflows can use
func (cm *CredentialManager) GetWorkflowCredentials() ([]WorkflowCredentials, error) {
	// This would integrate with n8n's credential system
	// For now, return example structure

	var workflowCreds []WorkflowCredentials

	// Example workflow credentials by environment
	switch cm.Environment {
	case "development":
		workflowCreds = []WorkflowCredentials{
			{
				ID:          "smtp_dev",
				Name:        "SMTP Development",
				Type:        "smtp",
				Environment: "development",
				Data: map[string]string{
					"host": "smtp.mailtrap.io",
					"port": "587",
					"user": "dev_smtp_user",
					// password would be retrieved securely
				},
			},
			{
				ID:          "db_dev",
				Name:        "Database Development",
				Type:        "postgres",
				Environment: "development",
				Data: map[string]string{
					"host":     "localhost",
					"database": "n8n_dev",
					"user":     "dev_user",
				},
			},
		}

	case "production":
		workflowCreds = []WorkflowCredentials{
			{
				ID:          "smtp_prod",
				Name:        "SMTP Production",
				Type:        "smtp",
				Environment: "production",
				Data: map[string]string{
					"host": "smtp.sendgrid.net",
					"port": "587",
					"user": "apikey",
				},
			},
			{
				ID:          "db_prod",
				Name:        "Database Production",
				Type:        "postgres",
				Environment: "production",
				Data: map[string]string{
					"host":     "prod-db.company.com",
					"database": "n8n_production",
					"user":     "prod_user",
				},
			},
		}
	}

	return workflowCreds, nil
}

// SyncCredentialsToN8n ensures workflow credentials exist in n8n
func (cm *CredentialManager) SyncCredentialsToN8n() error {
	// This would connect to n8n API and ensure credentials exist
	// Implementation would depend on n8n's credential management API

	workflowCreds, err := cm.GetWorkflowCredentials()
	if err != nil {
		return err
	}

	fmt.Printf("📋 Found %d workflow credentials for %s environment\n",
		len(workflowCreds), cm.Environment)

	for _, cred := range workflowCreds {
		fmt.Printf("  • %s (%s)\n", cred.Name, cred.Type)
	}

	return nil
}

// ExportCredentialsTemplate creates a template config file
func (cm *CredentialManager) ExportCredentialsTemplate(outputPath string) error {
	template := map[string]interface{}{
		"environments": map[string]interface{}{
			"development": map[string]interface{}{
				"n8n_url":      "http://localhost:5678",
				"n8n_api_key":  "your_development_api_key_here",
				"gitlab_token": "your_gitlab_token_here",
				"custom_credentials": map[string]string{
					"stripe_key": "sk_test_...",
					"aws_key":    "AKIA...",
				},
			},
			"staging": map[string]interface{}{
				"n8n_url":      "https://n8n-staging.company.com",
				"n8n_api_key":  "your_staging_api_key_here",
				"gitlab_token": "your_gitlab_token_here",
				"custom_credentials": map[string]string{
					"stripe_key": "sk_test_...",
					"aws_key":    "AKIA...",
				},
			},
			"production": map[string]interface{}{
				"n8n_url":      "https://n8n.company.com",
				"n8n_api_key":  "your_production_api_key_here",
				"gitlab_token": "your_gitlab_token_here",
				"custom_credentials": map[string]string{
					"stripe_key": "sk_live_...",
					"aws_key":    "AKIA...",
				},
			},
		},
	}

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0600) // Secure permissions
}
