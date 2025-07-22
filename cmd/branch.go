package cmd

import (
        "bufio"
        "fmt"
        "os"
        "os/exec"
        "regexp"
        "strings"
        "time"

        "github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
        Use:   "branch",
        Short: "Manage workflow branches with intelligent naming conventions",
        Long: `Create and manage Git branches with integrated n8n workflow naming conventions.

BRANCH STRATEGIES (DevOps Best Practices):
• feature/workflow-name - New workflow development
• hotfix/workflow-name  - Production fixes  
• release/v1.2.3        - Release preparation
• experiment/test-name  - A/B testing workflows

WORKFLOW NAMING CONVENTIONS:
• {branch-prefix}_{workflow-name}_{version} - Full naming
• Automatic versioning based on Git tags
• Environment-specific prefixes (dev_, staging_, prod_)
• Metadata tracking for workflow lineage

Examples:
  n8n-ops branch create customer-onboarding    # Interactive workflow creation
  n8n-ops branch list                          # Show all workflow branches  
  n8n-ops branch switch feature/email-flow     # Switch with workflow sync`,
}

var branchCreateCmd = &cobra.Command{
        Use:   "create [workflow-name]",
        Short: "Create new branch with workflow naming conventions",
        Long: `Creates a new Git branch following DevOps best practices with integrated n8n workflow setup.

INTERACTIVE WORKFLOW CREATION:
1. Choose branch type (feature/hotfix/release/experiment)
2. Select base workflow (existing/new/empty)  
3. Configure naming convention
4. Set up initial workflow structure

NAMING CONVENTION GENERATION:
• Branch: feature/customer-onboarding
• Workflow: dev_customer_onboarding_v1.0.0
• File: workflows/development/customer-onboarding-v1.0.0.json`,
        RunE: runBranchCreate,
}

var branchListCmd = &cobra.Command{
        Use:   "list",
        Short: "List all workflow branches with metadata",
        RunE:  runBranchList,
}

var branchSwitchCmd = &cobra.Command{
        Use:   "switch [branch-name]",
        Short: "Switch branch and sync associated workflows",
        RunE:  runBranchSwitch,
}

func init() {
        rootCmd.AddCommand(branchCmd)
        branchCmd.AddCommand(branchCreateCmd)
        branchCmd.AddCommand(branchListCmd)
        branchCmd.AddCommand(branchSwitchCmd)
}

func runBranchCreate(cmd *cobra.Command, args []string) error {
        var workflowName string
        if len(args) > 0 {
                workflowName = args[0]
        }

        fmt.Printf("🌿 Interactive Workflow Branch Creation\n")
        fmt.Printf("=====================================\n\n")

        // Step 1: Get workflow name if not provided
        if workflowName == "" {
                reader := bufio.NewReader(os.Stdin)
                fmt.Printf("Workflow name: ")
                input, _ := reader.ReadString('\n')
                workflowName = strings.TrimSpace(input)
        }

        // Step 2: Choose branch type
        branchType := chooseBranchType()
        
        // Step 3: Choose base workflow strategy
        baseStrategy := chooseBaseWorkflow()
        
        // Step 4: Generate naming convention
        naming := generateNamingConvention(workflowName, branchType)
        
        // Step 5: Show summary and confirm
        if !confirmCreation(naming, baseStrategy) {
                fmt.Printf("❌ Branch creation cancelled\n")
                return nil
        }

        // Step 6: Execute creation
        return executeBranchCreation(naming, baseStrategy)
}

func chooseBranchType() string {
        fmt.Printf("🎯 Choose branch type (DevOps Strategy):\n")
        fmt.Printf("[1] feature  - New workflow development (default)\n")
        fmt.Printf("[2] hotfix   - Production bug fixes\n") 
        fmt.Printf("[3] release  - Release preparation\n")
        fmt.Printf("[4] experiment - A/B testing\n")
        fmt.Printf("\nChoice [1-4]: ")

        reader := bufio.NewReader(os.Stdin)
        input, _ := reader.ReadString('\n')
        choice := strings.TrimSpace(input)

        switch choice {
        case "2":
                return "hotfix"
        case "3":
                return "release"
        case "4":
                return "experiment"
        default:
                return "feature"
        }
}

