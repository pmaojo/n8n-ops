package main

import (
        "os"

        "github.com/n8n-workflows/cli/cmd"
        "github.com/n8n-workflows/cli/internal/utils"
)

func main() {
        // Initialize logger
        logger := utils.NewLogger()
        
        // Execute CLI
        if err := cmd.Execute(); err != nil {
                logger.WithField("error", err).Error("Application failed to start")
                os.Exit(1)
        }
}
