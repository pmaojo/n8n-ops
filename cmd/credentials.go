package cmd

import (
        "fmt"
        "os"
        "strings"

        "github.com/spf13/cobra"
        "github.com/n8n-workflows/n8n-ops/internal/credentials"
)

var credentialsCmd = &cobra.Command{
        Use:   "credentials",
        Short: "Manage workflow credentials securely across environments",
        Long: `Secure credential management for n8n workflows using environment variables only.

SECURITY PRINCIPLES:
• No credentials stored in files or Git
• Environment-specific credential mapping
• Automatic credential substitution in workflows
• Secure variable naming conventions

WORKFLOW CREDENTIAL MAPPING:
• Development: Uses _DEV suffix variables
• Staging: Uses _STAGING suffix variables  
• Production: Uses _PROD suffix variables

Examples:
  n8n-ops credentials list                    # Show credential mappings
  n8n-ops credentials validate               # Validate all required credentials
  n8n-ops credentials map --workflow-id wf123 # Show workflow credential mapping`,
}

var credentialsListCmd = &cobra.Command{
        Use:   "list",
        Short: "List credential mappings for current environment",
        RunE:  runCredentialsList,
}

var credentialsValidateCmd = &cobra.Command{
        Use:   "validate",
        Short: "Validate all required credentials are set",
        RunE:  runCredentialsValidate,
}

var credentialsMapCmd = &cobra.Command{
        Use:   "map",
        Short: "Show credential mapping for workflows",
        RunE:  runCredentialsMap,
}

var credentialsTemplateCmd = &cobra.Command{
        Use:   "template",
        Short: "Generate environment variable template",
        RunE:  runCredentialsTemplate,
}

var (
        workflowID string
        showValues bool
)

func init() {
        rootCmd.AddCommand(credentialsCmd)
        credentialsCmd.AddCommand(credentialsListCmd)
        credentialsCmd.AddCommand(credentialsValidateCmd)
        credentialsCmd.AddCommand(credentialsMapCmd)
        credentialsCmd.AddCommand(credentialsTemplateCmd)

        credentialsMapCmd.Flags().StringVar(&workflowID, "workflow-id", "", "specific workflow ID to map")
        credentialsListCmd.Flags().BoolVar(&showValues, "show-values", false, "show credential values (use carefully)")
}

func runCredentialsList(cmd *cobra.Command, args []string) error {
        fmt.Printf("🔐 Credential Mappings - %s Environment\n", environment)
        fmt.Printf("=====================================\n\n")

        manager := credentials.NewSecureCredentialManager(environment)
        
        mappings, err := manager.GetCredentialMappings()
        if err != nil {
                return fmt.Errorf("failed to get credential mappings: %w", err)
        }

        fmt.Printf("📋 N8N API Credentials:\n")
        fmt.Printf("  N8N_URL_%s = %s\n", strings.ToUpper(environment), getEnvOrDefault(fmt.Sprintf("N8N_URL_%s", strings.ToUpper(environment)), "not_set"))
        fmt.Printf("  N8N_API_KEY_%s = %s\n", strings.ToUpper(environment), maskValue(os.Getenv(fmt.Sprintf("N8N_API_KEY_%s", strings.ToUpper(environment)))))
        fmt.Printf("\n")

        fmt.Printf("🔧 Workflow Node Credentials:\n")
        for nodeType, creds := range mappings {
                fmt.Printf("  %s:\n", nodeType)
                for credName, envVar := range creds {
                        value := os.Getenv(envVar)
                        if showValues {
                                fmt.Printf("    %s = %s (%s)\n", credName, value, envVar)
                        } else {
                                fmt.Printf("    %s = %s (%s)\n", credName, maskValue(value), envVar)
                        }
                }
                fmt.Printf("\n")
        }

        if !showValues {
                fmt.Printf("💡 Use --show-values to see actual values (be careful!)\n")
        }

        return nil
}

