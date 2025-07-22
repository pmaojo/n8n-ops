package cmd

import (
        "fmt"
        "os"
        "path/filepath"
        "strings"

        "github.com/spf13/cobra"
        "github.com/n8n-workflows/cli/internal/utils"
        "github.com/n8n-workflows/cli/internal/workflow"
)

var validateCmd = &cobra.Command{
        Use:   "validate [file-or-directory...]",
        Short: "Validate workflow JSON files",
        Long: `Validate workflow JSON files for structure, syntax, and n8n compatibility.
If no arguments are provided, validates all workflow files in the current directory.

Examples:
  n8n-cli validate workflow.json         # Validate single file
  n8n-cli validate workflows/            # Validate directory
  n8n-cli validate --strict              # Enable strict validation mode`,
        RunE: runValidate,
}

var (
        strict    bool
        recursive bool
)

func init() {
        rootCmd.AddCommand(validateCmd)
        
        validateCmd.Flags().BoolVar(&strict, "strict", false, "enable strict validation mode")
        validateCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "validate files recursively")
}

func runValidate(cmd *cobra.Command, args []string) error {
        logger.Info("Starting workflow validation", "strict", strict, "recursive", recursive)

        // Determine files to validate
        var files []string
        var err error

        if len(args) == 0 {
                // Validate current directory
                files, err = getValidationFiles(".", recursive)
                if err != nil {
                        return fmt.Errorf("failed to get files from current directory: %w", err)
                }
        } else {
                // Validate specified files/directories
                for _, arg := range args {
                        if isDirectory(arg) {
                                dirFiles, err := getValidationFiles(arg, recursive)
                                if err != nil {
                                        return fmt.Errorf("failed to get files from directory %s: %w", arg, err)
                                }
                                files = append(files, dirFiles...)
                        } else {
                                files = append(files, arg)
                        }
                }
        }

        if len(files) == 0 {
                logger.Info("No workflow files found to validate")
                return nil
        }

        logger.Info("Found workflow files to validate", "count", len(files))

        // Validate each file
        var validationErrors []ValidationError
        validCount := 0

        for _, file := range files {
                if err := validateWorkflowFile(file, strict); err != nil {
                        validationErrors = append(validationErrors, ValidationError{
                                File:  file,
                                Error: err,
                        })
                        logger.Error("Validation failed", "file", file, "error", err)
                } else {
                        validCount++
                        logger.Info("Validation passed", "file", file)
                }
        }

        // Print summary
        logger.Info("Validation completed",
                "total", len(files),
                "valid", validCount,
                "invalid", len(validationErrors),
        )

        if len(validationErrors) > 0 {
                fmt.Printf("\n❌ Validation Summary:\n")
                fmt.Printf("Total files: %d\n", len(files))
                fmt.Printf("Valid: %d\n", validCount)
                fmt.Printf("Invalid: %d\n\n", len(validationErrors))

                fmt.Println("Validation errors:")
                for _, ve := range validationErrors {
                        fmt.Printf("  %s: %v\n", ve.File, ve.Error)
                }

                return fmt.Errorf("validation failed for %d files", len(validationErrors))
        }

        fmt.Printf("\n✅ All %d workflow files are valid!\n", validCount)
        return nil
}

type ValidationError struct {
        File  string
        Error error
}

func getValidationFiles(dir string, recursive bool) ([]string, error) {
        var files []string

        if recursive {
                err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
                        if err != nil {
                                return err
                        }
                        if !info.IsDir() && isWorkflowFile(path) {
                                files = append(files, path)
                        }
                        return nil
                })
                return files, err
        } else {
                entries, err := os.ReadDir(dir)
                if err != nil {
                        return nil, err
                }

                for _, entry := range entries {
                        if !entry.IsDir() {
                                path := filepath.Join(dir, entry.Name())
                                if isWorkflowFile(path) {
                                        files = append(files, path)
                                }
                        }
                }
                return files, nil
        }
}

func isWorkflowFile(path string) bool {
        return strings.HasSuffix(strings.ToLower(path), ".json") && !strings.HasPrefix(filepath.Base(path), "_")
}

func isDirectory(path string) bool {
        info, err := os.Stat(path)
        return err == nil && info.IsDir()
}

func validateWorkflowFile(file string, strict bool) error {
        // Basic file validation
        if err := workflow.ValidateWorkflowFile(file); err != nil {
                return err
        }

        if strict {
                // Load workflow for additional validation
                wf, err := utils.LoadWorkflowFromFile(file)
                if err != nil {
                        return fmt.Errorf("failed to load workflow: %w", err)
                }

                // Strict validation checks
                if err := workflow.ValidateWorkflowStrict(wf); err != nil {
                        return fmt.Errorf("strict validation failed: %w", err)
                }
        }

        return nil
}
