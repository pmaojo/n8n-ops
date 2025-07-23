package credentials

import (
        "encoding/json"
        "fmt"
        "os"
        "strings"
)

// N8nCredentialManager handles n8n-specific credential operations
type N8nCredentialManager struct {
        Environment string
        BaseURL     string
        APIKey      string
}

// N8nCredential represents a credential in n8n
type N8nCredential struct {
        ID   string `json:"id"`
        Name string `json:"name"`
        Type string `json:"type"`
        Data struct {
                // This varies by credential type - n8n encrypts this data
                // We only work with credential references, not actual values
        } `json:"data"`
        NodesAccess []struct {
                NodeType string `json:"nodeType"`
        } `json:"nodesAccess"`
}

// WorkflowCredentialReference represents how credentials are referenced in workflows
type WorkflowCredentialReference struct {
        WorkflowID   string `json:"workflowId"`
        NodeID       string `json:"nodeId"`
        NodeType     string `json:"nodeType"`
        CredentialID string `json:"credentialId"`
        Environment  string `json:"environment"`
}

// CredentialMapping maps environment-specific credential names
type CredentialMapping struct {
        CredentialName string `json:"credentialName"`
        DevID          string `json:"devId"`
        StagingID      string `json:"stagingId"`
        ProductionID   string `json:"productionId"`
}

func NewN8nCredentialManager(environment string) *N8nCredentialManager {
        return &N8nCredentialManager{
                Environment: environment,
                BaseURL:     os.Getenv(fmt.Sprintf("N8N_URL_%s", strings.ToUpper(environment))),
                APIKey:      os.Getenv(fmt.Sprintf("N8N_API_KEY_%s", strings.ToUpper(environment))),
        }
}

// GetWorkflowCredentials analyzes a workflow and returns credential references
func (ncm *N8nCredentialManager) GetWorkflowCredentials(workflowJSON []byte) ([]WorkflowCredentialReference, error) {
        var workflow struct {
                ID    string `json:"id"`
                Name  string `json:"name"`
                Nodes []struct {
                        ID         string `json:"id"`
                        Type       string `json:"type"`
                        Credentials map[string]struct {
                                ID   string `json:"id"`
                                Name string `json:"name"`
                        } `json:"credentials"`
                } `json:"nodes"`
        }

        if err := json.Unmarshal(workflowJSON, &workflow); err != nil {
                return nil, fmt.Errorf("failed to parse workflow JSON: %w", err)
        }

        var credRefs []WorkflowCredentialReference

        for _, node := range workflow.Nodes {
                for _, credInfo := range node.Credentials {
                        credRefs = append(credRefs, WorkflowCredentialReference{
                                WorkflowID:   workflow.ID,
                                NodeID:       node.ID,
                                NodeType:     node.Type,
                                CredentialID: credInfo.ID,
                                Environment:  ncm.Environment,
                        })
                }
        }

        return credRefs, nil
}

// MapCredentialsForEnvironment updates workflow credential IDs for target environment
func (ncm *N8nCredentialManager) MapCredentialsForEnvironment(workflowJSON []byte, targetEnv string) ([]byte, error) {
        // Parse workflow
        var workflow map[string]interface{}
        if err := json.Unmarshal(workflowJSON, &workflow); err != nil {
                return nil, fmt.Errorf("failed to parse workflow: %w", err)
        }

        // Get credential mappings
        mappings := ncm.getCredentialMappings()

        // Update credential references in nodes
        if nodes, ok := workflow["nodes"].([]interface{}); ok {
                for _, nodeInterface := range nodes {
                        if node, ok := nodeInterface.(map[string]interface{}); ok {
                                if credentials, ok := node["credentials"].(map[string]interface{}); ok {
                                        for _, credInterface := range credentials {
                                                if cred, ok := credInterface.(map[string]interface{}); ok {
                                                        if credID, ok := cred["id"].(string); ok {
                                                                // Map credential ID to target environment
                                                                if newCredID := ncm.mapCredentialID(credID, targetEnv, mappings); newCredID != "" {
                                                                        cred["id"] = newCredID
                                                                }
                                                        }
                                                }
                                        }
                                }
                        }
                }
        }

        return json.MarshalIndent(workflow, "", "  ")
}

// getCredentialMappings returns environment-specific credential mappings
func (ncm *N8nCredentialManager) getCredentialMappings() map[string]CredentialMapping {
        // These mappings would typically be stored in a secure configuration
        // For now, we use a conventional naming pattern
        return map[string]CredentialMapping{
                "smtp": {
                        CredentialName: "SMTP Email",
                        DevID:          "smtp_dev_001",
                        StagingID:      "smtp_staging_001", 
                        ProductionID:   "smtp_prod_001",
                },
                "postgres": {
                        CredentialName: "PostgreSQL Database",
                        DevID:          "postgres_dev_001",
                        StagingID:      "postgres_staging_001",
                        ProductionID:   "postgres_prod_001",
                },
                "stripe": {
                        CredentialName: "Stripe Payment",
                        DevID:          "stripe_dev_001",
                        StagingID:      "stripe_staging_001", 
                        ProductionID:   "stripe_prod_001",
                },
                "aws": {
                        CredentialName: "AWS S3",
                        DevID:          "aws_dev_001",
                        StagingID:      "aws_staging_001",
                        ProductionID:   "aws_prod_001",
                },
        }
}

