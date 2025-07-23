package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "strconv"
        "time"

        "github.com/gorilla/mux"
)

type Workflow struct {
        ID          string            `json:"id"`
        Name        string            `json:"name"`
        Active      bool              `json:"active"`
        Nodes       []WorkflowNode    `json:"nodes"`
        Connections map[string]interface{} `json:"connections"`
        CreatedAt   time.Time         `json:"createdAt"`
        UpdatedAt   time.Time         `json:"updatedAt"`
        VersionId   int               `json:"versionId"`
        Tags        []string          `json:"tags"`
}

type WorkflowNode struct {
        ID         string                 `json:"id"`
        Name       string                 `json:"name"`
        Type       string                 `json:"type"`
        TypeVersion float64               `json:"typeVersion"`
        Position   []int                  `json:"position"`
        Parameters map[string]interface{} `json:"parameters"`
}

type WorkflowList struct {
        Data []Workflow `json:"data"`
}

type User struct {
        ID       string `json:"id"`
        Email    string `json:"email"`
        Settings map[string]interface{} `json:"settings"`
}

// Mock data store
var workflows = map[string]*Workflow{
        "1001": {
                ID:        "1001",
                Name:      "Customer Onboarding",
                Active:    true,
                CreatedAt: time.Now().Add(-72 * time.Hour),
                UpdatedAt: time.Now().Add(-30 * time.Minute),
                VersionId: 15,
                Tags:      []string{"customer", "onboarding"},
                Nodes: []WorkflowNode{
                        {
                                ID:         "webhook-1001",
                                Name:       "Webhook",
                                Type:       "n8n-nodes-base.webhook",
                                TypeVersion: 1.0,
                                Position:   []int{250, 300},
                                Parameters: map[string]interface{}{
                                        "httpMethod": "POST",
                                        "path":       "customer-signup",
                                },
                        },
                        {
                                ID:         "email-1001",
                                Name:       "Send Welcome Email",
                                Type:       "n8n-nodes-base.emailSend",
                                TypeVersion: 1.0,
                                Position:   []int{450, 300},
                                Parameters: map[string]interface{}{
                                        "subject": "Welcome to our platform!",
                                        "toEmail": "={{ $json.email }}",
                                },
                        },
                },
                Connections: map[string]interface{}{
                        "Webhook": map[string]interface{}{
                                "main": [][]map[string]interface{}{
                                        {
                                                {
                                                        "node": "Send Welcome Email",
                                                        "type": "main",
                                                        "index": 0,
                                                },
                                        },
                                },
                        },
                },
        },
        "1002": {
                ID:        "1002",
                Name:      "Payment Processing",
                Active:    true,
                CreatedAt: time.Now().Add(-48 * time.Hour),
                UpdatedAt: time.Now().Add(-10 * time.Minute),
                VersionId: 23,
                Tags:      []string{"payment", "stripe"},
                Nodes: []WorkflowNode{
                        {
                                ID:         "webhook-1002",
                                Name:       "Payment Webhook",
                                Type:       "n8n-nodes-base.webhook",
                                TypeVersion: 1.0,
                                Position:   []int{250, 200},
                                Parameters: map[string]interface{}{
                                        "httpMethod": "POST",
                                        "path":       "stripe-webhook",
                                },
                        },
                        {
                                ID:         "stripe-1002",
                                Name:       "Process Payment",
                                Type:       "n8n-nodes-base.stripe",
                                TypeVersion: 1.0,
                                Position:   []int{450, 200},
                                Parameters: map[string]interface{}{
                                        "operation": "charge",
                                },
                        },
                },
                Connections: map[string]interface{}{
                        "Payment Webhook": map[string]interface{}{
                                "main": [][]map[string]interface{}{
                                        {
                                                {
                                                        "node": "Process Payment",
                                                        "type": "main",
                                                        "index": 0,
                                                },
                                        },
                                },
                        },
                },
        },
        "1003": {
                ID:        "1003",
                Name:      "Order Fulfillment",
                Active:    false,
                CreatedAt: time.Now().Add(-24 * time.Hour),
                UpdatedAt: time.Now().Add(-5 * time.Minute),
                VersionId: 8,
                Tags:      []string{"orders", "fulfillment"},
                Nodes: []WorkflowNode{
                        {
                                ID:         "trigger-1003",
                                Name:       "Order Created",
                                Type:       "n8n-nodes-base.httpRequest",
                                TypeVersion: 1.0,
                                Position:   []int{250, 400},
                                Parameters: map[string]interface{}{
                                        "method": "POST",
                                },
                        },
                },
                Connections: map[string]interface{}{},
        },
}

