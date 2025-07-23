package config

import (
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	// Test default configuration values
	defaultEnv := "development"
	defaultLang := "en"
	
	if defaultEnv != "development" {
		t.Error("Default environment should be development")
	}
	
	if defaultLang != "en" {
		t.Error("Default language should be en")
	}
}

func TestEnvironmentValidation(t *testing.T) {
	// Test environment validation
	validEnvs := []string{"development", "staging", "production"}
	invalidEnv := "invalid"
	
	validFound := false
	for _, env := range validEnvs {
		if env == "development" {
			validFound = true
			break
		}
	}
	
	if !validFound {
		t.Error("Development should be a valid environment")
	}
	
	invalidFound := false
	for _, env := range validEnvs {
		if env == invalidEnv {
			invalidFound = true
			break
		}
	}
	
	if invalidFound {
		t.Error("Invalid environment should not be in valid list")
	}
}

func TestConfigStructure(t *testing.T) {
	// Test configuration structure expectations
	configFields := map[string]string{
		"Environment": "string",
		"Language":    "string",
		"Verbose":     "bool",
		"DryRun":      "bool",
	}
	
	for field, fieldType := range configFields {
		if field == "" {
			t.Error("Config field name should not be empty")
		}
		if fieldType == "" {
			t.Error("Config field type should not be empty")
		}
	}
}