// mapCredentialID maps a credential ID to the target environment
func (ncm *N8nCredentialManager) mapCredentialID(credID, targetEnv string, mappings map[string]CredentialMapping) string {
        // Find the credential mapping
        for _, mapping := range mappings {
                switch targetEnv {
                case "development":
                        if credID == mapping.StagingID || credID == mapping.ProductionID {
                                return mapping.DevID
                        }
                case "staging":
                        if credID == mapping.DevID || credID == mapping.ProductionID {
                                return mapping.StagingID
                        }
                case "production":
                        if credID == mapping.DevID || credID == mapping.StagingID {
                                return mapping.ProductionID
                        }
                }
        }

        // Return original ID if no mapping found
        return credID
}

// GenerateCredentialReport creates a report of credentials used across environments
func (ncm *N8nCredentialManager) GenerateCredentialReport() (*CredentialReport, error) {
        report := &CredentialReport{
                Environment: ncm.Environment,
                Credentials: make(map[string]CredentialInfo),
        }

        mappings := ncm.getCredentialMappings()

        for _, mapping := range mappings {
                var currentID string
                switch ncm.Environment {
                case "development":
                        currentID = mapping.DevID
                case "staging":
                        currentID = mapping.StagingID
                case "production":
                        currentID = mapping.ProductionID
                }

                report.Credentials[mapping.CredentialName] = CredentialInfo{
                        Name:        mapping.CredentialName,
                        Type:        mapping.CredentialName,
                        ID:          currentID,
                        Environment: ncm.Environment,
                        Status:      ncm.checkCredentialStatus(currentID),
                }
        }

        return report, nil
}

// checkCredentialStatus verifies if a credential exists in n8n
func (ncm *N8nCredentialManager) checkCredentialStatus(credID string) string {
        // In a real implementation, this would make an API call to n8n
        // For now, we simulate the check
        if credID != "" {
                return "active"
        }
        return "missing"
}

// SecureCredentialManager handles secure credential operations
type SecureCredentialManager struct {
        Environment string
}

func NewSecureCredentialManager(environment string) *SecureCredentialManager {
        return &SecureCredentialManager{
                Environment: environment,
        }
}

// GetCredentialMappings returns the mapping of node types to environment variables
func (scm *SecureCredentialManager) GetCredentialMappings() (map[string]map[string]string, error) {
        envSuffix := strings.ToUpper(scm.Environment)
        
        return map[string]map[string]string{
                "SMTP": {
                        "host":     fmt.Sprintf("SMTP_HOST_%s", envSuffix),
                        "port":     fmt.Sprintf("SMTP_PORT_%s", envSuffix),
                        "user":     fmt.Sprintf("SMTP_USER_%s", envSuffix),
                        "password": fmt.Sprintf("SMTP_PASSWORD_%s", envSuffix),
                },
                "PostgreSQL": {
                        "host":     fmt.Sprintf("POSTGRES_HOST_%s", envSuffix),
                        "port":     fmt.Sprintf("POSTGRES_PORT_%s", envSuffix),
                        "database": fmt.Sprintf("POSTGRES_DB_%s", envSuffix),
                        "username": fmt.Sprintf("POSTGRES_USER_%s", envSuffix),
                        "password": fmt.Sprintf("POSTGRES_PASSWORD_%s", envSuffix),
                },
                "Stripe": {
                        "secretKey":      fmt.Sprintf("STRIPE_SECRET_KEY_%s", envSuffix),
                        "publishableKey": fmt.Sprintf("STRIPE_PUBLISHABLE_KEY_%s", envSuffix),
                        "webhookSecret":  fmt.Sprintf("STRIPE_WEBHOOK_SECRET_%s", envSuffix),
                },
                "AWS": {
                        "accessKeyId":     fmt.Sprintf("AWS_ACCESS_KEY_ID_%s", envSuffix),
                        "secretAccessKey": fmt.Sprintf("AWS_SECRET_ACCESS_KEY_%s", envSuffix),
                        "region":          fmt.Sprintf("AWS_REGION_%s", envSuffix),
                },
                "Discord": {
                        "webhookUrl": fmt.Sprintf("DISCORD_WEBHOOK_URL_%s", envSuffix),
                        "botToken":   fmt.Sprintf("DISCORD_BOT_TOKEN_%s", envSuffix),
                },
                "Slack": {
                        "webhookUrl": fmt.Sprintf("SLACK_WEBHOOK_URL_%s", envSuffix),
                        "botToken":   fmt.Sprintf("SLACK_BOT_TOKEN_%s", envSuffix),
                },
        }, nil
}

