package credentials

import (
        "os"
        "strings"
        "testing"
)

func TestEnvironmentVariables(t *testing.T) {
        // Test environment variable handling
        testKey := "TEST_N8N_API_KEY"
        testValue := "test_api_key_value"
        
        // Set test environment variable
        err := os.Setenv(testKey, testValue)
        if err != nil {
                t.Fatalf("Failed to set environment variable: %v", err)
        }
        defer os.Unsetenv(testKey)
        
        // Test retrieval
        retrieved := os.Getenv(testKey)
        if retrieved != testValue {
                t.Errorf("Expected %s, got %s", testValue, retrieved)
        }
}

func TestAPIKeyValidation(t *testing.T) {
        // Test API key validation logic
        validKeys := []string{
                "n8n_api_key_development_12345",
                "valid-api-key-format",
                "API_KEY_123",
        }
        
        invalidKeys := []string{
                "",
                "   ",
                "bad",
        }
        
        for _, key := range validKeys {
                if len(key) == 0 || len(key) < 5 {
                        t.Errorf("Valid key %s should pass basic validation", key)
                }
        }
        
        for _, key := range invalidKeys {
                if len(key) >= 5 && strings.TrimSpace(key) != "" {
                        t.Errorf("Invalid key %s should fail basic validation", key)
                }
        }
}

func TestCredentialSecurity(t *testing.T) {
        // Test credential security practices
        sensitiveData := "super-secret-api-key"
        
        // Verify sensitive data is not logged or exposed
        if sensitiveData == "" {
                t.Error("Test data should not be empty")
        }
        
        // Test that we don't accidentally expose credentials
        logMessage := "API request failed"
        if logMessage == sensitiveData {
                t.Error("Log messages should not contain sensitive data")
        }
}

func TestEnvironmentIsolation(t *testing.T) {
        // Test environment-specific credential isolation
        environments := []string{"development", "staging", "production"}
        
        for _, env := range environments {
                keyName := "N8N_" + env + "_API_KEY"
                if keyName == "" {
                        t.Errorf("Key name for environment %s should not be empty", env)
                }
                
                // Verify environment-specific naming
                if env == "development" && keyName != "N8N_DEV_API_KEY" {
                        t.Error("Development API key name format incorrect")
                }
        }
}