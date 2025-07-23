package tutorial

import (
	"fmt"
	"time"

	"github.com/pmaojo/n8n-ops/internal/ascii"
	"github.com/pmaojo/n8n-ops/internal/utils"
)

// ShowIntro displays the animated tutorial introduction.
func ShowIntro() {
	utils.ClearTerminalScreen()
	fmt.Print(ascii.TutorialWelcome())
	time.Sleep(time.Second)
	fmt.Println("\n🚀 Welcome to the n8n-ops Interactive Tutorial!")
	fmt.Println("This tutorial will guide you through using n8n-ops step by step.")
	time.Sleep(time.Second)
}
