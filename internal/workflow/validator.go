package workflow

import (
        "encoding/json"
        "fmt"
        "os"
        "strings"

        "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
        validate = validator.New()
}

// ValidateWorkflowFile validates a workflow JSON file
func ValidateWorkflowFile(filepath string) error {
        // Check if file exists
        if _, err := os.Stat(filepath); os.IsNotExist(err) {
                return fmt.Errorf("workflow file does not exist: %s", filepath)
        }

        // Read file content
        content, err := os.ReadFile(filepath)
        if err != nil {
                return fmt.Errorf("failed to read workflow file: %w", err)
        }

        // Validate JSON syntax
        var rawWorkflow map[string]interface{}
        if err := json.Unmarshal(content, &rawWorkflow); err != nil {
                return fmt.Errorf("invalid JSON syntax: %w", err)
        }

        // Parse into workflow struct
        var workflow Workflow
        if err := json.Unmarshal(content, &workflow); err != nil {
                return fmt.Errorf("failed to parse workflow: %w", err)
        }

        // Validate workflow structure
        return ValidateWorkflow(&workflow)
}

// ValidateWorkflow validates a workflow struct
func ValidateWorkflow(wf *Workflow) error {
        // Basic validation using tags
        if err := validate.Struct(wf); err != nil {
                return fmt.Errorf("validation failed: %w", err)
        }

        // Custom validation rules
        if err := validateWorkflowName(wf.Name); err != nil {
                return err
        }

        if err := validateNodes(wf.Nodes); err != nil {
                return err
        }

        if err := validateConnections(wf.Connections); err != nil {
                return err
        }

        return nil
}

// ValidateWorkflowStrict performs strict validation with additional checks
func ValidateWorkflowStrict(wf *Workflow) error {
        // Run basic validation first
        if err := ValidateWorkflow(wf); err != nil {
                return err
        }

        // Strict validation rules
        if err := validateNodeTypes(wf.Nodes); err != nil {
                return err
        }

        if err := validateNodeConnectivity(wf.Nodes, wf.Connections); err != nil {
                return err
        }

        if err := validateWorkflowSettings(wf.Settings); err != nil {
                return err
        }

        return nil
}

// validateWorkflowName validates the workflow name
func validateWorkflowName(name string) error {
        if name == "" {
                return fmt.Errorf("workflow name cannot be empty")
        }

        if len(name) > 100 {
                return fmt.Errorf("workflow name too long (max 100 characters)")
        }

        // Check for invalid characters that might cause issues
        invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
        for _, char := range invalidChars {
                if strings.Contains(name, char) {
                        return fmt.Errorf("workflow name contains invalid character: %s", char)
                }
        }

        return nil
}

// validateNodes validates the workflow nodes
func validateNodes(nodes []Node) error {
        if len(nodes) == 0 {
                return fmt.Errorf("workflow must contain at least one node")
        }

        nodeNames := make(map[string]bool)
        for _, node := range nodes {
                // Check for duplicate node names
                if nodeNames[node.Name] {
                        return fmt.Errorf("duplicate node name: %s", node.Name)
                }
                nodeNames[node.Name] = true

                // Validate individual node
                if err := validateNode(node); err != nil {
                        return fmt.Errorf("invalid node '%s': %w", node.Name, err)
                }
        }

        return nil
}

// validateNode validates a single node
func validateNode(node Node) error {
        if node.Name == "" {
                return fmt.Errorf("node name cannot be empty")
        }

        if node.Type == "" {
                return fmt.Errorf("node type cannot be empty")
        }

        // Validate position
        if node.Position == nil || len(node.Position) != 2 {
                return fmt.Errorf("node position must be an array of two numbers")
        }

        return nil
}

// validateConnections validates workflow connections
func validateConnections(connections map[string]interface{}) error {
        // Basic validation - ensure connections is a valid object
        if connections == nil {
                return nil // Connections are optional
        }

        // More detailed connection validation could be added here
        // For now, we'll just ensure it's a valid map structure
        return nil
}

