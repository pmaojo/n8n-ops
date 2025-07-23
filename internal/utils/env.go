package utils

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnvFile loads environment variables from .env file if it exists
func LoadEnvFile() {
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		// Read .env file
		file, err := os.Open(envFile)
		if err != nil {
			return
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
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				
				// Remove quotes if present
				if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
					value = value[1 : len(value)-1]
				}
				
				// Set environment variable
				os.Setenv(key, value)
			}
		}
	}
}