package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Environments map[string]EnvironmentConfig `mapstructure:"environments"`
	Defaults     DefaultConfig                `mapstructure:"defaults"`
	Logging      LoggingConfig                `mapstructure:"logging"`
}

type EnvironmentConfig struct {
	URL       string `mapstructure:"url"`
	APIKeyEnv string `mapstructure:"api_key_env"`
	APIKey    string // Resolved from environment variable
}

type DefaultConfig struct {
	Environment string           `mapstructure:"environment"`
	Validation  ValidationConfig `mapstructure:"validation"`
	Sync        SyncConfig       `mapstructure:"sync"`
	Deploy      DeployConfig     `mapstructure:"deploy"`
}

type ValidationConfig struct {
	Strict bool `mapstructure:"strict"`
}

type SyncConfig struct {
	AutoBackup bool `mapstructure:"auto_backup"`
}

type DeployConfig struct {
	DryRun         bool `mapstructure:"dry_run"`
	SkipValidation bool `mapstructure:"skip_validation"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	File   string `mapstructure:"file"`
}

// GetLoggingConfig returns the logging configuration from the Config instance.
// It does not modify the state of the Config and adheres to the principle of
// providing read-only access to configuration values.
func (c *Config) GetLoggingConfig() *LoggingConfig {
	if c == nil {
		return &LoggingConfig{}
	}
	return &c.Logging
}

// NewConfig creates and initializes a new Config object
func NewConfig() (*Config, error) {
	config := &Config{}

	// Set defaults
	viper.SetDefault("defaults.environment", "development")
	viper.SetDefault("defaults.validation.strict", false)
	viper.SetDefault("defaults.sync.auto_backup", true)
	viper.SetDefault("defaults.deploy.dry_run", false)
	viper.SetDefault("defaults.deploy.skip_validation", false)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")

	// Unmarshal configuration
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Resolve API keys from environment variables
	for name, env := range config.Environments {
		if env.APIKeyEnv != "" {
			apiKey := os.Getenv(env.APIKeyEnv)
			if apiKey == "" {
				fmt.Printf("Warning: Environment variable %s not set for environment %s\n", env.APIKeyEnv, name)
			}
			env.APIKey = apiKey
			config.Environments[name] = env
		}
	}
	return config, nil
}

// GetEnvironmentConfig returns configuration for a specific environment
func (c *Config) GetEnvironmentConfig(environment string) (*EnvironmentConfig, error) {
	if c == nil {
		return nil, fmt.Errorf("nil config provided")
	}

	envConfig, exists := c.Environments[environment]
	if !exists {
		return nil, fmt.Errorf("environment '%s' not found in configuration", environment)
	}

	if envConfig.URL == "" {
		return nil, fmt.Errorf("URL not configured for environment '%s'", environment)
	}

	if envConfig.APIKey == "" {
		return nil, fmt.Errorf("API key not configured for environment '%s'", environment)
	}

	return &envConfig, nil
}

// GetDefaultEnvironment returns the default environment name
func (c *Config) GetDefaultEnvironment() string {
	if c == nil {
		return "development"
	}
	if c.Defaults.Environment != "" {
		return c.Defaults.Environment
	}
	return "development"
}

// IsStrictValidation returns whether strict validation is enabled
func (c *Config) IsStrictValidation() bool {
	if c == nil {
		return false
	}
	return c.Defaults.Validation.Strict
}

// IsAutoBackupEnabled returns whether auto backup is enabled for sync
func (c *Config) IsAutoBackupEnabled() bool {
	if c == nil {
		return false
	}
	return c.Defaults.Sync.AutoBackup
}
