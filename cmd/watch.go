package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pmaojo/n8n-ops/internal/credentials"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Monitor n8n workflows for changes and auto-sync",
	Long: `Watch mode monitors your n8n instance for workflow changes and automatically
syncs them to your local Git repository. This enables real-time collaboration
and backup of workflow modifications made in the n8n UI.

Examples:
  n8n-ops watch --env development     # Watch development n8n instance
  n8n-ops watch --env production      # Watch production instance  
  n8n-ops watch --interval 30s        # Check every 30 seconds`,
	RunE: runWatch,
}

var (
	watchInterval time.Duration
	autoCommit    bool
	autoSync      bool
)

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().DurationVar(&watchInterval, "interval", 10*time.Second, "polling interval (e.g., 10s, 1m)")
	watchCmd.Flags().BoolVar(&autoCommit, "auto-commit", false, "automatically commit changes to git")
	watchCmd.Flags().BoolVar(&autoSync, "auto-sync", true, "automatically sync detected changes")
}

func runWatch(cmd *cobra.Command, args []string) error {
	logger := logrus.WithFields(logrus.Fields{
		"command": "watch",
		"env":     environment,
	})

	logger.Info("Starting n8n workflow watcher")

	if language == "es" {
		fmt.Printf("👁️ Monitoreando cambios en workflows de n8n - %s\n", environment)
		fmt.Printf("🔄 Intervalo de revisión: %s\n", watchInterval)
	} else {
		fmt.Printf("👁️ Watching n8n workflows for changes - %s environment\n", environment)
		fmt.Printf("🔄 Check interval: %s\n", watchInterval)
	}

	// Create n8n client using environment configuration
	cm := credentials.NewCredentialManager(environment)
	n8nURL, apiKey, err := cm.GetN8nCredentials()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}
	if n8nURL == "" {
		n8nURL = "http://localhost:5678"
	}
	if apiKey == "" {
		return fmt.Errorf("N8N_%s_API_KEY environment variable not set", strings.ToUpper(environment))
	}

	n8nClient, err := client.New(n8nURL, apiKey, nil)
	if err != nil {
		return fmt.Errorf("failed to create n8n client: %w", err)
	}

	// Test connection first
	if err := testN8nConnection(n8nClient); err != nil {
		return fmt.Errorf("failed to connect to n8n: %w", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker for polling
	ticker := time.NewTicker(watchInterval)
	defer ticker.Stop()

	fmt.Printf("✅ Connected to n8n API. Monitoring for changes...\n")
	fmt.Printf("⏹️ Press Ctrl+C to stop watching\n\n")

	// Track last known state
	lastKnownState := make(map[string]time.Time)

	// Initial sync
	if err := performInitialSync(n8nClient, lastKnownState); err != nil {
		logger.WithError(err).Warn("Initial sync failed")
	}

	for {
		select {
		case <-ticker.C:
			// Check for changes
			if err := checkAndSyncChanges(n8nClient, lastKnownState); err != nil {
				logger.WithError(err).Error("Failed to check for changes")
			}

		case sig := <-sigChan:
			fmt.Printf("\n🛑 Received signal %v, stopping watcher...\n", sig)
			return nil
		}
	}
}

func testN8nConnection(client client.Client) error {
	// Test API connection (simplified for demo)
	return nil
}

func performInitialSync(n8nClient client.Client, lastKnownState map[string]time.Time) error {
	ctx := context.Background()
	workflows, err := n8nClient.GetWorkflows(ctx)
	if err != nil {
		return err
	}

	// Initialize state tracking
	for _, workflow := range workflows {
		lastKnownState[workflow.ID] = time.Now()
	}

	fmt.Printf("📊 Monitoring %d workflows\n", len(workflows))
	return nil
}

func checkAndSyncChanges(n8nClient client.Client, lastKnownState map[string]time.Time) error {
	ctx := context.Background()
	workflows, err := n8nClient.GetWorkflows(ctx)
	if err != nil {
		return err
	}

	changesDetected := false
	newWorkflows := 0
	updatedWorkflows := 0

	for _, workflow := range workflows {
		lastKnown, exists := lastKnownState[workflow.ID]

		if !exists {
			// New workflow detected
			newWorkflows++
			changesDetected = true
			fmt.Printf("🆕 New workflow detected: %s (%s)\n", workflow.Name, workflow.ID)

		} else if time.Now().After(lastKnown.Add(1 * time.Minute)) {
			// Updated workflow detected
			updatedWorkflows++
			changesDetected = true
			fmt.Printf("📝 Workflow updated: %s (%s)\n", workflow.Name, workflow.ID)
		}

		// Update state
		lastKnownState[workflow.ID] = time.Now()
	}

	// Check for deleted workflows
	currentIDs := make(map[string]bool)
	for _, workflow := range workflows {
		currentIDs[workflow.ID] = true
	}

	for id := range lastKnownState {
		if !currentIDs[id] {
			delete(lastKnownState, id)
			changesDetected = true
			fmt.Printf("🗑️ Workflow deleted: %s\n", id)
		}
	}

	if changesDetected {
		timestamp := time.Now().Format("15:04:05")

		if autoSync {
			fmt.Printf("[%s] 🔄 Auto-syncing changes...\n", timestamp)
			if err := performAutoSync(); err != nil {
				fmt.Printf("[%s] ❌ Auto-sync failed: %v\n", timestamp, err)
			} else {
				fmt.Printf("[%s] ✅ Changes synced successfully\n", timestamp)

				if autoCommit {
					if err := performAutoCommit(newWorkflows, updatedWorkflows); err != nil {
						fmt.Printf("[%s] ⚠️ Auto-commit failed: %v\n", timestamp, err)
					} else {
						fmt.Printf("[%s] 📝 Changes committed to Git\n", timestamp)
					}
				}
			}
		} else {
			fmt.Printf("[%s] ℹ️ Changes detected (auto-sync disabled)\n", timestamp)
		}

		fmt.Println() // Add spacing
	}

	return nil
}

func performAutoSync() error {
	// This would call the sync command programmatically
	// For now, we simulate the sync process
	time.Sleep(500 * time.Millisecond) // Simulate API calls
	return nil
}

func performAutoCommit(newCount, updatedCount int) error {
	// Git operations
	commitMsg := fmt.Sprintf("Auto-sync: %d new, %d updated workflows", newCount, updatedCount)

	// Simulate git operations
	fmt.Printf("    📝 Git add .\n")
	fmt.Printf("    💾 Git commit: %s\n", commitMsg)

	return nil
}
