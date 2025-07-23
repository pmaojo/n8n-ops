package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var ErrRecordNotFound = errors.New("record not found")

type SQLiteDB struct {
	db *sql.DB
}

type WorkflowRecord struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Environment  string    `db:"environment"`
	FilePath     string    `db:"file_path"`
	Hash         string    `db:"hash"`
	LastSync     time.Time `db:"last_sync"`
	LastDeploy   time.Time `db:"last_deploy"`
	Version      int       `db:"version"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type DeploymentRecord struct {
	ID            string    `db:"id"`
	Environment   string    `db:"environment"`
	Status        string    `db:"status"` // in_progress, success, failed
	StartTime     time.Time `db:"start_time"`
	EndTime       time.Time `db:"end_time"`
	WorkflowCount int       `db:"workflow_count"`
	GitCommit     string    `db:"git_commit"`
	DeployedBy    string    `db:"deployed_by"`
	ErrorMessage  string    `db:"error_message"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB() (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", ".n8n-ops.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqliteDB := &SQLiteDB{db: db}
	
	// Initialize tables
	if err := sqliteDB.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return sqliteDB, nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// initTables creates the required tables if they don't exist
func (s *SQLiteDB) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS workflows (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			environment TEXT NOT NULL,
			file_path TEXT NOT NULL,
			hash TEXT NOT NULL,
			last_sync DATETIME,
			last_deploy DATETIME,
			version INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(name, environment)
		)`,
		`CREATE TABLE IF NOT EXISTS deployments (
			id TEXT PRIMARY KEY,
			environment TEXT NOT NULL,
			status TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			workflow_count INTEGER DEFAULT 0,
			git_commit TEXT,
			deployed_by TEXT,
			error_message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_workflows_env ON workflows(environment)`,
		`CREATE INDEX IF NOT EXISTS idx_workflows_name ON workflows(name)`,
		`CREATE INDEX IF NOT EXISTS idx_deployments_env ON deployments(environment)`,
		`CREATE INDEX IF NOT EXISTS idx_deployments_time ON deployments(start_time)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// CreateWorkflowRecord creates a new workflow record
func CreateWorkflowRecord(db *SQLiteDB, record *WorkflowRecord) error {
	query := `
		INSERT INTO workflows (id, name, environment, file_path, hash, last_sync, last_deploy, version)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := db.db.Exec(query,
		record.ID,
		record.Name,
		record.Environment,
		record.FilePath,
		record.Hash,
		record.LastSync,
		record.LastDeploy,
		record.Version,
	)
	
	return err
}

// UpdateWorkflowRecord updates an existing workflow record
func UpdateWorkflowRecord(db *SQLiteDB, record *WorkflowRecord) error {
	query := `
		UPDATE workflows 
		SET name = ?, file_path = ?, hash = ?, last_sync = ?, last_deploy = ?, version = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND environment = ?
	`
	
	result, err := db.db.Exec(query,
		record.Name,
		record.FilePath,
		record.Hash,
		record.LastSync,
		record.LastDeploy,
		record.Version,
		record.ID,
		record.Environment,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	
	return nil
}

// GetWorkflowRecord retrieves a workflow record by ID and environment
func GetWorkflowRecord(db *SQLiteDB, id, environment string) (*WorkflowRecord, error) {
	query := `
		SELECT id, name, environment, file_path, hash, last_sync, last_deploy, version, created_at, updated_at
		FROM workflows 
		WHERE id = ? AND environment = ?
	`
	
	record := &WorkflowRecord{}
	err := db.db.QueryRow(query, id, environment).Scan(
		&record.ID,
		&record.Name,
		&record.Environment,
		&record.FilePath,
		&record.Hash,
		&record.LastSync,
		&record.LastDeploy,
		&record.Version,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	
	return record, nil
}

// GetWorkflowsByEnvironment retrieves all workflows for an environment
func GetWorkflowsByEnvironment(db *SQLiteDB, environment string) ([]*WorkflowRecord, error) {
	query := `
		SELECT id, name, environment, file_path, hash, last_sync, last_deploy, version, created_at, updated_at
		FROM workflows 
		WHERE environment = ?
		ORDER BY name
	`
	
	rows, err := db.db.Query(query, environment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*WorkflowRecord
	for rows.Next() {
		record := &WorkflowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Environment,
			&record.FilePath,
			&record.Hash,
			&record.LastSync,
			&record.LastDeploy,
			&record.Version,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	
	return records, nil
}

// CreateDeploymentRecord creates a new deployment record
func CreateDeploymentRecord(db *SQLiteDB, record *DeploymentRecord) error {
	query := `
		INSERT INTO deployments (id, environment, status, start_time, end_time, workflow_count, git_commit, deployed_by, error_message)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := db.db.Exec(query,
		record.ID,
		record.Environment,
		record.Status,
		record.StartTime,
		record.EndTime,
		record.WorkflowCount,
		record.GitCommit,
		record.DeployedBy,
		record.ErrorMessage,
	)
	
	return err
}

// UpdateDeploymentRecord updates an existing deployment record
func UpdateDeploymentRecord(db *SQLiteDB, record *DeploymentRecord) error {
	query := `
		UPDATE deployments 
		SET status = ?, end_time = ?, workflow_count = ?, error_message = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	
	_, err := db.db.Exec(query,
		record.Status,
		record.EndTime,
		record.WorkflowCount,
		record.ErrorMessage,
		record.ID,
	)
	
	return err
}

// GetDeploymentRecord retrieves a deployment record by ID
func GetDeploymentRecord(db *SQLiteDB, id string) (*DeploymentRecord, error) {
	query := `
		SELECT id, environment, status, start_time, end_time, workflow_count, git_commit, deployed_by, error_message, created_at, updated_at
		FROM deployments 
		WHERE id = ?
	`
	
	record := &DeploymentRecord{}
	err := db.db.QueryRow(query, id).Scan(
		&record.ID,
		&record.Environment,
		&record.Status,
		&record.StartTime,
		&record.EndTime,
		&record.WorkflowCount,
		&record.GitCommit,
		&record.DeployedBy,
		&record.ErrorMessage,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	
	return record, nil
}

// GetDeploymentHistory retrieves deployment history for an environment
func GetDeploymentHistory(db *SQLiteDB, environment string, limit int) ([]*DeploymentRecord, error) {
	query := `
		SELECT id, environment, status, start_time, end_time, workflow_count, git_commit, deployed_by, error_message, created_at, updated_at
		FROM deployments 
		WHERE environment = ?
		ORDER BY start_time DESC
		LIMIT ?
	`
	
	rows, err := db.db.Query(query, environment, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*DeploymentRecord
	for rows.Next() {
		record := &DeploymentRecord{}
		err := rows.Scan(
			&record.ID,
			&record.Environment,
			&record.Status,
			&record.StartTime,
			&record.EndTime,
			&record.WorkflowCount,
			&record.GitCommit,
			&record.DeployedBy,
			&record.ErrorMessage,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	
	return records, nil
}

// GetWorkflowsByDeployment retrieves workflows that were part of a specific deployment
// This is a simplified implementation - in a more complex system, you might have a deployment_workflows junction table
func GetWorkflowsByDeployment(db *SQLiteDB, deploymentID string) ([]*WorkflowRecord, error) {
	// For simplicity, we'll get the deployment record and then find workflows updated around that time
	deployment, err := GetDeploymentRecord(db, deploymentID)
	if err != nil {
		return nil, err
	}
	
	// Get workflows that were deployed within a time window around the deployment
	query := `
		SELECT id, name, environment, file_path, hash, last_sync, last_deploy, version, created_at, updated_at
		FROM workflows 
		WHERE environment = ? 
		AND last_deploy BETWEEN datetime(?, '-5 minutes') AND datetime(?, '+5 minutes')
		ORDER BY name
	`
	
	rows, err := db.db.Query(query, 
		deployment.Environment, 
		deployment.StartTime.Format("2006-01-02 15:04:05"),
		deployment.EndTime.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*WorkflowRecord
	for rows.Next() {
		record := &WorkflowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Environment,
			&record.FilePath,
			&record.Hash,
			&record.LastSync,
			&record.LastDeploy,
			&record.Version,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	
	return records, nil
}

// DeleteWorkflowRecord deletes a workflow record
func DeleteWorkflowRecord(db *SQLiteDB, id, environment string) error {
	query := `DELETE FROM workflows WHERE id = ? AND environment = ?`
	_, err := db.db.Exec(query, id, environment)
	return err
}

// DeleteDeploymentRecord deletes a deployment record
func DeleteDeploymentRecord(db *SQLiteDB, id string) error {
	query := `DELETE FROM deployments WHERE id = ?`
	_, err := db.db.Exec(query, id)
	return err
}