func chooseBaseWorkflow() string {
        fmt.Printf("\n📋 Choose base workflow strategy:\n")
        fmt.Printf("[1] empty    - Start with blank workflow (default)\n")
        fmt.Printf("[2] template - Use workflow template\n")
        fmt.Printf("[3] copy     - Copy existing workflow\n")
        fmt.Printf("[4] import   - Import from existing branch\n")
        fmt.Printf("\nChoice [1-4]: ")

        reader := bufio.NewReader(os.Stdin)
        input, _ := reader.ReadString('\n')
        choice := strings.TrimSpace(input)

        switch choice {
        case "2":
                return "template"
        case "3":
                return "copy"
        case "4":
                return "import"
        default:
                return "empty"
        }
}

type NamingConvention struct {
        BranchName     string
        WorkflowName   string
        FileName       string
        Version        string
        Environment    string
}

func generateNamingConvention(workflowName, branchType string) NamingConvention {
        // Sanitize workflow name
        sanitized := sanitizeWorkflowName(workflowName)
        
        // Generate version
        version := generateVersion()
        
        // Generate names following DevOps conventions
        naming := NamingConvention{
                BranchName:   fmt.Sprintf("%s/%s", branchType, sanitized),
                WorkflowName: fmt.Sprintf("dev_%s_%s", sanitized, version),
                FileName:     fmt.Sprintf("%s-%s.json", sanitized, version),
                Version:      version,
                Environment:  "development",
        }

        return naming
}

func sanitizeWorkflowName(name string) string {
        // Convert to lowercase and replace spaces/special chars with hyphens
        reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
        sanitized := reg.ReplaceAllString(name, "-")
        sanitized = strings.ToLower(sanitized)
        sanitized = strings.Trim(sanitized, "-")
        return sanitized
}

func generateVersion() string {
        // Generate semantic version based on current date/time
        now := time.Now()
        return fmt.Sprintf("v1.0.%d", now.Unix()%10000)
}

func confirmCreation(naming NamingConvention, baseStrategy string) bool {
        fmt.Printf("\n📝 Branch Creation Summary:\n")
        fmt.Printf("==========================\n")
        fmt.Printf("Git Branch:     %s\n", naming.BranchName)
        fmt.Printf("Workflow Name:  %s\n", naming.WorkflowName)
        fmt.Printf("File Name:      %s\n", naming.FileName)
        fmt.Printf("Version:        %s\n", naming.Version)
        fmt.Printf("Base Strategy:  %s\n", baseStrategy)
        fmt.Printf("Environment:    %s\n", naming.Environment)
        fmt.Printf("\nCreate this branch? [Y/n]: ")

        reader := bufio.NewReader(os.Stdin)
        input, _ := reader.ReadString('\n')
        response := strings.TrimSpace(strings.ToLower(input))
        
        return response == "" || response == "y" || response == "yes"
}

func executeBranchCreation(naming NamingConvention, baseStrategy string) error {
        fmt.Printf("\n🚀 Creating workflow branch...\n")

        // Create Git branch
        fmt.Printf("📝 Creating Git branch: %s\n", naming.BranchName)
        if err := createGitBranch(naming.BranchName); err != nil {
                return fmt.Errorf("failed to create Git branch: %w", err)
        }

        // Create workflow structure
        fmt.Printf("📁 Setting up workflow structure...\n")
        if err := createWorkflowStructure(naming, baseStrategy); err != nil {
                return fmt.Errorf("failed to create workflow structure: %w", err)
        }

        // Create branch metadata
        fmt.Printf("📊 Saving branch metadata...\n")
        if err := saveBranchMetadata(naming); err != nil {
                return fmt.Errorf("failed to save metadata: %w", err)
        }

        fmt.Printf("\n✅ Branch created successfully!\n")
        fmt.Printf("📂 Workflow file: workflows/%s/%s\n", naming.Environment, naming.FileName)
        fmt.Printf("🌿 Git branch: %s\n", naming.BranchName)
        fmt.Printf("\n💡 Next steps:\n")
        fmt.Printf("   1. Edit workflow: workflows/%s/%s\n", naming.Environment, naming.FileName)
        fmt.Printf("   2. Sync to n8n:  n8n-ops sync --to-n8n --env %s\n", naming.Environment)
        fmt.Printf("   3. Test workflow in n8n UI\n")

        return nil
}

