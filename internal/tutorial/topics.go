package tutorial

import (
	"fmt"

	"github.com/pmaojo/n8n-ops/internal/ascii"
	"github.com/pmaojo/n8n-ops/internal/utils"
)

// topic represents a selectable tutorial section.
type topic struct {
	title       string
	description string
	action      func()
}

func runMainTutorial(advanced bool) {
	topics := []topic{
		{
			title:       "Understanding n8n-ops",
			description: "Learn about the core concepts of n8n-ops",
			action:      ShowConceptsTutorial,
		},
		{
			title:       "Syncing Workflows",
			description: "Learn how to sync workflows between n8n and your local files",
			action:      ShowSyncTutorial,
		},
		{
			title:       "Deploying Workflows",
			description: "Learn how to deploy workflows to different environments",
			action:      ShowDeployTutorial,
		},
		{
			title:       "Validating Workflows",
			description: "Learn how to validate your workflow files",
			action:      ShowValidateTutorial,
		},
		{
			title:       "Git Integration",
			description: "Learn how to use n8n-ops with Git",
			action:      ShowGitTutorial,
		},
	}

	if advanced {
		advancedTopics := []topic{
			{
				title:       "Multi-Environment Setup",
				description: "Learn how to manage multiple environments",
				action:      ShowMultiEnvTutorial,
			},
			{
				title:       "CI/CD Integration",
				description: "Learn how to integrate with GitLab CI/CD",
				action:      ShowCICDTutorial,
			},
			{
				title:       "Monitoring & Alerts",
				description: "Learn how to monitor workflows and set up alerts",
				action:      ShowMonitoringTutorial,
			},
		}
		topics = append(topics, advancedTopics...)
	}

	currentTopic := 0
	for {
		utils.ClearTerminalScreen()
		fmt.Println(ascii.Banner("tutorial"))
		fmt.Printf("\n%s%s📚 TUTORIAL TOPICS%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

		for i, topic := range topics {
			if i == currentTopic {
				fmt.Printf("%s%s> %d. %s%s\n", ascii.Bold, ascii.Green, i+1, topic.title, ascii.Reset)
				fmt.Printf("   %s%s%s\n", ascii.Dim, topic.description, ascii.Reset)
			} else {
				fmt.Printf("  %d. %s\n", i+1, topic.title)
			}
		}

		fmt.Printf("\n%s%sUse arrow keys to navigate, Enter to select, Q to quit%s\n", ascii.Bold, ascii.Cyan, ascii.Reset)

		key := WaitForKey()
		switch key {
		case "up":
			if currentTopic > 0 {
				currentTopic--
			}
		case "down":
			if currentTopic < len(topics)-1 {
				currentTopic++
			}
		case "enter":
			utils.ClearTerminalScreen()
			topics[currentTopic].action()
			WaitForEnter()
		case "q":
			fmt.Println("\nExiting tutorial. Run 'n8n-ops tutorial' anytime to continue learning!")
			return
		}
	}
}
