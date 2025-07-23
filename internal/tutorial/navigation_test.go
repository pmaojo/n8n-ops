package tutorial

import (
	"strings"
	"testing"
)

func TestReadKey(t *testing.T) {
	cases := map[string]string{
		"q":     "q",
		"u":     "up",
		"down":  "down",
		"other": "enter",
	}

	for input, want := range cases {
		got := readKey(strings.NewReader(input + "\n"))
		if got != want {
			t.Errorf("readKey(%q)=%q, want %q", input, got, want)
		}
	}
}
