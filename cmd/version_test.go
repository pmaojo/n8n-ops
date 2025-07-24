package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	runVersion(nil, nil)
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "n8n-ops version") {
		t.Errorf("unexpected output: %s", out)
	}
}
