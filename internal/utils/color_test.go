package utils

import "testing"

func TestStatusColor(t *testing.T) {
	cases := map[string]string{
		"active":   "green",
		"error":    "red",
		"warning":  "yellow",
		"whatever": "cyan",
	}

	for status, expected := range cases {
		if got := StatusColor(status); got != expected {
			t.Fatalf("status %s expected %s, got %s", status, expected, got)
		}
	}
}
