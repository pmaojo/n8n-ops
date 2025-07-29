package cmd

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestRunStatusJSON(t *testing.T) {
	repo := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(repo)

	if err := runCmd("git", "init"); err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	os.WriteFile("dummy.txt", []byte("hi"), 0644)
	runCmd("git", "add", "dummy.txt")
	runCmd("git", "-c", "user.email=test@example.com", "-c", "user.name=test", "commit", "-m", "init")

	reader, writer, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = writer
	err := runStatusJSON("dev")
	writer.Close()
	os.Stdout = old
	out, _ := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		t.Fatalf("runStatusJSON failed: %v", err)
	}

	var data map[string]interface{}
	if json.Unmarshal(out, &data) != nil {
		t.Fatalf("invalid json output: %s", string(out))
	}
	if data["environment"] != "dev" {
		t.Fatalf("unexpected env: %v", data["environment"])
	}
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
