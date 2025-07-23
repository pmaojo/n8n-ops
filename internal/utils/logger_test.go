package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func countOpenFDsFor(path string) (int, error) {
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		return 0, err
	}
	count := 0
	for _, e := range entries {
		link, err := os.Readlink(filepath.Join("/proc/self/fd", e.Name()))
		if err != nil {
			continue
		}
		if link == path {
			count++
		}
	}
	return count, nil
}

func TestSetLogFileClosesPrevious(t *testing.T) {
	if _, err := os.Stat("/proc/self/fd"); err != nil {
		t.Skip("proc filesystem not available")
	}

	dir := t.TempDir()
	file1 := filepath.Join(dir, "first.log")
	file2 := filepath.Join(dir, "second.log")

	if err := SetLogFile(file1); err != nil {
		t.Fatalf("set first log file: %v", err)
	}
	c, err := countOpenFDsFor(file1)
	if err != nil || c != 1 {
		t.Fatalf("expected 1 open fd for first file, got %d (err=%v)", c, err)
	}

	if err := SetLogFile(file2); err != nil {
		t.Fatalf("set second log file: %v", err)
	}
	c1, _ := countOpenFDsFor(file1)
	if c1 != 0 {
		t.Fatalf("first log file descriptor should be closed, got %d", c1)
	}
	c2, _ := countOpenFDsFor(file2)
	if c2 != 1 {
		t.Fatalf("expected 1 open fd for second file, got %d", c2)
	}

	if err := CloseLogFile(); err != nil {
		t.Fatalf("close second log file: %v", err)
	}
	c2After, _ := countOpenFDsFor(file2)
	if c2After != 0 {
		t.Fatalf("log file should be closed after CloseLogFile, got %d", c2After)
	}
}