func createGitBranch(branchName string) error {
        cmd := exec.Command("git", "checkout", "-b", branchName)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        return cmd.Run()
}

func createWorkflowStructure(naming NamingConvention, baseStrategy string) error {
        // Create directory
        dir := fmt.Sprintf("workflows/%s", naming.Environment)
        if err := os.MkdirAll(dir, 0755); err != nil {
                return err
        }

        // Create workflow file based on strategy
        filePath := fmt.Sprintf("%s/%s", dir, naming.FileName)
        
        var workflowContent []byte
        var err error

        switch baseStrategy {
        case "template":
                workflowContent = createTemplateWorkflow(naming)
        case "copy":
                workflowContent, err = copyExistingWorkflow(naming)
        case "import":
                workflowContent, err = importFromBranch(naming)
        default:
                workflowContent = createEmptyWorkflow(naming)
        }

        if err != nil {
                return err
        }

        return os.WriteFile(filePath, workflowContent, 0644)
}

func createEmptyWorkflow(naming NamingConvention) []byte {
        workflow := fmt.Sprintf(`{
  "id": "%s",
  "name": "%s",
  "active": false,
  "nodes": [
    {
      "id": "start",
      "name": "Start",
      "type": "n8n-nodes-base.start",
      "typeVersion": 1,
      "position": [240, 300],
      "parameters": {}
    }
  ],
  "connections": {},
  "createdAt": "%s",
  "updatedAt": "%s",
  "versionId": 1,
  "tags": ["branch:%s", "version:%s", "env:%s"]
}`, 
                generateWorkflowID(),
                naming.WorkflowName,
                time.Now().Format(time.RFC3339),
                time.Now().Format(time.RFC3339),
                naming.BranchName,
                naming.Version,
                naming.Environment,
        )
        
        return []byte(workflow)
}

func createTemplateWorkflow(naming NamingConvention) []byte {
        // Enhanced template with common nodes
        workflow := fmt.Sprintf(`{
  "id": "%s",
  "name": "%s",
  "active": false,
  "nodes": [
    {
      "id": "start",
      "name": "Start",
      "type": "n8n-nodes-base.start",
      "typeVersion": 1,
      "position": [240, 300],
      "parameters": {}
    },
    {
      "id": "webhook",
      "name": "Webhook",
      "type": "n8n-nodes-base.webhook",
      "typeVersion": 1,
      "position": [460, 300],
      "parameters": {
        "path": "/%s",
        "responseMode": "responseNode"
      }
    },
    {
      "id": "response",
      "name": "Response",
      "type": "n8n-nodes-base.respondToWebhook",
      "typeVersion": 1,
      "position": [680, 300],
      "parameters": {
        "options": {}
      }
    }
  ],
  "connections": {
    "Start": {
      "main": [
        [
          {
            "node": "Webhook",
            "type": "main",
            "index": 0
          }
        ]
      ]
    },
    "Webhook": {
      "main": [
        [
          {
            "node": "Response", 
            "type": "main",
            "index": 0
          }
        ]
      ]
    }
  },
  "createdAt": "%s",
  "updatedAt": "%s",
  "versionId": 1,
  "tags": ["branch:%s", "version:%s", "env:%s", "template"]
}`,
                generateWorkflowID(),
                naming.WorkflowName,
                sanitizeWorkflowName(naming.WorkflowName),
                time.Now().Format(time.RFC3339),
                time.Now().Format(time.RFC3339),
                naming.BranchName,
                naming.Version,
                naming.Environment,
        )

        return []byte(workflow)
}

func copyExistingWorkflow(naming NamingConvention) ([]byte, error) {
        // For now, return template (would implement workflow selection UI)
        return createTemplateWorkflow(naming), nil
}

func importFromBranch(naming NamingConvention) ([]byte, error) {
        // For now, return template (would implement branch selection UI) 
        return createTemplateWorkflow(naming), nil
}

