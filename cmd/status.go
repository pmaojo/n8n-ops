package cmd

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "os"
        "path/filepath"
        "text/tabwriter"
        "time"

        "github.com/spf13/cobra"
        "github.com/n8n-workflows/n8n-ops/internal/ascii"
)

var statusCmd = &cobra.Command{
        Use:   "status",
        Short: "Show deployment status and workflow versions",
        Long: `Display the current status of workflows across environments, including:
- Deployed workflow versions
- Last deployment dates
- Git commit hashes
- Active/inactive status`,
        Run: runStatus,
}

var (
        showAll    bool
        outputJSON bool
)

func init() {
        rootCmd.AddCommand(statusCmd)
        statusCmd.Flags().BoolVar(&showAll, "all", false, "show all workflows, including inactive ones")
        statusCmd.Flags().BoolVar(&outputJSON, "json", false, "output in JSON format")
}

type WorkflowInfo struct {
        Name        string    `json:"name"`
        ID          string    `json:"id"`
        Active      bool      `json:"active"`
        Nodes       int       `json:"nodes"`
        Version     string    `json:"version"`
        LastDeploy  time.Time `json:"last_deploy"`
        GitCommit   string    `json:"git_commit"`
        Environment string    `json:"environment"`
}

func runStatus(cmd *cobra.Command, args []string) {
        logger.WithField("env", environment).Info("Getting workflow status")

        // Read local workflow files to get status
        workflowDir := fmt.Sprintf("./workflows/%s", environment)
        workflows, err := getLocalWorkflowStatus(workflowDir)
        if err != nil {
                fmt.Printf("❌ Error reading local workflows: %v\n", err)
                fmt.Printf("💡 Tip: Run 'n8n-ops sync --env %s' first to download workflows\n", environment)
                return
        }

        if outputJSON {
                printStatusJSON(workflows)
        } else {
                printStatusTable(workflows)
        }
}

func getLocalWorkflowStatus(workflowDir string) ([]WorkflowInfo, error) {
        var workflows []WorkflowInfo

        if _, err := os.Stat(workflowDir); os.IsNotExist(err) {
                return workflows, fmt.Errorf("workflow directory %s does not exist", workflowDir)
        }

        files, err := ioutil.ReadDir(workflowDir)
        if err != nil {
                return workflows, err
        }

        for _, file := range files {
                if filepath.Ext(file.Name()) != ".json" {
                        continue
                }

                filePath := filepath.Join(workflowDir, file.Name())
                content, err := ioutil.ReadFile(filePath)
                if err != nil {
                        continue
                }

                var workflowData map[string]interface{}
                if err := json.Unmarshal(content, &workflowData); err != nil {
                        continue
                }

                // Extract version from metadata or generate from git
                version := getStringField(workflowData, "version", "")
                if version == "" {
                        // Try to get git-based version
                        version = getGitVersion()
                }
                if version == "" {
                        version = "local"
                }

                workflow := WorkflowInfo{
                        Name:        getStringField(workflowData, "name", file.Name()),
                        ID:          getStringField(workflowData, "id", "unknown"),
                        Active:      getBoolField(workflowData, "active", true),
                        Environment: environment,
                        Version:     version,
                        LastDeploy:  file.ModTime(),
                        GitCommit:   getGitCommit(),
                }

                // Count nodes
                if nodes, ok := workflowData["nodes"].([]interface{}); ok {
                        workflow.Nodes = len(nodes)
                }

                workflows = append(workflows, workflow)
        }

        return workflows, nil
}

func getStringField(data map[string]interface{}, key, defaultValue string) string {
        if val, ok := data[key].(string); ok {
                return val
        }
        return defaultValue
}

func getBoolField(data map[string]interface{}, key string, defaultValue bool) bool {
        if val, ok := data[key].(bool); ok {
                return val
        }
        return defaultValue
}

func printStatusTable(workflows []WorkflowInfo) {
        fmt.Print(ascii.SmallLogo())
        fmt.Printf("\n%sWorkflow Status - %s Environment%s\n", ascii.Bold, environment, ascii.Reset)
        fmt.Printf("%sLast updated: %s%s\n\n", ascii.Dim, time.Now().Format("2006-01-02 15:04:05"), ascii.Reset)

        if len(workflows) == 0 {
                fmt.Printf("%sNo workflows found in %s environment%s\n", ascii.Dim, environment, ascii.Reset)
                fmt.Printf("%sTip: Run 'n8n-ops sync --env %s' to download workflows%s\n", ascii.Dim, environment, ascii.Reset)
                return
        }

        // Create table writer
        w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
        fmt.Fprintf(w, "%sNAME\tSTATUS\tVERSION\tLAST MODIFIED\tNODES%s\n", 
                ascii.Bold, ascii.Reset)
        fmt.Fprintf(w, "────\t──────\t───────\t─────────────\t─────\n")

        for _, workflow := range workflows {
                if !showAll && !workflow.Active {
                        continue
                }
                
                status := "🟢 Active"
                if !workflow.Active {
                        status = "🔴 Inactive"
                }

                lastModified := workflow.LastDeploy.Format("2006-01-02 15:04")
                nodeCount := fmt.Sprintf("%d", workflow.Nodes)

                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
                        workflow.Name, status, workflow.Version, lastModified, nodeCount)
        }
        w.Flush()

        fmt.Printf("\n%sTotal workflows: %d%s\n", ascii.Dim, len(workflows), ascii.Reset)
        if !showAll {
                fmt.Printf("%sUse --all to show inactive workflows%s\n", ascii.Dim, ascii.Reset)
        }
}

func getGitVersion() string {
        // Try to get version from git tags
        return "v1.0.0" // Placeholder - would implement git tag reading
}

func getGitCommit() string {
        // Try to get current git commit hash
        return "local" // Placeholder - would implement git commit reading
}

func printStatusJSON(workflows []WorkflowInfo) {
        output, err := json.MarshalIndent(workflows, "", "  ")
        if err != nil {
                fmt.Printf("❌ Error formatting JSON: %v\n", err)
                return
        }
        fmt.Println(string(output))
}