// validateNodeTypes validates that node types are known/supported
func validateNodeTypes(nodes []Node) error {
        // List of commonly known n8n node types
        // In a real implementation, this could be loaded from a configuration file
        knownTypes := map[string]bool{
                "n8n-nodes-base.start":              true,
                "n8n-nodes-base.webhook":            true,
                "n8n-nodes-base.httpRequest":        true,
                "n8n-nodes-base.function":           true,
                "n8n-nodes-base.if":                 true,
                "n8n-nodes-base.switch":             true,
                "n8n-nodes-base.set":                true,
                "n8n-nodes-base.merge":              true,
                "n8n-nodes-base.wait":               true,
                "n8n-nodes-base.spreadsheetFile":    true,
                "n8n-nodes-base.emailSend":          true,
                "n8n-nodes-base.slack":              true,
                "n8n-nodes-base.mysql":              true,
                "n8n-nodes-base.postgres":           true,
                "n8n-nodes-base.mongoDb":            true,
                "n8n-nodes-base.redis":              true,
                "n8n-nodes-base.cron":               true,
                "n8n-nodes-base.interval":           true,
                // Add more as needed
        }

        for _, node := range nodes {
                // Skip validation for custom/community nodes
                if strings.HasPrefix(node.Type, "@") {
                        continue
                }

                if !knownTypes[node.Type] {
                        // This is a warning, not an error in strict mode
                        // You might want to make this configurable
                        fmt.Printf("Warning: Unknown node type '%s' for node '%s'\n", node.Type, node.Name)
                }
        }

        return nil
}

// validateNodeConnectivity ensures nodes are properly connected
func validateNodeConnectivity(nodes []Node, connections map[string]interface{}) error {
        nodeMap := make(map[string]Node)
        for _, node := range nodes {
                nodeMap[node.Name] = node
        }

        // Check if all nodes (except start nodes) have incoming connections
        startNodeTypes := map[string]bool{
                "n8n-nodes-base.start":    true,
                "n8n-nodes-base.webhook":  true,
                "n8n-nodes-base.cron":     true,
                "n8n-nodes-base.interval": true,
        }

        hasIncomingConnection := make(map[string]bool)

        // Parse connections to build incoming connection map
        for _, sourceConnections := range connections {
                if sourceConns, ok := sourceConnections.(map[string]interface{}); ok {
                        for _, targets := range sourceConns {
                                if targetList, ok := targets.([]interface{}); ok {
                                        for _, target := range targetList {
                                                if targetConn, ok := target.(map[string]interface{}); ok {
                                                        if targetNode, ok := targetConn["node"].(string); ok {
                                                                hasIncomingConnection[targetNode] = true
                                                        }
                                                }
                                        }
                                }
                        }
                }
        }

        // Check for isolated nodes
        for _, node := range nodes {
                if !startNodeTypes[node.Type] && !hasIncomingConnection[node.Name] {
                        return fmt.Errorf("node '%s' has no incoming connections and is not a start node", node.Name)
                }
        }

        return nil
}

// validateWorkflowSettings validates workflow settings
func validateWorkflowSettings(settings map[string]interface{}) error {
        if settings == nil {
                return nil // Settings are optional
        }

        // Validate timezone if present
        if timezone, ok := settings["timezone"].(string); ok {
                if timezone == "" {
                        return fmt.Errorf("timezone cannot be empty if specified")
                }
                // You could add more timezone validation here
        }

        // Validate saveDataErrorExecution
        if saveDataError, ok := settings["saveDataErrorExecution"]; ok {
                if _, ok := saveDataError.(string); !ok {
                        return fmt.Errorf("saveDataErrorExecution must be a string")
                }
        }

        // Validate saveDataSuccessExecution  
        if saveDataSuccess, ok := settings["saveDataSuccessExecution"]; ok {
                if _, ok := saveDataSuccess.(string); !ok {
                        return fmt.Errorf("saveDataSuccessExecution must be a string")
                }
        }

        return nil
}

// ValidateWorkflowBatch validates multiple workflows
func ValidateWorkflowBatch(workflows []*Workflow) []error {
        var errors []error
        
        for i, wf := range workflows {
                if err := ValidateWorkflow(wf); err != nil {
                        errors = append(errors, fmt.Errorf("workflow %d (%s): %w", i, wf.Name, err))
                }
        }
        
        return errors
}

// ValidateWorkflowJSON validates raw JSON content
func ValidateWorkflowJSON(jsonContent []byte) error {
        var wf Workflow
        if err := json.Unmarshal(jsonContent, &wf); err != nil {
                return fmt.Errorf("invalid workflow JSON: %w", err)
        }
        
        return ValidateWorkflow(&wf)
}