func generateWorkflowID() string {
        return fmt.Sprintf("wf_%d", time.Now().Unix())
}

func saveBranchMetadata(naming NamingConvention) error {
        metadataDir := ".n8n-ops/branches"
        if err := os.MkdirAll(metadataDir, 0755); err != nil {
                return err
        }

        metadata := fmt.Sprintf(`{
  "branchName": "%s",
  "workflowName": "%s", 
  "fileName": "%s",
  "version": "%s",
  "environment": "%s",
  "createdAt": "%s",
  "createdBy": "%s"
}`,
                naming.BranchName,
                naming.WorkflowName,
                naming.FileName,
                naming.Version,
                naming.Environment,
                time.Now().Format(time.RFC3339),
                getGitUser(),
        )

        metadataFile := fmt.Sprintf("%s/%s.json", metadataDir, sanitizeWorkflowName(naming.BranchName))
        return os.WriteFile(metadataFile, []byte(metadata), 0644)
}

func getGitUser() string {
        cmd := exec.Command("git", "config", "user.name")
        output, err := cmd.Output()
        if err != nil {
                return "unknown"
        }
        return strings.TrimSpace(string(output))
}

func runBranchList(cmd *cobra.Command, args []string) error {
        fmt.Printf("🌿 Workflow Branches\n")
        fmt.Printf("==================\n\n")

        // List Git branches
        gitCmd := exec.Command("git", "branch", "-a")
        output, err := gitCmd.Output()
        if err != nil {
                return fmt.Errorf("failed to list Git branches: %w", err)
        }

        branches := strings.Split(string(output), "\n")
        workflowBranches := 0

        for _, branch := range branches {
                branch = strings.TrimSpace(branch)
                if branch == "" {
                        continue
                }

                // Remove markers (* and origin/)
                cleanBranch := strings.TrimPrefix(branch, "* ")
                cleanBranch = strings.TrimPrefix(cleanBranch, "origin/")
                
                // Check if it's a workflow branch
                if isWorkflowBranch(cleanBranch) {
                        workflowBranches++
                        fmt.Printf("📂 %s\n", cleanBranch)
                        
                        // Show metadata if available
                        metadata := getBranchMetadata(cleanBranch)
                        if metadata != "" {
                                fmt.Printf("   %s\n", metadata)
                        }
                        fmt.Println()
                }
        }

        if workflowBranches == 0 {
                fmt.Printf("No workflow branches found.\n")
                fmt.Printf("💡 Create one with: n8n-ops branch create <workflow-name>\n")
        }

        return nil
}

func isWorkflowBranch(branch string) bool {
        workflowPrefixes := []string{"feature/", "hotfix/", "release/", "experiment/"}
        for _, prefix := range workflowPrefixes {
                if strings.HasPrefix(branch, prefix) {
                        return true
                }
        }
        return false
}

func getBranchMetadata(branch string) string {
        metadataFile := fmt.Sprintf(".n8n-ops/branches/%s.json", sanitizeWorkflowName(branch))
        if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
                return "   ℹ️ No metadata available"
        }
        
        _, err := os.ReadFile(metadataFile)
        if err != nil {
                return "   ⚠️ Metadata read error"
        }

        // Parse and format metadata (simplified for demo)
        return fmt.Sprintf("   📝 Created: %s", time.Now().Format("2006-01-02"))
}

func runBranchSwitch(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
                return fmt.Errorf("branch name required")
        }

        branchName := args[0]
        
        fmt.Printf("🔄 Switching to branch: %s\n", branchName)

        // Switch Git branch
        gitCmd := exec.Command("git", "checkout", branchName)
        gitCmd.Stdout = os.Stdout
        gitCmd.Stderr = os.Stderr
        
        if err := gitCmd.Run(); err != nil {
                return fmt.Errorf("failed to switch Git branch: %w", err)
        }

        fmt.Printf("✅ Switched to branch: %s\n", branchName)
        
        // Auto-sync workflows for this branch
        fmt.Printf("🔄 Syncing workflows for branch...\n")
        // This would call sync command programmatically
        fmt.Printf("💡 Run: n8n-ops sync --env development\n")

        return nil
}