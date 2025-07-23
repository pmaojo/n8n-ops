package ascii

import (
	"fmt"
	"strings"
)

// OnboardingWelcome returns a visually appealing onboarding welcome screen
func OnboardingWelcome() string {
	return fmt.Sprintf(`%s%s
╔══════════════════════════════════════════════════════════════════════╗
║                                                                      ║
║  %s⚡ WELCOME TO N8N-OPS ONBOARDING ⚡%s                               ║
║                                                                      ║
║  %sLet's get you set up in just a few minutes!%s                      ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝%s
`, Purple, Bold, Yellow, Purple, Cyan, Purple, Reset)
}

// OnboardingStep returns a formatted step header for the onboarding process
func OnboardingStep(stepNumber int, title string) string {
	return fmt.Sprintf(`%s%s
╔══════════════════════════════════════════════════════════════════════╗
║  %sSTEP %d: %s%s                                                       ║
╚══════════════════════════════════════════════════════════════════════╝%s
`, Blue, Bold, Yellow, stepNumber, strings.ToUpper(title), Blue, Reset)
}

// OnboardingComplete returns a success message for completing onboarding
func OnboardingComplete() string {
	return fmt.Sprintf(`%s%s
╔══════════════════════════════════════════════════════════════════════╗
║                                                                      ║
║  %s✨ ONBOARDING COMPLETE! ✨%s                                        ║
║                                                                      ║
║  %sYou're ready to start using n8n-ops!%s                             ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝%s
`, Green, Bold, Yellow, Green, Cyan, Green, Reset)
}

// OnboardingVisualGuide returns a visual guide for the n8n-ops workflow
func OnboardingVisualGuide() string {
        return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  %sn8n-ops Workflow:%s                                                │
│                                                                     │
│  %s┌───────────┐         ┌───────────┐         ┌───────────┐%s         │
│  %s│           │  sync   │           │ deploy  │           │%s         │
│  %s│    n8n    │ ◄─────► │   Files   │ ◄─────► │    n8n    │%s         │
│  %s│    Dev    │         │   (Git)   │         │   Prod    │%s         │
│  %s└───────────┘         └───────────┘         └───────────┘%s         │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘%s
`, Cyan, Bold,
   Yellow, Cyan,
   Blue, Cyan,
   Blue, Cyan,
   Blue, Cyan,
   Blue, Cyan,
   Blue, Cyan,
   Reset)
}

// OnboardingAPIKeyInstructions returns instructions for getting API keys
func OnboardingAPIKeyInstructions() string {
        return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  %sHow to get your n8n API keys:%s                                    │
│                                                                     │
│  %s1. Log in to your n8n instance%s                                   │
│  %s2. Go to Settings → API Keys%s                                     │
│  %s3. Click "Create API Key"%s                                        │
│  %s4. Name it "n8n-ops-development"%s                                 │
│  %s5. Copy the generated key (starts with "n8n_api_")%s               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘%s
`, Yellow, Bold,
   Green, Yellow,
   Cyan, Yellow,
   Cyan, Yellow,
   Cyan, Yellow,
   Cyan, Yellow,
   Cyan, Yellow,
   Reset)
}

// OnboardingNextSteps returns a guide for next steps after onboarding
func OnboardingNextSteps() string {
        return fmt.Sprintf(`%s%s
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  %sNext Steps:%s                                                      │
│                                                                     │
│  %s1. Sync workflows:%s                                               │
│     %sn8n-ops sync --env development%s                                │
│                                                                     │
│  %s2. Make changes to workflows%s                                     │
│     %s(in n8n UI or edit JSON files)%s                                │
│                                                                     │
│  %s3. Deploy your changes:%s                                          │
│     %sn8n-ops deploy --env development%s                              │
│                                                                     │
│  %s4. Commit to Git:%s                                                │
│     %sgit add workflows/ && git commit -m "update workflows"%s         │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘%s
`, Green, Bold,
   Yellow, Green,
   Cyan, Green,
   Cyan, Green,
   Dim, Green,
   Cyan, Green,
   Dim, Green,
   Cyan, Green,
   Dim, Green,
   Cyan, Green,
   Reset)
}

