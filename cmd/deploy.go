package cmd

import (
        "encoding/json"
        "fmt"
        "io/fs"
        "os"
        "path/filepath"
        "strings"

        "github.com/n8n-workflows/n8n-ops/internal/client"
        "github.com/spf13/cobra"
)

// deployCmd sincroniza los JSON locales con la instancia n8n.
//
// Ejemplos:
//   n8n-ops deploy --env staging        # dry-run
//   n8n-ops deploy -f workflows/staging/foo.json --force
var deployCmd = &cobra.Command{
        Use:   "deploy",
        Short: "Sube workflows locales a n8n",
        RunE:  runDeploy,
}

var (
        deployFile  string
        deployForce bool
)

func init() {
        rootCmd.AddCommand(deployCmd)
        deployCmd.Flags().StringVarP(&deployFile, "file", "f", "", "JSON individual a desplegar")
        deployCmd.Flags().BoolVarP(&deployForce, "force", "", false, "ejecutar (por defecto dry-run)")
}

func runDeploy(cmd *cobra.Command, _ []string) error {
        // 1. Resolver lista de archivos
        files, err := collectFiles()
        if err != nil {
                return err
        }

        // 2. Instanciar cliente
        apiURL, apiKey := os.Getenv("N8N_URL"), os.Getenv("N8N_API_KEY")
        var cli client.N8nClientInterface
        if demoMode {
                cli = client.NewDemoN8nClient()
        } else {
                if apiURL == "" || apiKey == "" {
                        return fmt.Errorf("N8N_URL / N8N_API_KEY no configurados")
                }
                cli = client.NewRealN8nClient(apiURL, apiKey)
                if err := cli.TestConnection(); err != nil {
                        return fmt.Errorf("conexión n8n falló: %w", err)
                }
        }

        // 3. Mapear workflows remotos por nombre
        remote, _ := cli.GetWorkflows()
        remoteMap := map[string]client.Workflow{}
        for _, w := range remote {
                remoteMap[strings.ToLower(w.Name)] = w
        }

        // 4. Recorrer archivos locales
        created, updated := 0, 0
        for _, path := range files {
                wf, name, err := readWorkflow(path)
                if err != nil {
                        return err
                }
                if r, ok := remoteMap[strings.ToLower(name)]; ok {
                        wf.ID = r.ID
                        if deployForce {
                                if _, err := cli.UpdateWorkflow(r.ID, &wf); err != nil {
                                        return err
                                }
                                updated++
                        }
                        logAction("Actualizar", name, deployForce)
                } else {
                        if deployForce {
                                if _, err := cli.CreateWorkflow(&wf); err != nil {
                                        return err
                                }
                                created++
                        }
                        logAction("Crear", name, deployForce)
                }
        }

        fmt.Printf("✔️  %d creados, %d actualizados (%s)\n",
                created, updated, ternary(deployForce, "real", "dry-run"))
        return nil
}

// Helpers --------------------------------------------------------------

func collectFiles() ([]string, error) {
        if deployFile != "" {
                return []string{deployFile}, nil
        }
        dir := filepath.Join("workflows", environment)
        var list []string
        return list, filepath.WalkDir(dir, func(p string, d fs.DirEntry, e error) error {
                if e != nil {
                        return e
                }
                if !d.IsDir() && strings.HasSuffix(p, ".json") {
                        list = append(list, p)
                }
                return nil
        })
}

func readWorkflow(path string) (client.Workflow, string, error) {
        data, err := os.ReadFile(path)
        if err != nil {
                return client.Workflow{}, "", err
        }
        var raw map[string]interface{}
        if err := json.Unmarshal(data, &raw); err != nil {
                return client.Workflow{}, "", err
        }
        name := raw["name"]
        if name == nil {
                base := filepath.Base(path)
                name = strings.TrimSuffix(base, ".json")
        }
        var wf client.Workflow
        b, _ := json.Marshal(raw)
        _ = json.Unmarshal(b, &wf)
        wf.Name = fmt.Sprint(name)
        return wf, wf.Name, nil
}

func logAction(action, name string, real bool) {
        if real {
                fmt.Printf("🚀 %s → %s\n", action, name)
        } else {
                fmt.Printf("ℹ️  [dry-run] %s → %s\n", action, name)
        }
}

func ternary(cond bool, a, b string) string {
        if cond {
                return a
        }
        return b
}