// CredentialValidation represents validation results
type CredentialValidation struct {
        Required           int      `json:"required"`
        Present            int      `json:"present"`
        Missing            int      `json:"missing"`
        MissingCredentials []string `json:"missingCredentials"`
}

// ValidateAllCredentials checks if all required credentials are set
func (scm *SecureCredentialManager) ValidateAllCredentials() (*CredentialValidation, error) {
        mappings, err := scm.GetCredentialMappings()
        if err != nil {
                return nil, err
        }

        validation := &CredentialValidation{
                MissingCredentials: make([]string, 0),
        }

        for _, creds := range mappings {
                for _, envVar := range creds {
                        validation.Required++
                        if os.Getenv(envVar) != "" {
                                validation.Present++
                        } else {
                                validation.Missing++
                                validation.MissingCredentials = append(validation.MissingCredentials, envVar)
                        }
                }
        }

        return validation, nil
}

// GetWorkflowCredentialMapping returns credential mapping for a specific workflow
func (scm *SecureCredentialManager) GetWorkflowCredentialMapping(workflowID string) (map[string]map[string]string, error) {
        // This would analyze a specific workflow and return its credential requirements
        // For now, return a sample mapping
        return map[string]map[string]string{
                "smtp_node_1": {
                        "host":     fmt.Sprintf("SMTP_HOST_%s", strings.ToUpper(scm.Environment)),
                        "user":     fmt.Sprintf("SMTP_USER_%s", strings.ToUpper(scm.Environment)),
                        "password": fmt.Sprintf("SMTP_PASSWORD_%s", strings.ToUpper(scm.Environment)),
                },
                "postgres_node_1": {
                        "host":     fmt.Sprintf("POSTGRES_HOST_%s", strings.ToUpper(scm.Environment)),
                        "database": fmt.Sprintf("POSTGRES_DB_%s", strings.ToUpper(scm.Environment)),
                        "username": fmt.Sprintf("POSTGRES_USER_%s", strings.ToUpper(scm.Environment)),
                },
        }, nil
}

// EnvTemplate represents environment variable template
type EnvTemplate struct {
        N8nURL          string                       `json:"n8nUrl"`
        N8nAPIKey       string                       `json:"n8nApiKey"`
        NodeCredentials map[string]map[string]string `json:"nodeCredentials"`
}

// GenerateEnvTemplate creates a template for environment variables
func (scm *SecureCredentialManager) GenerateEnvTemplate() (*EnvTemplate, error) {
        envSuffix := strings.ToUpper(scm.Environment)
        
        var n8nURL, n8nAPIKey string
        switch scm.Environment {
        case "development":
                n8nURL = "http://localhost:5678"
                n8nAPIKey = "your_development_api_key_here"
        case "staging":
                n8nURL = "https://n8n-staging.company.com"
                n8nAPIKey = "your_staging_api_key_here"
        case "production":
                n8nURL = "https://n8n.company.com"
                n8nAPIKey = "your_production_api_key_here"
        }

        template := &EnvTemplate{
                N8nURL:    n8nURL,
                N8nAPIKey: n8nAPIKey,
                NodeCredentials: map[string]map[string]string{
                        "SMTP": {
                                fmt.Sprintf("SMTP_HOST_%s", envSuffix):     "smtp.example.com",
                                fmt.Sprintf("SMTP_USER_%s", envSuffix):     "user@example.com",
                                fmt.Sprintf("SMTP_PASSWORD_%s", envSuffix): "your_smtp_password",
                        },
                        "Database": {
                                fmt.Sprintf("POSTGRES_HOST_%s", envSuffix):     "db.example.com",
                                fmt.Sprintf("POSTGRES_DB_%s", envSuffix):       "myapp_" + scm.Environment,
                                fmt.Sprintf("POSTGRES_USER_%s", envSuffix):     "db_user",
                                fmt.Sprintf("POSTGRES_PASSWORD_%s", envSuffix): "your_db_password",
                        },
                        "Stripe": {
                                fmt.Sprintf("STRIPE_SECRET_KEY_%s", envSuffix):      "sk_test_...",
                                fmt.Sprintf("STRIPE_PUBLISHABLE_KEY_%s", envSuffix): "pk_test_...",
                        },
                },
        }

        return template, nil
}

// CredentialReport represents a credential status report
type CredentialReport struct {
        Environment string                    `json:"environment"`
        Credentials map[string]CredentialInfo `json:"credentials"`
}

// CredentialInfo represents information about a single credential
type CredentialInfo struct {
        Name        string `json:"name"`
        Type        string `json:"type"`
        ID          string `json:"id"`
        Environment string `json:"environment"`
        Status      string `json:"status"`
}