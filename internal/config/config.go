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

var globalConfig *Config

// InitConfig initializes the global configuration
func InitConfig() {
	globalConfig = &Config{}

	// Set defaults
	viper.SetDefault("defaults.environment", "development")
	viper.SetDefault("defaults.validation.strict", false)
	viper.SetDefault("defaults.sync.auto_backup", true)
	viper.SetDefault("defaults.deploy.dry_run", false)
	viper.SetDefault("defaults.deploy.skip_validation", false)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")

	// Unmarshal configuration
	if err := viper.Unmarshal(globalConfig); err != nil {
		fmt.Printf("Warning: Failed to parse configuration: %v\n", err)
		return
	}

	// Resolve API keys from environment variables
	for name, env := range globalConfig.Environments {
		if env.APIKeyEnv != "" {
			apiKey := os.Getenv(env.APIKeyEnv)
			if apiKey == "" {
				fmt.Printf("Warning: Environment variable %s not set for environment %s\n", env.APIKeyEnv, name)
			}
			env.APIKey = apiKey
			globalConfig.Environments[name] = env
		}
	}
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	if globalConfig == nil {
		InitConfig()
	}
	return globalConfig
}

// GetEnvironmentConfig returns configuration for a specific environment
func GetEnvironmentConfig(environment string) (*EnvironmentConfig, error) {
	config := GetConfig()

	envConfig, exists := config.Environments[environment]
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
func GetDefaultEnvironment() string {
	config := GetConfig()
	if config.Defaults.Environment != "" {
		return config.Defaults.Environment
	}
	return "development"
}

// IsStrictValidation returns whether strict validation is enabled
func IsStrictValidation() bool {
	config := GetConfig()
	return config.Defaults.Validation.Strict
}

// IsAutoBackupEnabled returns whether auto backup is enabled for sync
func IsAutoBackupEnabled() bool {
	config := GetConfig()
	return config.Defaults.Sync.AutoBackup
}

// GetLoggingConfig returns the logging configuration
func GetLoggingConfig() *LoggingConfig {
	config := GetConfig()
	return &config.Logging
}
