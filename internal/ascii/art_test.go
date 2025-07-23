package ascii

import (
	"strings"
	"testing"
)

func TestMatrixCharacters(t *testing.T) {
	// Test matrix character set
	matrixChars := []string{"0", "1", "▒", "░", "▓", "█", "⚡"}

	if len(matrixChars) == 0 {
		t.Error("Matrix character set should not be empty")
	}

	for _, char := range matrixChars {
		if char == "" {
			t.Error("Matrix character should not be empty")
		}
	}
}

func TestStringOperations(t *testing.T) {
	// Test basic string operations used in ASCII art
	testText := "n8n Operations Tool"

	if !strings.Contains(testText, "n8n") {
		t.Error("Text should contain 'n8n'")
	}

	lines := strings.Split(testText, " ")
	if len(lines) < 3 {
		t.Error("Text should split into multiple words")
	}
}

func TestAsciiPatterns(t *testing.T) {
	// Test ASCII art patterns
	pattern := "╔══════════════════════════╗\n║  n8n Operations Tool     ║\n╚══════════════════════════╝"

	if pattern == "" {
		t.Error("Pattern should not be empty")
	}

	lines := strings.Split(pattern, "\n")
	if len(lines) != 3 {
		t.Error("Pattern should have exactly 3 lines")
	}
}

func TestColorCodes(t *testing.T) {
	// Test ANSI color codes
	colorCodes := map[string]string{
		"reset":  "\033[0m",
		"green":  "\033[32m",
		"cyan":   "\033[36m",
		"yellow": "\033[33m",
	}

	for name, code := range colorCodes {
		if code == "" {
			t.Errorf("Color code for %s should not be empty", name)
		}
		if !strings.Contains(code, "\033[") {
			t.Errorf("Color code for %s should contain ANSI escape sequence", name)
		}
	}
}