func runCredentialsValidate(cmd *cobra.Command, args []string) error {
        fmt.Printf("✅ Validating Credentials - %s Environment\n", environment)
        fmt.Printf("=========================================\n\n")

        manager := credentials.NewSecureCredentialManager(environment)
        
        validation, err := manager.ValidateAllCredentials()
        if err != nil {
                return fmt.Errorf("validation failed: %w", err)
        }

        fmt.Printf("📊 Validation Results:\n")
        fmt.Printf("  Required: %d\n", validation.Required)
        fmt.Printf("  Present:  %d\n", validation.Present)
        fmt.Printf("  Missing:  %d\n", validation.Missing)
        fmt.Printf("\n")

        if len(validation.MissingCredentials) > 0 {
                fmt.Printf("❌ Missing Credentials:\n")
                for _, missing := range validation.MissingCredentials {
                        fmt.Printf("  • %s\n", missing)
                }
                fmt.Printf("\n")
                fmt.Printf("💡 Set missing credentials:\n")
                for _, missing := range validation.MissingCredentials {
                        fmt.Printf("  export %s=\"your_value_here\"\n", missing)
                }
        } else {
                fmt.Printf("✅ All required credentials are present!\n")
        }

        return nil
}

func runCredentialsMap(cmd *cobra.Command, args []string) error {
        fmt.Printf("🗺️ Workflow Credential Mapping\n")
        fmt.Printf("==============================\n\n")

        manager := credentials.NewSecureCredentialManager(environment)

        if workflowID != "" {
                // Show specific workflow mapping
                mapping, err := manager.GetWorkflowCredentialMapping(workflowID)
                if err != nil {
                        return fmt.Errorf("failed to get workflow mapping: %w", err)
                }

                fmt.Printf("📋 Workflow: %s\n", workflowID)
                fmt.Printf("Environment: %s\n\n", environment)

                for nodeID, nodeCreds := range mapping {
                        fmt.Printf("🔧 Node: %s\n", nodeID)
                        for credName, envVar := range nodeCreds {
                                status := "✅"
                                if os.Getenv(envVar) == "" {
                                        status = "❌"
                                }
                                fmt.Printf("  %s %s → %s\n", status, credName, envVar)
                        }
                        fmt.Printf("\n")
                }
        } else {
                // Show general mapping structure
                mappings, err := manager.GetCredentialMappings()
                if err != nil {
                        return fmt.Errorf("failed to get mappings: %w", err)
                }

                fmt.Printf("📋 Environment Variable Patterns:\n\n")
                
                for nodeType, creds := range mappings {
                        fmt.Printf("🔧 %s Node:\n", nodeType)
                        for credName, pattern := range creds {
                                fmt.Printf("  %s → %s\n", credName, pattern)
                        }
                        fmt.Printf("\n")
                }
        }

        return nil
}

func runCredentialsTemplate(cmd *cobra.Command, args []string) error {
        fmt.Printf("📝 Environment Variable Template\n")
        fmt.Printf("==============================\n\n")

        manager := credentials.NewSecureCredentialManager(environment)
        
        template, err := manager.GenerateEnvTemplate()
        if err != nil {
                return fmt.Errorf("failed to generate template: %w", err)
        }

        fmt.Printf("# n8n-ops Environment Variables - %s\n", strings.ToUpper(environment))
        fmt.Printf("# Copy these variables to your environment\n\n")

        fmt.Printf("# N8N API Configuration\n")
        fmt.Printf("export N8N_URL_%s=\"%s\"\n", strings.ToUpper(environment), template.N8nURL)
        fmt.Printf("export N8N_API_KEY_%s=\"%s\"\n\n", strings.ToUpper(environment), template.N8nAPIKey)

        fmt.Printf("# Workflow Node Credentials\n")
        for category, vars := range template.NodeCredentials {
                fmt.Printf("# %s\n", category)
                for varName, example := range vars {
                        fmt.Printf("export %s=\"%s\"\n", varName, example)
                }
                fmt.Printf("\n")
        }

        fmt.Printf("# GitLab Integration (optional)\n")
        fmt.Printf("export GITLAB_TOKEN_%s=\"your_gitlab_token\"\n", strings.ToUpper(environment))
        fmt.Printf("export GITLAB_PROJECT_ID_%s=\"your_project_id\"\n\n", strings.ToUpper(environment))

        fmt.Printf("💡 Tip: Save this to .env.%s (never commit to Git!)\n", environment)

        return nil
}

func maskValue(value string) string {
        if value == "" {
                return "not_set"
        }
        if len(value) <= 4 {
                return "***"
        }
        return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

func getEnvOrDefault(key, defaultValue string) string {
        if value := os.Getenv(key); value != "" {
                return value
        }
        return defaultValue
}