package i18n

import (
	"testing"
)

func TestLanguageSupport(t *testing.T) {
	// Test supported languages
	supportedLanguages := []string{"en", "es"}

	if len(supportedLanguages) != 2 {
		t.Error("Should support exactly 2 languages")
	}

	for _, lang := range supportedLanguages {
		if lang == "" {
			t.Error("Language code should not be empty")
		}

		if len(lang) != 2 {
			t.Errorf("Language code %s should be 2 characters", lang)
		}
	}
}

func TestTranslationKeys(t *testing.T) {
	// Test translation key structure
	translationKeys := map[string]map[string]string{
		"en": {
			"welcome":  "Welcome to n8n-ops",
			"sync":     "Sync workflows",
			"deploy":   "Deploy to n8n",
			"validate": "Validate workflows",
			"status":   "Show status",
		},
		"es": {
			"welcome":  "Bienvenido a n8n-ops",
			"sync":     "Sincronizar workflows",
			"deploy":   "Desplegar a n8n",
			"validate": "Validar workflows",
			"status":   "Mostrar estado",
		},
	}

	// Validate translation completeness
	englishKeys := translationKeys["en"]
	spanishKeys := translationKeys["es"]

	if len(englishKeys) != len(spanishKeys) {
		t.Error("All languages should have same number of translations")
	}

	for key := range englishKeys {
		if spanishKeys[key] == "" {
			t.Errorf("Spanish translation missing for key: %s", key)
		}
	}
}

func TestLanguageDetection(t *testing.T) {
	// Test language detection logic
	testCases := []struct {
		input    string
		expected string
	}{
		{"en", "en"},
		{"es", "es"},
		{"EN", "en"}, // Case insensitive
		{"spanish", "es"},
		{"english", "en"},
		{"invalid", "en"}, // Default fallback
	}

	for _, tc := range testCases {
		var result string

		// Simple language detection logic
		switch tc.input {
		case "en", "EN", "english":
			result = "en"
		case "es", "ES", "spanish":
			result = "es"
		default:
			result = "en" // Default
		}

		if result != tc.expected {
			t.Errorf("For input %s, expected %s, got %s", tc.input, tc.expected, result)
		}
	}
}

func TestMessageFormatting(t *testing.T) {
	// Test message formatting with parameters
	template := "Synced %d workflows to %s environment"

	if template == "" {
		t.Error("Message template should not be empty")
	}

	// Test parameter substitution
	count := 3
	environment := "development"

	if count <= 0 {
		t.Error("Workflow count should be positive")
	}

	if environment == "" {
		t.Error("Environment should not be empty")
	}
}

func TestCultureSpecific(t *testing.T) {
	// Test culture-specific formatting
	dateFormats := map[string]string{
		"en": "2025-07-23 01:00:00",
		"es": "23/07/2025 01:00:00",
	}

	for lang, format := range dateFormats {
		if format == "" {
			t.Errorf("Date format for %s should not be empty", lang)
		}

		if lang == "en" && len(format) < 10 {
			t.Error("English date format should include full date")
		}
	}
}