// API Key validation
func validateAPIKey(r *http.Request) bool {
        apiKey := r.Header.Get("X-N8N-API-KEY")
        // Mock validation - accept any key starting with "n8n_api_"
        return apiKey != "" && len(apiKey) > 8 && apiKey[:8] == "n8n_api_"
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                if !validateAPIKey(r) {
                        http.Error(w, "Unauthorized", http.StatusUnauthorized)
                        return
                }
                next.ServeHTTP(w, r)
        }
}

// GET /api/v1/me
func getMe(w http.ResponseWriter, r *http.Request) {
        user := User{
                ID:    "user-123",
                Email: "admin@company.com",
                Settings: map[string]interface{}{
                        "timezone": "UTC",
                },
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
}

// GET /api/v1/workflows
func getWorkflows(w http.ResponseWriter, r *http.Request) {
        var workflowList []Workflow
        for _, workflow := range workflows {
                workflowList = append(workflowList, *workflow)
        }
        
        response := WorkflowList{Data: workflowList}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}

// GET /api/v1/workflows/{id}
func getWorkflow(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        
        workflow, exists := workflows[id]
        if !exists {
                http.Error(w, "Workflow not found", http.StatusNotFound)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(workflow)
}

// PUT /api/v1/workflows/{id}
func updateWorkflow(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        
        var updatedWorkflow Workflow
        if err := json.NewDecoder(r.Body).Decode(&updatedWorkflow); err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        // Update workflow
        updatedWorkflow.ID = id
        updatedWorkflow.UpdatedAt = time.Now()
        if existingWorkflow, exists := workflows[id]; exists {
                updatedWorkflow.CreatedAt = existingWorkflow.CreatedAt
                updatedWorkflow.VersionId = existingWorkflow.VersionId + 1
        } else {
                updatedWorkflow.CreatedAt = time.Now()
                updatedWorkflow.VersionId = 1
        }
        
        workflows[id] = &updatedWorkflow
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updatedWorkflow)
}

// POST /api/v1/workflows
func createWorkflow(w http.ResponseWriter, r *http.Request) {
        var newWorkflow Workflow
        if err := json.NewDecoder(r.Body).Decode(&newWorkflow); err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        // Generate new ID
        newID := fmt.Sprintf("%d", len(workflows)+1001)
        newWorkflow.ID = newID
        newWorkflow.CreatedAt = time.Now()
        newWorkflow.UpdatedAt = time.Now()
        newWorkflow.VersionId = 1
        
        workflows[newID] = &newWorkflow
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(newWorkflow)
}

// POST /api/v1/workflows/{id}/activate
func activateWorkflow(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        
        workflow, exists := workflows[id]
        if !exists {
                http.Error(w, "Workflow not found", http.StatusNotFound)
                return
        }
        
        workflow.Active = true
        workflow.UpdatedAt = time.Now()
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(workflow)
}

// POST /api/v1/workflows/{id}/deactivate
func deactivateWorkflow(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        
        workflow, exists := workflows[id]
        if !exists {
                http.Error(w, "Workflow not found", http.StatusNotFound)
                return
        }
        
        workflow.Active = false
        workflow.UpdatedAt = time.Now()
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(workflow)
}

// DELETE /api/v1/workflows/{id}
func deleteWorkflow(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        
        if _, exists := workflows[id]; !exists {
                http.Error(w, "Workflow not found", http.StatusNotFound)
                return
        }
        
        delete(workflows, id)
        w.WriteHeader(http.StatusNoContent)
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
        response := map[string]interface{}{
                "status": "ok",
                "timestamp": time.Now(),
                "version": "mock-n8n-api-v1.0.0",
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}

func main() {
        r := mux.NewRouter()
        
        // Health check (no auth required)
        r.HandleFunc("/health", healthCheck).Methods("GET")
        
        // API routes (with auth)
        api := r.PathPrefix("/api/v1").Subrouter()
        api.HandleFunc("/me", authMiddleware(getMe)).Methods("GET")
        api.HandleFunc("/workflows", authMiddleware(getWorkflows)).Methods("GET")
        api.HandleFunc("/workflows", authMiddleware(createWorkflow)).Methods("POST")
        api.HandleFunc("/workflows/{id}", authMiddleware(getWorkflow)).Methods("GET")
        api.HandleFunc("/workflows/{id}", authMiddleware(updateWorkflow)).Methods("PUT")
        api.HandleFunc("/workflows/{id}", authMiddleware(deleteWorkflow)).Methods("DELETE")
        api.HandleFunc("/workflows/{id}/activate", authMiddleware(activateWorkflow)).Methods("POST")
        api.HandleFunc("/workflows/{id}/deactivate", authMiddleware(deactivateWorkflow)).Methods("POST")
        api.HandleFunc("/executions", authMiddleware(getExecutions)).Methods("GET")
        api.HandleFunc("/executions/{id}", authMiddleware(getExecution)).Methods("GET")
        
        // CORS middleware for all routes
        r.Use(func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        w.Header().Set("Access-Control-Allow-Origin", "*")
                        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-N8N-API-KEY")
                        
                        if r.Method == "OPTIONS" {
                                w.WriteHeader(http.StatusOK)
                                return
                        }
                        
                        next.ServeHTTP(w, r)
                })
        })
        
        port := "3001"
        fmt.Printf("🚀 Mock n8n API Server starting on port %s\n", port)
        fmt.Printf("📍 Health check: http://localhost:%s/health\n", port)
        fmt.Printf("🔑 API Base URL: http://localhost:%s/api/v1\n", port)
        fmt.Printf("🔐 Use API key: n8n_api_mock_development\n")
        fmt.Println("===============================================")
        
        log.Fatal(http.ListenAndServe(":"+port, r))
}

