package storage

import (
	"os"
	"testing"
	"time"
)

func TestSQLiteWorkflowCRUD(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	db, err := NewSQLiteDB()
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	w := &WorkflowRecord{
		ID:          "wf1",
		Name:        "Test WF",
		Environment: "dev",
		FilePath:    "wf.json",
		Hash:        "abc",
		LastSync:    time.Now(),
	}

	if err := CreateWorkflowRecord(db, w); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := GetWorkflowRecord(db, w.ID, w.Environment)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Name != w.Name {
		t.Fatalf("unexpected name: %s", got.Name)
	}

	w.Name = "Updated"
	if err := UpdateWorkflowRecord(db, w); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	list, err := GetWorkflowsByEnvironment(db, w.Environment)
	if err != nil || len(list) != 1 {
		t.Fatalf("list failed: %v", err)
	}

	if err := DeleteWorkflowRecord(db, w.ID, w.Environment); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}

func TestSQLiteDeploymentCRUD(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	db, err := NewSQLiteDB()
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	d := &DeploymentRecord{
		ID:          "dep1",
		Environment: "dev",
		Status:      "in_progress",
		StartTime:   time.Now(),
	}

	if err := CreateDeploymentRecord(db, d); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	d.Status = "success"
	d.EndTime = time.Now()
	if err := UpdateDeploymentRecord(db, d); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	got, err := GetDeploymentRecord(db, d.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Status != d.Status {
		t.Fatalf("unexpected status: %s", got.Status)
	}

	hist, err := GetDeploymentHistory(db, d.Environment, 10)
	if err != nil || len(hist) != 1 {
		t.Fatalf("history failed: %v", err)
	}

	if err := DeleteDeploymentRecord(db, d.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}
