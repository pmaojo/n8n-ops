package credentials

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CredentialManager handles environment-specific credentials
type CredentialManager struct {
	ConfigPath  string
	Environment string
	Logger      logrus.FieldLogger
}

// newViper returns a Viper instance configured for credential operations.
// It does not read the configuration file, allowing callers to decide when
// loading should occur.
func newViper(configPath string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvPrefix("N8N_OPS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
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

// NewCredentialManager initializes a CredentialManager for the given
// environment. The configuration file path is resolved from the environment and
// no filesystem operations are performed during construction.
func NewCredentialManager(environment string, logger logrus.FieldLogger) *CredentialManager {
	configPath := os.Getenv("N8N_OPS_CONFIG")
	if configPath == "" {
		configPath = filepath.Join(os.Getenv("HOME"), ".n8n-ops.yaml")
	}

	if logger == nil {
		logger = logrus.New()
	}

	return &CredentialManager{
		ConfigPath:  configPath,
		Environment: environment,
		Logger:      logger,
	}
}

// loadConfig creates a new Viper instance and reads the configuration file.
// It sets environment variable handling consistent with other credential methods.
func (cm *CredentialManager) loadConfig() (*viper.Viper, error) {
	v := newViper(cm.ConfigPath)
	if err := v.ReadInConfig(); err != nil {
		return v, err
	}
	return v, nil
}

// GetEnvironmentCredentials retrieves credentials for the current environment
func (cm *CredentialManager) GetEnvironmentCredentials() (*EnvironmentCredentials, error) {
	v := newViper(cm.ConfigPath)
	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return cm.getFromEnvironmentVariables(nil), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	envKey := fmt.Sprintf("environments.%s", cm.Environment)

	creds := &EnvironmentCredentials{
		CustomCreds: make(map[string]string),
	}

	// Get credentials with environment variable fallback
	creds.N8nURL = cm.getCredential(v, envKey+".n8n_url", "N8N_URL")
	creds.N8nAPIKey = cm.getCredential(v, envKey+".n8n_api_key", "N8N_API_KEY")
	creds.GitLabToken = cm.getCredential(v, envKey+".gitlab_token", "GITLAB_TOKEN")

	// Get custom credentials
	customCredsKey := envKey + ".custom_credentials"
	if v.IsSet(customCredsKey) {
		customCreds := v.GetStringMapString(customCredsKey)
		for key, value := range customCreds {
			creds.CustomCreds[key] = value
		}
	}

	return creds, nil
}

// getCredential tries config file first, then environment variables
func (cm *CredentialManager) getCredential(v *viper.Viper, configKey, envKey string) string {
	// Try config file first when available
	if v != nil && configKey != "" && v.IsSet(configKey) {
		return v.GetString(configKey)
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
func (cm *CredentialManager) getFromEnvironmentVariables(v *viper.Viper) *EnvironmentCredentials {
	return &EnvironmentCredentials{
		N8nURL:      cm.getCredential(v, "", "N8N_URL"),
		N8nAPIKey:   cm.getCredential(v, "", "N8N_API_KEY"),
		GitLabToken: cm.getCredential(v, "", "GITLAB_TOKEN"),
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

// GetN8nCredentials returns the n8n URL and API key for the current environment.
func (cm *CredentialManager) GetN8nCredentials() (string, string, error) {
	creds, err := cm.GetEnvironmentCredentials()
	if err != nil {
		return "", "", err
	}
	return creds.N8nURL, creds.N8nAPIKey, nil
}

// SetCredential stores a credential for the current environment
func (cm *CredentialManager) SetCredential(key, value string) error {
	v := newViper(cm.ConfigPath)

	// Read existing config or create new
	if err := v.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	envKey := fmt.Sprintf("environments.%s.%s", cm.Environment, key)
	v.Set(envKey, value)

	return v.WriteConfig()
}

// ListCredentials shows available credential sources
func (cm *CredentialManager) ListCredentials() (map[string]string, error) {
	creds, err := cm.GetEnvironmentCredentials()
	if err != nil {
		return nil, err
	}

	cfg, _ := cm.loadConfig()

	sources := make(map[string]string)

	// Check sources for each credential
	sources["n8n_url"] = cm.getCredentialSource(cfg, "n8n_url", creds.N8nURL)
	sources["n8n_api_key"] = cm.getCredentialSource(cfg, "n8n_api_key", creds.N8nAPIKey)
	sources["gitlab_token"] = cm.getCredentialSource(cfg, "gitlab_token", creds.GitLabToken)

	for key := range creds.CustomCreds {
		sources[key] = cm.getCredentialSource(cfg, key, creds.CustomCreds[key])
	}

	return sources, nil
}

func (cm *CredentialManager) getCredentialSource(v *viper.Viper, key, value string) string {
	if value == "" {
		return "not_set"
	}

	envKey := fmt.Sprintf("environments.%s.%s", cm.Environment, key)
	if v != nil && v.IsSet(envKey) {
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
	cfg, err := cm.loadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return []WorkflowCredentials{}, nil
		}
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	key := fmt.Sprintf("environments.%s.workflow_credentials", cm.Environment)
	if !cfg.IsSet(key) {
		return []WorkflowCredentials{}, nil
	}

	var creds []WorkflowCredentials
	if err := cfg.UnmarshalKey(key, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse workflow credentials: %w", err)
	}

	return creds, nil
}

// SyncCredentialsToN8n ensures workflow credentials exist in n8n
func (cm *CredentialManager) SyncCredentialsToN8n() error {
	// This would connect to n8n API and ensure credentials exist
	// Implementation would depend on n8n's credential management API

	workflowCreds, err := cm.GetWorkflowCredentials()
	if err != nil {
		return err
	}

	cm.Logger.WithFields(logrus.Fields{
		"count":       len(workflowCreds),
		"environment": cm.Environment,
	}).Info("found workflow credentials")

	for _, cred := range workflowCreds {
		cm.Logger.WithFields(logrus.Fields{
			"name": cred.Name,
			"type": cred.Type,
		}).Info("workflow credential")
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
				"workflow_credentials": []map[string]interface{}{
					{
						"id":   "smtp_dev",
						"name": "SMTP Development",
						"type": "smtp",
						"data": map[string]string{
							"host": "smtp.mailtrap.io",
							"port": "587",
							"user": "dev_smtp_user",
						},
					},
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
				"workflow_credentials": []map[string]interface{}{
					{
						"id":   "smtp_staging",
						"name": "SMTP Staging",
						"type": "smtp",
						"data": map[string]string{
							"host": "smtp-staging.company.com",
							"port": "587",
							"user": "staging_user",
						},
					},
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
				"workflow_credentials": []map[string]interface{}{
					{
						"id":   "smtp_prod",
						"name": "SMTP Production",
						"type": "smtp",
						"data": map[string]string{
							"host": "smtp.sendgrid.net",
							"port": "587",
							"user": "apikey",
						},
					},
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