// Execution represents an n8n workflow execution
type Execution struct {
        ID           string                 `json:"id"`
        WorkflowID   string                 `json:"workflowId"`
        WorkflowName string                 `json:"workflowName"`
        Status       string                 `json:"status"` // running, success, error, waiting
        StartedAt    time.Time              `json:"startedAt"`
        StoppedAt    *time.Time             `json:"stoppedAt,omitempty"`
        Mode         string                 `json:"mode"` // manual, trigger, retry
        Retries      int                    `json:"retries"`
        Data         map[string]interface{} `json:"data,omitempty"`
        Error        *ExecutionError        `json:"error,omitempty"`
}

type ExecutionError struct {
        Message string `json:"message"`
        Node    string `json:"node"`
        Stack   string `json:"stack,omitempty"`
}

type ExecutionsList struct {
        Data []Execution `json:"data"`
}

// Mock executions data
var executions = map[string]*Execution{
        "exec_1001": {
                ID:           "exec_1001",
                WorkflowID:   "1001",
                WorkflowName: "Customer Onboarding",
                Status:       "success",
                StartedAt:    time.Now().Add(-2 * time.Hour),
                StoppedAt:    timePtr(time.Now().Add(-2*time.Hour).Add(5 * time.Minute)),
                Mode:         "trigger",
                Retries:      0,
        },
        "exec_1002": {
                ID:           "exec_1002",
                WorkflowID:   "1001",
                WorkflowName: "Customer Onboarding",
                Status:       "error",
                StartedAt:    time.Now().Add(-1 * time.Hour),
                StoppedAt:    timePtr(time.Now().Add(-1 * time.Hour).Add(2 * time.Minute)),
                Mode:         "trigger",
                Retries:      2,
                Error: &ExecutionError{
                        Message: "Database connection timeout",
                        Node:    "PostgreSQL",
                        Stack:   "Error: Database connection timeout\n    at PostgreSQL.execute()",
                },
        },
        "exec_1003": {
                ID:           "exec_1003",
                WorkflowID:   "1002",
                WorkflowName: "Payment Processing",
                Status:       "error",
                StartedAt:    time.Now().Add(-30 * time.Minute),
                StoppedAt:    timePtr(time.Now().Add(-25 * time.Minute)),
                Mode:         "trigger",
                Retries:      3,
                Error: &ExecutionError{
                        Message: "API rate limit exceeded",
                        Node:    "Stripe",
                        Stack:   "Error: API rate limit exceeded\n    at Stripe.makeRequest()",
                },
        },
}

func timePtr(t time.Time) *time.Time {
        return &t
}

// GET /api/v1/executions
func getExecutions(w http.ResponseWriter, r *http.Request) {
        workflowID := r.URL.Query().Get("workflowId")
        status := r.URL.Query().Get("status")
        limitStr := r.URL.Query().Get("limit")
        
        limit := 50 // default
        if limitStr != "" {
                if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
                        limit = l
                }
        }
        
        var filteredExecutions []Execution
        count := 0
        
        for _, exec := range executions {
                if count >= limit {
                        break
                }
                
                // Filter by workflowID if provided
                if workflowID != "" && exec.WorkflowID != workflowID {
                        continue
                }
                
                // Filter by status if provided
                if status != "" && exec.Status != status {
                        continue
                }
                
                filteredExecutions = append(filteredExecutions, *exec)
                count++
        }
        
        response := ExecutionsList{Data: filteredExecutions}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}

// GET /api/v1/executions/{id}
func getExecution(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        
        execution, exists := executions[id]
        if !exists {
                http.Error(w, "Execution not found", http.StatusNotFound)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(execution)
}