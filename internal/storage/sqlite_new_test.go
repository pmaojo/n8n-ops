package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/utils"
)

func TestNewSQLiteDBCreatesFile(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := NewSQLiteDB(dbPath)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(dbPath); err != nil {
		t.Errorf("expected db file at %s", dbPath)
	}
}

func TestNewSQLiteDBDefaultPath(t *testing.T) {
	dir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(oldWD)
	os.Chdir(dir)

	db, err := NewSQLiteDB("")
	if err != nil {
		t.Fatalf("create db: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(DefaultSQLiteDBPath); err != nil {
		t.Errorf("expected db file at %s", DefaultSQLiteDBPath)
	}
}

func TestNewSQLiteDBFromEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.db")
	provider := utils.OSProvider{}
	os.Setenv(EnvSQLiteDBPath, envPath)
	defer os.Unsetenv(EnvSQLiteDBPath)

	db, err := NewSQLiteDBFromEnv(provider)
	if err != nil {
		t.Fatalf("create db from env: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(envPath); err != nil {
		t.Errorf("expected db file at %s", envPath)
	}
}
