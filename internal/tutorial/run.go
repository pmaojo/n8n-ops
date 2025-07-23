package tutorial

import (
	"fmt"

	"github.com/pmaojo/n8n-ops/internal/ascii"
)

// Run executes the interactive tutorial.
func Run(skipIntro, advanced bool) {
	if !skipIntro {
		ShowIntro()
	}

	fmt.Println(ascii.Banner("tutorial"))

	if !ConfigExists() {
		fmt.Printf("%s\n", ascii.ErrorMessage("No configuration found. Let's set up your environment first."))
		fmt.Println("Run 'n8n-ops onboard' to create your configuration.")
		return
	}

	runMainTutorial(advanced)
}
