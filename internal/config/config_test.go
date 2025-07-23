package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// prepareConfig loads the provided YAML into viper and resets global state.
func prepareConfig(yaml string) {
	viper.Reset()
	viper.SetConfigType("yaml")
	_ = viper.ReadConfig(strings.NewReader(yaml))
}

func TestInitConfigLoadsConfiguration(t *testing.T) {
	cfgYAML := `
    environments:
      test:
        url: http://example.com
        api_key_env: TEST_KEY
    defaults:
      environment: test
      validation:
        strict: true
      sync:
        auto_backup: false
      deploy:
        dry_run: true
        skip_validation: true
    logging:
      level: debug
      format: json
    `
	prepareConfig(cfgYAML)
	t.Setenv("TEST_KEY", "secret")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("init config: %v", err)
	}

	if cfg.Defaults.Environment != "test" {
		t.Fatalf("expected default environment 'test', got '%s'", cfg.Defaults.Environment)
	}
	if !cfg.Defaults.Validation.Strict {
		t.Error("strict validation should be true")
	}
	if cfg.Defaults.Sync.AutoBackup {
		t.Error("auto backup should be false")
	}
	envCfg, err := cfg.GetEnvironmentConfig("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if envCfg.APIKey != "secret" {
		t.Errorf("expected API key 'secret', got '%s'", envCfg.APIKey)
	}
	logCfg := cfg.GetLoggingConfig()
	if logCfg.Level != "debug" || logCfg.Format != "json" {
		t.Errorf("unexpected logging config: %+v", logCfg)
	}
}

func TestGetEnvironmentConfigErrors(t *testing.T) {
	prepareConfig("environments: {}")
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("init config: %v", err)
	}
	if _, err := cfg.GetEnvironmentConfig("missing"); err == nil {
		t.Error("expected error for unknown environment")
	}
}
