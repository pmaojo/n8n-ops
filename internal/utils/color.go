package utils

import "strings"

// StatusColor returns a canonical color name for the provided status. Similar
// statuses map to the same color to keep the UI consistent.
func StatusColor(name string) string {
	switch strings.ToLower(name) {
	case "success", "completed", "active", "healthy", "online", "ready":
		return "green"
	case "error", "failed", "inactive", "critical", "offline":
		return "red"
	case "warning", "detecting", "pending":
		return "yellow"
	default:
		return "cyan"
	}
}
