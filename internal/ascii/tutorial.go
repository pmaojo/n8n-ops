package ascii

import (
	"fmt"
	"strings"
)

// TutorialWelcome returns a visually appealing tutorial welcome screen
func TutorialWelcome() string {
	return fmt.Sprintf(`%s%s
╔══════════════════════════════════════════════════════════════════════╗
║                                                                      ║
║  %s⚡ WELCOME TO THE N8N-OPS INTERACTIVE TUTORIAL ⚡%s                  ║
║                                                                      ║
║  %sLearn how to use n8n-ops like a pro!%s                             ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝%s
`, Blue, Bold, Yellow, Blue, Cyan, Blue, Reset)
}

// TutorialHeader returns a formatted header for tutorial sections
func TutorialHeader(title string) string {
	return fmt.Sprintf(`%s%s
╔══════════════════════════════════════════════════════════════════════╗
║  %s%s%s                                                               ║
╚══════════════════════════════════════════════════════════════════════╝%s
`, Blue, Bold, Yellow, strings.ToUpper(title), Blue, Reset)
}

// WorkflowDiagram returns a visual diagram of the n8n-ops workflow
func WorkflowDiagram() string {
	return fmt.Sprintf(`%s
┌───────────────┐         ┌───────────────┐         ┌───────────────┐
│               │         │               │         │               │
│  %sn8n Instance%s │  %ssync%s   │  %sLocal Files%s  │ %sdeploy%s │  %sn8n Instance%s │
│  %s(Development)%s │ ◄─────► │    %s(Git)%s     │ ◄─────► │  %s(Production)%s │
│               │         │               │         │               │
└───────────────┘         └───────────────┘         └───────────────┘
                                  │
                                  │ %svalidate%s
                                  ▼
                          ┌───────────────┐
                          │               │
                          │  %sCI/CD Tests%s  │
                          │               │
                          └───────────────┘%s
`, Dim,
		Green, Dim, Yellow, Dim, Cyan, Dim, Yellow, Dim, Green, Dim,
		Blue, Dim, Purple, Dim, Blue, Dim,
		Yellow, Dim,
		Red, Dim,
		Reset)
}

// QuickStartGuide returns a visual quick start guide
func QuickStartGuide() string {
	return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  %s🚀 QUICK START GUIDE%s                                             │
│                                                                     │
│  %s1. Configure your environment%s                                    │
│     %sn8n-ops onboard%s                                               │
│                                                                     │
│  %s2. Sync workflows from n8n%s                                       │
│     %sn8n-ops sync --env development%s                                │
│                                                                     │
│  %s3. Validate your workflows%s                                       │
│     %sn8n-ops validate ./workflows/development/%s                     │
│                                                                     │
│  %s4. Deploy to another environment%s                                 │
│     %sn8n-ops deploy --env staging%s                                  │
│                                                                     │
│  %s5. Monitor workflow status%s                                       │
│     %sn8n-ops status --env production%s                               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘%s
`, Green, Bold,
		Yellow, Green,
		Cyan, Green,
		Dim, Green,
		Cyan, Green,
		Dim, Green,
		Cyan, Green,
		Dim, Green,
		Cyan, Green,
		Dim, Green,
		Cyan, Green,
		Dim, Green,
		Reset)
}

// TutorialTip returns a formatted tip box
func TutorialTip(tip string) string {
	return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│  %s💡 TIP%s                                                           │
│                                                                     │
│  %s%s%s                                                               │
└─────────────────────────────────────────────────────────────────────┘%s
`, Yellow, Bold,
		Green, Yellow,
		Dim, tip, Yellow,
		Reset)
}

// TutorialWarning returns a formatted warning box
func TutorialWarning(warning string) string {
	return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│  %s⚠️  WARNING%s                                                       │
│                                                                     │
│  %s%s%s                                                               │
└─────────────────────────────────────────────────────────────────────┘%s
`, Red, Bold,
		Yellow, Red,
		Dim, warning, Red,
		Reset)
}

// TutorialExample returns a formatted example box
func TutorialExample(title, content string) string {
	return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│  %s📝 EXAMPLE: %s%s                                                    │
│                                                                     │
%s%s%s                                                                   │
└─────────────────────────────────────────────────────────────────────┘%s
`, Cyan, Bold,
		Yellow, title, Cyan,
		Dim, formatExampleContent(content), Cyan,
		Reset)
}

// Helper function to format example content with proper indentation
func formatExampleContent(content string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder

	for _, line := range lines {
		result.WriteString("│  " + line + "\n")
	}

	return result.String()
}

// TutorialProgress returns a progress indicator
func TutorialProgress(current, total int) string {
	percentage := float64(current) / float64(total) * 100
	filled := int(percentage / 5) // 20 characters total
	bar := strings.Repeat("█", filled) + strings.Repeat("░", 20-filled)

	return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  %sTutorial Progress%s                                                │
│                                                                     │
│  [%s%s%s] %s%.0f%%%s                                                    │
│  %sStep %d of %d%s                                                      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘%s
`, Purple, Bold,
		Cyan, Purple,
		Green, bar, Purple, Yellow, percentage, Purple,
		Dim, current, total, Purple,
		Reset)
}

// TutorialComplete returns a completion message
func TutorialComplete() string {
	return fmt.Sprintf(`%s%s
╔══════════════════════════════════════════════════════════════════════╗
║                                                                      ║
║  %s✨ TUTORIAL COMPLETE! ✨%s                                          ║
║                                                                      ║
║  %sYou're now ready to use n8n-ops like a pro!%s                      ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝%s
`, Green, Bold, Yellow, Green, Cyan, Green, Reset)
}

// TutorialCommandReference returns a command reference box
func TutorialCommandReference() string {
	return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  %s📚 COMMAND REFERENCE%s                                             │
│                                                                     │
│  %sCore Commands:%s                                                   │
│  %s• n8n-ops sync --env <env>%s - Sync workflows from n8n             │
│  %s• n8n-ops deploy --env <env>%s - Deploy workflows to n8n           │
│  %s• n8n-ops validate <path>%s - Validate workflow files              │
│  %s• n8n-ops status --env <env>%s - Check workflow status             │
│                                                                     │
│  %sAdvanced Commands:%s                                               │
│  %s• n8n-ops monitor --env <env>%s - Monitor workflows                │
│  %s• n8n-ops branch list%s - List branch mappings                     │
│  %s• n8n-ops check --env <env>%s - Check for changes                  │
│  %s• n8n-ops credentials --env <env>%s - Manage credentials           │
│                                                                     │
│  %sGet help for any command:%s                                        │
│  %s• n8n-ops <command> --help%s                                       │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘%s
`, Blue, Bold,
		Yellow, Blue,
		Cyan, Blue,
		Green, Blue,
		Green, Blue,
		Green, Blue,
		Green, Blue,
		Cyan, Blue,
		Green, Blue,
		Green, Blue,
		Green, Blue,
		Green, Blue,
		Cyan, Blue,
		Green, Blue,
		Reset)
}
