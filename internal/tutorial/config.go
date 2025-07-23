package tutorial

import "os"

// ConfigExists checks for a .n8n-ops.yaml configuration file in the current
// directory or the user's home directory.
func ConfigExists() bool {
	if _, err := os.Stat(".n8n-ops.yaml"); err == nil {
		return true
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	_, err = os.Stat(homeDir + "/.n8n-ops.yaml")
	return err == nil
}
